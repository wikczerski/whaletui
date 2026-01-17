package swarm

import (
	"context"
	"fmt"
	"strings"

	"github.com/wikczerski/whaletui/internal/shared"
)

// handleAction handles action key presses for swarm nodes
func (v *NodesView) handleAction(key rune, _ *shared.SwarmNode) {
	ctx := context.Background()

	switch key {
	case 'i':
		v.handleInspectAction(ctx)
	case 'a':
		v.handleUpdateAvailabilityAction(ctx)
	case 'r':
		if _, err := v.handleRemove(ctx); err != nil {
			v.log.Error("Failed to handle remove action", "error", err)
		}
	case 'f':
		v.Refresh()
	default:
		v.log.Warn("Unknown action key", "key", string(key))
	}
}

// handleInspectAction handles the inspect action
func (v *NodesView) handleInspectAction(ctx context.Context) {
	selectedNode := v.GetSelectedItem()
	if selectedNode == nil {
		v.log.Warn("No node selected for inspection")
		return
	}
	v.inspectNode(ctx, selectedNode)
}

// handleUpdateAvailabilityAction handles the update availability action
func (v *NodesView) handleUpdateAvailabilityAction(ctx context.Context) {
	selectedNode := v.GetSelectedItem()
	if selectedNode == nil {
		v.log.Warn("No node selected for availability update")
		return
	}

	nodeService := v.GetUI().GetSwarmNodeService()
	if nodeService == nil {
		v.log.Warn("Swarm node service not available")
		return
	}

	swarmNodeService, ok := nodeService.(*NodeService)
	if !ok {
		v.log.Warn("Swarm node service type assertion failed")
		return
	}

	v.showAvailabilityUpdateModal(ctx, selectedNode, swarmNodeService)
}

// inspectNode performs node inspection and shows details
func (v *NodesView) inspectNode(ctx context.Context, node *shared.SwarmNode) {
	nodeInfo, err := v.nodeService.InspectNode(ctx, node.ID)
	v.ShowItemDetails(*node, nodeInfo, err)
}

// handleRemove handles node removal
func (v *NodesView) handleRemove(ctx context.Context) (any, error) {
	// Get the currently selected node
	selectedNode := v.GetSelectedItem()
	if selectedNode == nil {
		return v, fmt.Errorf("no node selected")
	}

	// Show confirmation dialog with more detailed information
	message := v.buildRemoveNodeConfirmationMessage(selectedNode)
	v.GetUI().ShowConfirm(message, func(confirmed bool) {
		if confirmed {
			v.performNodeRemoval(ctx, selectedNode)
		}
	})

	return v, nil
}

// buildRemoveNodeConfirmationMessage builds the confirmation message for node removal
func (v *NodesView) buildRemoveNodeConfirmationMessage(selectedNode *shared.SwarmNode) string {
	return fmt.Sprintf(
		"⚠️  Remove Node Confirmation\n\n"+
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
		selectedNode.Hostname,
		shared.TruncName(selectedNode.ID, 12),
		selectedNode.Role,
		selectedNode.Status,
	)
}

// performNodeRemoval performs the actual node removal operation
func (v *NodesView) performNodeRemoval(ctx context.Context, selectedNode *shared.SwarmNode) {
	// Get the node service to access removal functionality
	nodeService := v.GetUI().GetSwarmNodeService()
	if nodeService == nil {
		v.GetUI().ShowError(fmt.Errorf("swarm node service is not available"))
		return
	}

	// Cast to the correct type
	if swarmNodeService, ok := nodeService.(*NodeService); ok {
		// Remove the node (force removal)
		err := swarmNodeService.RemoveNode(ctx, selectedNode.ID, true)
		if err != nil {
			v.handleNodeRemovalError(selectedNode, err)
		} else {
			v.handleNodeRemovalSuccess(selectedNode)
		}
	}
}

// handleNodeRemovalError handles errors during node removal
func (v *NodesView) handleNodeRemovalError(selectedNode *shared.SwarmNode, err error) {
	errorMsg := fmt.Sprintf(
		"failed to remove node '%s': %v\n\nPlease check:\n"+
			"• Docker daemon is running\n"+
			"• You have sufficient permissions\n"+
			"• Node is not the last manager node",
		selectedNode.Hostname,
		err,
	)
	v.GetUI().ShowError(fmt.Errorf("%s", errorMsg))
}

// handleNodeRemovalSuccess handles successful node removal
func (v *NodesView) handleNodeRemovalSuccess(selectedNode *shared.SwarmNode) {
	v.GetUI().ShowInfo(fmt.Sprintf("Node '%s' successfully removed from swarm", selectedNode.Hostname))
	v.Refresh()
}

// showAvailabilityUpdateModal shows the availability update modal
func (v *NodesView) showAvailabilityUpdateModal(
	ctx context.Context,
	selectedNode *shared.SwarmNode,
	swarmNodeService *NodeService,
) {
	v.GetUI().ShowNodeAvailabilityModal(
		selectedNode.Hostname,
		selectedNode.Availability,
		func(newAvailability string) {
			v.handleAvailabilityUpdate(ctx, selectedNode, newAvailability, swarmNodeService)
		},
	)
}

// handleAvailabilityUpdate handles the actual availability update
func (v *NodesView) handleAvailabilityUpdate(
	ctx context.Context,
	selectedNode *shared.SwarmNode,
	newAvailability string,
	swarmNodeService *NodeService,
) {
	err := swarmNodeService.UpdateNodeAvailability(ctx, selectedNode.ID, newAvailability)
	if err != nil {
		v.handleAvailabilityUpdateError(ctx, selectedNode, newAvailability, err, swarmNodeService)
	} else {
		v.Refresh()
	}
}

// handleAvailabilityUpdateError handles availability update errors
func (v *NodesView) handleAvailabilityUpdateError(
	ctx context.Context,
	node *shared.SwarmNode,
	newAvailability string,
	err error,
	swarmNodeService *NodeService,
) {
	if v.isRetryableError(err) {
		v.showRetryDialog(ctx, node, newAvailability, err, swarmNodeService)
	} else {
		v.showFallbackDialog(ctx, node, newAvailability, err, swarmNodeService)
	}
}

// showRetryDialog shows the retry dialog
func (v *NodesView) showRetryDialog(
	ctx context.Context,
	node *shared.SwarmNode,
	newAvailability string,
	err error,
	swarmNodeService *NodeService,
) {
	v.GetUI().ShowRetryDialog(
		fmt.Sprintf("update availability for node '%s' to '%s'", node.Hostname, newAvailability),
		err,
		func() error {
			return swarmNodeService.UpdateNodeAvailability(ctx, node.ID, newAvailability)
		},
		func() {
			v.GetUI().ShowInfo(fmt.Sprintf("Node '%s' availability updated", node.Hostname))
			v.Refresh()
		},
	)
}

// showFallbackDialog shows the fallback dialog
func (v *NodesView) showFallbackDialog(
	ctx context.Context,
	node *shared.SwarmNode,
	newAvailability string,
	err error,
	swarmNodeService *NodeService,
) {
	fallbackOptions := []string{
		"Check Node Status",
		"View Node Details",
		"Try Different Availability",
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

// isRetryableError determines if an error is retryable
func (v *NodesView) isRetryableError(err error) bool {
	errStr := strings.ToLower(err.Error())
	retryablePatterns := []string{
		"connection refused", "timeout", "temporary failure",
		"network unreachable", "service unavailable",
	}

	for _, pattern := range retryablePatterns {
		if strings.Contains(errStr, pattern) {
			return true
		}
	}
	return false
}

// executeAvailabilityFallback executes fallback operations
func (v *NodesView) executeAvailabilityFallback(
	ctx context.Context,
	node *shared.SwarmNode,
	fallbackOption, _ string,
	swarmNodeService *NodeService,
) {
	switch fallbackOption {
	case "Check Node Status":
		v.handleCheckNodeStatus(ctx, node, swarmNodeService)
	case "View Node Details":
		v.handleViewNodeDetails(ctx, node, swarmNodeService)
	case "Try Different Availability":
		v.handleTryDifferentAvailability(ctx, node, swarmNodeService)
	}
}

// handleCheckNodeStatus handles checking node status
func (v *NodesView) handleCheckNodeStatus(ctx context.Context, node *shared.SwarmNode, swarmNodeService *NodeService) {
	nodeInfo, err := swarmNodeService.InspectNode(ctx, node.ID)
	if err != nil {
		v.GetUI().ShowError(err)
		return
	}
	v.GetUI().ShowInfo(fmt.Sprintf("Node %s: %v", node.Hostname, nodeInfo["Status"]))
}

// handleViewNodeDetails handles viewing node details
func (v *NodesView) handleViewNodeDetails(ctx context.Context, node *shared.SwarmNode, swarmNodeService *NodeService) {
	nodeInfo, err := swarmNodeService.InspectNode(ctx, node.ID)
	if err != nil {
		v.GetUI().ShowError(err)
		return
	}
	v.GetUI().ShowInfo(fmt.Sprintf("Details for %s: %v", node.Hostname, nodeInfo))
}

// handleTryDifferentAvailability handles trying different availability
func (v *NodesView) handleTryDifferentAvailability(ctx context.Context, node *shared.SwarmNode, swarmNodeService *NodeService) {
	v.GetUI().ShowNodeAvailabilityModal(node.Hostname, node.Availability, func(availability string) {
		err := swarmNodeService.UpdateNodeAvailability(ctx, node.ID, availability)
		if err != nil {
			v.GetUI().ShowError(err)
		} else {
			v.Refresh()
		}
	})
}
