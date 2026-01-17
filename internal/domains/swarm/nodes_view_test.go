package swarm

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wikczerski/whaletui/internal/shared"
	uishared "github.com/wikczerski/whaletui/internal/ui/shared"
)

// TestNodesView_FormatNodeRow tests the formatNodeRow method without requiring full view construction
func TestNodesView_FormatNodeRow(t *testing.T) {
	// Create a test view instance
	view := &NodesView{}

	// Create test data
	node := shared.SwarmNode{
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
	node := shared.SwarmNode{
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
	node := shared.SwarmNode{
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

// TestNodesView_SearchState tests search state management
func TestNodesView_SearchState(t *testing.T) {
	// Create a test view instance with proper initialization
	view := &NodesView{}
	view.BaseView = uishared.NewBaseView[shared.SwarmNode](
		nil, "Test Nodes", []string{"ID", "Hostname"},
	)

	// Test initial state
	assert.False(t, view.IsSearchActive())
	assert.Empty(t, view.GetSearchTerm())

	// Test setting search term
	view.Search("test")
	assert.True(t, view.IsSearchActive())
	assert.Equal(t, "test", view.GetSearchTerm())

	// Test clearing search
	view.ClearSearch()
	assert.False(t, view.IsSearchActive())
	assert.Empty(t, view.GetSearchTerm())
}

// TestNodesView_SearchWithMockData tests search functionality with mock data
func TestNodesView_SearchWithMockData(t *testing.T) {
	// Create a test view instance with proper initialization
	view := &NodesView{}
	view.BaseView = uishared.NewBaseView[shared.SwarmNode](
		nil, "Test Nodes", []string{"ID", "Hostname", "Role", "Status"},
	)

	// Set up the callbacks that are needed for search to work
	view.FormatRow = func(n shared.SwarmNode) []string { return view.formatNodeRow(&n) }
	view.GetItemID = func(n shared.SwarmNode) string { return view.getNodeID(&n) }
	view.GetItemName = func(n shared.SwarmNode) string { return view.getNodeName(&n) }

	// Set up the view with mock data
	testNodes := []shared.SwarmNode{
		{
			ID:            "node1",
			Hostname:      "manager-node",
			Role:          "manager",
			Availability:  "active",
			Status:        "ready",
			ManagerStatus: "leader",
			EngineVersion: "20.10.0",
			Address:       "192.168.1.100",
		},
		{
			ID:            "node2",
			Hostname:      "worker-node-1",
			Role:          "worker",
			Availability:  "active",
			Status:        "ready",
			ManagerStatus: "",
			EngineVersion: "20.10.0",
			Address:       "192.168.1.101",
		},
		{
			ID:            "node3",
			Hostname:      "worker-node-2",
			Role:          "worker",
			Availability:  "drain",
			Status:        "ready",
			ManagerStatus: "",
			EngineVersion: "20.10.1",
			Address:       "192.168.1.102",
		},
	}

	// Set up the view with mock data by directly setting the items
	view.SetItems(testNodes)

	// Test search by hostname
	view.Search("manager")
	assert.True(t, view.IsSearchActive())
	assert.Equal(t, "manager", view.GetSearchTerm())

	// Get the filtered items to verify search worked
	filteredItems := view.GetItems()
	assert.Len(t, filteredItems, 1)
	assert.Equal(t, "manager-node", view.getNodeName(&filteredItems[0]))

	// Test search by role
	view.Search("worker")
	assert.True(t, view.IsSearchActive())
	assert.Equal(t, "worker", view.GetSearchTerm())

	filteredItems = view.GetItems()
	assert.Len(t, filteredItems, 2)
	assert.Equal(t, "worker-node-1", view.getNodeName(&filteredItems[0]))
	assert.Equal(t, "worker-node-2", view.getNodeName(&filteredItems[1]))

	// Test search by availability
	view.Search("drain")
	assert.True(t, view.IsSearchActive())
	assert.Equal(t, "drain", view.GetSearchTerm())

	filteredItems = view.GetItems()
	assert.Len(t, filteredItems, 1)
	assert.Equal(t, "worker-node-2", view.getNodeName(&filteredItems[0]))

	// Test case insensitive search
	view.Search("MANAGER")
	assert.True(t, view.IsSearchActive())
	assert.Equal(t, "MANAGER", view.GetSearchTerm())

	filteredItems = view.GetItems()
	assert.Len(t, filteredItems, 1)
	assert.Equal(t, "manager-node", view.getNodeName(&filteredItems[0]))

	// Test partial match
	view.Search("worker")
	assert.True(t, view.IsSearchActive())
	assert.Equal(t, "worker", view.GetSearchTerm())

	filteredItems = view.GetItems()
	assert.Len(t, filteredItems, 2)

	// Test no match
	view.Search("nonexistent")
	assert.True(t, view.IsSearchActive())
	assert.Equal(t, "nonexistent", view.GetSearchTerm())

	filteredItems = view.GetItems()
	assert.Len(t, filteredItems, 0)

	// Test clearing search
	view.ClearSearch()
	assert.False(t, view.IsSearchActive())
	assert.Empty(t, view.GetSearchTerm())

	filteredItems = view.GetItems()
	assert.Len(t, filteredItems, 3) // Should show all items again
}
