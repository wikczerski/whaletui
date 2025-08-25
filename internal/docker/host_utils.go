package docker

import (
	"errors"
	"strings"
)

// formatRemoteHost formats and validates a remote host URL
func formatRemoteHost(host string) (string, error) {
	// Handle SSH URLs - they should be passed through as-is
	if strings.HasPrefix(host, "ssh://") {
		return host, nil
	}

	if err := validateRemoteHost(host); err != nil {
		return "", err
	}

	// Automatically add tcp:// prefix if user didn't provide a scheme
	if !strings.Contains(host, "://") {
		host = "tcp://" + host
	}

	return host, nil
}

// validateRemoteHost validates a remote host string
func validateRemoteHost(host string) error {
	if host == "" {
		return errors.New("remote host cannot be empty")
	}

	// Basic validation - could be expanded based on your needs
	if strings.Contains(host, " ") {
		return errors.New("remote host cannot contain spaces")
	}

	return nil
}
