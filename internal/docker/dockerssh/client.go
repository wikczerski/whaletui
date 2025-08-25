// Package dockerssh provides SSH client functionality for Docker operations.
package dockerssh

import (
	"errors"
	"fmt"
	"log/slog"
	"net"
	"os"
	"strings"
	"time"

	"github.com/wikczerski/whaletui/internal/logger"
	"golang.org/x/crypto/ssh"
)

// SSHClient represents an SSH client for remote Docker connections
type SSHClient struct {
	config   *ssh.ClientConfig
	host     string
	port     string
	user     string
	log      *slog.Logger
	listener net.Listener
}

// NewSSHClient creates a new SSH client for the specified host
func NewSSHClient(host string, port int) (*SSHClient, error) {
	if host == "" {
		return nil, errors.New("SSH host cannot be empty")
	}

	user, hostname, sshPort, err := parseSSHHost(host)
	if err != nil {
		return nil, fmt.Errorf("failed to parse SSH host: %w", err)
	}

	config, err := createSSHConfigWithKey(user)
	if err != nil {
		return nil, err
	}

	return createSSHClient(config, hostname, sshPort, user), nil
}

// createSSHConfigWithKey creates SSH config with the user's key
func createSSHConfigWithKey(user string) (*ssh.ClientConfig, error) {
	sshKeyPath, err := getSSHKeyPath()
	if err != nil {
		return nil, fmt.Errorf("failed to get SSH key path: %w", err)
	}

	return createSSHConfig(user, sshKeyPath)
}

// createSSHClient creates the SSH client with the given configuration
func createSSHClient(config *ssh.ClientConfig, hostname, sshPort, user string) *SSHClient {
	return &SSHClient{
		config: config,
		host:   hostname,
		port:   sshPort,
		user:   user,
		log:    logger.GetLogger(),
	}
}

// GetConnectionInfo returns diagnostic information about the SSH client configuration
func (s *SSHClient) GetConnectionInfo() string {
	if s == nil {
		return "SSH Client: <nil>"
	}

	user := s.getUserInfo()
	host := s.getHostInfo()
	port := s.getPortInfo()

	return fmt.Sprintf("SSH Client: %s@%s:%s", user, host, port)
}

// DiagnoseConnection performs comprehensive diagnostics on the SSH connection
func (s *SSHClient) DiagnoseConnection() error {
	var errors []string

	s.diagnoseHostname(&errors)
	s.diagnoseSSHKey(&errors)

	if s.config == nil {
		errors = append(errors, "SSH configuration not available")
		return s.buildDiagnosticError(errors)
	}

	return s.diagnoseSSHConnection(&errors)
}

// Connect establishes an SSH connection without socat
func (s *SSHClient) Connect() (*SSHContext, error) {
	if err := s.validateHostname(); err != nil {
		return nil, fmt.Errorf("hostname validation failed: %w", err)
	}

	client, err := ssh.Dial("tcp", net.JoinHostPort(s.host, s.port), s.config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SSH server %s:%s: %w", s.host, s.port, err)
	}

	if err := s.checkDockerSocketAccess(client); err != nil {
		if closeErr := client.Close(); closeErr != nil {
			s.log.Warn("Failed to close SSH client", "error", closeErr)
		}
		return nil, fmt.Errorf("docker socket not accessible: %w", err)
	}

	return &SSHContext{
		client:     client,
		remoteHost: s.host,
		log:        s.log,
	}, nil
}

// ConnectWithSocat establishes an SSH connection and sets up port forwarding to the remote Docker socket
func (s *SSHClient) ConnectWithSocat(remotePort int) (*SSHConnection, error) {
	if err := s.validateHostname(); err != nil {
		return nil, fmt.Errorf("hostname validation failed: %w", err)
	}

	client, err := s.establishSSHConnection()
	if err != nil {
		return nil, err
	}

	return s.setupPortForwardingAndConnection(client, remotePort)
}

// establishSSHConnection establishes the initial SSH connection
func (s *SSHClient) establishSSHConnection() (*ssh.Client, error) {
	client, err := ssh.Dial("tcp", net.JoinHostPort(s.host, s.port), s.config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SSH server %s:%s: %w", s.host, s.port, err)
	}
	return client, nil
}

// setupPortForwardingAndConnection sets up port forwarding and creates the connection
func (s *SSHClient) setupPortForwardingAndConnection(
	client *ssh.Client,
	remotePort int,
) (*SSHConnection, error) {
	localPort, err := s.setupPortForwarding(client)
	if err != nil {
		s.closeSSHClient(client)
		return nil, err
	}

	s.setupAndTestPortForwarding(client, localPort)

	s.log.Info("SSH port forwarding established successfully", "localPort", localPort)

	return s.createSSHConnection(client, localPort), nil
}

// setupAndTestPortForwarding sets up and tests the port forwarding
func (s *SSHClient) setupAndTestPortForwarding(client *ssh.Client, localPort int) {
	s.startPortForwarding(client, localPort)
	s.waitForListenerReady()

	if err := s.testPortForwarding(localPort); err != nil {
		s.closeSSHClient(client)
	}
}

// setupPortForwarding sets up the local port and listener
func (s *SSHClient) setupPortForwarding(client *ssh.Client) (int, error) {
	localPort, err := findAvailablePort()
	if err != nil {
		return 0, fmt.Errorf("failed to find available port: %w", err)
	}

	s.log.Info(
		"Setting up SSH port forwarding",
		"localPort",
		localPort,
		"remoteSocket",
		"/var/run/docker.sock",
	)

	listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", localPort))
	if err != nil {
		return 0, fmt.Errorf("failed to bind to local port %d: %w", localPort, err)
	}

	s.listener = listener
	return localPort, nil
}

// startPortForwarding starts the port forwarding goroutine
func (s *SSHClient) startPortForwarding(client *ssh.Client, localPort int) {
	go func() {
		defer func() {
			if err := s.listener.Close(); err != nil {
				s.log.Warn("Failed to close listener", "error", err)
			}
		}()
		for {
			localConn, err := s.listener.Accept()
			if err != nil {
				s.log.Warn("Failed to accept local connection", "error", err)
				return
			}

			s.handlePortForwardingConnection(client, localConn)
		}
	}()
}

// handlePortForwardingConnection handles a single port forwarding connection
func (s *SSHClient) handlePortForwardingConnection(client *ssh.Client, localConn net.Conn) {
	session, err := s.createSSHSession(client)
	if err != nil {
		s.cleanupConnection(localConn, nil)
		return
	}

	s.setupSessionIO(session, localConn)

	if err := s.executeSocatCommand(session); err != nil {
		s.cleanupConnection(localConn, session)
		return
	}

	s.cleanupConnection(localConn, session)
}

// createSSHSession creates a new SSH session
func (s *SSHClient) createSSHSession(client *ssh.Client) (*ssh.Session, error) {
	session, err := client.NewSession()
	if err != nil {
		s.log.Warn("Failed to create SSH session for forwarding", "error", err)
		return nil, err
	}
	return session, nil
}

// executeSocatCommand executes the socat command on the remote machine
func (s *SSHClient) executeSocatCommand(session *ssh.Session) error {
	if err := session.Run("socat STDIO UNIX-CONNECT:/var/run/docker.sock"); err != nil {
		s.log.Warn("Failed to execute socat on remote machine", "error", err)
		return err
	}
	return nil
}

// cleanupConnection safely closes the connection and session
func (s *SSHClient) cleanupConnection(localConn net.Conn, session *ssh.Session) {
	if session != nil {
		if err := session.Close(); err != nil {
			s.log.Warn("Failed to close SSH session", "error", err)
		}
	}
	if localConn != nil {
		if err := localConn.Close(); err != nil {
			s.log.Warn("Failed to close local connection", "error", err)
		}
	}
}

// setupSessionIO sets up the session input/output
func (s *SSHClient) setupSessionIO(session *ssh.Session, localConn net.Conn) {
	session.Stdin = localConn
	session.Stdout = localConn
	session.Stderr = os.Stderr
}

// waitForListenerReady waits for the listener to be ready
func (s *SSHClient) waitForListenerReady() {
	time.Sleep(1 * time.Second)
}

// testPortForwarding tests if the port forwarding is working
func (s *SSHClient) testPortForwarding(localPort int) error {
	testConn, err := net.DialTimeout("tcp", fmt.Sprintf("127.0.0.1:%d", localPort), 2*time.Second)
	if err != nil {
		return fmt.Errorf("local port forwarding test failed: %w", err)
	}
	if err := testConn.Close(); err != nil {
		s.log.Warn("Failed to close test connection", "error", err)
	}
	return nil
}

// createSSHConnection creates the SSH connection struct
func (s *SSHClient) createSSHConnection(client *ssh.Client, localPort int) *SSHConnection {
	return &SSHConnection{
		client:     client,
		session:    nil, // Not using a single session for this approach
		localPort:  localPort,
		remoteHost: s.host,
		remotePort: 0,  // Not using remote port for this approach
		socatPID:   "", // Not using socat PID for this approach
		log:        s.log,
		listener:   s.listener, // Store the listener for cleanup
	}
}

// getUserInfo returns the user info with fallback
func (s *SSHClient) getUserInfo() string {
	if s.user == "" {
		return "<unknown>"
	}
	return s.user
}

// getHostInfo returns the host info with fallback
func (s *SSHClient) getHostInfo() string {
	if s.host == "" {
		return "<unknown>"
	}
	return s.host
}

// getPortInfo returns the port info with fallback
func (s *SSHClient) getPortInfo() string {
	if s.port == "" {
		return "<unknown>"
	}
	return s.port
}

// diagnoseHostname checks hostname resolution
func (s *SSHClient) diagnoseHostname(errors *[]string) {
	if err := s.validateHostname(); err != nil {
		*errors = append(*errors, fmt.Sprintf("Hostname resolution: %v", err))
	}
}

// diagnoseSSHKey checks SSH key configuration and permissions
func (s *SSHClient) diagnoseSSHKey(errors *[]string) {
	sshKeyPath, err := getSSHKeyPath()
	if err != nil {
		*errors = append(*errors, fmt.Sprintf("SSH key: %v", err))
		return
	}

	s.checkSSHKeyPermissions(sshKeyPath, errors)
}

// checkSSHKeyPermissions checks SSH key file permissions
func (s *SSHClient) checkSSHKeyPermissions(sshKeyPath string, errors *[]string) {
	if info, err := os.Stat(sshKeyPath); err == nil {
		if info.Mode().Perm()&0o077 != 0 {
			*errors = append(
				*errors,
				fmt.Sprintf(
					"SSH key has overly permissive permissions %v (should be 600)",
					info.Mode().Perm(),
				),
			)
		}
	}
}

// diagnoseSSHConnection tests the actual SSH connection
func (s *SSHClient) diagnoseSSHConnection(errors *[]string) error {
	client, err := ssh.Dial("tcp", net.JoinHostPort(s.host, s.port), s.config)
	if err != nil {
		*errors = append(*errors, fmt.Sprintf("SSH connection: %v", err))
		return s.buildDiagnosticError(*errors)
	}
	defer s.closeSSHClient(client)

	if err := s.checkDockerSocketAccess(client); err != nil {
		*errors = append(*errors, fmt.Sprintf("Docker socket access: %v", err))
	}

	return s.buildDiagnosticError(*errors)
}

// closeSSHClient safely closes the SSH client
func (s *SSHClient) closeSSHClient(client *ssh.Client) {
	if err := client.Close(); err != nil {
		s.log.Warn("Failed to close SSH client", "error", err)
	}
}

// buildDiagnosticError builds the final diagnostic error message
func (s *SSHClient) buildDiagnosticError(errors []string) error {
	if len(errors) > 0 {
		return fmt.Errorf("SSH connection diagnostics failed:\n%s", strings.Join(errors, "\n"))
	}
	return nil
}
