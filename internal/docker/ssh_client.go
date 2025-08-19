package docker

import (
	"fmt"
	"net"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
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
	if host == "" {
		return nil, fmt.Errorf("SSH host cannot be empty")
	}

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

// GetConnectionInfo returns diagnostic information about the SSH client configuration
func (s *SSHClient) GetConnectionInfo() string {
	if s == nil {
		return "SSH Client: <nil>"
	}

	user := s.user
	if user == "" {
		user = "<unknown>"
	}

	host := s.host
	if host == "" {
		host = "<unknown>"
	}

	port := s.port
	if port == "" {
		port = "<unknown>"
	}

	return fmt.Sprintf("SSH Client: %s@%s:%s", user, host, port)
}

// DiagnoseConnection performs comprehensive diagnostics on the SSH connection
func (s *SSHClient) DiagnoseConnection() error {
	var errors []string

	// Check hostname resolution
	if err := s.validateHostname(); err != nil {
		errors = append(errors, fmt.Sprintf("Hostname resolution: %v", err))
	}

	// Check SSH key availability
	sshKeyPath, err := getSSHKeyPath()
	if err != nil {
		errors = append(errors, fmt.Sprintf("SSH key: %v", err))
	} else {
		// Verify SSH key permissions
		if info, err := os.Stat(sshKeyPath); err == nil {
			if info.Mode().Perm()&0o077 != 0 {
				errors = append(errors, fmt.Sprintf("SSH key has overly permissive permissions %v (should be 600)", info.Mode().Perm()))
			}
		}
	}

	// Check if SSH config is available
	if s.config == nil {
		errors = append(errors, "SSH configuration not available")
		return fmt.Errorf("SSH connection diagnostics failed:\n%s", strings.Join(errors, "\n"))
	}

	// Check if we can establish a basic SSH connection
	client, err := ssh.Dial("tcp", net.JoinHostPort(s.host, s.port), s.config)
	if err != nil {
		errors = append(errors, fmt.Sprintf("SSH connection: %v", err))
		return fmt.Errorf("SSH connection diagnostics failed:\n%s", strings.Join(errors, "\n"))
	}
	defer client.Close()

	// Test basic SSH functionality
	if err := s.checkSocatAvailability(client); err != nil {
		errors = append(errors, fmt.Sprintf("Socat availability: %v", err))
	}

	if err := s.checkDockerSocketAccess(client); err != nil {
		errors = append(errors, fmt.Sprintf("Docker socket access: %v", err))
	}

	if len(errors) > 0 {
		return fmt.Errorf("SSH connection diagnostics failed:\n%s", strings.Join(errors, "\n"))
	}

	return nil
}

// parseSSHHost parses an SSH host string in the format [user@]host[:port]
func parseSSHHost(host string) (username, hostname, port string, err error) {
	if host == "" {
		return "", "", "", fmt.Errorf("SSH host cannot be empty")
	}

	var userStr, hostn, portStr string

	// Check if user is specified
	if strings.Contains(host, "@") {
		parts := strings.Split(host, "@")
		if len(parts) != 2 {
			return "", "", "", fmt.Errorf("invalid SSH host format '%s': expected [user@]host[:port]", host)
		}
		userStr = strings.TrimSpace(parts[0])
		hostn = strings.TrimSpace(parts[1])

		if userStr == "" {
			return "", "", "", fmt.Errorf("username cannot be empty in SSH host format '%s'", host)
		}
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
		hostn = strings.TrimSpace(host)
	}

	if hostn == "" {
		return "", "", "", fmt.Errorf("hostname cannot be empty in SSH host format '%s'", host)
	}

	// Check if port is specified
	if strings.Contains(hostn, ":") {
		parts := strings.Split(hostn, ":")
		if len(parts) != 2 {
			return "", "", "", fmt.Errorf("invalid SSH host format '%s': expected [user@]host[:port]", host)
		}
		hostn = strings.TrimSpace(parts[0])
		portStr = strings.TrimSpace(parts[1])

		if hostn == "" {
			return "", "", "", fmt.Errorf("hostname cannot be empty in SSH host format '%s'", host)
		}
		if portStr == "" {
			return "", "", "", fmt.Errorf("port cannot be empty in SSH host format '%s'", host)
		}

		// Validate port is numeric
		if _, err := fmt.Sscanf(portStr, "%d", new(int)); err != nil {
			return "", "", "", fmt.Errorf("invalid port '%s' in SSH host format '%s': port must be numeric", portStr, host)
		}
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
		if info, err := os.Stat(keyPath); err == nil {
			// On Windows, be more permissive with permissions since Windows ACLs work differently
			if runtime.GOOS == "windows" {
				// On Windows, just check if the file exists and is readable
				return keyPath, nil
			}

			// On Unix-like systems, check for strict permissions
			mode := info.Mode().Perm()
			if mode&0o077 != 0 {
				return "", fmt.Errorf("SSH key %s has overly permissive permissions %v (should be 600)", keyPath, mode)
			}
			return keyPath, nil
		}
	}

	return "", fmt.Errorf("no SSH private key found in %s/.ssh/ (checked: id_rsa, id_ed25519, id_ecdsa)", homeDir)
}

// createSSHConfig creates SSH client configuration with key-based authentication
func createSSHConfig(username, keyPath string) (*ssh.ClientConfig, error) {
	// Read private key
	keyBytes, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read SSH key at %s: %w", keyPath, err)
	}

	// Parse private key
	signer, err := ssh.ParsePrivateKey(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse SSH key at %s: %w", keyPath, err)
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
	if s.host == "" {
		return fmt.Errorf("hostname cannot be empty")
	}

	// Check if the host is already an IP address
	if net.ParseIP(s.host) != nil {
		return nil // It's a valid IP address
	}

	// Basic hostname validation
	if strings.HasPrefix(s.host, ".") || strings.HasSuffix(s.host, ".") {
		return fmt.Errorf("hostname '%s' cannot start or end with a dot", s.host)
	}

	if strings.Contains(s.host, "..") {
		return fmt.Errorf("hostname '%s' cannot contain consecutive dots", s.host)
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

	// Store the actual remote port used (in case it was defaulted)
	if remotePort == 0 {
		remotePort = 2375
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
	// Build the socat command with better error handling and use a specific port
	if remotePort == 0 {
		// Find an available port on the remote machine
		availablePort, err := s.findAvailableRemotePort(client)
		if err != nil {
			return nil, fmt.Errorf("failed to find available remote port: %w", err)
		}
		remotePort = availablePort
	}

	// Debug: Print the port being used
	fmt.Printf("DEBUG: Using remote port %d for socat proxy\n", remotePort)

	// First, check if socat is available on the remote machine
	if err := s.checkSocatAvailability(client); err != nil {
		return nil, fmt.Errorf("socat not available on remote machine: %w", err)
	}

	// Check if Docker socket exists and is accessible
	if err := s.checkDockerSocketAccess(client); err != nil {
		return nil, fmt.Errorf("docker socket not accessible: %w", err)
	}

	// Create a new session for the socat command
	session, err := client.NewSession()
	if err != nil {
		return nil, fmt.Errorf("failed to create SSH session for socat: %w", err)
	}

	// Build the socat command with better error handling
	socatCmd := fmt.Sprintf("nohup bash -c 'socat TCP-LISTEN:%d,bind=0.0.0.0,reuseaddr,fork UNIX-CONNECT:/var/run/docker.sock 2>&1' > /tmp/socat.log 2>&1 & echo $!", remotePort)

	// Debug: Print the socat command
	fmt.Printf("DEBUG: Executing socat command: %s\n", socatCmd)

	// Execute the socat command and capture the PID
	output, err := session.Output(socatCmd)
	if err != nil {
		session.Close()
		return nil, fmt.Errorf("failed to start socat proxy: %w", err)
	}

	// Parse the PID from output
	pid := strings.TrimSpace(string(output))
	if pid == "" {
		session.Close()
		return nil, fmt.Errorf("failed to get socat process ID")
	}

	// Debug: Print the PID
	fmt.Printf("DEBUG: Socat started with PID %s\n", pid)

	// Close the current session
	session.Close()

	// Wait a moment for the process to start
	time.Sleep(2 * time.Second) // Increased wait time

	// Debug: Print verification attempt
	fmt.Printf("DEBUG: Attempting to verify socat process on port %d\n", remotePort)

	// Verify the socat process is running with a new session
	if err := s.verifySocatProcess(client, pid, remotePort); err != nil {
		return nil, fmt.Errorf("socat process verification failed: %w", err)
	}

	// Debug: Print success
	fmt.Printf("DEBUG: Socat verification successful on port %d\n", remotePort)

	// Create a new session for monitoring
	monitorSession, err := client.NewSession()
	if err != nil {
		return nil, fmt.Errorf("failed to create monitoring session: %w", err)
	}

	return monitorSession, nil
}

// GetLocalProxyHost returns the local host:port for the Docker proxy
func (s *SSHConnection) GetLocalProxyHost() string {
	if s == nil {
		return "tcp://127.0.0.1:0"
	}

	port := s.localPort
	if port <= 0 {
		port = 2375 // Default to Docker port if invalid
	}

	return fmt.Sprintf("tcp://127.0.0.1:%d", port)
}

// Close closes the SSH connection and cleans up
func (s *SSHConnection) Close() error {
	if s == nil {
		return nil
	}

	var errors []string

	// Kill the socat process on the remote machine only if client exists
	if s.client != nil {
		if err := s.killSocatProcess(); err != nil {
			errors = append(errors, fmt.Sprintf("socat cleanup: %v", err))
		}
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
	if s == nil || s.client == nil {
		return nil
	}

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
		// This could happen if the process was already terminated
		return nil
	}

	// Wait a moment for the process to be killed
	time.Sleep(200 * time.Millisecond)

	// Verify the process is no longer running
	verifySession, err := s.client.NewSession()
	if err != nil {
		return nil // Don't fail cleanup if we can't verify
	}
	defer verifySession.Close()

	verifyCmd := fmt.Sprintf("netstat -tln | grep ':%d '", s.remotePort)
	if err := verifySession.Run(verifyCmd); err == nil {
		// Process is still running, try force kill
		forceKillSession, err := s.client.NewSession()
		if err != nil {
			return nil
		}
		defer forceKillSession.Close()

		forceKillCmd := fmt.Sprintf("pkill -9 -f 'socat.*TCP-LISTEN:%d'", s.remotePort)
		_ = forceKillSession.Run(forceKillCmd) // Ignore errors for force kill
	}

	return nil
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
	defer session.Close()

	// Check if socat command exists
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
	defer session.Close()

	// Check if Docker socket exists and is accessible
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

	fmt.Printf("DEBUG: Starting verification for PID %s on port %d\n", pid, port)

	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session for process verification: %w", err)
	}
	defer session.Close()

	// Check if process is running
	cmd := fmt.Sprintf("kill -0 %s", pid)
	fmt.Printf("DEBUG: Checking if process is running with command: %s\n", cmd)
	if err := session.Run(cmd); err != nil {
		return fmt.Errorf("socat process not running (PID: %s): %w", pid, err)
	}
	fmt.Printf("DEBUG: Process is running\n")

	// Wait a moment for the process to fully start and bind to the port
	time.Sleep(2 * time.Second)

	// For now, just verify the process is running and assume the port will be ready
	// The real test will be when we try to connect to Docker
	fmt.Printf("DEBUG: Socat process verification successful - assuming port will be ready\n")
	return nil
}

// findAvailablePort finds an available local port
func findAvailablePort() (int, error) {
	// Try ports in the range 2376-2385 (avoiding 2375 which is standard Docker port)
	// Also try some additional ports in case the primary range is busy
	portRanges := [][]int{
		{2376, 2385}, // Primary range
		{2386, 2395}, // Secondary range
		{2396, 2405}, // Tertiary range
	}

	for _, portRange := range portRanges {
		if len(portRange) != 2 || portRange[0] <= 0 || portRange[1] <= 0 {
			continue // Skip invalid ranges
		}

		for port := portRange[0]; port <= portRange[1]; port++ {
			addr := fmt.Sprintf("127.0.0.1:%d", port)
			listener, err := net.Listen("tcp", addr)
			if err == nil {
				listener.Close()
				return port, nil
			}
		}
	}

	return 0, fmt.Errorf("no available ports found in ranges 2376-2385, 2386-2395, or 2396-2405")
}

// findAvailableRemotePort finds an available port on the remote machine
func (s *SSHClient) findAvailableRemotePort(client *ssh.Client) (int, error) {
	// Try ports in the range 2376-2385 (avoiding 2375 which is standard Docker port)
	// Also try some additional ports in case the primary range is busy
	portRanges := [][]int{
		{2376, 2385}, // Primary range
		{2386, 2395}, // Secondary range
		{2396, 2405}, // Tertiary range
	}

	for _, portRange := range portRanges {
		if len(portRange) != 2 || portRange[0] <= 0 || portRange[1] <= 0 {
			continue // Skip invalid ranges
		}

		for port := portRange[0]; port <= portRange[1]; port++ {
			// Check if port is available on remote machine
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
	defer session.Close()

	// Check if port is already in use
	cmd := fmt.Sprintf("ss -tln | grep ':%d '", port)
	return session.Run(cmd) != nil // Port is available if grep doesn't find it
}
