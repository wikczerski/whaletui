package swarm

import (
	"context"
	"log/slog"

	"github.com/wikczerski/whaletui/internal/logger"
	"github.com/wikczerski/whaletui/internal/shared"
	"github.com/wikczerski/whaletui/internal/ui/interfaces"
	uishared "github.com/wikczerski/whaletui/internal/ui/shared"
	"github.com/wikczerski/whaletui/internal/ui/utils"
)

// NodesView represents the swarm nodes view
type NodesView struct {
	*uishared.BaseView[shared.SwarmNode]
	nodeService   *NodeService
	modalManager  interfaces.ModalManagerInterface
	headerManager interfaces.HeaderManagerInterface
	log           *slog.Logger
}

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
	baseView := uishared.NewBaseView[shared.SwarmNode](ui, "swarm nodes", headers)

	view := &NodesView{
		BaseView:      baseView,
		nodeService:   nodeService,
		modalManager:  modalManager,
		headerManager: headerManager,
		log:           logger.GetLogger(),
	}

	view.setupCallbacks()
	view.setupCharacterLimits(ui)

	return view
}

// setupCallbacks sets up the callbacks for the base view
func (v *NodesView) setupCallbacks() {
	v.setupBasicCallbacks()
	v.setupActionCallbacks()
}

// setupBasicCallbacks sets up the basic view callbacks
func (v *NodesView) setupBasicCallbacks() {
	v.ListItems = v.listNodes
	v.FormatRow = func(n shared.SwarmNode) []string { return v.formatNodeRow(&n) }
	v.GetItemID = func(n shared.SwarmNode) string { return v.getNodeID(&n) }
	v.GetItemName = func(n shared.SwarmNode) string { return v.getNodeName(&n) }
}

// setupActionCallbacks sets up the action-related callbacks
func (v *NodesView) setupActionCallbacks() {
	v.HandleKeyPress = func(key rune, n shared.SwarmNode) { v.handleAction(key, &n) }
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
