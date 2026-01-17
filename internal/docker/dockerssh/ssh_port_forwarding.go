package dockerssh

import (
	"errors"
	"fmt"
	"log/slog"
	"net"
	"time"

	"golang.org/x/crypto/ssh"
)

// SSHTunnelClient handles SSH tunneling with local TCP port creation on remote machine
type SSHTunnelClient struct {
	client           *ssh.Client
	log              *slog.Logger
	listener         net.Listener
	remoteLocalPort  int
	localPort        int
	remotePID        string
	connectionMethod string // Track which method was used (primary or fallback)
}

// NewSSHTunnelClient creates a new SSH tunnel client
func NewSSHTunnelClient(client *ssh.Client, log *slog.Logger) *SSHTunnelClient {
	return &SSHTunnelClient{
		client: client,
		log:    log,
	}
}

// SetupSSHTunnel sets up SSH tunneling with local TCP port creation on remote machine
func (s *SSHTunnelClient) SetupSSHTunnel(localPort int) error {
	if s.client == nil {
		return errors.New("SSH client is nil")
	}

	s.localPort = localPort
	s.log.Info("Setting up SSH tunnel with local TCP port creation", "localPort", localPort)

	// Step 1: Check if Docker TCP port exists on remote
	if err := s.checkExistingDockerTCPPort(); err != nil {
		s.log.Info("No existing Docker TCP port found, creating local TCP port on remote")

		// Step 2: Create local TCP port on remote using socat/nc/netcat
		if err := s.createLocalTCPPortOnRemote(); err != nil {
			return fmt.Errorf("failed to create local TCP port on remote: %w", err)
		}
	} else {
		s.log.Info("Found existing Docker TCP port on remote", "port", s.remoteLocalPort)
	}

	// Step 3: Create local listener
	listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", localPort))
	if err != nil {
		return fmt.Errorf("failed to create local listener on port %d: %w", localPort, err)
	}

	s.listener = listener
	s.log.Info("SSH tunnel established",
		"localPort", localPort,
		"remoteLocalPort", s.remoteLocalPort)

	// Step 4: Start handling connections in a goroutine
	go s.handleTunnelConnections(listener, localPort)

	return nil
}

// TestSSHTunnel tests if the SSH tunnel is working
func (s *SSHTunnelClient) TestSSHTunnel(localPort int) error {
	s.log.Info("Testing SSH tunnel", "localPort", localPort)

	// Give the tunnel a moment to establish
	time.Sleep(1 * time.Second)

	conn, err := s.connectToLocalPort(localPort)
	if err != nil {
		return err
	}
	defer func() {
		if err := conn.Close(); err != nil {
			s.log.Warn("Failed to close test connection", "error", err)
		}
	}()

	response, err := s.sendTestRequest(conn)
	if err != nil {
		return err
	}

	return s.validateResponse(response, localPort)
}

// Close closes the SSH tunnel and cleans up remote processes
func (s *SSHTunnelClient) Close() error {
	if s.listener != nil {
		if err := s.listener.Close(); err != nil {
			s.log.Warn("Failed to close tunnel listener", "error", err)
		}
	}

	// Clean up remote process if we created one
	if s.remotePID != "" && s.remotePID != "unknown" {
		if err := s.cleanupRemoteProcess(); err != nil {
			s.log.Warn("Failed to cleanup remote process", "error", err)
		}
	}

	return nil
}

// GetConnectionMethod returns the connection method name
func (s *SSHTunnelClient) GetConnectionMethod() string {
	if s.connectionMethod != "" {
		return s.connectionMethod
	}
	return "SSH Tunnel"
}
