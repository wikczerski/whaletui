package swarm

import (
	"testing"

	"github.com/stretchr/testify/assert"
	sharedTypes "github.com/wikczerski/whaletui/internal/shared"
)

// TestNodesView_FormatNodeRow tests the formatNodeRow method without requiring full view construction
func TestNodesView_FormatNodeRow(t *testing.T) {
	// Create a test view instance
	view := &NodesView{}

	// Create test data
	node := sharedTypes.SwarmNode{
		ID:            "node1234567890abcdef",
		Hostname:      "test-node-1",
		Role:          "manager",
		Availability:  "active",
		Status:        "ready",
		ManagerStatus: "leader",
		EngineVersion: "20.10.0",
		Address:       "192.168.1.100",
	}

	// Test formatNodeRow method
	result := view.formatNodeRow(&node)

	assert.Len(t, result, 8)                    // Should have 8 columns
	assert.Equal(t, "node12345678", result[0])  // Truncated ID (12 chars)
	assert.Equal(t, "test-node-1", result[1])   // Hostname
	assert.Equal(t, "manager", result[2])       // Role
	assert.Equal(t, "active", result[3])        // Availability
	assert.Equal(t, "ready", result[4])         // Status
	assert.Equal(t, "leader", result[5])        // ManagerStatus
	assert.Equal(t, "20.10.0", result[6])       // EngineVersion
	assert.Equal(t, "192.168.1.100", result[7]) // Address
}

// TestNodesView_GetNodeID tests the getNodeID method
func TestNodesView_GetNodeID(t *testing.T) {
	// Create a test view instance
	view := &NodesView{}

	// Create test data
	node := sharedTypes.SwarmNode{
		ID:       "node1",
		Hostname: "test-node-1",
	}

	// Test getNodeID method
	result := view.getNodeID(&node)

	assert.Equal(t, "node1", result)
}

// TestNodesView_GetNodeName tests the getNodeName method
func TestNodesView_GetNodeName(t *testing.T) {
	// Create a test view instance
	view := &NodesView{}

	// Create test data
	node := sharedTypes.SwarmNode{
		ID:       "node1",
		Hostname: "test-node-1",
	}

	// Test getNodeName method
	result := view.getNodeName(&node)

	assert.Equal(t, "test-node-1", result)
}

// TestNodesView_GetActions tests the getActions method
func TestNodesView_GetActions(t *testing.T) {
	// Create a test view instance
	view := &NodesView{}

	// Test getActions method
	actions := view.getActions()

	assert.NotNil(t, actions)
	assert.Contains(t, actions, 'i') // Inspect
	assert.Contains(t, actions, 'a') // Update Availability
	assert.Contains(t, actions, 'r') // Remove
	assert.Contains(t, actions, 'f') // Refresh
}

// TestNodesView_Constructor tests that we can create a view (skipped due to complex dependencies)
func TestNodesView_Constructor(t *testing.T) {
	// Skip this test as it requires real UI, services, and managers
	// The actual constructor is tested indirectly through integration tests
	t.Skip("Constructor requires complex dependencies - tested through integration")
}
