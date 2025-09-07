package swarm

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/wikczerski/whaletui/internal/logger"
	"github.com/wikczerski/whaletui/internal/shared"
	"github.com/wikczerski/whaletui/internal/ui/interfaces"
	"github.com/wikczerski/whaletui/internal/ui/utils"
)

// NodesView represents the swarm nodes view
type NodesView struct {
	*shared.BaseView[shared.SwarmNode]
	nodeService   *NodeService
	modalManager  interfaces.ModalManagerInterface
	headerManager interfaces.HeaderManagerInterface
	log           *slog.Logger
}

// getUI safely gets the UI interface
// NewNodesView creates a new swarm nodes view
func NewNodesView(
	ui interfaces.UIInterface,
	nodeService *NodeService,
	modalManager interfaces.ModalManagerInterface,
	headerManager interfaces.HeaderManagerInterface,
) *NodesView {
	headers := []string{
		"ID", "Hostname", "Role", "Availability", "Status",
		"Manager Status", "Engine Version", "Address",
	}

	view := &NodesView{
		BaseView:      shared.NewBaseView[shared.SwarmNode](ui, "Swarm Nodes", headers),
		nodeService:   nodeService,
		modalManager:  modalManager,
		headerManager: headerManager,
		log:           logger.GetLogger(),
	}

	view.setupCallbacks()
	view.setupCharacterLimits(ui)

	return view
}

func (v *NodesView) getUI() shared.SharedUIInterface {
	return v.GetUI()
}

// setupCallbacks sets up the callbacks for the base view
func (v *NodesView) setupCallbacks() {
	v.ListItems = v.listNodes
	v.FormatRow = func(n shared.SwarmNode) []string { return v.formatNodeRow(&n) }
	v.GetItemID = func(n shared.SwarmNode) string { return v.getNodeID(&n) }
	v.GetItemName = func(n shared.SwarmNode) string { return v.getNodeName(&n) }
	v.GetActions = v.getActions
	v.HandleKeyPress = func(key rune, n shared.SwarmNode) { v.handleAction(key, &n) }
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
		// Enhanced error handling with retry and fallback options
		v.handleAvailabilityUpdateError(ctx, selectedNode, newAvailability, err, swarmNodeService)
	} else {
		// Show success feedback and refresh
		v.Refresh()
	}
}

// handleAvailabilityUpdateError handles availability update errors with advanced recovery options
func (v *NodesView) handleAvailabilityUpdateError(
	ctx context.Context,
	node *shared.SwarmNode,
	newAvailability string,
	err error,
	swarmNodeService *NodeService,
) {
	// Check if this is a retryable error
	if v.isRetryableError(err) {
		v.showRetryDialog(ctx, node, newAvailability, err, swarmNodeService)
	} else {
		v.showFallbackDialog(ctx, node, newAvailability, err, swarmNodeService)
	}
}

// showRetryDialog shows the retry dialog for retryable errors
func (v *NodesView) showRetryDialog(
	ctx context.Context,
	node *shared.SwarmNode,
	newAvailability string,
	err error,
	swarmNodeService *NodeService,
) {
	if ui, ok := v.GetUI().(interface {
		ShowRetryDialog(string, error, func() error, func())
	}); ok {
		ui.ShowRetryDialog(
			fmt.Sprintf("update availability for node '%s' to '%s'",
				node.Hostname, newAvailability),
			err,
			func() error {
				// Retry function
				return swarmNodeService.UpdateNodeAvailability(ctx, node.ID, newAvailability)
			},
			func() {
				// Success callback
				if infoUI, ok := v.GetUI().(interface{ ShowInfo(string) }); ok {
					message := fmt.Sprintf("Node '%s' availability successfully updated to '%s'",
						node.Hostname, newAvailability)
					infoUI.ShowInfo(message)
				}
				v.Refresh()
			},
		)
	}
}

// showFallbackDialog shows the fallback dialog for non-retryable errors
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
		"Check Swarm Health",
	}

	if ui, ok := v.GetUI().(interface {
		ShowFallbackDialog(string, error, []string, func(string))
	}); ok {
		ui.ShowFallbackDialog(
			fmt.Sprintf("update availability for node '%s' to '%s'",
				node.Hostname, newAvailability),
			err,
			fallbackOptions,
			func(fallbackOption string) {
				v.executeAvailabilityFallback(ctx, node, fallbackOption,
					newAvailability, swarmNodeService)
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
	case "Check Swarm Health":
		v.handleCheckSwarmHealth(node)
	}
}

// handleCheckNodeStatus handles checking node status
func (v *NodesView) handleCheckNodeStatus(
	ctx context.Context,
	node *shared.SwarmNode,
	swarmNodeService *NodeService,
) {
	nodeInfo, err := swarmNodeService.InspectNode(ctx, node.ID)
	if err != nil {
		v.GetUI().ShowError(fmt.Errorf("failed to check node status: %v", err))
		return
	}

	statusText := fmt.Sprintf("Node '%s' Status:\n\n"+
		"Current Availability: %v\n"+
		"Status: %v\n"+
		"Role: %v\n"+
		"Manager Status: %v",
		node.Hostname, nodeInfo["Availability"], nodeInfo["Status"],
		nodeInfo["Role"], nodeInfo["ManagerStatus"])
	v.GetUI().ShowInfo(statusText)
}

// handleViewNodeDetails handles viewing node details
func (v *NodesView) handleViewNodeDetails(
	ctx context.Context,
	node *shared.SwarmNode,
	swarmNodeService *NodeService,
) {
	nodeInfo, err := swarmNodeService.InspectNode(ctx, node.ID)
	if err != nil {
		if ui := v.getUI(); ui != nil {
			ui.ShowError(fmt.Errorf("failed to get node details: %v", err))
		}
		return
	}

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
	if ui := v.getUI(); ui != nil {
		ui.ShowInfo(detailsText)
	}
}

// handleTryDifferentAvailability handles trying different availability
func (v *NodesView) handleTryDifferentAvailability(
	ctx context.Context,
	node *shared.SwarmNode,
	swarmNodeService *NodeService,
) {
	if ui := v.getUI(); ui != nil {
		ui.ShowNodeAvailabilityModal(node.Hostname, node.Availability,
			func(differentAvailability string) {
				err := swarmNodeService.UpdateNodeAvailability(ctx, node.ID, differentAvailability)
				if err != nil {
					if ui := v.getUI(); ui != nil {
						ui.ShowError(fmt.Errorf("failed to update availability to '%s': %v",
							differentAvailability, err))
					}
				} else {
					if ui := v.getUI(); ui != nil {
						ui.ShowInfo(fmt.Sprintf("Node '%s' availability successfully updated to '%s'",
							node.Hostname, differentAvailability))
					}
					v.Refresh()
				}
			})
	}
}

// handleCheckSwarmHealth handles checking swarm health
func (v *NodesView) handleCheckSwarmHealth(node *shared.SwarmNode) {
	v.GetUI().ShowInfo(fmt.Sprintf("Swarm Health Check for Node '%s':\n\n"+
		"Current Availability: %s\n"+
		"Status: %s\n"+
		"Role: %s\n\n"+
		"Note: Detailed swarm health information requires additional swarm service methods.",
		node.Hostname, node.Availability, node.Status, node.Role))
}

// handleRemove handles node removal
func (v *NodesView) handleRemove(ctx context.Context) (any, error) {
	// Get the currently selected node
	selectedNode := v.GetSelectedItem()
	if selectedNode == nil {
		v.GetUI().ShowError(errors.New("please select a node first"))
		return v, errors.New("no node selected")
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
		v.GetUI().
			ShowError(errors.New("swarm node service is not available - please check your Docker connection"))
		return
	}

	// Cast to the correct type
	if swarmNodeService, ok := nodeService.(NodeService); ok {
		// Remove the node (force removal)
		err := swarmNodeService.RemoveNode(ctx, selectedNode.ID, true)
		if err != nil {
			v.handleNodeRemovalError(selectedNode, err)
		} else {
			v.handleNodeRemovalSuccess(selectedNode)
		}
	} else {
		v.GetUI().ShowError(errors.New("swarm node service is not properly configured"))
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
	v.GetUI().
		ShowInfo(fmt.Sprintf("Node '%s' successfully removed from swarm", selectedNode.Hostname))
	v.Refresh()
}

// handleAction handles action key presses for swarm nodes
func (v *NodesView) handleAction(key rune, node *shared.SwarmNode) {
	ctx := context.Background()

	switch key {
	case 'i':
		// TODO: Implement inspect functionality
		v.log.Info("Inspect action not yet implemented")
	case 'a':
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

// setupCharacterLimits sets up character limits for table columns
func (v *NodesView) setupCharacterLimits(ui interfaces.UIInterface) {
	// Define column types for swarm nodes table:
	// ID, Hostname, Role, Availability, Status, Manager Status, Engine Version, Address
	columnTypes := []string{
		"id", "name", "hostname", "role", "availability",
		"status", "manager_status", "engine_version", "address",
	}
	v.SetColumnTypes(columnTypes)

	// Create formatter from theme manager
	formatter := utils.NewTableFormatterFromTheme(ui.GetThemeManager())
	v.SetFormatter(formatter)
}
