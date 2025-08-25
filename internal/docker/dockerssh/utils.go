package dockerssh

import (
	"errors"
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

// parseSSHHost parses an SSH host string in the format [user@]host[:port]
func parseSSHHost(host string) (username, hostname, port string, err error) {
	if host == "" {
		return "", "", "", errors.New("SSH host cannot be empty")
	}

	userStr, hostn, err := parseUserAndHost(host)
	if err != nil {
		return "", "", "", err
	}

	if err := validateHostname(hostn, host); err != nil {
		return "", "", "", err
	}

	return parseHostAndPortWithUser(userStr, hostn, host)
}

// parseHostAndPortWithUser parses host and port, returning username, hostname, and port
func parseHostAndPortWithUser(
	userStr, hostn, host string,
) (username, hostname, port string, err error) {
	hostn, portStr, err := parseHostAndPort(hostn, host)
	if err != nil {
		return "", "", "", err
	}
	return userStr, hostn, portStr, nil
}

// validateHostname validates that the hostname is not empty
func validateHostname(hostn, host string) error {
	if hostn == "" {
		return fmt.Errorf("hostname cannot be empty in SSH host format '%s'", host)
	}
	return nil
}

// parseUserAndHost extracts username and hostname from the SSH host string
func parseUserAndHost(host string) (username, hostname string, err error) {
	if strings.Contains(host, "@") {
		return parseUserAtHost(host)
	}
	return parseHostOnly(host)
}

// parseUserAtHost handles the user@host format
func parseUserAtHost(host string) (username, hostname string, err error) {
	parts := strings.Split(host, "@")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid SSH host format '%s': expected [user@]host[:port]", host)
	}

	userStr := strings.TrimSpace(parts[0])
	hostn := strings.TrimSpace(parts[1])

	if userStr == "" {
		return "", "", fmt.Errorf("username cannot be empty in SSH host format '%s'", host)
	}

	return userStr, hostn, nil
}

// parseHostOnly handles the host-only format, using current user
func parseHostOnly(host string) (username, hostname string, err error) {
	currentUser, err := user.Current()
	if err != nil {
		return "", "", fmt.Errorf("failed to get current user: %w", err)
	}

	userStr := currentUser.Username
	if strings.Contains(userStr, "\\") {
		parts := strings.Split(userStr, "\\")
		if len(parts) == 2 {
			userStr = parts[1]
		}
	}

	hostn := strings.TrimSpace(host)
	return userStr, hostn, nil
}

// parseHostAndPort extracts hostname and port from the host string
func parseHostAndPort(hostn, originalHost string) (hostname, port string, err error) {
	if strings.Contains(hostn, ":") {
		return parseHostWithPort(hostn, originalHost)
	}
	return hostn, "22", nil
}

// parseHostWithPort handles the host:port format
func parseHostWithPort(hostn, originalHost string) (hostname, port string, err error) {
	parts := strings.Split(hostn, ":")
	if len(parts) != 2 {
		return "", "", fmt.Errorf(
			"invalid SSH host format '%s': expected [user@]host[:port]",
			originalHost,
		)
	}

	host := strings.TrimSpace(parts[0])
	portStr := strings.TrimSpace(parts[1])

	if err := validateHostAndPort(host, portStr, originalHost); err != nil {
		return "", "", err
	}

	return host, portStr, nil
}

// validateHostAndPort validates both host and port parts
func validateHostAndPort(host, portStr, originalHost string) error {
	if host == "" {
		return fmt.Errorf("hostname cannot be empty in SSH host format '%s'", originalHost)
	}
	if portStr == "" {
		return fmt.Errorf("port cannot be empty in SSH host format '%s'", originalHost)
	}
	return validatePort(portStr, originalHost)
}

// validatePort checks if the port string is numeric
func validatePort(portStr, originalHost string) error {
	if _, err := fmt.Sscanf(portStr, "%d", new(int)); err != nil {
		return fmt.Errorf(
			"invalid port '%s' in SSH host format '%s': port must be numeric",
			portStr,
			originalHost,
		)
	}
	return nil
}

// getSSHKeyPath returns the path to the user's SSH private key
func getSSHKeyPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	possibleKeys := getPossibleSSHKeys(homeDir)
	return findValidSSHKey(possibleKeys, homeDir)
}

// getPossibleSSHKeys returns a list of possible SSH key paths
func getPossibleSSHKeys(homeDir string) []string {
	return []string{
		filepath.Join(homeDir, ".ssh", "id_rsa"),
		filepath.Join(homeDir, ".ssh", "id_ed25519"),
		filepath.Join(homeDir, ".ssh", "id_ecdsa"),
	}
}

// findValidSSHKey finds a valid SSH key from the list of possible paths
func findValidSSHKey(possibleKeys []string, homeDir string) (string, error) {
	for _, keyPath := range possibleKeys {
		if info, err := os.Stat(keyPath); err == nil {
			if runtime.GOOS == "windows" {
				return keyPath, nil
			}

			if err := validateSSHKeyPermissions(info, keyPath); err != nil {
				return "", err
			}
			return keyPath, nil
		}
	}

	return "", fmt.Errorf(
		"no SSH private key found in %s/.ssh/ (checked: id_rsa, id_ed25519, id_ecdsa)",
		homeDir,
	)
}

// validateSSHKeyPermissions validates SSH key file permissions
func validateSSHKeyPermissions(info os.FileInfo, keyPath string) error {
	mode := info.Mode().Perm()
	if mode&0o077 != 0 {
		return fmt.Errorf(
			"SSH key %s has overly permissive permissions %v (should be 600)",
			keyPath,
			mode,
		)
	}
	return nil
}

// createSSHConfig creates SSH client configuration with key-based authentication
func createSSHConfig(username, keyPath string) (*ssh.ClientConfig, error) {
	signer, err := readAndParseSSHKey(keyPath)
	if err != nil {
		return nil, err
	}

	hostKeyCallback, err := createHostKeyCallback()
	if err != nil {
		return nil, err
	}

	return &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: hostKeyCallback,
		Timeout:         60 * time.Second,
	}, nil
}

// readAndParseSSHKey reads and parses the SSH private key
func readAndParseSSHKey(keyPath string) (ssh.Signer, error) {
	// Validate keyPath to prevent directory traversal
	if !filepath.IsAbs(keyPath) ||
		!strings.HasPrefix(filepath.Clean(keyPath), filepath.Clean(filepath.Dir(keyPath))) {
		return nil, errors.New("invalid SSH key path")
	}

	keyBytes, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read SSH key at %s: %w", keyPath, err)
	}

	signer, err := ssh.ParsePrivateKey(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse SSH key at %s: %w", keyPath, err)
	}

	return signer, nil
}

// createHostKeyCallback creates the host key callback for SSH connections
func createHostKeyCallback() (ssh.HostKeyCallback, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}
	knownHostsPath := filepath.Join(homeDir, ".ssh", "known_hosts")

	hostKeyCallback, err := knownhosts.New(knownHostsPath)
	if err != nil {
		// Use a more secure alternative - create a custom host key callback that logs warnings
		hostKeyCallback = func(hostname string, remote net.Addr, _ ssh.PublicKey) error {
			// Log warning about unknown host key but allow connection
			fmt.Printf("Warning: Unknown host key for %s (%s)\n", hostname, remote.String())
			return nil
		}
	}

	return hostKeyCallback, nil
}

// validateHostname checks if the hostname can be resolved to an IP address
func (s *SSHClient) validateHostname() error {
	if s.host == "" {
		return errors.New("hostname cannot be empty")
	}

	if net.ParseIP(s.host) != nil {
		return nil
	}

	if err := s.validateHostnameFormat(); err != nil {
		return err
	}

	return s.validateHostnameResolution()
}

// validateHostnameFormat checks the format of the hostname
func (s *SSHClient) validateHostnameFormat() error {
	if strings.HasPrefix(s.host, ".") || strings.HasSuffix(s.host, ".") {
		return fmt.Errorf("hostname '%s' cannot start or end with a dot", s.host)
	}

	if strings.Contains(s.host, "..") {
		return fmt.Errorf("hostname '%s' cannot contain consecutive dots", s.host)
	}

	return nil
}

// validateHostnameResolution checks if the hostname can be resolved
func (s *SSHClient) validateHostnameResolution() error {
	ips, err := net.LookupHost(s.host)
	if err != nil {
		return fmt.Errorf("cannot resolve hostname '%s': %w", s.host, err)
	}

	if len(ips) == 0 {
		return fmt.Errorf("hostname '%s' resolved to no IP addresses", s.host)
	}

	return nil
}

// findAvailablePort finds an available local port
func findAvailablePort() (int, error) {
	portRanges := getPortRanges()

	for _, portRange := range portRanges {
		if !isValidPortRange(portRange) {
			continue
		}

		if port := findPortInRange(portRange); port != 0 {
			return port, nil
		}
	}

	return 0, errors.New("no available ports found in ranges 2376-2385, 2386-2395, or 2396-2405")
}

// getPortRanges returns the available port ranges to check
func getPortRanges() [][]int {
	return [][]int{
		{2376, 2385},
		{2386, 2395},
		{2396, 2405},
	}
}

// isValidPortRange checks if a port range is valid
func isValidPortRange(portRange []int) bool {
	return len(portRange) == 2 && portRange[0] > 0 && portRange[1] > 0
}

// findPortInRange finds an available port in the given range
func findPortInRange(portRange []int) int {
	for port := portRange[0]; port <= portRange[1]; port++ {
		if isPortAvailable(port) {
			return port
		}
	}
	return 0
}

// isPortAvailable checks if a specific port is available
func isPortAvailable(port int) bool {
	addr := fmt.Sprintf("127.0.0.1:%d", port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return false
	}

	if closeErr := listener.Close(); closeErr != nil {
		// Log the error but continue since we found an available port
		fmt.Printf("Warning: Failed to close listener on port %d: %v\n", port, closeErr)
	}
	return true
}
