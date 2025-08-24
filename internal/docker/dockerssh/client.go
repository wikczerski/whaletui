// Package dockerssh provides SSH client functionality for Docker operations.
package dockerssh

import (
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
	config *ssh.ClientConfig
	host   string
	port   string
	user   string
	log    *slog.Logger
}

// NewSSHClient creates a new SSH client for the specified host
func NewSSHClient(host string, port int) (*SSHClient, error) {
	if host == "" {
		return nil, fmt.Errorf("SSH host cannot be empty")
	}

	user, hostname, sshPort, err := parseSSHHost(host)
	if err != nil {
		return nil, fmt.Errorf("failed to parse SSH host: %w", err)
	}

	sshKeyPath, err := getSSHKeyPath()
	if err != nil {
		return nil, fmt.Errorf("failed to get SSH key path: %w", err)
	}

	config, err := createSSHConfig(user, sshKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create SSH config: %w", err)
	}

	return &SSHClient{
		config: config,
		host:   hostname,
		port:   sshPort,
		user:   user,
		log:    logger.GetLogger(),
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

	if err := s.validateHostname(); err != nil {
		errors = append(errors, fmt.Sprintf("Hostname resolution: %v", err))
	}

	sshKeyPath, err := getSSHKeyPath()
	if err != nil {
		errors = append(errors, fmt.Sprintf("SSH key: %v", err))
	} else {
		if info, err := os.Stat(sshKeyPath); err == nil {
			if info.Mode().Perm()&0o077 != 0 {
				errors = append(errors, fmt.Sprintf("SSH key has overly permissive permissions %v (should be 600)", info.Mode().Perm()))
			}
		}
	}

	if s.config == nil {
		errors = append(errors, "SSH configuration not available")
		return fmt.Errorf("SSH connection diagnostics failed:\n%s", strings.Join(errors, "\n"))
	}

	client, err := ssh.Dial("tcp", net.JoinHostPort(s.host, s.port), s.config)
	if err != nil {
		errors = append(errors, fmt.Sprintf("SSH connection: %v", err))
		return fmt.Errorf("SSH connection diagnostics failed:\n%s", strings.Join(errors, "\n"))
	}
	defer func() {
		if err := client.Close(); err != nil {
			s.log.Warn("Failed to close SSH client", "error", err)
		}
	}()

	if err := s.checkDockerSocketAccess(client); err != nil {
		errors = append(errors, fmt.Sprintf("Docker socket access: %v", err))
	}

	if len(errors) > 0 {
		return fmt.Errorf("SSH connection diagnostics failed:\n%s", strings.Join(errors, "\n"))
	}

	return nil
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

	client, err := ssh.Dial("tcp", net.JoinHostPort(s.host, s.port), s.config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SSH server %s:%s: %w", s.host, s.port, err)
	}

	// Find an available local port
	localPort, err := findAvailablePort()
	if err != nil {
		if closeErr := client.Close(); closeErr != nil {
			s.log.Warn("Failed to close SSH client", "error", closeErr)
		}
		return nil, fmt.Errorf("failed to find available port: %w", err)
	}

	s.log.Info("Setting up SSH port forwarding", "localPort", localPort, "remoteSocket", "/var/run/docker.sock")

	// Set up SSH port forwarding from local port to remote Docker socket
	listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", localPort))
	if err != nil {
		if closeErr := client.Close(); closeErr != nil {
			s.log.Warn("Failed to close SSH client", "error", closeErr)
		}
		return nil, fmt.Errorf("failed to bind to local port %d: %w", localPort, err)
	}

	// Start the port forwarding in a goroutine
	go func() {
		defer listener.Close()
		for {
			localConn, err := listener.Accept()
			if err != nil {
				s.log.Warn("Failed to accept local connection", "error", err)
				return
			}

			// Create a new SSH session for each connection
			session, err := client.NewSession()
			if err != nil {
				s.log.Warn("Failed to create SSH session for forwarding", "error", err)
				localConn.Close()
				continue
			}

			// Set up stdin/stdout for the session
			session.Stdin = localConn
			session.Stdout = localConn
			session.Stderr = os.Stderr

			// Execute socat on the remote machine to forward to Docker socket
			if err := session.Run(fmt.Sprintf("socat STDIO UNIX-CONNECT:/var/run/docker.sock")); err != nil {
				s.log.Warn("Failed to execute socat on remote machine", "error", err)
				session.Close()
				localConn.Close()
				continue
			}

			session.Close()
			localConn.Close()
		}
	}()

	// Wait a moment for the listener to be ready
	time.Sleep(1 * time.Second)

	// Test the local port to make sure it's working
	testConn, err := net.DialTimeout("tcp", fmt.Sprintf("127.0.0.1:%d", localPort), 2*time.Second)
	if err != nil {
		if closeErr := client.Close(); closeErr != nil {
			s.log.Warn("Failed to close SSH client", "error", closeErr)
		}
		listener.Close()
		return nil, fmt.Errorf("local port forwarding test failed: %w", err)
	}
	testConn.Close()

	s.log.Info("SSH port forwarding established successfully", "localPort", localPort)

	return &SSHConnection{
		client:     client,
		session:    nil, // Not using a single session for this approach
		localPort:  localPort,
		remoteHost: s.host,
		remotePort: 0,  // Not using remote port for this approach
		socatPID:   "", // Not using socat PID for this approach
		log:        s.log,
		listener:   listener, // Store the listener for cleanup
	}, nil
}
