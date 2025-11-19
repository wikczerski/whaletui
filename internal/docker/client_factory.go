package docker

import (
	"context"
	"fmt"
	"log/slog"
	"runtime"
	"strings"
	"time"

	"github.com/docker/docker/client"

	"github.com/wikczerski/whaletui/internal/config"
	"github.com/wikczerski/whaletui/internal/docker/services"
	"github.com/wikczerski/whaletui/internal/docker/utils"
)

// createDockerClient creates a Docker client with the given configuration
func createDockerClient(cfg *config.Config, log *slog.Logger) (*Client, error) {
	// For SSH connections, try direct connection first
	if cfg.RemoteHost != "" && strings.HasPrefix(cfg.RemoteHost, "ssh://") {
		return createSSHDockerClient(cfg, log)
	}

	opts, err := buildClientOptions(cfg)
	if err != nil {
		return nil, err
	}

	cli, err := client.NewClientWithOpts(opts...)
	if err != nil {
		return handleClientCreationError(cfg, log, err)
	}

	return testAndCreateClient(cfg, log, cli)
}

// buildClientOptions builds the client options based on configuration
func buildClientOptions(cfg *config.Config) ([]client.Opt, error) {
	var opts []client.Opt

	if cfg.RemoteHost != "" {
		// Check if this is an SSH connection
		if strings.HasPrefix(cfg.RemoteHost, "ssh://") {
			// For SSH connections, we'll handle this specially
			// The Docker client will handle ssh:// URLs natively
			opts = append(opts, client.WithHost(cfg.RemoteHost), client.WithTimeout(30*time.Second))
		} else {
			dockerHost, err := formatRemoteHost(cfg.RemoteHost)
			if err != nil {
				return nil, fmt.Errorf("invalid remote host format: %w", err)
			}
			opts = append(opts, client.WithHost(dockerHost), client.WithTimeout(30*time.Second))
		}
	} else if cfg.DockerHost != "" {
		opts = append(opts, client.WithHost(cfg.DockerHost))
	}

	opts = append(opts,
		client.WithAPIVersionNegotiation(),
		client.WithTimeout(10*time.Second),
	)

	return opts, nil
}

// handleClientCreationError handles errors during client creation, including Windows auto-detection
func handleClientCreationError(cfg *config.Config, log *slog.Logger, err error) (*Client, error) {
	if runtime.GOOS == "windows" && cfg.RemoteHost == "" {
		return tryWindowsAutoDetection(cfg, log)
	}
	return nil, fmt.Errorf("failed to create Docker client: %w", err)
}

// tryWindowsAutoDetection attempts to auto-detect the correct Docker host on Windows
func tryWindowsAutoDetection(cfg *config.Config, log *slog.Logger) (*Client, error) {
	log.Warn("Docker client creation failed, attempting to auto-detect correct host...")

	detectedHost, detectErr := utils.DetectWindowsDockerHost(log)
	if detectErr != nil {
		return nil, fmt.Errorf("failed to create Docker client: %w", detectErr)
	}

	log.Info("Detected working Docker host", "host", detectedHost)
	return createClientWithHost(cfg, log, detectedHost)
}

// createClientWithHost creates a client with a specific host
func createClientWithHost(cfg *config.Config, log *slog.Logger, host string) (*Client, error) {
	detectedOpts := []client.Opt{
		client.WithHost(host),
		client.WithAPIVersionNegotiation(),
		client.WithTimeout(10 * time.Second),
	}

	detectedCli, detectedErr := client.NewClientWithOpts(detectedOpts...)
	if detectedErr != nil {
		return nil, fmt.Errorf("failed to create Docker client with detected host: %w", detectedErr)
	}

	return testAndCreateClient(cfg, log, detectedCli)
}

// testAndCreateClient tests the connection and creates the client
func testAndCreateClient(
	cfg *config.Config,
	log *slog.Logger,
	cli *client.Client,
) (*Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if _, err := cli.Ping(ctx); err != nil {
		return handleConnectionFailure(cfg, log, cli, err)
	}

	client := createDockerClientInstance(cli, cfg, log)
	logConnectionSuccess(cfg, log)
	return client, nil
}

// handleConnectionFailure handles connection failure and attempts SSH connection
func handleConnectionFailure(
	cfg *config.Config,
	log *slog.Logger,
	cli *client.Client,
	err error,
) (*Client, error) {
	if closeErr := cli.Close(); closeErr != nil {
		log.Warn("Failed to close Docker client", "error", closeErr)
	}

	// If this is a remote host and direct connection failed, try SSH connection
	// Only try SSH connection if it's not already an SSH connection
	if cfg.RemoteHost != "" && !strings.HasPrefix(cfg.RemoteHost, "ssh://") {
		return trySSHConnectionFromConfig(cfg, log)
	}

	return nil, fmt.Errorf("failed to connect to Docker: %w", err)
}

// createDockerClientInstance creates a new client instance
func createDockerClientInstance(cli *client.Client, cfg *config.Config, log *slog.Logger) *Client {
	return &Client{
		cli:       cli,
		cfg:       cfg,
		log:       log,
		Container: services.NewContainerService(cli, log),
		Image:     services.NewImageService(cli, log),
		Volume:    services.NewVolumeService(cli, log),
		Network:   services.NewNetworkService(cli, log),
		Swarm:     services.NewSwarmService(cli, log),
	}
}

// logConnectionSuccess logs the successful connection
func logConnectionSuccess(cfg *config.Config, log *slog.Logger) {
	if cfg.RemoteHost != "" {
		log.Info("Successfully connected to remote Docker host", "host", cfg.RemoteHost)
	} else {
		log.Info("Successfully connected to local Docker instance")
	}
}
