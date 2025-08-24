package dockerssh

import (
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

// setupSocatProxy sets up a socat proxy on the remote machine
func (s *SSHClient) setupSocatProxy(client *ssh.Client, _, remotePort int) (session *ssh.Session, pid string, port int, err error) {
	determinedPort, err := s.determineRemotePort(client, remotePort)
	if err != nil {
		return nil, "", 0, err
	}

	s.log.Info("Using remote port for socat proxy", "port", determinedPort)

	if err := s.validateRemoteEnvironment(client); err != nil {
		return nil, "", 0, err
	}

	pid, err = s.startSocatProxy(client, determinedPort)
	if err != nil {
		return nil, "", 0, err
	}

	if err := s.verifySocatProcess(client, pid, determinedPort); err != nil {
		return nil, "", 0, fmt.Errorf("socat process verification failed: %w", err)
	}

	s.log.Info("Socat verification successful on port", "port", determinedPort)

	session, err = client.NewSession()
	if err != nil {
		return nil, "", 0, fmt.Errorf("failed to create monitoring session: %w", err)
	}

	return session, pid, determinedPort, nil
}

// determineRemotePort finds an available port if none specified
func (s *SSHClient) determineRemotePort(client *ssh.Client, remotePort int) (int, error) {
	if remotePort == 0 {
		availablePort, err := s.findAvailableRemotePort(client)
		if err != nil {
			return 0, fmt.Errorf("failed to find available remote port: %w", err)
		}
		return availablePort, nil
	}
	return remotePort, nil
}

// validateRemoteEnvironment checks socat availability and Docker socket access
func (s *SSHClient) validateRemoteEnvironment(client *ssh.Client) error {
	if err := s.checkSocatAvailability(client); err != nil {
		return fmt.Errorf("socat not available on remote machine: %w", err)
	}

	if err := s.checkDockerSocketAccess(client); err != nil {
		return fmt.Errorf("docker socket not accessible: %w", err)
	}

	return nil
}

// startSocatProxy starts the socat proxy and returns the process ID
func (s *SSHClient) startSocatProxy(client *ssh.Client, remotePort int) (string, error) {
	// Start socat with a much simpler approach to avoid timeout issues
	s.log.Info("Starting socat proxy with simplified approach", "port", remotePort)

	// First, check if socat is available
	versionSession, err := client.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create version check session: %w", err)
	}
	defer versionSession.Close()

	if output, err := versionSession.Output("socat -V 2>/dev/null || echo 'socat not found'"); err != nil {
		s.log.Warn("Could not check socat version", "error", err)
	} else {
		s.log.Info("Socat version info", "version", strings.TrimSpace(string(output)))
	}

	// Clean up any existing socat processes
	cleanupSession, err := client.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create cleanup session: %w", err)
	}
	defer cleanupSession.Close()

	if err := cleanupSession.Run(fmt.Sprintf("pkill -f 'socat.*:%d' 2>/dev/null || true", remotePort)); err != nil {
		s.log.Warn("Failed to cleanup existing socat processes", "error", err)
	}

	// Start socat with a simple command
	startSession, err := client.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create socat start session: %w", err)
	}
	defer startSession.Close()

	// Use a very simple socat command
	socatCmd := fmt.Sprintf("nohup socat TCP-LISTEN:%d,bind=0.0.0.0,reuseaddr,fork UNIX-CONNECT:/var/run/docker.sock > /tmp/socat.log 2>&1 & echo $!", remotePort)

	s.log.Debug("Executing simple socat command", "command", socatCmd)
	output, err := startSession.Output(socatCmd)

	// Add debug logging to see what's happening
	s.log.Debug("Socat command execution result", "error", err, "outputLength", len(output))

	if err != nil {
		// Try to get more detailed error information
		s.log.Warn("Socat startup failed, checking logs", "error", err)

		logSession, logSessionErr := client.NewSession()
		if logSessionErr == nil {
			defer logSession.Close()
			if logOutput, logErr := logSession.Output("cat /tmp/socat.log 2>/dev/null || echo 'No socat log found'"); logErr == nil {
				s.log.Warn("Socat startup log", "log", string(logOutput))
			}
		}

		return "", fmt.Errorf("failed to start socat proxy: %w", err)
	}

	// Check the output for any error indicators
	outputStr := string(output)
	s.log.Debug("Socat startup output", "output", outputStr)

	pid := strings.TrimSpace(outputStr)
	if pid == "" {
		return "", fmt.Errorf("failed to get socat process ID from output: %s", outputStr)
	}

	s.log.Info("Socat started with PID", "pid", pid, "port", remotePort)

	// Give socat more time to fully start and bind to port
	s.log.Debug("Waiting for socat to fully initialize...")
	time.Sleep(5 * time.Second)

	// Verify the process is still running
	verifySession, err := client.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create verification session: %w", err)
	}
	defer verifySession.Close()

	if err := verifySession.Run(fmt.Sprintf("kill -0 %s", pid)); err != nil {
		return "", fmt.Errorf("socat process not running after startup (PID: %s): %w", pid, err)
	}

	// Check if port is listening
	portSession, err := client.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create port check session: %w", err)
	}
	defer portSession.Close()

	portCmd := fmt.Sprintf("ss -tln | grep -q ':%d ' || netstat -tln 2>/dev/null | grep -q ':%d ' || lsof -i :%d", remotePort, remotePort, remotePort)
	if err := portSession.Run(portCmd); err != nil {
		// Port is not listening, let's check what happened
		s.log.Warn("Port is not listening, checking socat logs", "port", remotePort, "pid", pid)

		// Check socat logs for more details
		logSession, err := client.NewSession()
		if err == nil {
			defer logSession.Close()
			if logOutput, logErr := logSession.Output("cat /tmp/socat.log 2>/dev/null || echo 'No socat log found'"); logErr == nil {
				s.log.Warn("Socat startup log", "log", string(logOutput))
			}
		}

		return "", fmt.Errorf("socat process verification failed: socat process not listening on port %d (PID: %s)", remotePort, pid)
	}

	s.log.Info("Socat verification successful - port is listening", "port", remotePort, "pid", pid)
	return pid, nil
}

// checkSocatAvailability checks if socat is available on the remote machine
func (s *SSHClient) checkSocatAvailability(client *ssh.Client) error {
	if client == nil {
		return fmt.Errorf("SSH client is nil")
	}

	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session for socat check: %w", err)
	}
	defer func() {
		if err := session.Close(); err != nil {
			s.log.Warn("Failed to close SSH session", "error", err)
		}
	}()

	cmd := "which socat"
	if err := session.Run(cmd); err != nil {
		return fmt.Errorf("socat not found on remote machine: %w", err)
	}

	return nil
}

// checkDockerSocketAccess checks if Docker socket is accessible on the remote machine
func (s *SSHClient) checkDockerSocketAccess(client *ssh.Client) error {
	if client == nil {
		return fmt.Errorf("SSH client is nil")
	}

	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session for Docker socket check: %w", err)
	}
	defer func() {
		if err := session.Close(); err != nil {
			s.log.Warn("Failed to close SSH session", "error", err)
		}
	}()

	cmd := "test -S /var/run/docker.sock && test -r /var/run/docker.sock"
	if err := session.Run(cmd); err != nil {
		return fmt.Errorf("docker socket not accessible at /var/run/docker.sock: %w", err)
	}

	return nil
}

// verifySocatProcess verifies that the socat process is running and listening
func (s *SSHClient) verifySocatProcess(client *ssh.Client, pid string, port int) error {
	if client == nil {
		return fmt.Errorf("SSH client is nil")
	}

	if pid == "" {
		return fmt.Errorf("process ID cannot be empty")
	}

	if port <= 0 {
		return fmt.Errorf("port must be positive, got %d", port)
	}

	s.log.Debug("Starting verification for PID", "pid", pid, "port", port)

	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session for process verification: %w", err)
	}
	defer func() {
		if err := session.Close(); err != nil {
			s.log.Warn("Failed to close SSH session", "error", err)
		}
	}()

	// First check if the process is running
	cmd := fmt.Sprintf("kill -0 %s", pid)
	s.log.Debug("Checking if process is running with command", "command", cmd)
	if err := session.Run(cmd); err != nil {
		// Process is not running, let's check what happened
		s.log.Warn("Process is not running, checking socat logs", "pid", pid)

		// Check socat logs to see what went wrong
		logSession, err := client.NewSession()
		if err == nil {
			defer logSession.Close()
			if output, logErr := logSession.Output("cat /tmp/socat.log 2>/dev/null || echo 'No socat log found'"); logErr == nil {
				s.log.Warn("Socat log output", "log", string(output))
			}
		}

		return fmt.Errorf("socat process not running (PID: %s): %w", pid, err)
	}
	s.log.Debug("Process is running")

	// Wait a bit more for socat to fully start
	time.Sleep(3 * time.Second)

	// Now check if the port is actually listening
	portSession, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session for port verification: %w", err)
	}
	defer func() {
		if err := portSession.Close(); err != nil {
			s.log.Warn("Failed to close port verification session", "error", err)
		}
	}()

	// Check if the port is listening
	portCmd := fmt.Sprintf("ss -tln | grep ':%d ' || netstat -tln | grep ':%d ' || lsof -i :%d", port, port, port)
	s.log.Debug("Checking if port is listening with command", "command", portCmd)
	if err := portSession.Run(portCmd); err != nil {
		// Port is not listening, let's check what happened
		s.log.Warn("Port is not listening, checking process status", "port", port, "pid", pid)

		// Check process status
		statusSession, err := client.NewSession()
		if err == nil {
			defer statusSession.Close()
			if output, statusErr := statusSession.Output(fmt.Sprintf("ps -p %s -o pid,ppid,cmd,stat 2>/dev/null || echo 'Process not found'", pid)); statusErr == nil {
				s.log.Warn("Process status", "status", string(output))
			}
		}

		return fmt.Errorf("socat process verification failed: socat process not listening on port %d (PID: %s): %w", port, pid, err)
	}

	s.log.Debug("Port is listening, socat process verification successful")
	return nil
}

// findAvailableRemotePort finds an available port on the remote machine
func (s *SSHClient) findAvailableRemotePort(client *ssh.Client) (int, error) {
	portRanges := [][]int{
		{2376, 2385},
		{2386, 2395},
		{2396, 2405},
	}

	for _, portRange := range portRanges {
		if len(portRange) != 2 || portRange[0] <= 0 || portRange[1] <= 0 {
			continue
		}

		for port := portRange[0]; port <= portRange[1]; port++ {
			if s.isRemotePortAvailable(client, port) {
				return port, nil
			}
		}
	}

	return 0, fmt.Errorf("no available ports found in ranges 2376-2385, 2386-2395, or 2396-2405")
}

// isRemotePortAvailable checks if a port is available on the remote machine
func (s *SSHClient) isRemotePortAvailable(client *ssh.Client, port int) bool {
	session, err := client.NewSession()
	if err != nil {
		return false
	}
	defer func() {
		if err := session.Close(); err != nil {
			s.log.Warn("Failed to close SSH session", "error", err)
		}
	}()

	cmd := fmt.Sprintf("ss -tln | grep ':%d '", port)
	return session.Run(cmd) != nil
}
