package dockerssh

import (
	"log/slog"

	"golang.org/x/crypto/ssh"
)

// SSHContext represents the SSH connection context without socat
type SSHContext struct {
	client     *ssh.Client
	remoteHost string
	log        *slog.Logger
}

// GetRemoteHost returns the remote host address
func (s *SSHContext) GetRemoteHost() string {
	return s.remoteHost
}

// Close closes the SSH context
func (s *SSHContext) Close() error {
	if s == nil || s.client == nil {
		return nil
	}

	return s.client.Close()
}
