package docker

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"github.com/wikczerski/D5r/internal/config"
	"github.com/wikczerski/D5r/internal/logger"
)

// Client represents a Docker client wrapper
type Client struct {
	cli *client.Client
	cfg *config.Config
	log *logger.Logger
}

// detectWindowsDockerHost attempts to find the correct Docker host on Windows
func detectWindowsDockerHost(log *logger.Logger) (string, error) {
	if runtime.GOOS != "windows" {
		return "", fmt.Errorf("not on Windows")
	}

	// Common Windows Docker Desktop pipe paths
	possibleHosts := []string{
		"npipe:////./pipe/dockerDesktopLinuxEngine", // Linux containers
		"npipe:////./pipe/docker_engine",            // Windows containers
		"npipe:////./pipe/dockerDesktopEngine",      // Legacy Docker Desktop
	}

	for _, host := range possibleHosts {
		opts := []client.Opt{
			client.WithHost(host),
			client.WithAPIVersionNegotiation(),
			client.WithTimeout(5 * time.Second),
		}

		cli, err := client.NewClientWithOpts(opts...)
		if err != nil {
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		if _, err = cli.Ping(ctx); err == nil {
			cancel()
			if closeErr := cli.Close(); closeErr != nil {
				log.Warn("Failed to close Docker client during host detection: %v", closeErr)
			}
			return host, nil
		}
		cancel()
		if closeErr := cli.Close(); closeErr != nil {
			log.Warn("Failed to close Docker client during host detection: %v", closeErr)
		}
	}

	return "", fmt.Errorf("no working Docker host found")
}

// New creates a new Docker client
func New(cfg *config.Config) (*Client, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	log := logger.GetLogger()
	log.SetPrefix("Docker")

	var opts []client.Opt

	// Handle remote host connection
	if cfg.RemoteHost != "" {
		log.Info("Connecting to remote Docker host: %s", cfg.RemoteHost)

		// Validate remote host format
		if err := validateRemoteHost(cfg.RemoteHost); err != nil {
			return nil, fmt.Errorf("invalid remote host format: %w", err)
		}

		opts = append(opts, client.WithHost(cfg.RemoteHost), client.WithTimeout(30*time.Second))
	} else if cfg.DockerHost != "" {
		opts = append(opts, client.WithHost(cfg.DockerHost))
	}

	opts = append(opts,
		client.WithAPIVersionNegotiation(),
		client.WithTimeout(10*time.Second),
	)

	cli, err := client.NewClientWithOpts(opts...)
	if err != nil {
		// On Windows, try to auto-detect the correct Docker host if client creation fails
		if runtime.GOOS == "windows" && cfg.RemoteHost == "" {
			log.Warn("Docker client creation failed, attempting to auto-detect correct host...")

			detectedHost, detectErr := detectWindowsDockerHost(log)
			if detectErr == nil {
				log.Info("Detected working Docker host: %s", detectedHost)

				// Try to connect with the detected host
				detectedOpts := []client.Opt{
					client.WithHost(detectedHost),
					client.WithAPIVersionNegotiation(),
					client.WithTimeout(10 * time.Second),
				}

				detectedCli, detectedErr := client.NewClientWithOpts(detectedOpts...)
				if detectedErr == nil {
					detectedCtx, detectedCancel := context.WithTimeout(context.Background(), 10*time.Second)
					if _, detectedPingErr := detectedCli.Ping(detectedCtx); detectedPingErr == nil {
						detectedCancel()
						log.Info("Successfully connected with auto-detected host")

						client := &Client{
							cli: detectedCli,
							cfg: cfg,
							log: log,
						}

						return client, nil
					}
					detectedCancel()
					if closeErr := detectedCli.Close(); closeErr != nil {
						log.Warn("Failed to close detected Docker client: %v", closeErr)
					}
				}
			}
		}

		return nil, fmt.Errorf("failed to create Docker client: %w", err)
	}

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if _, err := cli.Ping(ctx); err != nil {
		if closeErr := cli.Close(); closeErr != nil {
			log.Warn("Failed to close Docker client: %v", closeErr)
		}
		if cfg.RemoteHost != "" {
			return nil, fmt.Errorf("failed to connect to remote Docker host '%s': %w", cfg.RemoteHost, err)
		}
		return nil, fmt.Errorf("failed to connect to Docker: %w", err)
	}

	client := &Client{
		cli: cli,
		cfg: cfg,
		log: log,
	}

	if cfg.RemoteHost != "" {
		log.Info("Successfully connected to remote Docker host: %s", cfg.RemoteHost)
	} else {
		log.Info("Successfully connected to local Docker instance")
	}

	return client, nil
}

// validateRemoteHost validates the format of a remote Docker host
func validateRemoteHost(host string) error {
	if host == "" {
		return fmt.Errorf("host cannot be empty")
	}

	// Check for valid URL schemes
	validSchemes := []string{"tcp://", "http://", "https://"}
	hasValidScheme := false

	for _, scheme := range validSchemes {
		if len(host) > len(scheme) && host[:len(scheme)] == scheme {
			hasValidScheme = true
			break
		}
	}

	if !hasValidScheme {
		return fmt.Errorf("host must use a valid scheme (tcp://, http://, or https://)")
	}

	// Basic format validation for tcp://host:port
	if host[:4] == "tcp:" {
		// Remove scheme
		hostPart := host[6:] // "tcp://" is 6 characters

		// Check if it contains a colon (for port)
		if !strings.Contains(hostPart, ":") {
			return fmt.Errorf("tcp:// host must include port (e.g., tcp://192.168.1.100:2375)")
		}

		// Split host and port
		parts := strings.Split(hostPart, ":")
		if len(parts) != 2 {
			return fmt.Errorf("invalid host:port format")
		}

		// Validate port is numeric
		if parts[1] == "" {
			return fmt.Errorf("port cannot be empty")
		}
	}

	return nil
}

// GetInfo retrieves Docker system information
func (c *Client) GetInfo(ctx context.Context) (map[string]any, error) {
	info, err := c.cli.Info(ctx)
	if err != nil {
		return nil, fmt.Errorf("docker info failed: %w", err)
	}
	return marshalToMap(info)
}

// InspectContainer inspects a container
func (c *Client) InspectContainer(ctx context.Context, id string) (map[string]any, error) {
	container, err := c.cli.ContainerInspect(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("container inspect failed %s: %w", id, err)
	}

	return marshalToMap(container)
}

// GetContainerLogs retrieves container logs
func (c *Client) GetContainerLogs(ctx context.Context, id string) (string, error) {
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
	defer logs.Close()

	var logLines []string
	buffer := make([]byte, 1024)
	for {
		n, err := logs.Read(buffer)
		if n > 0 {
			line := string(buffer[:n])
			if len(line) >= 8 {
				logLines = append(logLines, line[8:])
			} else {
				logLines = append(logLines, line)
			}
		}
		if err != nil {
			break
		}
	}

	return strings.Join(logLines, ""), nil
}

// InspectImage inspects an image
func (c *Client) InspectImage(ctx context.Context, id string) (map[string]any, error) {
	imageInfo, err := c.cli.ImageInspect(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("image inspect failed %s: %w", id, err)
	}
	return marshalToMap(imageInfo)
}

// InspectVolume inspects a volume
func (c *Client) InspectVolume(ctx context.Context, name string) (map[string]any, error) {
	volumeInfo, err := c.cli.VolumeInspect(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("volume inspect failed %s: %w", name, err)
	}
	return marshalToMap(volumeInfo)
}

// RemoveVolume removes a volume
func (c *Client) RemoveVolume(ctx context.Context, name string, force bool) error {
	if name == "" {
		return fmt.Errorf("volume name cannot be empty")
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
	return marshalToMap(networkInfo)
}

// Close closes the Docker client connection
func (c *Client) Close() error {
	return c.cli.Close()
}

// ListContainers lists all containers
func (c *Client) ListContainers(ctx context.Context, all bool) ([]Container, error) {
	opts := container.ListOptions{
		All: all,
	}

	containers, err := c.cli.ContainerList(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to list containers: %w", err)
	}

	result := make([]Container, 0, len(containers))
	for i := range containers {
		cont := &containers[i]
		// Format ports
		ports := ""
		for _, p := range cont.Ports {
			if p.PublicPort > 0 {
				if ports != "" {
					ports += ", "
				}
				ports += fmt.Sprintf("%d:%d/%s", p.PublicPort, p.PrivatePort, p.Type)
			} else if p.PrivatePort > 0 {
				if ports != "" {
					ports += ", "
				}
				ports += fmt.Sprintf("%d/%s", p.PrivatePort, p.Type)
			}
		}

		result = append(result, Container{
			ID:      cont.ID[:12],
			Name:    cont.Names[0][1:], // Remove leading slash
			Image:   cont.Image,
			Status:  cont.Status,
			State:   cont.State,
			Created: time.Unix(cont.Created, 0),
			Ports:   ports,
			Size:    "", // Size is not available in ContainerList
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Created.After(result[j].Created)
	})

	return result, nil
}

// GetContainerStats gets container stats
func (c *Client) GetContainerStats(ctx context.Context, id string) (map[string]any, error) {
	stats, err := c.cli.ContainerStats(ctx, id, false)
	if err != nil {
		return nil, fmt.Errorf("failed to get container stats: %w", err)
	}
	defer stats.Body.Close()

	var statsJSON map[string]any
	if err := json.NewDecoder(stats.Body).Decode(&statsJSON); err != nil {
		return nil, fmt.Errorf("failed to decode container stats: %w", err)
	}

	return statsJSON, nil
}

// ListImages lists all images
func (c *Client) ListImages(ctx context.Context) ([]Image, error) {
	opts := image.ListOptions{}

	images, err := c.cli.ImageList(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to list images: %w", err)
	}

	result := make([]Image, 0, len(images))
	for i := range images {
		img := &images[i]
		// Format repository and tag
		repo := "<none>"
		tag := "<none>"
		if len(img.RepoTags) > 0 && img.RepoTags[0] != "<none>:<none>" {
			parts := strings.Split(img.RepoTags[0], ":")
			if len(parts) >= 2 {
				repo = parts[0]
				tag = parts[1]
			}
		}

		size := formatSize(img.Size)

		result = append(result, Image{
			ID:         img.ID[7:19], // Remove "sha256:" prefix and truncate
			Repository: repo,
			Tag:        tag,
			Size:       size,
			Created:    time.Unix(img.Created, 0),
			Containers: int(img.Containers),
		})
	}

	// Sort by creation time (newest first)
	sort.Slice(result, func(i, j int) bool {
		return result[i].Created.After(result[j].Created)
	})

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
		created := time.Time{}
		if vol.CreatedAt != "" {
			created, _ = time.Parse(time.RFC3339, vol.CreatedAt)
		}

		result = append(result, Volume{
			Name:       vol.Name,
			Driver:     vol.Driver,
			Mountpoint: vol.Mountpoint,
			Created:    created,
			Size:       "", // Size is not available in VolumeList
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})

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
	opts := container.StopOptions{}
	if timeout != nil {
		opts.Signal = "SIGTERM"
		seconds := int(timeout.Seconds())
		opts.Timeout = &seconds
	}
	return c.cli.ContainerStop(ctx, id, opts)
}

// RestartContainer restarts a container
func (c *Client) RestartContainer(ctx context.Context, id string, timeout *time.Duration) error {
	opts := container.StopOptions{}
	if timeout != nil {
		opts.Signal = "SIGTERM"
		seconds := int(timeout.Seconds())
		opts.Timeout = &seconds
	}
	return c.cli.ContainerRestart(ctx, id, opts)
}

// RemoveContainer removes a container
func (c *Client) RemoveContainer(ctx context.Context, id string, force bool) error {
	opts := container.RemoveOptions{
		Force: force,
	}

	return c.cli.ContainerRemove(ctx, id, opts)
}

// ExecContainer executes a command in a running container and returns the output
func (c *Client) ExecContainer(ctx context.Context, id string, command []string, _ bool) (string, error) {
	if c.cli == nil {
		return "", fmt.Errorf("docker client is not initialized")
	}

	execConfig := container.ExecOptions{
		Cmd:          command,
		Tty:          false, // Set to false to capture output
		AttachStdin:  false, // We don't need stdin for command execution
		AttachStdout: true,
		AttachStderr: true,
		Detach:       false,
	}

	execResp, err := c.cli.ContainerExecCreate(ctx, id, execConfig)
	if err != nil {
		return "", fmt.Errorf("failed to create exec instance: %w", err)
	}

	attachConfig := container.ExecStartOptions{
		Tty: false,
	}

	hijackedResp, err := c.cli.ContainerExecAttach(ctx, execResp.ID, attachConfig)
	if err != nil {
		return "", fmt.Errorf("failed to attach to exec instance: %w", err)
	}
	defer hijackedResp.Close()

	err = c.cli.ContainerExecStart(ctx, execResp.ID, attachConfig)
	if err != nil {
		return "", fmt.Errorf("failed to start exec instance: %w", err)
	}

	var output strings.Builder
	buffer := make([]byte, 1024)

	for {
		n, err := hijackedResp.Reader.Read(buffer)
		if n > 0 {
			output.Write(buffer[:n])
		}
		if err != nil {
			break
		}
	}

	return output.String(), nil
}

// AttachContainer attaches to a running container
func (c *Client) AttachContainer(ctx context.Context, id string) (types.HijackedResponse, error) {
	if c.cli == nil {
		return types.HijackedResponse{}, fmt.Errorf("docker client is not initialized")
	}

	attachConfig := container.AttachOptions{
		Stream: true,
		Stdin:  true,
		Stdout: true,
		Stderr: true,
		Logs:   false,
	}

	response, err := c.cli.ContainerAttach(ctx, id, attachConfig)
	if err != nil {
		return types.HijackedResponse{}, fmt.Errorf("failed to attach to container: %w", err)
	}

	return response, nil
}

// RemoveImage removes an image
func (c *Client) RemoveImage(ctx context.Context, id string, force bool) error {
	if id == "" {
		return fmt.Errorf("image ID cannot be empty")
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
	if id == "" {
		return fmt.Errorf("network ID cannot be empty")
	}

	if err := c.cli.NetworkRemove(ctx, id); err != nil {
		return fmt.Errorf("failed to remove network %s: %w", id, err)
	}

	return nil
}
