package dockerssh

import (
	"fmt"
	"log/slog"
	"net"
	"strings"

	"golang.org/x/crypto/ssh"
)

// SSHConnection represents an active SSH connection
type SSHConnection struct {
	client         *ssh.Client
	session        *ssh.Session
	localPort      int
	remoteHost     string
	remotePort     int
	log            *slog.Logger
	listener       net.Listener     // Add listener for port forwarding
	connectionType string           // Type of connection method used
	tunnelClient   *SSHTunnelClient // SSH tunnel client for tunnel connections
}

// GetLocalProxyHost returns the local host:port for the Docker proxy
func (s *SSHConnection) GetLocalProxyHost() string {
	if s == nil {
		return "tcp://127.0.0.1:0"
	}
	port := s.localPort
	if port <= 0 {
		port = 2375 // Default Docker port
	}
	return fmt.Sprintf("tcp://127.0.0.1:%d", port)
}

// GetConnectionMethod returns the connection method used
func (s *SSHConnection) GetConnectionMethod() string {
	if s == nil {
		return "Unknown"
	}

	if s.connectionType != "" {
		return s.connectionType
	}

	// Default to SSH Tunnel
	return "SSH Tunnel"
}

// Close closes the SSH connection and cleans up resources
func (s *SSHConnection) Close() error {
	if s == nil {
		return nil
	}

	var errors []string
	s.closeListener(&errors)

	// Handle SSH tunnel connection
	if s.tunnelClient != nil {
		// SSH tunnel connection - close the tunnel client
		if err := s.tunnelClient.Close(); err != nil {
			errors = append(errors, fmt.Sprintf("tunnel client cleanup: %v", err))
		}
	}

	s.closeSSHClient(&errors)
	s.closeSession(&errors)

	return s.buildCloseError(errors)
}

// closeListener closes the port forwarding listener
func (s *SSHConnection) closeListener(errors *[]string) {
	if s.listener == nil {
		return
	}

	s.log.Info("Closing port forwarding listener")
	if err := s.listener.Close(); err != nil {
		*errors = append(*errors, fmt.Sprintf("listener close: %v", err))
	} else {
		s.log.Info("Port forwarding listener closed successfully")
	}
}

// closeSSHClient closes the SSH client
func (s *SSHConnection) closeSSHClient(errors *[]string) {
	if s.client == nil {
		return
	}

	s.log.Info("SSH client closing")
	if err := s.client.Close(); err != nil {
		*errors = append(*errors, fmt.Sprintf("SSH client close: %v", err))
	} else {
		s.log.Info("SSH client closed successfully")
	}
}

// closeSession closes the SSH session
func (s *SSHConnection) closeSession(errors *[]string) {
	if s.session == nil {
		return
	}

	s.log.Info("SSH session closing")
	if err := s.session.Close(); err != nil {
		*errors = append(*errors, fmt.Sprintf("SSH session close: %v", err))
	} else {
		s.log.Info("SSH session closed successfully")
	}
}

// buildCloseError builds the final error message from collected errors
func (s *SSHConnection) buildCloseError(errors []string) error {
	if len(errors) == 0 {
		s.log.Info("SSH connection closed successfully")
		return nil
	}

	s.log.Error("Failed to close SSH connection", "errors", strings.Join(errors, "; "))
	return fmt.Errorf("failed to close SSH connection: %s", strings.Join(errors, "; "))
}
