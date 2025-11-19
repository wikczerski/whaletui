package e2e

import (
	"testing"
	"time"

	"github.com/docker/docker/api/types/swarm"
	"github.com/stretchr/testify/assert"
	"github.com/wikczerski/whaletui/e2e/framework"
)

// TestSwarmNodeList tests listing swarm nodes.
func TestSwarmNodeList(t *testing.T) {
	fw := framework.NewTestFramework(t)
	dh := fw.GetDockerHelper()

	// Get nodes
	nodes := dh.GetSwarmNodes()

	// Verify we have at least one node (the manager)
	assert.NotEmpty(t, nodes, "Should have at least one swarm node")
}

// TestSwarmNodeInspect tests inspecting a swarm node.
func TestSwarmNodeInspect(t *testing.T) {
	fw := framework.NewTestFramework(t)
	dh := fw.GetDockerHelper()

	// Get nodes to find an ID
	nodes := dh.GetSwarmNodes()
	assert.NotEmpty(t, nodes, "Should have at least one swarm node")
	nodeID := nodes[0].ID

	// Inspect node
	node := dh.GetSwarmNode(nodeID)

	// Verify inspect data
	assert.Equal(t, nodeID, node.ID, "Node ID should match")
	assert.NotEmpty(t, node.Description.Hostname, "Node should have a hostname")
}

// TestSwarmNodeUpdateAvailability tests updating a swarm node's availability.
func TestSwarmNodeUpdateAvailability(t *testing.T) {
	fw := framework.NewTestFramework(t)
	dh := fw.GetDockerHelper()

	// Get nodes to find an ID
	nodes := dh.GetSwarmNodes()
	assert.NotEmpty(t, nodes, "Should have at least one swarm node")
	nodeID := nodes[0].ID

	// Ensure initial state is active
	dh.UpdateNodeAvailability(nodeID, swarm.NodeAvailabilityActive)
	time.Sleep(1 * time.Second)

	// Update to Drain
	dh.UpdateNodeAvailability(nodeID, swarm.NodeAvailabilityDrain)
	time.Sleep(1 * time.Second)

	// Verify Drain state
	node := dh.GetSwarmNode(nodeID)
	assert.Equal(
		t,
		swarm.NodeAvailabilityDrain,
		node.Spec.Availability,
		"Node should be in Drain state",
	)

	// Update back to Active
	dh.UpdateNodeAvailability(nodeID, swarm.NodeAvailabilityActive)
	time.Sleep(1 * time.Second)

	// Verify Active state
	node = dh.GetSwarmNode(nodeID)
	assert.Equal(
		t,
		swarm.NodeAvailabilityActive,
		node.Spec.Availability,
		"Node should be in Active state",
	)
}
