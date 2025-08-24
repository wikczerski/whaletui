// Package swarm provides Swarm-related UI views for WhaleTUI.
package swarm

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/wikczerski/whaletui/internal/domains/swarm"
	"github.com/wikczerski/whaletui/internal/logger"
	"github.com/wikczerski/whaletui/internal/shared"
	uiinterfaces "github.com/wikczerski/whaletui/internal/ui/interfaces"
	"github.com/wikczerski/whaletui/internal/ui/managers"
)

// NodesView represents the swarm nodes view
type NodesView struct {
	*shared.BaseView[shared.SwarmNode]
	nodeService   *swarm.NodeService
	modalManager  *managers.ModalManager
	headerManager *managers.HeaderManager
	log           *slog.Logger
}

// NewNodesView creates a new swarm nodes view
func NewNodesView(
	ui uiinterfaces.UIInterface,
	nodeService *swarm.NodeService,
	modalManager *managers.ModalManager,
	headerManager *managers.HeaderManager,
) *NodesView {
	headers := []string{"ID", "Hostname", "Role", "Availability", "Status", "Manager Status", "Engine Version", "Address"}

	view := &NodesView{
		BaseView:      shared.NewBaseView[shared.SwarmNode](ui, "Swarm Nodes", headers),
		nodeService:   nodeService,
		modalManager:  modalManager,
		headerManager: headerManager,
		log:           logger.GetLogger(),
	}

	view.setupCallbacks()

	return view
}

// Render renders the swarm nodes view
func (v *NodesView) Render(_ context.Context) error {
	// The base view handles rendering automatically through the callbacks
	// Just refresh the data
	v.Refresh()
	return nil
}

// HandleInput handles user input for the nodes view
func (v *NodesView) HandleInput(ctx context.Context, input rune) (interface{}, error) {
	switch input {
	case 'i':
		return v.handleInspect(ctx)
	case 'a':
		return v.handleUpdateAvailability(ctx)
	case 'r':
		return v.handleRemove(ctx)
	case 'f':
		return v, nil // Refresh current view
	case 's':
		return v.handleNavigateToServices(ctx)
	case 'q':
		return v.handleBackToMain(ctx)
	case 'h':
		v.handleHelp()
		return v, nil
	default:
		return v, nil
	}
}

// setupCallbacks sets up the callbacks for the base view
func (v *NodesView) setupCallbacks() {
	v.ListItems = v.listNodes
	v.FormatRow = func(n shared.SwarmNode) []string { return v.formatNodeRow(&n) }
	v.GetItemID = func(n shared.SwarmNode) string { return v.getNodeID(&n) }
	v.GetItemName = func(n shared.SwarmNode) string { return v.getNodeName(&n) }
	v.GetActions = v.getActions
}

// listNodes lists all swarm nodes
func (v *NodesView) listNodes(ctx context.Context) ([]shared.SwarmNode, error) {
	return v.nodeService.ListNodes(ctx)
}

// formatNodeRow formats a node row for display
func (v *NodesView) formatNodeRow(node *shared.SwarmNode) []string {
	return []string{
		shared.TruncName(node.ID, 12),
		node.Hostname,
		node.Role,
		node.Availability,
		node.Status,
		node.ManagerStatus,
		node.EngineVersion,
		node.Address,
	}
}

// getNodeID returns the node ID
func (v *NodesView) getNodeID(node *shared.SwarmNode) string {
	return node.ID
}

// getNodeName returns the node name
func (v *NodesView) getNodeName(node *shared.SwarmNode) string {
	return node.Hostname
}

// getActions returns the available actions for swarm nodes
func (v *NodesView) getActions() map[rune]string {
	return map[rune]string{
		'i': "Inspect",
		'a': "Update Availability",
		'r': "Remove",
		'f': "Refresh",
	}
}

// handleInspect handles node inspection
func (v *NodesView) handleInspect(ctx context.Context) (interface{}, error) {
	// Get the currently selected node
	selectedNode := v.GetSelectedItem()
	if selectedNode == nil {
		v.GetUI().ShowError(fmt.Errorf("please select a node first"))
		return v, fmt.Errorf("no node selected")
	}

	// Get the node service to access inspection functionality
	nodeService := v.GetUI().GetSwarmNodeService()
	if nodeService == nil {
		v.GetUI().ShowError(fmt.Errorf("swarm node service is not available - please check your Docker connection"))
		return v, fmt.Errorf("swarm node service not available")
	}

	// Cast to the correct type
	if swarmNodeService, ok := nodeService.(*swarm.NodeService); ok {
		// Get detailed node information
		nodeInfo, err := swarmNodeService.InspectNode(ctx, selectedNode.ID)
		if err != nil {
			errorMsg := fmt.Sprintf("failed to inspect node '%s': %v\n\nPlease check:\n• Node is accessible\n• You have sufficient permissions\n• Docker daemon is running", selectedNode.Hostname, err)
			v.GetUI().ShowError(fmt.Errorf("%s", errorMsg))
			return v, fmt.Errorf("failed to inspect node: %w", err)
		}

		// Show node information in a modal
		if nodeInfo == nil {
			v.GetUI().ShowInfo(fmt.Sprintf("Node '%s' has no detailed information available", selectedNode.Hostname))
		} else {
			// Format node information for display
			infoText := fmt.Sprintf("Node Details: %s\n\n"+
				"ID: %s\n"+
				"Role: %s\n"+
				"Availability: %s\n"+
				"Status: %s\n"+
				"Manager Status: %s\n"+
				"Engine Version: %s\n"+
				"Address: %s\n"+
				"CPUs: %d\n"+
				"Memory: %d",
				selectedNode.Hostname, shared.TruncName(selectedNode.ID, 12), selectedNode.Role, selectedNode.Availability,
				selectedNode.Status, selectedNode.ManagerStatus, selectedNode.EngineVersion,
				selectedNode.Address, selectedNode.CPUs, selectedNode.Memory)

			v.GetUI().ShowInfo(infoText)
		}

		return v, nil
	}

	v.GetUI().ShowError(fmt.Errorf("swarm node service is not properly configured"))
	return v, fmt.Errorf("swarm node service not available")
}

// handleUpdateAvailability handles node availability updates
func (v *NodesView) handleUpdateAvailability(ctx context.Context) (interface{}, error) {
	// Get the currently selected node
	selectedNode := v.GetSelectedItem()
	if selectedNode == nil {
		v.GetUI().ShowError(fmt.Errorf("please select a node first"))
		return v, fmt.Errorf("no node selected")
	}

	// Get the node service to access availability update functionality
	nodeService := v.GetUI().GetSwarmNodeService()
	if nodeService == nil {
		v.GetUI().ShowError(fmt.Errorf("swarm node service is not available - please check your Docker connection"))
		return v, fmt.Errorf("swarm node service not available")
	}

	// Cast to the correct type
	if swarmNodeService, ok := nodeService.(*swarm.NodeService); ok {
		// Show availability update modal
		v.GetUI().ShowNodeAvailabilityModal(selectedNode.Hostname, selectedNode.Availability, func(newAvailability string) {
			// Callback when user confirms availability update
			err := swarmNodeService.UpdateNodeAvailability(ctx, selectedNode.ID, newAvailability)
			if err != nil {
				// Enhanced error handling with retry and fallback options
				v.handleAvailabilityUpdateError(ctx, selectedNode, newAvailability, err, swarmNodeService)
			} else {
				// Show success feedback and refresh
				v.Refresh()
			}
		})

		return v, nil
	}

	v.GetUI().ShowError(fmt.Errorf("swarm node service is not properly configured"))
	return v, fmt.Errorf("swarm node service not available")
}

// handleAvailabilityUpdateError handles availability update errors with advanced recovery options
func (v *NodesView) handleAvailabilityUpdateError(ctx context.Context, node *shared.SwarmNode, newAvailability string, err error, swarmNodeService *swarm.NodeService) {
	// Check if this is a retryable error
	if v.isRetryableError(err) {
		// Show retry dialog with automatic retry option
		v.GetUI().ShowRetryDialog(
			fmt.Sprintf("update availability for node '%s' to '%s'", node.Hostname, newAvailability),
			err,
			func() error {
				// Retry function
				return swarmNodeService.UpdateNodeAvailability(ctx, node.ID, newAvailability)
			},
			func() {
				// Success callback
				v.GetUI().ShowInfo(fmt.Sprintf("Node '%s' availability successfully updated to '%s'", node.Hostname, newAvailability))
				v.Refresh()
			},
		)
	} else {
		// Show fallback options for non-retryable errors
		fallbackOptions := []string{
			"Check Node Status",
			"View Node Details",
			"Try Different Availability",
			"Check Swarm Health",
		}

		v.GetUI().ShowFallbackDialog(
			fmt.Sprintf("update availability for node '%s' to '%s'", node.Hostname, newAvailability),
			err,
			fallbackOptions,
			func(fallbackOption string) {
				v.executeAvailabilityFallback(ctx, node, fallbackOption, newAvailability, swarmNodeService)
			},
		)
	}
}

// isRetryableError determines if an error is retryable
func (v *NodesView) isRetryableError(err error) bool {
	errStr := err.Error()

	// Common retryable errors
	retryablePatterns := []string{
		"connection refused",
		"timeout",
		"temporary failure",
		"network unreachable",
		"service unavailable",
		"too many requests",
		"rate limit exceeded",
		"node is busy",
		"operation in progress",
	}

	for _, pattern := range retryablePatterns {
		if strings.Contains(strings.ToLower(errStr), pattern) {
			return true
		}
	}

	return false
}

// executeAvailabilityFallback executes fallback operations for availability update failures
func (v *NodesView) executeAvailabilityFallback(ctx context.Context, node *shared.SwarmNode, fallbackOption, _ string, swarmNodeService *swarm.NodeService) {
	switch fallbackOption {
	case "Check Node Status":
		// Show current node status
		nodeInfo, err := swarmNodeService.InspectNode(ctx, node.ID)
		if err != nil {
			v.GetUI().ShowError(fmt.Errorf("failed to check node status: %v", err))
		} else {
			statusText := fmt.Sprintf("Node '%s' Status:\n\n"+
				"Current Availability: %v\n"+
				"Status: %v\n"+
				"Role: %v\n"+
				"Manager Status: %v",
				node.Hostname, nodeInfo["Availability"], nodeInfo["Status"], nodeInfo["Role"], nodeInfo["ManagerStatus"])
			v.GetUI().ShowInfo(statusText)
		}

	case "View Node Details":
		// Show detailed node information
		nodeInfo, err := swarmNodeService.InspectNode(ctx, node.ID)
		if err != nil {
			v.GetUI().ShowError(fmt.Errorf("failed to get node details: %v", err))
		} else {
			detailsText := fmt.Sprintf("Node '%s' Details:\n\n"+
				"ID: %s\n"+
				"Role: %s\n"+
				"Availability: %v\n"+
				"Status: %v\n"+
				"Manager Status: %v\n"+
				"Engine Version: %v\n"+
				"Address: %s\n"+
				"CPUs: %d\n"+
				"Memory: %d",
				node.Hostname, shared.TruncName(node.ID, 12), node.Role, nodeInfo["Availability"],
				nodeInfo["Status"], nodeInfo["ManagerStatus"], nodeInfo["EngineVersion"],
				node.Address, node.CPUs, node.Memory)
			v.GetUI().ShowInfo(detailsText)
		}

	case "Try Different Availability":
		// Show availability modal again with current availability
		v.GetUI().ShowNodeAvailabilityModal(node.Hostname, node.Availability, func(differentAvailability string) {
			err := swarmNodeService.UpdateNodeAvailability(ctx, node.ID, differentAvailability)
			if err != nil {
				v.GetUI().ShowError(fmt.Errorf("failed to update availability to '%s': %v", differentAvailability, err))
			} else {
				v.GetUI().ShowInfo(fmt.Sprintf("Node '%s' availability successfully updated to '%s'", node.Hostname, differentAvailability))
				v.Refresh()
			}
		})

	case "Check Swarm Health":
		// Show swarm health information
		// This would require additional swarm service methods
		v.GetUI().ShowInfo(fmt.Sprintf("Swarm Health Check for Node '%s':\n\n"+
			"Current Availability: %s\n"+
			"Status: %s\n"+
			"Role: %s\n\n"+
			"Note: Detailed swarm health information requires additional swarm service methods.",
			node.Hostname, node.Availability, node.Status, node.Role))
	}
}

// handleRemove handles node removal
func (v *NodesView) handleRemove(ctx context.Context) (interface{}, error) {
	// Get the currently selected node
	selectedNode := v.GetSelectedItem()
	if selectedNode == nil {
		v.GetUI().ShowError(fmt.Errorf("please select a node first"))
		return v, fmt.Errorf("no node selected")
	}

	// Show confirmation dialog with more detailed information
	message := fmt.Sprintf("⚠️  Remove Node Confirmation\n\n"+
		"Node: %s\n"+
		"ID: %s\n"+
		"Role: %s\n"+
		"Status: %s\n\n"+
		"This action will:\n"+
		"• Force removal of the node\n"+
		"• Stop all tasks running on this node\n"+
		"• Cannot be undone\n\n"+
		"⚠️  Warning: Removing a manager node may affect swarm stability\n\n"+
		"Are you sure you want to continue?",
		selectedNode.Hostname, shared.TruncName(selectedNode.ID, 12), selectedNode.Role, selectedNode.Status)

	v.GetUI().ShowConfirm(message, func(confirmed bool) {
		if confirmed {
			// Get the node service to access removal functionality
			nodeService := v.GetUI().GetSwarmNodeService()
			if nodeService == nil {
				v.GetUI().ShowError(fmt.Errorf("swarm node service is not available - please check your Docker connection"))
				return
			}

			// Cast to the correct type
			if swarmNodeService, ok := nodeService.(swarm.NodeService); ok {
				// Remove the node (force removal)
				err := swarmNodeService.RemoveNode(ctx, selectedNode.ID, true)
				if err != nil {
					// Show detailed error modal with recovery suggestions
					errorMsg := fmt.Sprintf("failed to remove node '%s': %v\n\nPlease check:\n• Docker daemon is running\n• You have sufficient permissions\n• Node is not the last manager node", selectedNode.Hostname, err)
					v.GetUI().ShowError(fmt.Errorf("%s", errorMsg))
				} else {
					// Show success feedback and refresh
					v.GetUI().ShowInfo(fmt.Sprintf("Node '%s' successfully removed from swarm", selectedNode.Hostname))
					v.Refresh()
				}
			} else {
				v.GetUI().ShowError(fmt.Errorf("swarm node service is not properly configured"))
			}
		}
	})

	return v, nil
}

// handleNavigateToServices handles navigation to swarm services view
func (v *NodesView) handleNavigateToServices(_ context.Context) (interface{}, error) {
	// This would return a services view - placeholder for now
	return v, fmt.Errorf("services view not implemented yet")
}

// handleBackToMain handles navigation back to main menu
func (v *NodesView) handleBackToMain(_ context.Context) (interface{}, error) {
	// This would return the main menu view - placeholder for now
	return v, fmt.Errorf("main menu view not implemented yet")
}

// handleHelp shows contextual help for swarm nodes
func (v *NodesView) handleHelp() {
	// Show general swarm nodes help
	v.GetUI().ShowContextualHelp("swarm_nodes", "")
}

// showOperationHelp shows contextual help for a specific operation
// This function is intentionally unused - it's a placeholder for future help functionality
// nolint:unused // Intentionally unused - placeholder for future help functionality
func (v *NodesView) showOperationHelp(operation string) {
	v.GetUI().ShowContextualHelp("swarm_nodes", operation)
}
