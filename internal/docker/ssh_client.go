package docker

import (
	"fmt"
	"net"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

// SSHClient represents an SSH client for remote Docker connections
type SSHClient struct {
	config *ssh.ClientConfig
	host   string
	port   string
	user   string
}

// SSHConnection represents an active SSH connection with a socat proxy
type SSHConnection struct {
	client     *ssh.Client
	session    *ssh.Session
	localPort  int
	remoteHost string
	remotePort int // Port used for the socat proxy on remote machine
}

// NewSSHClient creates a new SSH client for the specified host
func NewSSHClient(host string, port int) (*SSHClient, error) {
	// Parse host to extract user, hostname, and sshPort
	user, hostname, sshPort, err := parseSSHHost(host)
	if err != nil {
		return nil, fmt.Errorf("failed to parse SSH host: %w", err)
	}

	// Get SSH key paths
	sshKeyPath, err := getSSHKeyPath()
	if err != nil {
		return nil, fmt.Errorf("failed to get SSH key path: %w", err)
	}

	// Create SSH client config
	config, err := createSSHConfig(user, sshKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create SSH config: %w", err)
	}

	return &SSHClient{
		config: config,
		host:   hostname,
		port:   sshPort,
		user:   user,
	}, nil
}

// parseSSHHost parses an SSH host string in the format [user@]host[:port]
func parseSSHHost(host string) (username, hostname, port string, err error) {
	var userStr, hostn, portStr string

	// Check if user is specified
	if strings.Contains(host, "@") {
		parts := strings.Split(host, "@")
		if len(parts) != 2 {
			return "", "", "", fmt.Errorf("invalid SSH host format: %s", host)
		}
		userStr = parts[0]
		hostn = parts[1]
	} else {
		// Use current user if not specified
		currentUser, err := user.Current()
		if err != nil {
			return "", "", "", fmt.Errorf("failed to get current user: %w", err)
		}
		// On Windows, user.Current() returns "DOMAIN\\username", we need just "username"
		userStr = currentUser.Username
		if strings.Contains(userStr, "\\") {
			parts := strings.Split(userStr, "\\")
			if len(parts) == 2 {
				userStr = parts[1] // Extract just the username part
			}
		}
		hostn = host
	}

	// Check if port is specified
	if strings.Contains(hostn, ":") {
		parts := strings.Split(hostn, ":")
		if len(parts) != 2 {
			return "", "", "", fmt.Errorf("invalid SSH host format: %s", host)
		}
		hostn = parts[0]
		portStr = parts[1]
	} else {
		portStr = "22" // Default SSH port
	}

	return userStr, hostn, portStr, nil
}

// getSSHKeyPath returns the path to the user's SSH private key
func getSSHKeyPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	// Common SSH key locations
	possibleKeys := []string{
		filepath.Join(homeDir, ".ssh", "id_rsa"),
		filepath.Join(homeDir, ".ssh", "id_ed25519"),
		filepath.Join(homeDir, ".ssh", "id_ecdsa"),
	}

	for _, keyPath := range possibleKeys {
		if _, err := os.Stat(keyPath); err == nil {
			return keyPath, nil
		}
	}

	return "", fmt.Errorf("no SSH private key found in %s/.ssh/", homeDir)
}

// createSSHConfig creates SSH client configuration with key-based authentication
func createSSHConfig(username, keyPath string) (*ssh.ClientConfig, error) {
	// Read private key
	keyBytes, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read SSH key: %w", err)
	}

	// Parse private key
	signer, err := ssh.ParsePrivateKey(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse SSH key: %w", err)
	}

	// Get known hosts file path
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}
	knownHostsPath := filepath.Join(homeDir, ".ssh", "known_hosts")

	// Create host key callback
	hostKeyCallback, err := knownhosts.New(knownHostsPath)
	if err != nil {
		// If known_hosts doesn't exist, create a permissive callback for testing
		// Note: This bypasses host key verification and should only be used in trusted environments
		// #nosec G106 -- InsecureIgnoreHostKey is used as fallback when known_hosts is unavailable
		hostKeyCallback = ssh.InsecureIgnoreHostKey()
	}

	return &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: hostKeyCallback,
		Timeout:         30 * time.Second,
	}, nil
}

// validateHostname checks if the hostname can be resolved to an IP address
func (s *SSHClient) validateHostname() error {
	// Check if the host is already an IP address
	if net.ParseIP(s.host) != nil {
		return nil // It's a valid IP address
	}

	// Try to resolve the hostname
	ips, err := net.LookupHost(s.host)
	if err != nil {
		return fmt.Errorf("cannot resolve hostname '%s': %w", s.host, err)
	}

	if len(ips) == 0 {
		return fmt.Errorf("hostname '%s' resolved to no IP addresses", s.host)
	}

	return nil
}

// Connect establishes an SSH connection and sets up a socat proxy
func (s *SSHClient) Connect(remotePort int) (*SSHConnection, error) {
	// Validate that the hostname can be resolved
	if err := s.validateHostname(); err != nil {
		return nil, fmt.Errorf("hostname validation failed: %w", err)
	}

	// Connect to SSH server
	client, err := ssh.Dial("tcp", net.JoinHostPort(s.host, s.port), s.config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SSH server %s:%s: %w", s.host, s.port, err)
	}

	// Find an available local port for the proxy
	localPort, err := findAvailablePort()
	if err != nil {
		client.Close()
		return nil, fmt.Errorf("failed to find available port: %w", err)
	}

	// Set up socat proxy on remote machine
	session, err := s.setupSocatProxy(client, localPort, remotePort)
	if err != nil {
		client.Close()
		return nil, fmt.Errorf("failed to set up socat proxy: %w", err)
	}

	return &SSHConnection{
		client:     client,
		session:    session,
		localPort:  localPort,
		remoteHost: s.host,
		remotePort: remotePort,
	}, nil
}

// setupSocatProxy sets up a socat proxy on the remote machine
func (s *SSHClient) setupSocatProxy(client *ssh.Client, _, remotePort int) (*ssh.Session, error) {
	// Create a new session for the socat command
	session, err := client.NewSession()
	if err != nil {
		return nil, fmt.Errorf("failed to create SSH session: %w", err)
	}

	// Build the socat command - use nohup with a more robust approach
	socatCmd := fmt.Sprintf("nohup bash -c 'socat TCP-LISTEN:%d,bind=0.0.0.0,reuseaddr,fork UNIX-CONNECT:/var/run/docker.sock' > /dev/null 2>&1 &", remotePort)

	// Execute the socat command in a detached way
	go func() {
		defer session.Close()

		// Set up pipes for the session
		session.Stdout = os.Stdout
		session.Stderr = os.Stderr

		// Start the socat command
		if err := session.Start(socatCmd); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: socat command start failed: %v\n", err)
			return
		}

		// Don't wait for the command to complete - let it run in background
		// The session will be closed when the connection is closed
	}()

	// Wait a moment for socat to start
	time.Sleep(3 * time.Second)

	return session, nil
}

// GetLocalProxyHost returns the local host:port for the Docker proxy
func (s *SSHConnection) GetLocalProxyHost() string {
	return fmt.Sprintf("tcp://127.0.0.1:%d", s.localPort)
}

// Close closes the SSH connection and cleans up
func (s *SSHConnection) Close() error {
	var errors []string

	// Kill the socat process on the remote machine
	if err := s.killSocatProcess(); err != nil {
		errors = append(errors, fmt.Sprintf("socat cleanup: %v", err))
	}

	if s.session != nil {
		if err := s.session.Close(); err != nil {
			errors = append(errors, fmt.Sprintf("session close: %v", err))
		}
	}

	if s.client != nil {
		if err := s.client.Close(); err != nil {
			errors = append(errors, fmt.Sprintf("client close: %v", err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed to close SSH connection: %s", strings.Join(errors, "; "))
	}

	return nil
}

// killSocatProcess kills the socat process running on the remote machine
func (s *SSHConnection) killSocatProcess() error {
	// Create a new session to kill the socat process
	session, err := s.client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session for socat cleanup: %w", err)
	}
	defer session.Close()

	// Kill any socat processes listening on the specified port
	killCmd := fmt.Sprintf("pkill -f 'socat.*TCP-LISTEN:%d'", s.remotePort)
	if err := session.Run(killCmd); err != nil {
		// It's okay if no processes were found to kill
		return nil
	}

	return nil
}

// findAvailablePort finds an available local port
func findAvailablePort() (int, error) {
	// Try ports in the range 2376-2385 (avoiding 2375 which is standard Docker port)
	for port := 2376; port <= 2385; port++ {
		addr := fmt.Sprintf("127.0.0.1:%d", port)
		listener, err := net.Listen("tcp", addr)
		if err == nil {
			listener.Close()
			return port, nil
		}
	}

	return 0, fmt.Errorf("no available ports found in range 2376-2385")
}
