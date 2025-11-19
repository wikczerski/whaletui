package framework

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/api/types/volume"
	"github.com/stretchr/testify/require"
)

// DockerHelper provides Docker operation utilities for testing.
type DockerHelper struct {
	fw *TestFramework
}

// NewDockerHelper creates a new Docker helper.
func NewDockerHelper(fw *TestFramework) *DockerHelper {
	return &DockerHelper{fw: fw}
}

// CreateTestContainer creates a container for testing.
func (dh *DockerHelper) CreateTestContainer(
	name, imageName string,
	config *container.Config,
	hostConfig *container.HostConfig,
) string {
	dh.fw.t.Helper()

	ctx := dh.fw.GetContext()
	client := dh.fw.GetDockerClient()

	// Pull image if not exists
	dh.EnsureImage(imageName)

	// Create container config if not provided
	if config == nil {
		config = &container.Config{
			Image: imageName,
		}
	} else if config.Image == "" {
		config.Image = imageName
	}

	// Create container
	resp, err := client.ContainerCreate(ctx, config, hostConfig, nil, nil, name)
	require.NoError(dh.fw.t, err, "Failed to create container")

	// Register for cleanup
	dh.fw.RegisterTestContainer(resp.ID)

	return resp.ID
}

// StartContainer starts a container.
func (dh *DockerHelper) StartContainer(containerID string) {
	dh.fw.t.Helper()

	ctx := dh.fw.GetContext()
	client := dh.fw.GetDockerClient()

	err := client.ContainerStart(ctx, containerID, container.StartOptions{})
	require.NoError(dh.fw.t, err, "Failed to start container")
}

// StopContainer stops a container.
func (dh *DockerHelper) StopContainer(containerID string) {
	dh.fw.t.Helper()

	ctx := dh.fw.GetContext()
	client := dh.fw.GetDockerClient()

	timeout := 10
	err := client.ContainerStop(ctx, containerID, container.StopOptions{Timeout: &timeout})
	require.NoError(dh.fw.t, err, "Failed to stop container")
}

// RemoveContainer removes a container.
func (dh *DockerHelper) RemoveContainer(containerID string, force bool) {
	dh.fw.t.Helper()

	ctx := dh.fw.GetContext()
	client := dh.fw.GetDockerClient()

	err := client.ContainerRemove(ctx, containerID, container.RemoveOptions{Force: force})
	require.NoError(dh.fw.t, err, "Failed to remove container")
}

// GetContainerState returns the current state of a container.
func (dh *DockerHelper) GetContainerState(containerID string) string {
	dh.fw.t.Helper()

	ctx := dh.fw.GetContext()
	client := dh.fw.GetDockerClient()

	inspect, err := client.ContainerInspect(ctx, containerID)
	require.NoError(dh.fw.t, err, "Failed to inspect container")

	return inspect.State.Status
}

// WaitForContainerState waits for a container to reach a specific state.
func (dh *DockerHelper) WaitForContainerState(
	containerID, expectedState string,
	timeout time.Duration,
) {
	dh.fw.t.Helper()

	dh.fw.WaitForCondition(func() bool {
		state := dh.GetContainerState(containerID)
		return state == expectedState
	}, timeout, fmt.Sprintf("container %s to reach state %s", containerID, expectedState))
}

// EnsureImage ensures an image exists locally, pulling if necessary.
func (dh *DockerHelper) EnsureImage(imageName string) {
	dh.fw.t.Helper()

	ctx := dh.fw.GetContext()
	client := dh.fw.GetDockerClient()

	// Check if image exists
	_, err := client.ImageInspect(ctx, imageName)
	if err == nil {
		return // Image already exists
	}

	// Pull image
	reader, err := client.ImagePull(ctx, imageName, image.PullOptions{})
	require.NoError(dh.fw.t, err, "Failed to pull image")
	defer func() { _ = reader.Close() }()

	// Wait for pull to complete
	_, err = io.Copy(io.Discard, reader)
	require.NoError(dh.fw.t, err, "Failed to read image pull response")
}

// CreateTestVolume creates a volume for testing.
func (dh *DockerHelper) CreateTestVolume(name string) string {
	dh.fw.t.Helper()

	ctx := dh.fw.GetContext()
	client := dh.fw.GetDockerClient()

	vol, err := client.VolumeCreate(ctx, volume.CreateOptions{
		Name: name,
	})
	require.NoError(dh.fw.t, err, "Failed to create volume")

	// Register for cleanup
	dh.fw.RegisterTestVolume(vol.Name)

	return vol.Name
}

// RemoveVolume removes a volume.
func (dh *DockerHelper) RemoveVolume(name string, force bool) {
	dh.fw.t.Helper()

	ctx := dh.fw.GetContext()
	client := dh.fw.GetDockerClient()

	err := client.VolumeRemove(ctx, name, force)
	require.NoError(dh.fw.t, err, "Failed to remove volume")
}

// CreateTestNetwork creates a network for testing.
func (dh *DockerHelper) CreateTestNetwork(name, driver string) string {
	dh.fw.t.Helper()

	ctx := dh.fw.GetContext()
	client := dh.fw.GetDockerClient()

	resp, err := client.NetworkCreate(ctx, name, network.CreateOptions{
		Driver: driver,
	})
	require.NoError(dh.fw.t, err, "Failed to create network")

	// Register for cleanup
	dh.fw.RegisterTestNetwork(resp.ID)

	return resp.ID
}

// RemoveNetwork removes a network.
func (dh *DockerHelper) RemoveNetwork(networkID string) {
	dh.fw.t.Helper()

	ctx := dh.fw.GetContext()
	client := dh.fw.GetDockerClient()

	err := client.NetworkRemove(ctx, networkID)
	require.NoError(dh.fw.t, err, "Failed to remove network")
}

// CreateTestService creates a swarm service for testing.
func (dh *DockerHelper) CreateTestService(name, imageName string, replicas uint64) string {
	dh.fw.t.Helper()

	ctx := dh.fw.GetContext()
	client := dh.fw.GetDockerClient()

	// Ensure image exists
	dh.EnsureImage(imageName)

	// Create service spec
	serviceSpec := swarm.ServiceSpec{
		Annotations: swarm.Annotations{
			Name: name,
		},
		TaskTemplate: swarm.TaskSpec{
			ContainerSpec: &swarm.ContainerSpec{
				Image: imageName,
			},
		},
		Mode: swarm.ServiceMode{
			Replicated: &swarm.ReplicatedService{
				Replicas: &replicas,
			},
		},
	}

	resp, err := client.ServiceCreate(ctx, serviceSpec, swarm.ServiceCreateOptions{})
	require.NoError(dh.fw.t, err, "Failed to create service")

	// Register for cleanup
	dh.fw.RegisterTestService(resp.ID)

	return resp.ID
}

// RemoveService removes a swarm service.
func (dh *DockerHelper) RemoveService(serviceID string) {
	dh.fw.t.Helper()

	ctx := dh.fw.GetContext()
	client := dh.fw.GetDockerClient()

	err := client.ServiceRemove(ctx, serviceID)
	require.NoError(dh.fw.t, err, "Failed to remove service")
}

// ScaleService scales a swarm service.
func (dh *DockerHelper) ScaleService(serviceID string, replicas uint64) {
	dh.fw.t.Helper()

	ctx := dh.fw.GetContext()
	client := dh.fw.GetDockerClient()

	// Get current service spec
	service, _, err := client.ServiceInspectWithRaw(ctx, serviceID, swarm.ServiceInspectOptions{})
	require.NoError(dh.fw.t, err, "Failed to inspect service")

	// Update replicas
	service.Spec.Mode.Replicated.Replicas = &replicas

	// Update service
	_, err = client.ServiceUpdate(
		ctx,
		serviceID,
		service.Version,
		service.Spec,
		swarm.ServiceUpdateOptions{},
	)
	require.NoError(dh.fw.t, err, "Failed to scale service")
}

// GetServiceReplicas returns the current replica count of a service.
func (dh *DockerHelper) GetServiceReplicas(serviceID string) uint64 {
	dh.fw.t.Helper()

	ctx := dh.fw.GetContext()
	client := dh.fw.GetDockerClient()

	service, _, err := client.ServiceInspectWithRaw(ctx, serviceID, swarm.ServiceInspectOptions{})
	require.NoError(dh.fw.t, err, "Failed to inspect service")

	if service.Spec.Mode.Replicated != nil && service.Spec.Mode.Replicated.Replicas != nil {
		return *service.Spec.Mode.Replicated.Replicas
	}

	return 0
}

// WaitForServiceReplicas waits for a service to reach a specific replica count.
func (dh *DockerHelper) WaitForServiceReplicas(
	serviceID string,
	expectedReplicas uint64,
	timeout time.Duration,
) {
	dh.fw.t.Helper()

	dh.fw.WaitForCondition(func() bool {
		replicas := dh.GetServiceReplicas(serviceID)
		return replicas == expectedReplicas
	}, timeout, fmt.Sprintf("service %s to reach %d replicas", serviceID, expectedReplicas))
}

// ListContainers lists all containers.
func (dh *DockerHelper) ListContainers(all bool) []container.Summary {
	dh.fw.t.Helper()

	ctx := dh.fw.GetContext()
	client := dh.fw.GetDockerClient()

	containers, err := client.ContainerList(ctx, container.ListOptions{All: all})
	require.NoError(dh.fw.t, err, "Failed to list containers")

	return containers
}

// CleanupAll cleans up all registered test resources.
func (dh *DockerHelper) CleanupAll() {
	ctx := context.Background()
	client := dh.fw.GetDockerClient()

	// Cleanup services
	for _, serviceID := range dh.fw.testServices {
		_ = client.ServiceRemove(ctx, serviceID)
	}

	// Cleanup containers
	for _, containerID := range dh.fw.testContainers {
		timeout := 5
		_ = client.ContainerStop(ctx, containerID, container.StopOptions{Timeout: &timeout})
		_ = client.ContainerRemove(ctx, containerID, container.RemoveOptions{Force: true})
	}

	// Cleanup volumes
	for _, volumeName := range dh.fw.testVolumes {
		_ = client.VolumeRemove(ctx, volumeName, true)
	}

	// Cleanup networks
	for _, networkID := range dh.fw.testNetworks {
		_ = client.NetworkRemove(ctx, networkID)
	}
}

// FindContainerByName finds a container by name.
func (dh *DockerHelper) FindContainerByName(name string) *container.Summary {
	dh.fw.t.Helper()

	ctx := dh.fw.GetContext()
	client := dh.fw.GetDockerClient()

	filterArgs := filters.NewArgs()
	filterArgs.Add("name", name)

	containers, err := client.ContainerList(ctx, container.ListOptions{
		All:     true,
		Filters: filterArgs,
	})
	require.NoError(dh.fw.t, err, "Failed to list containers")

	if len(containers) == 0 {
		return nil
	}

	return &containers[0]
}

// GetSwarmNodes returns a list of all swarm nodes.
func (dh *DockerHelper) GetSwarmNodes() []swarm.Node {
	dh.fw.t.Helper()

	ctx := dh.fw.GetContext()
	client := dh.fw.GetDockerClient()

	nodes, err := client.NodeList(ctx, swarm.NodeListOptions{})
	require.NoError(dh.fw.t, err, "Failed to list swarm nodes")

	return nodes
}

// GetSwarmNode returns a swarm node by ID.
func (dh *DockerHelper) GetSwarmNode(nodeID string) swarm.Node {
	dh.fw.t.Helper()

	ctx := dh.fw.GetContext()
	client := dh.fw.GetDockerClient()

	node, _, err := client.NodeInspectWithRaw(ctx, nodeID)
	require.NoError(dh.fw.t, err, "Failed to inspect swarm node")

	return node
}

// UpdateNodeAvailability updates the availability of a swarm node.
func (dh *DockerHelper) UpdateNodeAvailability(nodeID string, availability swarm.NodeAvailability) {
	dh.fw.t.Helper()

	ctx := dh.fw.GetContext()
	client := dh.fw.GetDockerClient()

	// Get current node spec
	node, _, err := client.NodeInspectWithRaw(ctx, nodeID)
	require.NoError(dh.fw.t, err, "Failed to inspect swarm node")

	// Update availability
	node.Spec.Availability = availability

	// Update node
	err = client.NodeUpdate(ctx, nodeID, node.Version, node.Spec)
	require.NoError(dh.fw.t, err, "Failed to update swarm node availability")
}
