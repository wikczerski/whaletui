// Package docker provides Docker client functionality for WhaleTUI.
//
//nolint:revive // Main client package legitimately needs multiple public types
package docker

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sort"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/volume"
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
	if c.cli == nil {
		return nil, errors.New("docker client not initialized")
	}

	container, err := c.cli.ContainerInspect(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("container inspect failed %s: %w", id, err)
	}

	return utils.MarshalToMap(container)
}

// GetContainerLogs retrieves container logs
func (c *Client) GetContainerLogs(ctx context.Context, id string) (string, error) {
	if c.cli == nil {
		return "", errors.New("docker client not initialized")
	}

	logs, err := c.cli.ContainerLogs(ctx, id, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Timestamps: true,
		Follow:     false,
		Tail:       "100",
	})
	if err != nil {
		return "", fmt.Errorf("container logs failed %s: %w", id, err)
	}
	defer func() {
		if err := logs.Close(); err != nil {
			c.log.Warn("Failed to close logs", "error", err)
		}
	}()

	return c.readAndFormatLogs(logs), nil
}

// InspectImage inspects an image
func (c *Client) InspectImage(ctx context.Context, id string) (map[string]any, error) {
	imageInfo, err := c.cli.ImageInspect(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("image inspect failed %s: %w", id, err)
	}
	return utils.MarshalToMap(imageInfo)
}

// InspectVolume inspects a volume
func (c *Client) InspectVolume(ctx context.Context, name string) (map[string]any, error) {
	volumeInfo, err := c.cli.VolumeInspect(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("volume inspect failed %s: %w", name, err)
	}
	return utils.MarshalToMap(volumeInfo)
}

// RemoveVolume removes a volume
func (c *Client) RemoveVolume(ctx context.Context, name string, force bool) error {
	if err := utils.ValidateID(name, "volume name"); err != nil {
		return err
	}

	if err := c.cli.VolumeRemove(ctx, name, force); err != nil {
		return fmt.Errorf("failed to remove volume %s: %w", name, err)
	}

	return nil
}

// InspectNetwork inspects a network
func (c *Client) InspectNetwork(ctx context.Context, id string) (map[string]any, error) {
	networkInfo, err := c.cli.NetworkInspect(ctx, id, network.InspectOptions{})
	if err != nil {
		return nil, fmt.Errorf("network inspect failed %s: %w", id, err)
	}
	return utils.MarshalToMap(networkInfo)
}

// ListContainers lists all containers
func (c *Client) ListContainers(ctx context.Context, all bool) ([]Container, error) {
	if c.cli == nil {
		return nil, errors.New("docker client not initialized")
	}

	containers, err := c.getContainerList(ctx, all)
	if err != nil {
		return nil, err
	}

	result := c.convertToContainers(containers)
	utils.SortContainersByCreationTime(result)
	return result, nil
}

// GetContainerStats gets container stats
func (c *Client) GetContainerStats(ctx context.Context, id string) (map[string]any, error) {
	stats, err := c.cli.ContainerStats(ctx, id, false)
	if err != nil {
		return nil, fmt.Errorf("failed to get container stats: %w", err)
	}
	defer func() {
		if err := stats.Body.Close(); err != nil {
			c.log.Warn("Failed to close stats body", "error", err)
		}
	}()

	return c.decodeStatsResponse(stats.Body)
}

// ListImages lists all images
func (c *Client) ListImages(ctx context.Context) ([]Image, error) {
	images, err := c.getImageList(ctx)
	if err != nil {
		return nil, err
	}

	result := c.convertToImages(images)
	utils.SortImagesByCreationTime(result)
	return result, nil
}

// ListVolumes lists all volumes
func (c *Client) ListVolumes(ctx context.Context) ([]Volume, error) {
	vols, err := c.cli.VolumeList(ctx, volume.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list volumes: %w", err)
	}

	result := make([]Volume, 0, len(vols.Volumes))
	for _, vol := range vols.Volumes {
		volume := c.createVolumeFromAPI(vol)
		result = append(result, volume)
	}

	utils.SortVolumesByName(result)
	return result, nil
}

// ListNetworks lists all networks
func (c *Client) ListNetworks(ctx context.Context) ([]Network, error) {
	networks, err := c.cli.NetworkList(ctx, network.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list networks: %w", err)
	}

	result := make([]Network, 0, len(networks))
	for i := range networks {
		net := &networks[i]
		result = append(result, Network{
			ID:         net.ID[:12],
			Name:       net.Name,
			Driver:     net.Driver,
			Scope:      net.Scope,
			Created:    net.Created,
			Internal:   net.Internal,
			Attachable: net.Attachable,
			Ingress:    net.Ingress,
			IPv6:       net.EnableIPv6,
			EnableIPv6: net.EnableIPv6,
			Labels:     net.Labels,
			Containers: len(net.Containers),
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})

	return result, nil
}

// StartContainer starts a container
func (c *Client) StartContainer(ctx context.Context, id string) error {
	return c.cli.ContainerStart(ctx, id, container.StartOptions{})
}

// StopContainer stops a container
func (c *Client) StopContainer(ctx context.Context, id string, timeout *time.Duration) error {
	opts := utils.BuildStopOptions(timeout)
	return c.cli.ContainerStop(ctx, id, opts)
}

// RestartContainer restarts a container
func (c *Client) RestartContainer(ctx context.Context, id string, timeout *time.Duration) error {
	opts := utils.BuildStopOptions(timeout)
	return c.cli.ContainerRestart(ctx, id, opts)
}

// RemoveContainer removes a container
func (c *Client) RemoveContainer(ctx context.Context, id string, force bool) error {
	opts := container.RemoveOptions{
		Force: force,
	}

	if err := utils.ValidateID(id, "container ID"); err != nil {
		return err
	}

	return c.cli.ContainerRemove(ctx, id, opts)
}

// ExecContainer executes a command in a running container and returns the output
func (c *Client) ExecContainer(
	ctx context.Context,
	id string,
	command []string,
	_ bool,
) (string, error) {
	if err := c.validateClient(); err != nil {
		return "", err
	}

	execResp, err := c.createExecInstance(ctx, id, command)
	if err != nil {
		return "", err
	}

	output, err := c.executeAndCollectOutput(ctx, execResp.ID)
	if err != nil {
		return "", err
	}

	return output, nil
}

// AttachContainer attaches to a running container
func (c *Client) AttachContainer(ctx context.Context, id string) (types.HijackedResponse, error) {
	if err := c.validateClient(); err != nil {
		return types.HijackedResponse{}, err
	}

	attachConfig := container.AttachOptions{
		Stream: true,
		Stdin:  true,
		Stdout: true,
		Stderr: true,
		Logs:   false,
	}

	if err := utils.ValidateID(id, "container ID"); err != nil {
		return types.HijackedResponse{}, err
	}

	response, err := c.cli.ContainerAttach(ctx, id, attachConfig)
	if err != nil {
		return types.HijackedResponse{}, fmt.Errorf("failed to attach to container: %w", err)
	}

	return response, nil
}

// RemoveImage removes an image
func (c *Client) RemoveImage(ctx context.Context, id string, force bool) error {
	if err := utils.ValidateID(id, "image ID"); err != nil {
		return err
	}

	opts := image.RemoveOptions{
		Force:         force,
		PruneChildren: true, // Remove dependent images by default
	}

	_, err := c.cli.ImageRemove(ctx, id, opts)
	if err != nil {
		return fmt.Errorf("failed to remove image %s: %w", id, err)
	}

	return nil
}

// RemoveNetwork removes a network
func (c *Client) RemoveNetwork(ctx context.Context, id string) error {
	if err := utils.ValidateID(id, "network ID"); err != nil {
		return err
	}

	if err := c.cli.NetworkRemove(ctx, id); err != nil {
		return fmt.Errorf("failed to remove network %s: %w", id, err)
	}
	return nil
}

// validateClient validates that the Docker client is initialized
func (c *Client) validateClient() error {
	if c.cli != nil {
		return nil
	}
	return errors.New("docker client is not initialized")
}
