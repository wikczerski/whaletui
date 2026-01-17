// Package docker provides Docker client functionality for WhaleTUI.
//
//nolint:revive // Main client package legitimately needs multiple public types
package docker

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/wikczerski/whaletui/internal/config"
	"github.com/wikczerski/whaletui/internal/docker/dockerssh"
	"github.com/wikczerski/whaletui/internal/docker/services"
	domaintypes "github.com/wikczerski/whaletui/internal/docker/types"
	"github.com/wikczerski/whaletui/internal/docker/utils"
	"github.com/wikczerski/whaletui/internal/logger"
)

// Type aliases for backward compatibility
type (
	Container = domaintypes.Container
	Image     = domaintypes.Image
	Volume    = domaintypes.Volume
	Network   = domaintypes.Network
)

// Client represents a Docker client wrapper
type Client struct {
	cli     *client.Client
	cfg     *config.Config
	log     *slog.Logger
	sshConn *dockerssh.SSHConnection
	sshCtx  *dockerssh.SSHContext

	// Services for domain-specific operations
	Container *services.ContainerService
	Image     *services.ImageService
	Volume    *services.VolumeService
	Network   *services.NetworkService
	Swarm     *services.SwarmService
}

// New creates a new Docker client
func New(cfg *config.Config) (*Client, error) {
	if cfg == nil {
		return nil, errors.New("config cannot be nil")
	}

	log := logger.GetLogger()

	client, err := createDockerClient(cfg, log)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// createSSHDockerClient creates a Docker client via SSH connection with authentication
func createSSHDockerClient(cfg *config.Config, log *slog.Logger) (*Client, error) {
	log.Info("Attempting SSH connection for Docker client", "host", cfg.RemoteHost)

	// Extract host from SSH URL
	host, err := extractHostFromURL(cfg.RemoteHost)
	if err != nil {
		return nil, fmt.Errorf("failed to extract host from SSH URL: %w", err)
	}

	// Use authentication options if provided
	if cfg.SSHKeyPath != "" || cfg.SSHPassword != "" {
		sshConn, err := createSSHConnectionWithAuth(host, cfg.SSHKeyPath, cfg.SSHPassword, log)
		if err != nil {
			return nil, err
		}

		return tryCreateDockerClient(cfg, log, host, sshConn)
	}

	// Use SSH tunneling for remote connections
	log.Info("Using SSH tunneling for remote connection")
	return trySSHConnection(cfg, log, host)
}

// GetConnectionMethod returns the connection method used for this Docker client
func (c *Client) GetConnectionMethod() string {
	if c.sshConn != nil {
		return c.sshConn.GetConnectionMethod()
	}
	if c.sshCtx != nil {
		return "SSH Context"
	}
	return "Local Docker"
}

// Close closes the Docker client and cleans up SSH connections
func (c *Client) Close() error {
	var errors []string

	c.log.Info("Docker client closing, starting cleanup")

	c.closeDockerClient(&errors)
	c.closeSSHConnection(&errors)
	c.closeSSHContext(&errors)

	return c.buildCloseError(errors)
}

// GetInfo retrieves Docker system information
func (c *Client) GetInfo(ctx context.Context) (map[string]any, error) {
	if c.cli == nil {
		return nil, errors.New("docker client not initialized")
	}

	info, err := c.cli.Info(ctx)
	if err != nil {
		return nil, fmt.Errorf("docker info failed: %w", err)
	}
	return utils.MarshalToMap(info)
}

// InspectContainer inspects a container
func (c *Client) InspectContainer(ctx context.Context, id string) (map[string]any, error) {
	if c.Container == nil {
		return nil, fmt.Errorf("container service not initialized")
	}
	return c.Container.InspectContainer(ctx, id)
}

// GetContainerLogs retrieves container logs
func (c *Client) GetContainerLogs(ctx context.Context, id string) (string, error) {
	if c.Container == nil {
		return "", fmt.Errorf("container service not initialized")
	}
	return c.Container.GetContainerLogs(ctx, id)
}

// InspectImage inspects an image
func (c *Client) InspectImage(ctx context.Context, id string) (map[string]any, error) {
	if c.Image == nil {
		return nil, fmt.Errorf("image service not initialized")
	}
	return c.Image.InspectImage(ctx, id)
}

// InspectVolume inspects a volume
func (c *Client) InspectVolume(ctx context.Context, name string) (map[string]any, error) {
	if c.Volume == nil {
		return nil, fmt.Errorf("volume service not initialized")
	}
	return c.Volume.InspectVolume(ctx, name)
}

// RemoveVolume removes a volume
func (c *Client) RemoveVolume(ctx context.Context, name string, force bool) error {
	if c.Volume == nil {
		return fmt.Errorf("volume service not initialized")
	}
	return c.Volume.RemoveVolume(ctx, name, force)
}

// InspectNetwork inspects a network
func (c *Client) InspectNetwork(ctx context.Context, id string) (map[string]any, error) {
	if c.Network == nil {
		return nil, fmt.Errorf("network service not initialized")
	}
	return c.Network.InspectNetwork(ctx, id)
}

// ListContainers lists all containers
func (c *Client) ListContainers(ctx context.Context, all bool) ([]Container, error) {
	if c.Container == nil {
		return nil, fmt.Errorf("container service not initialized")
	}
	return c.Container.ListContainers(ctx, all)
}

// GetContainerStats gets container stats
func (c *Client) GetContainerStats(ctx context.Context, id string) (map[string]any, error) {
	if c.Container == nil {
		return nil, fmt.Errorf("container service not initialized")
	}
	return c.Container.GetContainerStats(ctx, id)
}

// ListImages lists all images
func (c *Client) ListImages(ctx context.Context) ([]Image, error) {
	if c.Image == nil {
		return nil, fmt.Errorf("image service not initialized")
	}
	return c.Image.ListImages(ctx)
}

// ListVolumes lists all volumes
func (c *Client) ListVolumes(ctx context.Context) ([]Volume, error) {
	if c.Volume == nil {
		return nil, fmt.Errorf("volume service not initialized")
	}
	return c.Volume.ListVolumes(ctx)
}

// ListNetworks lists all networks
func (c *Client) ListNetworks(ctx context.Context) ([]Network, error) {
	if c.Network == nil {
		return nil, fmt.Errorf("network service not initialized")
	}
	return c.Network.ListNetworks(ctx)
}

// StartContainer starts a container
func (c *Client) StartContainer(ctx context.Context, id string) error {
	if c.Container == nil {
		return fmt.Errorf("container service not initialized")
	}
	return c.Container.StartContainer(ctx, id)
}

// StopContainer stops a container
func (c *Client) StopContainer(ctx context.Context, id string, timeout *time.Duration) error {
	if c.Container == nil {
		return fmt.Errorf("container service not initialized")
	}
	return c.Container.StopContainer(ctx, id, timeout)
}

// RestartContainer restarts a container
func (c *Client) RestartContainer(ctx context.Context, id string, timeout *time.Duration) error {
	if c.Container == nil {
		return fmt.Errorf("container service not initialized")
	}
	return c.Container.RestartContainer(ctx, id, timeout)
}

// RemoveContainer removes a container
func (c *Client) RemoveContainer(ctx context.Context, id string, force bool) error {
	if c.Container == nil {
		return fmt.Errorf("container service not initialized")
	}
	return c.Container.RemoveContainer(ctx, id, force)
}

// ExecContainer executes a command in a running container and returns the output
func (c *Client) ExecContainer(
	ctx context.Context,
	id string,
	command []string,
	tty bool,
) (string, error) {
	if c.Container == nil {
		return "", fmt.Errorf("container service not initialized")
	}
	return c.Container.ExecContainer(ctx, id, command, tty)
}

// AttachContainer attaches to a running container
func (c *Client) AttachContainer(ctx context.Context, id string) (types.HijackedResponse, error) {
	if c.Container == nil {
		return types.HijackedResponse{}, fmt.Errorf("container service not initialized")
	}
	return c.Container.AttachContainer(ctx, id)
}

// RemoveImage removes an image
func (c *Client) RemoveImage(ctx context.Context, id string, force bool) error {
	if c.Image == nil {
		return fmt.Errorf("image service not initialized")
	}
	return c.Image.RemoveImage(ctx, id, force)
}

// RemoveNetwork removes a network
func (c *Client) RemoveNetwork(ctx context.Context, id string) error {
	if c.Network == nil {
		return fmt.Errorf("network service not initialized")
	}
	return c.Network.RemoveNetwork(ctx, id)
}
