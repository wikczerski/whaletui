package docker

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/docker/docker/client"
	"github.com/wikczerski/whaletui/internal/config"
	"github.com/wikczerski/whaletui/internal/docker/dockerssh"
)

// extractHostFromURL extracts the host from a Docker URL or SSH connection string
func extractHostFromURL(url string) (string, error) {
	if url == "" {
		return "", errors.New("URL cannot be empty")
	}

	if host := extractSSHHost(url); host != "" {
		return host, nil
	}

	if host := extractTCPHost(url); host != "" {
		return host, nil
	}

	if host := extractUnixHost(url); host != "" {
		return host, nil
	}

	// If no prefix, assume it's a host string
	return url, nil
}

// extractSSHHost extracts host from SSH URLs
func extractSSHHost(url string) string {
	if !strings.HasPrefix(url, "ssh://") {
		return ""
	}

	hostPart := strings.TrimPrefix(url, "ssh://")
	// Extract host part (before any path)
	if slashIndex := strings.Index(hostPart, "/"); slashIndex != -1 {
		hostPart = hostPart[:slashIndex]
	}
	return hostPart
}

// extractTCPHost extracts host from TCP URLs
func extractTCPHost(url string) string {
	if !strings.HasPrefix(url, "tcp://") {
		return ""
	}
	return strings.TrimPrefix(url, "tcp://")
}

// extractUnixHost extracts host from Unix socket URLs
func extractUnixHost(url string) string {
	if !strings.HasPrefix(url, "unix://") {
		return ""
	}
	return "localhost"
}

// These functions are currently unused and have been removed to satisfy the linter

// trySSHConnection attempts to create a Docker client using SSH tunneling
func trySSHConnection(cfg *config.Config, log *slog.Logger, host string) (*Client, error) {
	if host == "" {
		return nil, errors.New("host cannot be empty")
	}

	sshConn, err := createSSHConnection(host, log)
	if err != nil {
		return nil, err
	}

	return tryCreateDockerClient(cfg, log, host, sshConn)
}

// tryCreateDockerClient attempts to create a Docker client with the given SSH connection
func tryCreateDockerClient(
	cfg *config.Config,
	log *slog.Logger,
	host string,
	sshConn *dockerssh.SSHConnection,
) (*Client, error) {
	dockerCli, err := createDockerClientForProxy(sshConn)
	if err != nil {
		cleanupSSHConnection(sshConn, log)
		return nil, err
	}

	if err := testDockerConnection(dockerCli); err != nil {
		cleanupDockerClient(dockerCli, sshConn, log)
		return nil, err
	}

	return createClient(cfg, log, host, sshConn, dockerCli), nil
}

// createSSHConnection creates an SSH connection using SSH tunneling
func createSSHConnection(host string, log *slog.Logger) (*dockerssh.SSHConnection, error) {
	// Parse the host to extract username, hostname, and port
	username, hostname, port, err := dockerssh.ParseSSHHost(host)
	if err != nil {
		return nil, fmt.Errorf("failed to parse SSH host: %w", err)
	}

	sshClient := dockerssh.NewSSHClient(hostname, port, username, log)

	// Try SSH tunneling connection
	sshConn, err := sshClient.ConnectWithFallback(2375)
	if err != nil {
		return nil, fmt.Errorf("failed to establish SSH connection: %w", err)
	}

	localProxyHost := sshConn.GetLocalProxyHost()
	connectionMethod := sshConn.GetConnectionMethod()
	log.Info("Created SSH connection",
		"localProxy", localProxyHost,
		"method", connectionMethod)
	log.Info("🔗 Connection Method", "method", connectionMethod)
	return sshConn, nil
}

// createSSHConnectionWithAuth creates an SSH connection with authentication options
func createSSHConnectionWithAuth(
	host string, keyPath, password string, log *slog.Logger,
) (*dockerssh.SSHConnection, error) {
	// Parse the host to extract username, hostname, and port
	username, hostname, port, err := dockerssh.ParseSSHHost(host)
	if err != nil {
		return nil, fmt.Errorf("failed to parse SSH host: %w", err)
	}

	sshClient := dockerssh.NewSSHClientWithAuth(hostname, port, username, keyPath, password, log)

	// Try SSH tunneling connection
	sshConn, err := sshClient.ConnectWithFallback(2375)
	if err != nil {
		return nil, fmt.Errorf("failed to establish SSH connection: %w", err)
	}

	localProxyHost := sshConn.GetLocalProxyHost()
	connectionMethod := sshConn.GetConnectionMethod()
	log.Info("Created SSH connection with authentication",
		"localProxy", localProxyHost,
		"method", connectionMethod,
		"keyPath", keyPath,
		"hasPassword", password != "")
	log.Info("🔗 Connection Method", "method", connectionMethod)
	return sshConn, nil
}

// createDockerClientForProxy creates a Docker client for the local proxy
func createDockerClientForProxy(sshConn *dockerssh.SSHConnection) (*client.Client, error) {
	localProxyHost := sshConn.GetLocalProxyHost()
	opts := []client.Opt{
		client.WithHost(localProxyHost),
		client.WithAPIVersionNegotiation(),
		client.WithTimeout(30 * time.Second),
	}

	dockerCli, err := client.NewClientWithOpts(opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create Docker client for local proxy: %w", err)
	}

	return dockerCli, nil
}

// testDockerConnection tests the Docker connection
func testDockerConnection(dockerCli *client.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if _, err := dockerCli.Ping(ctx); err != nil {
		return fmt.Errorf("failed to connect to Docker via SSH tunnel: %w", err)
	}

	return nil
}

// cleanupSSHConnection cleans up the SSH connection
func cleanupSSHConnection(sshConn *dockerssh.SSHConnection, log *slog.Logger) {
	if closeErr := sshConn.Close(); closeErr != nil {
		log.Warn("Failed to close SSH connection after Docker client creation failure",
			"error", closeErr)
	}
}

// cleanupDockerClient cleans up the Docker client and SSH connection
func cleanupDockerClient(
	dockerCli *client.Client,
	sshConn *dockerssh.SSHConnection,
	log *slog.Logger,
) {
	if closeErr := dockerCli.Close(); closeErr != nil {
		log.Warn("Failed to close Docker client after ping failure",
			"error", closeErr)
	}
	if closeErr := sshConn.Close(); closeErr != nil {
		log.Warn("Failed to close SSH connection after ping failure",
			"error", closeErr)
	}
}

// createClient creates the final client
func createClient(
	cfg *config.Config,
	log *slog.Logger,
	host string,
	sshConn *dockerssh.SSHConnection,
	dockerCli *client.Client,
) *Client {
	localProxyHost := sshConn.GetLocalProxyHost()
	log.Info("Successfully created Docker client via SSH tunnel",
		"host", host,
		"localProxy", localProxyHost)

	return &Client{
		cli:     dockerCli,
		cfg:     cfg,
		log:     log,
		sshConn: sshConn,
		sshCtx:  nil, // Not using direct SSH context
	}
}

// trySSHConnectionFromConfig is a simplified version for the client factory
func trySSHConnectionFromConfig(cfg *config.Config, log *slog.Logger) (*Client, error) {
	if cfg.RemoteHost == "" {
		return nil, errors.New("no remote host configured")
	}

	host, err := extractHostFromURL(cfg.RemoteHost)
	if err != nil {
		return nil, fmt.Errorf("failed to extract host: %w", err)
	}

	return trySSHConnection(cfg, log, host)
}
