package docker

import (
	"context"
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
		return "", fmt.Errorf("URL cannot be empty")
	}

	// Handle SSH URLs
	if strings.HasPrefix(url, "ssh://") {
		// Remove ssh:// prefix
		hostPart := strings.TrimPrefix(url, "ssh://")
		// Extract host part (before any path)
		if slashIndex := strings.Index(hostPart, "/"); slashIndex != -1 {
			hostPart = hostPart[:slashIndex]
		}
		return hostPart, nil
	}

	// Handle TCP URLs
	if strings.HasPrefix(url, "tcp://") {
		hostPart := strings.TrimPrefix(url, "tcp://")
		return hostPart, nil
	}

	// Handle Unix socket URLs
	if strings.HasPrefix(url, "unix://") {
		return "localhost", nil
	}

	// If no prefix, assume it's a host string
	return url, nil
}

// establishSSHConnection establishes an SSH connection using the dockerssh package
func establishSSHConnection(host, _ string, _ *slog.Logger) (*dockerssh.SSHContext, error) {
	if host == "" {
		return nil, fmt.Errorf("host cannot be empty")
	}

	// Create SSH client
	sshClient, err := dockerssh.NewSSHClient(host, 22)
	if err != nil {
		return nil, fmt.Errorf("failed to create SSH client: %w", err)
	}

	// Connect to establish the connection
	sshCtx, err := sshClient.Connect()
	if err != nil {
		return nil, fmt.Errorf("failed to establish SSH connection: %w", err)
	}

	return sshCtx, nil
}

// createDockerClientViaSSH creates a Docker client via SSH connection
func createDockerClientViaSSH(sshCtx *dockerssh.SSHContext, log *slog.Logger) (interface{}, error) {
	if sshCtx == nil {
		return nil, fmt.Errorf("SSH context cannot be nil")
	}

	// For now, return nil as this would require implementing the actual Docker client creation
	// This is a placeholder - the actual implementation would depend on how you want to
	// create Docker clients over SSH
	log.Warn("Direct SSH Docker client creation not yet implemented")
	return nil, fmt.Errorf("direct SSH Docker client creation not yet implemented")
}

// trySSHFallbackWithSocat attempts to create a Docker client using SSH with socat proxy
func trySSHFallbackWithSocat(cfg *config.Config, log *slog.Logger, host string) (*Client, error) {
	if host == "" {
		return nil, fmt.Errorf("host cannot be empty")
	}

	// Create SSH client
	sshClient, err := dockerssh.NewSSHClient(host, 22)
	if err != nil {
		return nil, fmt.Errorf("failed to create SSH client: %w", err)
	}

	// Try to connect with socat
	sshConn, err := sshClient.ConnectWithSocat(2375)
	if err != nil {
		return nil, fmt.Errorf("failed to establish SSH connection with socat: %w", err)
	}

	// Create Docker client using the local proxy
	localProxyHost := sshConn.GetLocalProxyHost()
	log.Info("Created SSH connection with socat", "localProxy", localProxyHost)

	// Create Docker client options for the local proxy
	opts := []client.Opt{
		client.WithHost(localProxyHost),
		client.WithAPIVersionNegotiation(),
		client.WithTimeout(30 * time.Second),
	}

	// Create the actual Docker client
	dockerCli, err := client.NewClientWithOpts(opts...)
	if err != nil {
		// Clean up SSH connection if Docker client creation fails
		if closeErr := sshConn.Close(); closeErr != nil {
			log.Warn("Failed to close SSH connection after Docker client creation failure", "error", closeErr)
		}
		return nil, fmt.Errorf("failed to create Docker client for local proxy: %w", err)
	}

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if _, err := dockerCli.Ping(ctx); err != nil {
		// Clean up if connection test fails
		if closeErr := dockerCli.Close(); closeErr != nil {
			log.Warn("Failed to close Docker client after ping failure", "error", closeErr)
		}
		if closeErr := sshConn.Close(); closeErr != nil {
			log.Warn("Failed to close SSH connection after ping failure", "error", closeErr)
		}
		return nil, fmt.Errorf("failed to connect to Docker via socat proxy: %w", err)
	}

	// Create and return the client with all fields properly initialized
	client := &Client{
		cli:     dockerCli,
		cfg:     cfg,
		log:     log,
		sshConn: sshConn,
		sshCtx:  nil, // Not using direct SSH context
	}

	log.Info("Successfully created Docker client via SSH socat proxy", "host", host, "localProxy", localProxyHost)
	return client, nil
}

// trySSHFallback is a simplified version for the client factory
func trySSHFallback(cfg *config.Config, log *slog.Logger) (*Client, error) {
	if cfg.RemoteHost == "" {
		return nil, fmt.Errorf("no remote host configured")
	}

	host, err := extractHostFromURL(cfg.RemoteHost)
	if err != nil {
		return nil, fmt.Errorf("failed to extract host: %w", err)
	}

	return trySSHFallbackWithSocat(cfg, log, host)
}
