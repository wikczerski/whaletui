// Package docker provides Docker client functionality for WhaleTUI.
package docker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"sort"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"github.com/wikczerski/whaletui/internal/config"
	"github.com/wikczerski/whaletui/internal/docker/dockerssh"
	"github.com/wikczerski/whaletui/internal/logger"
)

// Client represents a Docker client wrapper
type Client struct {
	cli     *client.Client
	cfg     *config.Config
	log     *slog.Logger
	sshConn *dockerssh.SSHConnection
	sshCtx  *dockerssh.SSHContext
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
	return marshalToMap(info)
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

	return marshalToMap(container)
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
	if err := validateID(name, "volume name"); err != nil {
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
	return marshalToMap(networkInfo)
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
	sortContainersByCreationTime(result)
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
	sortImagesByCreationTime(result)
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

	sortVolumesByName(result)
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
	opts := buildStopOptions(timeout)
	return c.cli.ContainerStop(ctx, id, opts)
}

// RestartContainer restarts a container
func (c *Client) RestartContainer(ctx context.Context, id string, timeout *time.Duration) error {
	opts := buildStopOptions(timeout)
	return c.cli.ContainerRestart(ctx, id, opts)
}

// RemoveContainer removes a container
func (c *Client) RemoveContainer(ctx context.Context, id string, force bool) error {
	opts := container.RemoveOptions{
		Force: force,
	}

	if err := validateID(id, "container ID"); err != nil {
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

	if err := validateID(id, "container ID"); err != nil {
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
	if err := validateID(id, "image ID"); err != nil {
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
	if err := validateID(id, "network ID"); err != nil {
		return err
	}

	if err := c.cli.NetworkRemove(ctx, id); err != nil {
		return fmt.Errorf("failed to remove network %s: %w", id, err)
	}
	return nil
}

// ListSwarmServices lists all swarm services
func (c *Client) ListSwarmServices(ctx context.Context) ([]swarm.Service, error) {
	if err := c.validateClient(); err != nil {
		return nil, err
	}

	services, err := c.cli.ServiceList(ctx, swarm.ServiceListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list swarm services: %w", err)
	}

	return services, nil
}

// InspectSwarmService inspects a swarm service
func (c *Client) InspectSwarmService(ctx context.Context, id string) (swarm.Service, error) {
	if err := c.validateClient(); err != nil {
		return swarm.Service{}, err
	}

	if err := validateID(id, "service ID"); err != nil {
		return swarm.Service{}, err
	}

	service, _, err := c.cli.ServiceInspectWithRaw(ctx, id, swarm.ServiceInspectOptions{})
	if err != nil {
		return swarm.Service{}, fmt.Errorf("failed to inspect swarm service %s: %w", id, err)
	}

	return service, nil
}

// UpdateSwarmService updates a swarm service
// nolint:gocritic // Docker API requires value parameter for compatibility
func (c *Client) UpdateSwarmService(
	ctx context.Context,
	id string,
	version swarm.Version,
	spec swarm.ServiceSpec,
) error {
	if err := c.validateClient(); err != nil {
		return err
	}

	if err := validateID(id, "service ID"); err != nil {
		return err
	}

	response, err := c.cli.ServiceUpdate(ctx, id, version, spec, swarm.ServiceUpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to update swarm service %s: %w", id, err)
	}

	if len(response.Warnings) > 0 {
		c.log.Warn("Service update warnings", "service_id", id, "warnings", response.Warnings)
	}

	return nil
}

// RemoveSwarmService removes a swarm service
func (c *Client) RemoveSwarmService(ctx context.Context, id string) error {
	if err := c.validateClient(); err != nil {
		return err
	}

	if err := validateID(id, "service ID"); err != nil {
		return err
	}

	if err := c.cli.ServiceRemove(ctx, id); err != nil {
		return fmt.Errorf("failed to remove swarm service %s: %w", id, err)
	}

	return nil
}

// GetSwarmServiceLogs gets logs for a swarm service
func (c *Client) GetSwarmServiceLogs(ctx context.Context, id string) (string, error) {
	if err := c.validateServiceLogsRequest(id); err != nil {
		return "", err
	}

	response, err := c.getServiceLogsResponse(ctx, id)
	if err != nil {
		return "", err
	}
	defer c.closeServiceLogsResponse(response)

	return c.readServiceLogs(response), nil
}

// ListSwarmNodes lists all swarm nodes
func (c *Client) ListSwarmNodes(ctx context.Context) ([]swarm.Node, error) {
	if err := c.validateClient(); err != nil {
		return nil, err
	}

	nodes, err := c.cli.NodeList(ctx, swarm.NodeListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list swarm nodes: %w", err)
	}

	return nodes, nil
}

// InspectSwarmNode inspects a swarm node
func (c *Client) InspectSwarmNode(_ context.Context, id string) (swarm.Node, error) {
	if err := c.validateClient(); err != nil {
		return swarm.Node{}, err
	}

	if err := validateID(id, "node ID"); err != nil {
		return swarm.Node{}, err
	}

	return swarm.Node{}, errors.New(
		"NodeInspect method not available in this Docker client version",
	)
}

// UpdateSwarmNode updates a swarm node
func (c *Client) UpdateSwarmNode(
	ctx context.Context,
	id string,
	version swarm.Version,
	spec swarm.NodeSpec,
) error {
	if err := c.validateClient(); err != nil {
		return err
	}

	if err := validateID(id, "node ID"); err != nil {
		return err
	}

	if err := c.cli.NodeUpdate(ctx, id, version, spec); err != nil {
		return fmt.Errorf("failed to update swarm node %s: %w", id, err)
	}

	return nil
}

// RemoveSwarmNode removes a swarm node
func (c *Client) RemoveSwarmNode(ctx context.Context, id string, force bool) error {
	if err := c.validateClient(); err != nil {
		return err
	}

	if err := validateID(id, "node ID"); err != nil {
		return err
	}

	options := swarm.NodeRemoveOptions{
		Force: force,
	}

	if err := c.cli.NodeRemove(ctx, id, options); err != nil {
		return fmt.Errorf("failed to remove swarm node %s: %w", id, err)
	}

	return nil
}

// validateClient validates that the Docker client is initialized
func (c *Client) validateClient() error {
	if c.cli == nil {
		return errors.New("docker client is not initialized")
	}
	return nil
}

// readAndFormatLogs reads logs from the response and formats them
func (c *Client) readAndFormatLogs(logs io.ReadCloser) string {
	var logLines []string
	buffer := make([]byte, 1024)

	c.readLogsIntoBuffer(logs, buffer, &logLines)

	return strings.Join(logLines, "")
}

// readLogsIntoBuffer reads logs into the buffer and formats them
func (c *Client) readLogsIntoBuffer(logs io.ReadCloser, buffer []byte, logLines *[]string) {
	for {
		n, err := logs.Read(buffer)
		if n > 0 {
			line := string(buffer[:n])
			formattedLine := c.formatLogLine(line)
			*logLines = append(*logLines, formattedLine)
		}
		if err != nil {
			break
		}
	}
}

// formatLogLine formats a single log line by removing timestamp prefix if present
func (c *Client) formatLogLine(line string) string {
	if len(line) >= 8 {
		return line[8:] // Remove timestamp prefix
	}
	return line
}

// decodeStatsResponse decodes the stats response body into a map
func (c *Client) decodeStatsResponse(body io.ReadCloser) (map[string]any, error) {
	var statsJSON map[string]any
	if err := json.NewDecoder(body).Decode(&statsJSON); err != nil {
		return nil, fmt.Errorf("failed to decode container stats: %w", err)
	}
	return statsJSON, nil
}

// createVolumeFromAPI creates a Volume from the API response
func (c *Client) createVolumeFromAPI(vol *volume.Volume) Volume {
	created := time.Time{}
	if vol.CreatedAt != "" {
		created, _ = time.Parse(time.RFC3339, vol.CreatedAt)
	}

	return Volume{
		Name:       vol.Name,
		Driver:     vol.Driver,
		Mountpoint: vol.Mountpoint,
		Created:    created,
		Size:       "", // Size is not available in VolumeList
	}
}

// createExecInstance creates an exec instance in the container
func (c *Client) createExecInstance(
	ctx context.Context,
	id string,
	command []string,
) (*container.ExecCreateResponse, error) {
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
		return nil, fmt.Errorf("failed to create exec instance: %w", err)
	}

	return &execResp, nil
}

// executeAndCollectOutput executes the exec instance and collects the output
func (c *Client) executeAndCollectOutput(ctx context.Context, execID string) (string, error) {
	attachConfig := container.ExecStartOptions{
		Tty: false,
	}

	hijackedResp, err := c.cli.ContainerExecAttach(ctx, execID, attachConfig)
	if err != nil {
		return "", fmt.Errorf("failed to attach to exec instance: %w", err)
	}
	defer hijackedResp.Close()

	if err := c.cli.ContainerExecStart(ctx, execID, attachConfig); err != nil {
		return "", fmt.Errorf("failed to start exec instance: %w", err)
	}

	return c.readExecOutput(hijackedResp), nil
}

// readExecOutput reads the output from the hijacked response
func (c *Client) readExecOutput(hijackedResp types.HijackedResponse) string {
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

	return output.String()
}

// getImageList retrieves the raw image list from Docker
func (c *Client) getImageList(ctx context.Context) ([]image.Summary, error) {
	opts := image.ListOptions{}
	return c.cli.ImageList(ctx, opts)
}

// convertToImages converts Docker API images to our Image type
func (c *Client) convertToImages(images []image.Summary) []Image {
	result := make([]Image, 0, len(images))
	for i := range images {
		img := &images[i]
		result = append(result, c.convertImage(img))
	}
	return result
}

// convertImage converts a single Docker API image to our Image type
func (c *Client) convertImage(img *image.Summary) Image {
	repo, tag := parseImageRepository(img.RepoTags)
	size := formatSize(img.Size)
	return Image{
		ID:         img.ID[7:19], // Remove "sha256:" prefix and truncate
		Repository: repo,
		Tag:        tag,
		Size:       size,
		Created:    time.Unix(img.Created, 0),
		Containers: int(img.Containers),
	}
}

// validateServiceLogsRequest validates the service logs request
func (c *Client) validateServiceLogsRequest(id string) error {
	if err := c.validateClient(); err != nil {
		return err
	}
	return validateID(id, "service ID")
}

// getServiceLogsResponse gets the service logs response from Docker
func (c *Client) getServiceLogsResponse(ctx context.Context, id string) (io.ReadCloser, error) {
	options := container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Timestamps: true,
		Follow:     false,
	}

	response, err := c.cli.ServiceLogs(ctx, id, options)
	if err != nil {
		return nil, fmt.Errorf("failed to get logs for swarm service %s: %w", id, err)
	}
	return response, nil
}

// closeServiceLogsResponse safely closes the service logs response
func (c *Client) closeServiceLogsResponse(response io.ReadCloser) {
	if err := response.Close(); err != nil {
		c.log.Warn("Failed to close response", "error", err)
	}
}

// readServiceLogs reads and formats the service logs
func (c *Client) readServiceLogs(response io.ReadCloser) string {
	output := &strings.Builder{}
	buffer := make([]byte, 1024)

	for {
		n, err := response.Read(buffer)
		if n > 0 {
			output.Write(buffer[:n])
		}
		if err != nil {
			break
		}
	}

	return output.String()
}

// getContainerList retrieves the raw container list from Docker
func (c *Client) getContainerList(ctx context.Context, all bool) ([]container.Summary, error) {
	opts := container.ListOptions{All: all}
	return c.cli.ContainerList(ctx, opts)
}

// convertToContainers converts Docker API containers to our Container type
func (c *Client) convertToContainers(containers []container.Summary) []Container {
	result := make([]Container, 0, len(containers))
	for i := range containers {
		cont := &containers[i]
		result = append(result, c.convertContainer(cont))
	}
	return result
}

// convertContainer converts a single Docker API container to our Container type
func (c *Client) convertContainer(cont *container.Summary) Container {
	ports := formatContainerPorts(cont.Ports)
	return Container{
		ID:      cont.ID[:12],
		Name:    cont.Names[0][1:], // Remove leading slash
		Image:   cont.Image,
		Status:  cont.Status,
		State:   cont.State,
		Created: time.Unix(cont.Created, 0),
		Ports:   ports,
		Size:    "", // Size is not available in ContainerList
	}
}

// closeDockerClient closes the Docker client and logs the result
func (c *Client) closeDockerClient(errors *[]string) {
	if c.cli != nil {
		if err := c.cli.Close(); err != nil {
			c.log.Error("Failed to close Docker client", "error", err)
			*errors = append(*errors, fmt.Sprintf("Docker client close: %v", err))
		} else {
			c.log.Info("Docker client closed successfully")
		}
	}
}

// closeSSHConnection closes the SSH connection and logs the result
func (c *Client) closeSSHConnection(errors *[]string) {
	if c.sshConn != nil {
		c.log.Info("Closing SSH connection with socat")
		if err := c.sshConn.Close(); err != nil {
			c.log.Error("Failed to close SSH connection", "error", err)
			*errors = append(*errors, fmt.Sprintf("SSH connection close: %v", err))
		} else {
			c.log.Info("SSH connection closed successfully")
		}
	}
}

// closeSSHContext closes the SSH context and logs the result
func (c *Client) closeSSHContext(errors *[]string) {
	if c.sshCtx != nil {
		c.log.Info("Closing SSH context")
		if err := c.sshCtx.Close(); err != nil {
			c.log.Error("Failed to close SSH context", "error", err)
			*errors = append(*errors, fmt.Sprintf("SSH context close: %v", err))
		} else {
			c.log.Info("SSH context closed successfully")
		}
	}
}

// buildCloseError builds the final error message if any errors occurred
func (c *Client) buildCloseError(errors []string) error {
	if len(errors) > 0 {
		return fmt.Errorf("failed to close client: %s", strings.Join(errors, "; "))
	}
	return nil
}
