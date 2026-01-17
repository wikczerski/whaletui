package dockerssh

import (
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
	"time"
)

// checkExistingDockerTCPPort checks if Docker is already listening on a TCP port
func (s *SSHTunnelClient) checkExistingDockerTCPPort() error {
	s.log.Info("Checking for existing Docker TCP port on remote machine")

	// Check for common Docker TCP ports
	commonPorts := []int{2375, 2376, 2377, 2378, 2379}
	for _, port := range commonPorts {
		// Try standard tools first
		if err := s.checkPortWithStandardTools(port); err == nil {
			s.remoteLocalPort = port
			s.log.Info("Found existing Docker TCP port", "port", port)
			return nil
		}

		// Fallback: Check /proc/net/tcp
		if s.checkPortInProcNetTCP(port) {
			s.remoteLocalPort = port
			s.log.Info("Found existing Docker TCP port via /proc/net/tcp", "port", port)
			return nil
		}
	}

	return errors.New("no existing Docker TCP port found")
}

// checkPortWithStandardTools checks if a port is listening using standard tools
func (s *SSHTunnelClient) checkPortWithStandardTools(port int) error {
	session, err := s.client.NewSession()
	if err != nil {
		return err
	}
	defer func() { _ = session.Close() }()

	cmd := fmt.Sprintf("netstat -tln | grep ':%d ' || ss -tln | grep ':%d '", port, port)
	output, err := session.Output(cmd)
	if err == nil && len(output) > 0 {
		return nil
	}
	return errors.New("port not found")
}

// findAvailableRemotePort finds an available port on the remote machine
func (s *SSHTunnelClient) findAvailableRemotePort() (int, error) {
	// Try to find an available port starting from 2375, expanding the range
	for port := 2375; port <= 2400; port++ {
		// Check if port is NOT in use
		// Try standard tools first
		if s.isPortAvailableWithStandardTools(port) {
			s.log.Info("Found available port on remote machine", "port", port)
			return port, nil
		}

		// Fallback: Check /proc/net/tcp
		// If checkPortInProcNetTCP returns false, it means the port is NOT in use (or we couldn't read the file)
		if !s.checkPortInProcNetTCP(port) {
			// Double check that we can actually read /proc/net/tcp
			if s.canReadProcNetTCP() {
				s.log.Info("Found available port on remote machine via /proc/net/tcp", "port", port)
				return port, nil
			}
		}
	}

	return 0, errors.New("no available ports found in range 2375-2400")
}

// isPortAvailableWithStandardTools checks if a port is available using standard tools
func (s *SSHTunnelClient) isPortAvailableWithStandardTools(port int) bool {
	session, err := s.client.NewSession()
	if err != nil {
		return false
	}
	defer func() { _ = session.Close() }()

	cmd := fmt.Sprintf(
		"! (netstat -tln | grep -q ':%d ' || ss -tln | grep -q ':%d ') && echo 'available'",
		port, port)
	output, err := session.Output(cmd)
	return err == nil && strings.Contains(string(output), "available")
}

// verifyTCPPortCreation verifies that the TCP port was created successfully
func (s *SSHTunnelClient) verifyTCPPortCreation() error {
	// Try multiple verification methods
	if pid, err := s.verifyWithStandardTools(); err == nil {
		s.remotePID = pid
		return nil
	}

	// Fallback: Check /proc/net/tcp
	if s.checkPortInProcNetTCP(s.remoteLocalPort) {
		s.log.Info("TCP port verification successful via /proc/net/tcp",
			"port", s.remoteLocalPort)

		// Try to get PID if possible, but don't fail if we can't
		if pid, err := s.getRemoteProcessPID(); err == nil {
			s.remotePID = pid
		}
		return nil
	}

	return fmt.Errorf(
		"failed to verify TCP port creation - port %d not listening",
		s.remoteLocalPort)
}

// verifyWithStandardTools verifies TCP port creation using standard tools
func (s *SSHTunnelClient) verifyWithStandardTools() (string, error) {
	verifySession, err := s.client.NewSession()
	if err != nil {
		return "", err
	}
	defer func() { _ = verifySession.Close() }()

	verifyCmd := fmt.Sprintf(
		"netstat -tln | grep ':%d ' || ss -tln | grep ':%d ' || lsof -i :%d",
		s.remoteLocalPort, s.remoteLocalPort, s.remoteLocalPort)
	output, err := verifySession.Output(verifyCmd)
	if err == nil && len(output) > 0 {
		s.log.Info("TCP port verification successful",
			"port", s.remoteLocalPort,
			"output", string(output))
		return s.getRemoteProcessPID()
	}
	return "", errors.New("verification failed")
}

// getRemoteProcessPID gets the PID of the remote process for cleanup
func (s *SSHTunnelClient) getRemoteProcessPID() (string, error) {
	// Try multiple methods to find the process PID
	pidMethods := []string{
		fmt.Sprintf("lsof -ti:%d", s.remoteLocalPort),
		fmt.Sprintf("fuser %d/tcp 2>/dev/null", s.remoteLocalPort),
		fmt.Sprintf("ps aux | grep -E '(python|python3|bash).*%d' | "+
			"grep -v grep | awk '{print $2}' | head -1", s.remoteLocalPort),
		fmt.Sprintf("netstat -tlnp | grep ':%d ' | awk '{print $7}' | "+
			"cut -d'/' -f1 | head -1", s.remoteLocalPort),
		fmt.Sprintf("ss -tlnp | grep ':%d ' | awk '{print $6}' | "+
			"cut -d',' -f2 | cut -d'=' -f2 | head -1", s.remoteLocalPort),
	}

	for _, pidCmd := range pidMethods {
		cmdSession, err := s.client.NewSession()
		if err != nil {
			continue
		}

		pidOutput, err := cmdSession.Output(pidCmd)
		_ = cmdSession.Close()

		if err == nil && len(strings.TrimSpace(string(pidOutput))) > 0 {
			pid := strings.TrimSpace(string(pidOutput))
			if pid != "" && pid != "unknown" {
				s.log.Info("Found remote process PID", "pid", pid, "method", pidCmd)
				return pid, nil
			}
		}
	}

	return "unknown", nil
}

// handleTunnelConnections handles incoming connections and forwards them via SSH tunnel
func (s *SSHTunnelClient) handleTunnelConnections(listener net.Listener, localPort int) {
	defer func() {
		if err := listener.Close(); err != nil {
			s.log.Warn("Failed to close tunnel listener", "error", err)
		}
	}()

	s.log.Info("SSH tunnel listener started", "localPort", localPort)

	for {
		localConn, err := listener.Accept()
		if err != nil {
			s.log.Warn("Failed to accept local connection", "error", err)
			return
		}

		// Handle each connection in a separate goroutine
		go s.forwardTunnelConnection(localConn)
	}
}

// forwardTunnelConnection forwards a single connection through SSH tunnel to remote local port
func (s *SSHTunnelClient) forwardTunnelConnection(localConn net.Conn) {
	defer func() {
		if err := localConn.Close(); err != nil {
			s.log.Warn("Failed to close local connection", "error", err)
		}
	}()

	s.log.Debug("Forwarding connection through SSH tunnel", "remoteLocalPort", s.remoteLocalPort)

	// Connect to the remote local port through SSH tunnel
	remoteConn, err := s.client.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", s.remoteLocalPort))
	if err != nil {
		s.log.Error("Failed to connect to remote local port through SSH tunnel",
			"error", err,
			"remoteLocalPort", s.remoteLocalPort)
		return
	}
	defer func() {
		if err := remoteConn.Close(); err != nil {
			s.log.Warn("Failed to close remote connection", "error", err)
		}
	}()

	// Copy data bidirectionally between local and remote connections
	s.copyDataBidirectionally(localConn, remoteConn)
}

// copyDataBidirectionally copies data between local and remote connections
func (s *SSHTunnelClient) copyDataBidirectionally(localConn, remoteConn net.Conn) {
	// Copy data from remote to local in a goroutine
	go func() {
		defer func() {
			if err := localConn.Close(); err != nil {
				s.log.Warn("Failed to close local connection in goroutine", "error", err)
			}
		}()
		defer func() {
			if err := remoteConn.Close(); err != nil {
				s.log.Warn("Failed to close remote connection in goroutine", "error", err)
			}
		}()
		if _, err := io.Copy(localConn, remoteConn); err != nil {
			s.log.Warn("Error copying data from remote to local", "error", err)
		}
	}()

	// Copy data from local to remote in the main goroutine
	if _, err := io.Copy(remoteConn, localConn); err != nil {
		s.log.Warn("Error copying data from local to remote", "error", err)
	}
}

// cleanupRemoteProcess cleans up the remote process that was created
func (s *SSHTunnelClient) cleanupRemoteProcess() error {
	session, err := s.client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create cleanup session: %w", err)
	}
	defer func() {
		if err := session.Close(); err != nil {
			s.log.Warn("Failed to close cleanup session", "error", err)
		}
	}()

	// If we have a specific PID, try to kill it
	if s.remotePID != "" && s.remotePID != "unknown" {
		cmd := fmt.Sprintf("kill %s 2>/dev/null || true", s.remotePID)
		if err := session.Run(cmd); err != nil {
			s.log.Warn("Failed to kill remote process by PID", "pid", s.remotePID, "error", err)
		} else {
			s.log.Info("Cleaned up remote process by PID", "pid", s.remotePID)
		}
	}

	// Also try to kill any processes listening on our port using multiple methods
	cleanupCommands := []string{
		fmt.Sprintf("lsof -ti:%d | xargs kill -9 2>/dev/null || true", s.remoteLocalPort),
		fmt.Sprintf("fuser -k %d/tcp 2>/dev/null || true", s.remoteLocalPort),
		fmt.Sprintf("pkill -f 'python.*%d' 2>/dev/null || true", s.remoteLocalPort),
		fmt.Sprintf("pkill -f 'python3.*%d' 2>/dev/null || true", s.remoteLocalPort),
		fmt.Sprintf("pkill -f 'bash.*%d' 2>/dev/null || true", s.remoteLocalPort),
	}

	for _, cmd := range cleanupCommands {
		if err := session.Run(cmd); err != nil {
			s.log.Debug("Cleanup command failed", "cmd", cmd, "error", err)
		}
	}

	s.log.Info("Completed cleanup of remote processes", "port", s.remoteLocalPort)
	return nil
}

// checkPortInProcNetTCP checks if a port is listening in /proc/net/tcp
func (s *SSHTunnelClient) checkPortInProcNetTCP(port int) bool {
	session, err := s.client.NewSession()
	if err != nil {
		return false
	}
	defer func() { _ = session.Close() }()

	// Convert port to hex string (uppercase)
	hexPort := fmt.Sprintf("%04X", port)

	cmd := fmt.Sprintf(
		"grep -F ':%s ' /proc/net/tcp | awk '{print $4}' | grep -q '0A' && echo 'found'",
		hexPort,
	)

	output, err := session.Output(cmd)
	return err == nil && strings.Contains(string(output), "found")
}

// canReadProcNetTCP checks if /proc/net/tcp is readable
func (s *SSHTunnelClient) canReadProcNetTCP() bool {
	session, err := s.client.NewSession()
	if err != nil {
		return false
	}
	defer func() { _ = session.Close() }()

	cmd := "test -r /proc/net/tcp && echo 'readable'"
	output, err := session.Output(cmd)
	return err == nil && strings.Contains(string(output), "readable")
}

// connectToLocalPort connects to the local port for testing
func (s *SSHTunnelClient) connectToLocalPort(localPort int) (net.Conn, error) {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("127.0.0.1:%d", localPort), 5*time.Second)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to local port %d: %w", localPort, err)
	}
	return conn, nil
}

// sendTestRequest sends a test HTTP request and reads the response
func (s *SSHTunnelClient) sendTestRequest(conn net.Conn) (string, error) {
	testRequest := "GET /_ping HTTP/1.1\r\nHost: localhost\r\n\r\n"
	if _, err := conn.Write([]byte(testRequest)); err != nil {
		return "", fmt.Errorf("failed to write test request: %w", err)
	}

	// Read response with timeout
	if err := conn.SetReadDeadline(time.Now().Add(10 * time.Second)); err != nil {
		return "", fmt.Errorf("failed to set read deadline: %w", err)
	}
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("failed to read test response: %w", err)
	}

	response := string(buffer[:n])
	if len(response) == 0 {
		return "", errors.New("no response received from Docker socket")
	}

	return response, nil
}

// validateResponse validates the response from the Docker socket
func (s *SSHTunnelClient) validateResponse(response string, localPort int) error {
	// Check if we got a valid HTTP response (Docker API should return HTTP)
	if !strings.Contains(response, "HTTP/") {
		responsePreview := response
		if len(response) > 50 {
			responsePreview = response[:50]
		}
		return fmt.Errorf("invalid response from Docker socket: %s", responsePreview)
	}

	s.log.Info("SSH tunnel test successful",
		"localPort", localPort,
		"responseLength", len(response))
	return nil
}
