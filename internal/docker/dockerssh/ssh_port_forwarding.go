package dockerssh

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

// SSHTunnelClient handles SSH tunneling with local TCP port creation on remote machine
type SSHTunnelClient struct {
	client           *ssh.Client
	log              *slog.Logger
	listener         net.Listener
	remoteLocalPort  int
	localPort        int
	remotePID        string
	connectionMethod string // Track which method was used (primary or fallback)
}

// NewSSHTunnelClient creates a new SSH tunnel client
func NewSSHTunnelClient(client *ssh.Client, log *slog.Logger) *SSHTunnelClient {
	return &SSHTunnelClient{
		client: client,
		log:    log,
	}
}

// SetupSSHTunnel sets up SSH tunneling with local TCP port creation on remote machine
func (s *SSHTunnelClient) SetupSSHTunnel(localPort int) error {
	if s.client == nil {
		return errors.New("SSH client is nil")
	}

	s.localPort = localPort
	s.log.Info("Setting up SSH tunnel with local TCP port creation", "localPort", localPort)

	// Step 1: Check if Docker TCP port exists on remote
	if err := s.checkExistingDockerTCPPort(); err != nil {
		s.log.Info("No existing Docker TCP port found, creating local TCP port on remote")

		// Step 2: Create local TCP port on remote using socat/nc/netcat
		if err := s.createLocalTCPPortOnRemote(); err != nil {
			return fmt.Errorf("failed to create local TCP port on remote: %w", err)
		}
	} else {
		s.log.Info("Found existing Docker TCP port on remote", "port", s.remoteLocalPort)
	}

	// Step 3: Create local listener
	listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", localPort))
	if err != nil {
		return fmt.Errorf("failed to create local listener on port %d: %w", localPort, err)
	}

	s.listener = listener
	s.log.Info("SSH tunnel established",
		"localPort", localPort,
		"remoteLocalPort", s.remoteLocalPort)

	// Step 4: Start handling connections in a goroutine
	go s.handleTunnelConnections(listener, localPort)

	return nil
}

// TestSSHTunnel tests if the SSH tunnel is working
func (s *SSHTunnelClient) TestSSHTunnel(localPort int) error {
	s.log.Info("Testing SSH tunnel", "localPort", localPort)

	// Give the tunnel a moment to establish
	time.Sleep(1 * time.Second)

	conn, err := s.connectToLocalPort(localPort)
	if err != nil {
		return err
	}
	defer func() {
		if err := conn.Close(); err != nil {
			s.log.Warn("Failed to close test connection", "error", err)
		}
	}()

	response, err := s.sendTestRequest(conn)
	if err != nil {
		return err
	}

	return s.validateResponse(response, localPort)
}

// Close closes the SSH tunnel and cleans up remote processes
func (s *SSHTunnelClient) Close() error {
	if s.listener != nil {
		if err := s.listener.Close(); err != nil {
			s.log.Warn("Failed to close tunnel listener", "error", err)
		}
	}

	// Clean up remote process if we created one
	if s.remotePID != "" && s.remotePID != "unknown" {
		if err := s.cleanupRemoteProcess(); err != nil {
			s.log.Warn("Failed to cleanup remote process", "error", err)
		}
	}

	return nil
}

// GetConnectionMethod returns the connection method name
func (s *SSHTunnelClient) GetConnectionMethod() string {
	if s.connectionMethod != "" {
		return s.connectionMethod
	}
	return "SSH Tunnel"
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

// checkExistingDockerTCPPort checks if Docker is already listening on a TCP port
func (s *SSHTunnelClient) checkExistingDockerTCPPort() error {
	s.log.Info("Checking for existing Docker TCP port on remote machine")

	session, err := s.client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create SSH session: %w", err)
	}
	defer func() {
		if err := session.Close(); err != nil {
			s.log.Warn("Failed to close SSH session", "error", err)
		}
	}()

	// Check for common Docker TCP ports
	commonPorts := []int{2375, 2376, 2377, 2378, 2379}
	for _, port := range commonPorts {
		cmd := fmt.Sprintf("netstat -tln | grep ':%d ' || ss -tln | grep ':%d '", port, port)
		output, err := session.Output(cmd)
		if err == nil && len(output) > 0 {
			s.remoteLocalPort = port
			s.log.Info("Found existing Docker TCP port", "port", port)
			return nil
		}
	}

	return errors.New("no existing Docker TCP port found")
}

// createLocalTCPPortOnRemote creates a local TCP port on remote machine using socat/nc/netcat
func (s *SSHTunnelClient) createLocalTCPPortOnRemote() error {
	s.log.Info("Creating local TCP port on remote machine")

	// Find an available port on remote machine
	remotePort, err := s.findAvailableRemotePort()
	if err != nil {
		return fmt.Errorf("failed to find available remote port: %w", err)
	}

	s.remoteLocalPort = remotePort

	// Try to create TCP port using different methods
	return s.tryCreateTCPPortWithMethods(remotePort)
}

// tryCreateTCPPortWithMethods tries to create TCP port using different methods
func (s *SSHTunnelClient) tryCreateTCPPortWithMethods(remotePort int) error {
	methods := s.getPrimaryMethods(remotePort)

	// Try the primary methods first
	for _, method := range methods {
		if err := s.tryCreateTCPPort(method.name, method.cmd); err != nil {
			s.log.Warn("Failed to create TCP port with method",
				"method", method.name, "error", err)
			continue
		}

		s.log.Info("Successfully created local TCP port on remote",
			"method", method.name, "port", remotePort)
		s.connectionMethod = fmt.Sprintf("SSH Tunnel (%s)", method.name)
		return nil
	}

	// If primary methods fail, try fallback methods using built-in tools
	s.log.Info("Primary methods failed, trying fallback methods using built-in tools")
	return s.tryCreateTCPPortWithFallbackMethods(remotePort)
}

// getPrimaryMethods returns the primary methods for creating TCP ports
func (s *SSHTunnelClient) getPrimaryMethods(remotePort int) []struct {
	name string
	cmd  string
} {
	return []struct {
		name string
		cmd  string
	}{
		{
			"socat",
			fmt.Sprintf(
				"nohup socat TCP-LISTEN:%d,reuseaddr,fork "+
					"UNIX-CONNECT:/var/run/docker.sock >/dev/null 2>&1 &",
				remotePort),
		},
		{
			"nc",
			fmt.Sprintf(
				"nohup nc -l -p %d -e 'cat /var/run/docker.sock' >/dev/null 2>&1 &",
				remotePort),
		},
		{
			"netcat",
			fmt.Sprintf(
				"nohup netcat -l -p %d -e 'cat /var/run/docker.sock' >/dev/null 2>&1 &",
				remotePort),
		},
	}
}

// tryCreateTCPPortWithFallbackMethods tries to create TCP port using built-in tools
//
//nolint:revive // Function is complex but necessary for fallback methods
func (s *SSHTunnelClient) tryCreateTCPPortWithFallbackMethods(remotePort int) error {
	fallbackMethods := []struct {
		name string
		cmd  string
	}{
		{
			"python3",
			fmt.Sprintf(
				`nohup python3 -c "
import socket, os, sys, threading, signal, time
import sys

def handle_client(client_socket, addr):
    try:
        with open('/var/run/docker.sock', 'rb') as docker_sock:
            while True:
                data = client_socket.recv(4096)
                if not data:
                    break
                docker_sock.write(data)
                docker_sock.flush()
                response = docker_sock.read(4096)
                if response:
                    client_socket.send(response)
    except Exception:
        pass
    finally:
        client_socket.close()

def signal_handler(signum, frame):
    sys.exit(0)

signal.signal(signal.SIGTERM, signal_handler)
signal.signal(signal.SIGINT, signal_handler)

server = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
server.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
server.bind(('127.0.0.1', %d))
server.listen(1)

while True:
    try:
        client_socket, addr = server.accept()
        thread = threading.Thread(target=handle_client, args=(client_socket, addr))
        thread.daemon = True
        thread.start()
    except Exception:
        break
" >/dev/null 2>&1 &`, remotePort),
		},
		{
			"python",
			fmt.Sprintf(
				`nohup python -c "
import socket, os, sys, threading, signal, time
import sys

def handle_client(client_socket, addr):
    try:
        with open('/var/run/docker.sock', 'rb') as docker_sock:
            while True:
                data = client_socket.recv(4096)
                if not data:
                    break
                docker_sock.write(data)
                docker_sock.flush()
                response = docker_sock.read(4096)
                if response:
                    client_socket.send(response)
    except Exception:
        pass
    finally:
        client_socket.close()

def signal_handler(signum, frame):
    sys.exit(0)

signal.signal(signal.SIGTERM, signal_handler)
signal.signal(signal.SIGINT, signal_handler)

server = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
server.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
server.bind(('127.0.0.1', %d))
server.listen(1)

while True:
    try:
        client_socket, addr = server.accept()
        thread = threading.Thread(target=handle_client, args=(client_socket, addr))
        thread.daemon = True
        thread.start()
    except Exception:
        break
" >/dev/null 2>&1 &`, remotePort),
		},
		{
			"bash",
			fmt.Sprintf(
				`nohup bash -c '
set -e
exec 3<>/var/run/docker.sock
while true; do
    timeout 1 bash -c "exec 4<>/dev/tcp/127.0.0.1/%d" 2>/dev/null || break
    {
        while IFS= read -r line <&4; do
            echo "$line" >&3
            read -r response <&3
            echo "$response" >&4
        done
    } &
    wait
    exec 4<&-
    exec 4>&-
done
' >/dev/null 2>&1 &`, remotePort),
		},
	}

	for _, method := range fallbackMethods {
		if err := s.tryCreateTCPPort(method.name, method.cmd); err != nil {
			s.log.Warn("Failed to create TCP port with fallback method",
				"method", method.name, "error", err)
			continue
		}

		s.log.Info("Successfully created local TCP port on remote using fallback method",
			"method", method.name,
			"port", remotePort)
		s.connectionMethod = fmt.Sprintf("SSH Tunnel (fallback %s)", method.name)
		return nil
	}

	return errors.New("failed to create TCP port with any fallback method")
}

// findAvailableRemotePort finds an available port on the remote machine
func (s *SSHTunnelClient) findAvailableRemotePort() (int, error) {
	session, err := s.client.NewSession()
	if err != nil {
		return 0, fmt.Errorf("failed to create SSH session: %w", err)
	}
	defer func() {
		if err := session.Close(); err != nil {
			s.log.Warn("Failed to close SSH session", "error", err)
		}
	}()

	// Try to find an available port starting from 2375, expanding the range
	for port := 2375; port <= 2400; port++ {
		// Check if port is NOT in use
		cmd := fmt.Sprintf(
			"! (netstat -tln | grep -q ':%d ' || ss -tln | grep -q ':%d ') && echo 'available'",
			port, port)
		output, err := session.Output(cmd)
		if err == nil && strings.Contains(string(output), "available") {
			s.log.Info("Found available port on remote machine", "port", port)
			return port, nil
		}
	}

	return 0, errors.New("no available ports found in range 2375-2400")
}

// tryCreateTCPPort attempts to create a TCP port using the specified method
func (s *SSHTunnelClient) tryCreateTCPPort(method, cmd string) error {
	if err := s.executeTCPPortCommand(method, cmd); err != nil {
		return err
	}

	// Give the process a moment to start
	time.Sleep(2 * time.Second)

	return s.verifyTCPPortCreation()
}

// executeTCPPortCommand executes the command to create the TCP port
func (s *SSHTunnelClient) executeTCPPortCommand(method, cmd string) error {
	session, err := s.client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create SSH session: %w", err)
	}
	defer func() {
		if err := session.Close(); err != nil {
			s.log.Warn("Failed to close SSH session", "error", err)
		}
	}()

	// Execute the command to create the TCP port
	if err := session.Run(cmd); err != nil {
		return fmt.Errorf("failed to execute %s command: %w", method, err)
	}

	return nil
}

// verifyTCPPortCreation verifies that the TCP port was created successfully
func (s *SSHTunnelClient) verifyTCPPortCreation() error {
	verifySession, err := s.client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create verification session: %w", err)
	}
	defer func() {
		if err := verifySession.Close(); err != nil {
			s.log.Warn("Failed to close verification session", "error", err)
		}
	}()

	// Try multiple verification methods
	verifyCmd := fmt.Sprintf(
		"netstat -tln | grep ':%d ' || ss -tln | grep ':%d ' || lsof -i :%d",
		s.remoteLocalPort, s.remoteLocalPort, s.remoteLocalPort)
	output, err := verifySession.Output(verifyCmd)
	if err != nil || len(output) == 0 {
		return fmt.Errorf(
			"failed to verify TCP port creation - port %d not listening",
			s.remoteLocalPort)
	}

	s.log.Info("TCP port verification successful",
		"port", s.remoteLocalPort,
		"output", string(output))

	return s.getRemoteProcessPID(verifySession)
}

// getRemoteProcessPID gets the PID of the remote process for cleanup
func (s *SSHTunnelClient) getRemoteProcessPID(session *ssh.Session) error {
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
		pidOutput, err := session.Output(pidCmd)
		if err == nil && len(strings.TrimSpace(string(pidOutput))) > 0 {
			s.remotePID = strings.TrimSpace(string(pidOutput))
			if s.remotePID != "" && s.remotePID != "unknown" {
				s.log.Info("Found remote process PID", "pid", s.remotePID, "method", pidCmd)
				return nil
			}
		}
	}

	s.remotePID = "unknown"
	s.log.Warn("Could not find remote process PID")
	return nil
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
