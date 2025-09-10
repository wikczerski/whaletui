package dockerssh

import (
	"errors"
	"fmt"
	"log/slog"
	"net"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

// SSHClient handles SSH connections for Docker access
type SSHClient struct {
	host     string
	port     string
	user     string
	config   *ssh.ClientConfig
	log      *slog.Logger
	keyPath  string
	password string
}

// NewSSHClient creates a new SSH client
func NewSSHClient(host, port, user string, log *slog.Logger) *SSHClient {
	return &SSHClient{
		host:   host,
		port:   port,
		user:   user,
		config: createSSHConfig(user, log, "", ""),
		log:    log,
	}
}

// NewSSHClientWithAuth creates a new SSH client with authentication options
func NewSSHClientWithAuth(
	host, port, user, keyPath, password string,
	log *slog.Logger,
) *SSHClient {
	return &SSHClient{
		host:     host,
		port:     port,
		user:     user,
		config:   createSSHConfig(user, log, keyPath, password),
		log:      log,
		keyPath:  keyPath,
		password: password,
	}
}

// createSSHConfig creates the SSH client configuration
func createSSHConfig(
	username string,
	log *slog.Logger,
	keyPath, password string,
) *ssh.ClientConfig {
	config := &ssh.ClientConfig{
		User:            username,
		Auth:            []ssh.AuthMethod{},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), //nolint:gosec // Acceptable for development/testing
		Timeout:         30 * time.Second,
	}

	// Try to use SSH key authentication
	if err := addSSHKeyAuth(config, log, keyPath); err != nil {
		log.Warn("Failed to add SSH key authentication", "error", err)
	}

	// Add password authentication
	if password != "" {
		config.Auth = append(config.Auth, ssh.Password(password))
		log.Info("Added password authentication")
	} else {
		// Add interactive password callback as fallback
		config.Auth = append(config.Auth, ssh.PasswordCallback(func() (string, error) {
			return promptForSSHPassword()
		}))
	}

	return config
}

// addSSHKeyAuth adds SSH key authentication to the config
func addSSHKeyAuth(config *ssh.ClientConfig, log *slog.Logger, customKeyPath string) error {
	keyPaths := getSSHKeyPaths(customKeyPath)

	for _, keyPath := range keyPaths {
		if err := trySSHKeyAuth(config, log, keyPath); err == nil {
			return nil
		}
	}

	if customKeyPath != "" {
		return fmt.Errorf("SSH key not found at specified path: %s", customKeyPath)
	}
	return errors.New("no SSH keys found")
}

// getSSHKeyPaths returns the list of SSH key paths to try
func getSSHKeyPaths(customKeyPath string) []string {
	if customKeyPath != "" {
		return []string{customKeyPath}
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return []string{}
	}

	return []string{
		filepath.Join(homeDir, ".ssh", "id_rsa"),
		filepath.Join(homeDir, ".ssh", "id_ed25519"),
		filepath.Join(homeDir, ".ssh", "id_ecdsa"),
	}
}

// trySSHKeyAuth attempts to use a specific SSH key for authentication
func trySSHKeyAuth(config *ssh.ClientConfig, log *slog.Logger, keyPath string) error {
	if _, err := os.Stat(keyPath); err != nil {
		return err
	}

	key, err := os.ReadFile(keyPath) //nolint:gosec // SSH key file reading is necessary
	if err != nil {
		log.Warn("Failed to read SSH key", "path", keyPath, "error", err)
		return err
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Warn("Failed to parse SSH key", "path", keyPath, "error", err)
		return err
	}

	config.Auth = append(config.Auth, ssh.PublicKeys(signer))
	log.Info("Added SSH key authentication", "path", keyPath)
	return nil
}

// ConnectWithFallback establishes an SSH connection using SSH tunneling only
// Only SSH tunnel method is supported - all fallback methods are disabled
func (s *SSHClient) ConnectWithFallback(remotePort int) (*SSHConnection, error) {
	if err := s.validateHostname(); err != nil {
		return nil, fmt.Errorf("hostname validation failed: %w", err)
	}

	client, err := s.establishSSHConnection()
	if err != nil {
		return nil, err
	}

	// Try SSH tunnel only
	connection, err := s.trySSHTunnel(client, remotePort)
	if err == nil {
		s.log.Info("✅ Successfully established SSH tunnel connection")
		s.log.Info("🔗 Connection Method: SSH Tunnel (local TCP port on remote machine)")
		return connection, nil
	}
	s.log.Warn("SSH tunnel failed", "error", err)

	// Clean up any local TCP ports that might have been created for SSH tunnel
	s.cleanupLocalTCPPorts(client)

	// SSH tunnel connection failed
	return nil, fmt.Errorf("SSH tunnel connection failed: %w", err)
}

// trySSHTunnel attempts to establish SSH tunnel connection
func (s *SSHClient) trySSHTunnel(client *ssh.Client, remotePort int) (*SSHConnection, error) {
	s.log.Info("Attempting SSH tunnel connection")

	// Create SSH tunnel client
	tunnelClient := NewSSHTunnelClient(client, s.log)

	// Check if Docker socket is accessible
	if err := s.checkDockerSocketAccess(client); err != nil {
		return nil, fmt.Errorf("docker socket not accessible for SSH tunnel: %w", err)
	}

	// Find an available local port
	localPort, err := findAvailablePort()
	if err != nil {
		return nil, fmt.Errorf("failed to find available local port: %w", err)
	}

	// Set up and test SSH tunnel
	if err := s.setupAndTestTunnel(tunnelClient, localPort); err != nil {
		return nil, err
	}

	// Create and return the SSH connection
	return s.createSSHConnection(client, localPort, remotePort, tunnelClient), nil
}

// setupAndTestTunnel sets up and tests the SSH tunnel
func (s *SSHClient) setupAndTestTunnel(tunnelClient *SSHTunnelClient, localPort int) error {
	// Set up SSH tunnel
	if err := tunnelClient.SetupSSHTunnel(localPort); err != nil {
		return fmt.Errorf("failed to setup SSH tunnel: %w", err)
	}

	// Test the SSH tunnel
	if err := tunnelClient.TestSSHTunnel(localPort); err != nil {
		// Clean up the tunnel client
		if err := tunnelClient.Close(); err != nil {
			s.log.Warn("Failed to close tunnel client", "error", err)
		}
		return fmt.Errorf("SSH tunnel test failed: %w", err)
	}

	return nil
}

// createSSHConnection creates the SSH connection struct
func (s *SSHClient) createSSHConnection(
	client *ssh.Client,
	localPort, remotePort int,
	tunnelClient *SSHTunnelClient,
) *SSHConnection {
	return &SSHConnection{
		client:         client,
		localPort:      localPort,
		remoteHost:     s.host,
		remotePort:     remotePort,
		log:            s.log,
		connectionType: "SSH Tunnel",
		tunnelClient:   tunnelClient,
	}
}

// checkDockerSocketAccess checks if the Docker socket is accessible via SSH
func (s *SSHClient) checkDockerSocketAccess(client *ssh.Client) error {
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session for Docker socket check: %w", err)
	}
	defer func() {
		if err := session.Close(); err != nil {
			s.log.Warn("Failed to close Docker socket check session", "error", err)
		}
	}()

	// Check if the Docker socket exists and is accessible
	cmd := "test -S /var/run/docker.sock && echo 'Docker socket accessible' || echo 'Docker socket not accessible'"
	output, err := session.Output(cmd)
	if err != nil {
		return fmt.Errorf("failed to check Docker socket access: %w", err)
	}

	outputStr := string(output)
	if outputStr == "Docker socket not accessible\n" {
		return errors.New("docker socket is not accessible on remote host")
	}

	s.log.Info("Docker socket access verified", "output", outputStr)
	return nil
}

// cleanupLocalTCPPorts cleans up any local TCP ports that might have been created for SSH tunnel
func (s *SSHClient) cleanupLocalTCPPorts(client *ssh.Client) {
	s.log.Info("Cleaning up any local TCP ports created for SSH tunnel")

	session, err := client.NewSession()
	if err != nil {
		s.log.Warn("Failed to create cleanup session", "error", err)
		return
	}
	defer func() {
		if err := session.Close(); err != nil {
			s.log.Warn("Failed to close cleanup session", "error", err)
		}
	}()

	// Kill any processes listening on common Docker ports
	ports := []int{2375, 2376, 2377, 2378, 2379}
	for _, port := range ports {
		cmd := fmt.Sprintf(
			"lsof -ti:%d | xargs kill -9 2>/dev/null || fuser -k %d/tcp 2>/dev/null || true",
			port, port)
		if err := session.Run(cmd); err != nil {
			s.log.Debug("Failed to cleanup port", "port", port, "error", err)
		}
	}

	s.log.Info("Local TCP port cleanup completed")
}

// establishSSHConnection establishes the initial SSH connection
func (s *SSHClient) establishSSHConnection() (*ssh.Client, error) {
	client, err := ssh.Dial("tcp", net.JoinHostPort(s.host, s.port), s.config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SSH server %s:%s: %w", s.host, s.port, err)
	}
	return client, nil
}

// validateHostname validates the hostname
func (s *SSHClient) validateHostname() error {
	if s.host == "" {
		return errors.New("hostname cannot be empty")
	}
	return nil
}

// findAvailablePort finds an available local port
func findAvailablePort() (int, error) {
	// Try to find an available port in the range 2376-2390
	for port := 2376; port <= 2390; port++ {
		listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
		if err == nil {
			_ = listener.Close() // Ignore close error
			return port, nil
		}
	}

	// Fallback to system-assigned port if range is full
	listener, err := net.Listen("tcp", ":0") //nolint:gosec // Safe for finding available ports
	if err != nil {
		return 0, err
	}
	defer func() {
		_ = listener.Close() // Ignore close error
	}()

	addr, ok := listener.Addr().(*net.TCPAddr)
	if !ok {
		return 0, errors.New("failed to get TCP address")
	}
	return addr.Port, nil
}

// promptForSSHPassword securely prompts the user for a password
func promptForSSHPassword() (string, error) {
	fmt.Print("SSH key authentication failed. Enter SSH password: ")

	// Read password from stdin without echoing
	passwordBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return "", fmt.Errorf("failed to read password: %w", err)
	}

	fmt.Println() // Add newline after password input

	return string(passwordBytes), nil
}
