package dockerssh

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

// setupSocatProxy sets up a socat proxy on the remote machine
func (s *SSHClient) setupSocatProxy(
	client *ssh.Client,
	_, remotePort int,
) (session *ssh.Session, pid string, port int, err error) {
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

	return s.createMonitoringSession(client, pid, determinedPort)
}

// createMonitoringSession creates the monitoring session after socat setup
func (s *SSHClient) createMonitoringSession(
	client *ssh.Client,
	pid string,
	determinedPort int,
) (*ssh.Session, string, int, error) {
	if err := s.verifySocatProcess(client, pid, determinedPort); err != nil {
		return nil, "", 0, fmt.Errorf("socat process verification failed: %w", err)
	}

	s.log.Info("Socat verification successful on port", "port", determinedPort)

	session, err := client.NewSession()
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
	s.log.Info("Starting socat proxy with simplified approach", "port", remotePort)

	s.prepareSocatEnvironment(client, remotePort)

	pid, err := s.startSocatProcess(client, remotePort)
	if err != nil {
		return "", err
	}

	return s.verifyAndReturnPid(client, pid, remotePort)
}

// prepareSocatEnvironment prepares the environment for socat startup
func (s *SSHClient) prepareSocatEnvironment(client *ssh.Client, remotePort int) {
	if err := s.checkSocatVersion(client); err != nil {
		s.log.Warn("Could not check socat version", "error", err)
	}

	if err := s.cleanupExistingSocatProcesses(client, remotePort); err != nil {
		s.log.Warn("Failed to cleanup existing socat processes", "error", err)
	}
}

// verifyAndReturnPid verifies the socat process and returns the PID
func (s *SSHClient) verifyAndReturnPid(
	client *ssh.Client,
	pid string,
	remotePort int,
) (string, error) {
	if err := s.verifySocatProcess(client, pid, remotePort); err != nil {
		return "", err
	}

	s.log.Info("Socat verification successful - port is listening", "port", remotePort, "pid", pid)
	return pid, nil
}

// checkSocatVersion checks if socat is available and its version
func (s *SSHClient) checkSocatVersion(client *ssh.Client) error {
	versionSession, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create version check session: %w", err)
	}
	defer func() {
		if err := versionSession.Close(); err != nil {
			s.log.Warn("Failed to close version check session", "error", err)
		}
	}()

	if output, err := versionSession.Output("socat -V 2>/dev/null || echo 'socat not found'"); err != nil {
		return err
	} else {
		s.log.Info("Socat version info", "version", strings.TrimSpace(string(output)))
	}
	return nil
}

// cleanupExistingSocatProcesses cleans up any existing socat processes
func (s *SSHClient) cleanupExistingSocatProcesses(client *ssh.Client, remotePort int) error {
	cleanupSession, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create cleanup session: %w", err)
	}
	defer func() {
		if err := cleanupSession.Close(); err != nil {
			s.log.Warn("Failed to close cleanup session", "error", err)
		}
	}()

	cleanupCmd := fmt.Sprintf("pkill -f 'socat.*:%d' 2>/dev/null || true", remotePort)
	return cleanupSession.Run(cleanupCmd)
}

// startSocatProcess starts the socat process and returns the PID
func (s *SSHClient) startSocatProcess(client *ssh.Client, remotePort int) (string, error) {
	startSession, err := client.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create socat start session: %w", err)
	}
	defer func() {
		if err := startSession.Close(); err != nil {
			s.log.Warn("Failed to close socat start session", "error", err)
		}
	}()

	return s.executeSocatStartupCommand(startSession, remotePort, client)
}

// executeSocatStartupCommand executes the socat command and returns the PID
func (s *SSHClient) executeSocatStartupCommand(
	startSession *ssh.Session,
	remotePort int,
	client *ssh.Client,
) (string, error) {
	socatCmd := s.buildSocatCommand(remotePort)

	s.log.Debug("Executing simple socat command", "command", socatCmd)
	output, err := startSession.Output(socatCmd)

	s.log.Debug("Socat command execution result", "error", err, "outputLength", len(output))

	if err != nil {
		s.logSocatStartupFailure(client)
		return "", fmt.Errorf("failed to start socat proxy: %w", err)
	}

	return s.extractAndValidatePid(output, remotePort)
}

// buildSocatCommand builds the socat command string
func (s *SSHClient) buildSocatCommand(remotePort int) string {
	return fmt.Sprintf(
		"nohup socat TCP-LISTEN:%d,bind=0.0.0.0,reuseaddr,fork "+
			"UNIX-CONNECT:/var/run/docker.sock > /tmp/socat.log 2>&1 & echo $!",
		remotePort)
}

// extractAndValidatePid extracts and validates the PID from the command output
func (s *SSHClient) extractAndValidatePid(output []byte, remotePort int) (string, error) {
	pid := strings.TrimSpace(string(output))
	if pid == "" {
		return "", fmt.Errorf("failed to get socat process ID from output: %s", string(output))
	}

	s.log.Info("Socat started with PID", "pid", pid, "port", remotePort)
	return pid, nil
}

// logSocatStartupFailure logs socat startup failure details
func (s *SSHClient) logSocatStartupFailure(client *ssh.Client) {
	s.log.Warn("Socat startup failed, checking logs")

	logSession, logSessionErr := client.NewSession()
	if logSessionErr == nil {
		defer func() {
			if err := logSession.Close(); err != nil {
				s.log.Warn("Failed to close log session", "error", err)
			}
		}()
		logCmd := "cat /tmp/socat.log 2>/dev/null || echo 'No socat log found'"
		if logOutput, logErr := logSession.Output(logCmd); logErr == nil {
			s.log.Warn("Socat startup log", "log", string(logOutput))
		}
	}
}

// verifySocatProcess verifies that the socat process is running and listening
func (s *SSHClient) verifySocatProcess(client *ssh.Client, pid string, remotePort int) error {
	// Give socat more time to fully start and bind to port
	s.log.Debug("Waiting for socat to fully initialize...")
	time.Sleep(5 * time.Second)

	if err := s.verifyProcessRunning(client, pid); err != nil {
		return err
	}

	return s.verifyPortListening(client, pid, remotePort)
}

// verifyProcessRunning verifies that the socat process is still running
func (s *SSHClient) verifyProcessRunning(client *ssh.Client, pid string) error {
	verifySession, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create verification session: %w", err)
	}
	defer func() {
		if err := verifySession.Close(); err != nil {
			s.log.Warn("Failed to close verification session", "error", err)
		}
	}()

	return s.runProcessVerification(verifySession, pid)
}

// runProcessVerification runs the actual process verification command
func (s *SSHClient) runProcessVerification(session *ssh.Session, pid string) error {
	verifyCmd := fmt.Sprintf("kill -0 %s", pid)
	if err := session.Run(verifyCmd); err != nil {
		return fmt.Errorf("socat process not running after startup (PID: %s): %w", pid, err)
	}

	s.log.Debug("Socat process is running", "pid", pid)
	return nil
}

// verifyPortListening verifies that the socat process is listening on the port
func (s *SSHClient) verifyPortListening(client *ssh.Client, pid string, remotePort int) error {
	portSession, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create port verification session: %w", err)
	}
	defer func() {
		if err := portSession.Close(); err != nil {
			s.log.Warn("Failed to close port verification session", "error", err)
		}
	}()

	return s.runPortListeningVerification(portSession, pid, remotePort)
}

// runPortListeningVerification runs the actual port listening verification
func (s *SSHClient) runPortListeningVerification(
	session *ssh.Session,
	pid string,
	remotePort int,
) error {
	portCheckCmd := fmt.Sprintf(
		"lsof -i:%d 2>/dev/null || netstat -tlnp 2>/dev/null | grep ':%d ' || ss -tlnp 2>/dev/null | grep ':%d '",
		remotePort,
		remotePort,
		remotePort,
	)
	if err := session.Run(portCheckCmd); err != nil {
		return fmt.Errorf(
			"socat process not listening on port %d (PID: %s): %w",
			remotePort,
			pid,
			err,
		)
	}

	s.log.Debug("Socat process is listening on port", "pid", pid, "port", remotePort)
	return nil
}

// checkSocatAvailability checks if socat is available on the remote machine
func (s *SSHClient) checkSocatAvailability(client *ssh.Client) error {
	if client == nil {
		return errors.New("socat not found on remote machine")
	}

	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create socat availability check session: %w", err)
	}
	defer func() {
		if err := session.Close(); err != nil {
			s.log.Warn("Failed to close socat availability check session", "error", err)
		}
	}()

	return session.Run("which socat >/dev/null 2>&1")
}

// checkDockerSocketAccess checks if the Docker socket is accessible
func (s *SSHClient) checkDockerSocketAccess(client *ssh.Client) error {
	if client == nil {
		return errors.New("SSH client is nil")
	}

	session, err := s.createDockerSocketCheckSession(client)
	if err != nil {
		return err
	}

	return s.runDockerSocketCheck(session)
}

// createDockerSocketCheckSession creates a session to check Docker socket access
func (s *SSHClient) createDockerSocketCheckSession(client *ssh.Client) (*ssh.Session, error) {
	session, err := client.NewSession()
	if err != nil {
		return nil, fmt.Errorf("failed to create Docker socket check session: %w", err)
	}
	return session, nil
}

// runDockerSocketCheck runs the Docker socket accessibility check
func (s *SSHClient) runDockerSocketCheck(session *ssh.Session) error {
	defer func() {
		if err := session.Close(); err != nil {
			s.log.Warn("Failed to close Docker socket check session", "error", err)
		}
	}()

	cmd := "test -S /var/run/docker.sock && test -r /var/run/docker.sock"
	if err := session.Run(cmd); err != nil {
		return fmt.Errorf("docker socket not accessible at /var/run/docker.sock: %w", err)
	}
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

	return 0, errors.New("no available ports found in ranges 2376-2385, 2386-2395, or 2396-2405")
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
