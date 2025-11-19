package e2e

import (
	"testing"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wikczerski/whaletui/e2e/framework"
)

// TestNetworkList tests listing Docker networks.
func TestNetworkList(t *testing.T) {
	fw := framework.NewTestFramework(t)
	dh := fw.GetDockerHelper()
	ctx := fw.GetContext()
	client := fw.GetDockerClient()

	// Create test networks
	net1 := dh.CreateTestNetwork(framework.NetworkFixtures.Bridge1.Name, "bridge")
	net2 := dh.CreateTestNetwork(framework.NetworkFixtures.Bridge2.Name, "bridge")

	// List networks
	networks, err := client.NetworkList(ctx, network.ListOptions{})
	require.NoError(t, err, "Failed to list networks")

	// Verify networks exist
	networkIDs := make([]string, 0)
	for _, n := range networks {
		networkIDs = append(networkIDs, n.ID)
	}

	assert.Contains(t, networkIDs, net1, "Network 1 should be in list")
	assert.Contains(t, networkIDs, net2, "Network 2 should be in list")
}

// TestNetworkCreate tests creating a network.
func TestNetworkCreate(t *testing.T) {
	fw := framework.NewTestFramework(t)
	dh := fw.GetDockerHelper()

	// Create network
	networkID := dh.CreateTestNetwork("e2e-test-network-create", "bridge")

	// Verify network exists
	ctx := fw.GetContext()
	client := fw.GetDockerClient()
	inspect, err := client.NetworkInspect(ctx, networkID, network.InspectOptions{})
	require.NoError(t, err, "Failed to inspect network")
	assert.Equal(t, networkID, inspect.ID, "Network ID should match")
}

// TestNetworkDelete tests deleting a network.
func TestNetworkDelete(t *testing.T) {
	fw := framework.NewTestFramework(t)
	dh := fw.GetDockerHelper()

	// Create network
	networkID := dh.CreateTestNetwork("e2e-test-network-delete", "bridge")

	// Delete network
	dh.RemoveNetwork(networkID)

	// Verify network is gone
	ctx := fw.GetContext()
	client := fw.GetDockerClient()
	_, err := client.NetworkInspect(ctx, networkID, network.InspectOptions{})
	assert.Error(t, err, "Network should be deleted")
}

// TestNetworkDeleteWithContainers tests deleting a network with connected containers.
func TestNetworkDeleteWithContainers(t *testing.T) {
	fw := framework.NewTestFramework(t)
	dh := fw.GetDockerHelper()
	ctx := fw.GetContext()
	client := fw.GetDockerClient()

	// Create network
	networkID := dh.CreateTestNetwork("e2e-test-network-in-use", "bridge")

	// Create container connected to network
	containerID := dh.CreateTestContainer(
		"e2e-test-network-container",
		framework.ImageFixtures.Alpine,
		&container.Config{
			Image: framework.ImageFixtures.Alpine,
			Cmd:   []string{"sleep", "3600"},
		},
		&container.HostConfig{
			NetworkMode: container.NetworkMode(networkID),
		},
	)
	dh.StartContainer(containerID)

	// Try to delete network (should fail)
	err := client.NetworkRemove(ctx, networkID)
	assert.Error(t, err, "Should fail to delete network with connected containers")

	// Cleanup
	dh.StopContainer(containerID)
	dh.RemoveContainer(containerID, true)
	dh.RemoveNetwork(networkID)
}

// TestNetworkDeleteDefault tests attempting to delete default networks.
func TestNetworkDeleteDefault(t *testing.T) {
	fw := framework.NewTestFramework(t)
	ctx := fw.GetContext()
	client := fw.GetDockerClient()

	// Try to delete bridge network (should fail)
	err := client.NetworkRemove(ctx, "bridge")
	assert.Error(t, err, "Should fail to delete default bridge network")

	// Try to delete host network (should fail)
	err = client.NetworkRemove(ctx, "host")
	assert.Error(t, err, "Should fail to delete default host network")
}

// TestNetworkInspect tests inspecting a network.
func TestNetworkInspect(t *testing.T) {
	fw := framework.NewTestFramework(t)
	dh := fw.GetDockerHelper()
	ctx := fw.GetContext()
	client := fw.GetDockerClient()

	// Create network
	networkID := dh.CreateTestNetwork("e2e-test-network-inspect", "bridge")

	// Inspect network
	inspect, err := client.NetworkInspect(ctx, networkID, network.InspectOptions{})
	require.NoError(t, err, "Failed to inspect network")

	// Verify inspect data
	assert.Equal(t, networkID, inspect.ID, "Network ID should match")
	assert.Equal(t, "bridge", inspect.Driver, "Driver should be bridge")
	assert.NotNil(t, inspect.IPAM, "IPAM config should not be nil")
}

// TestNetworkEmptyList tests handling of empty network list scenario.
func TestNetworkEmptyList(t *testing.T) {
	fw := framework.NewTestFramework(t)
	ctx := fw.GetContext()
	client := fw.GetDockerClient()

	// List networks (will always have default networks)
	networks, err := client.NetworkList(ctx, network.ListOptions{})
	require.NoError(t, err, "Failed to list networks")

	// Verify operation works and default networks exist
	assert.NotNil(t, networks, "Network list should not be nil")
	assert.GreaterOrEqual(
		t,
		len(networks),
		3,
		"Should have at least default networks (bridge, host, none)",
	)
}
