package dockerssh

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
)

// ParseSSHHost parses an SSH host string in the format [user@]host[:port]
func ParseSSHHost(host string) (username, hostname, port string, err error) {
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
