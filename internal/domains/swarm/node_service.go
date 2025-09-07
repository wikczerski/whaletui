// Package swarm provides Docker Swarm functionality for WhaleTUI.
package swarm

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/docker/docker/api/types/swarm"
	"github.com/wikczerski/whaletui/internal/docker"
	"github.com/wikczerski/whaletui/internal/logger"
	"github.com/wikczerski/whaletui/internal/shared"
	"github.com/wikczerski/whaletui/internal/shared/interfaces"
)

// NodeService implements the SwarmNodeService interface
type NodeService struct {
	client *docker.Client
	log    *slog.Logger
}

// NewNodeService creates a new swarm node service
func NewNodeService(client *docker.Client) interfaces.SwarmNodeService {
	return &NodeService{
		client: client,
		log:    logger.GetLogger(),
	}
}

// ListNodes lists all swarm nodes
func (n *NodeService) ListNodes(ctx context.Context) ([]shared.SwarmNode, error) {
	dockerNodes, err := n.client.ListSwarmNodes(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list swarm nodes: %w", err)
	}

	nodes := make([]shared.SwarmNode, len(dockerNodes))
	for i := range dockerNodes {
		nodes[i] = n.convertToSharedNode(dockerNodes[i])
	}

	return nodes, nil
}

// InspectNode inspects a swarm node
func (n *NodeService) InspectNode(ctx context.Context, id string) (map[string]any, error) {
	node, err := n.client.InspectSwarmNode(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to inspect swarm node: %w", err)
	}

	// Convert to map for inspection
	result := map[string]any{
		"ID":            node.ID,
		"Hostname":      node.Description.Hostname,
		"Role":          string(node.Spec.Role),
		"Availability":  string(node.Spec.Availability),
		"Status":        string(node.Status.State),
		"ManagerStatus": getManagerStatus(node.ManagerStatus),
		"EngineVersion": node.Description.Engine.EngineVersion,
		"Address":       node.Status.Addr,
		"CPUs":          node.Description.Resources.NanoCPUs,
		"Memory":        node.Description.Resources.MemoryBytes,
		"Labels":        node.Spec.Labels,
	}

	return result, nil
}

// UpdateNodeAvailability updates a node's availability
func (n *NodeService) UpdateNodeAvailability(ctx context.Context, id, availability string) error {
	node, err := n.client.InspectSwarmNode(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to inspect node for availability update: %w", err)
	}

	if err := n.validateAvailability(availability); err != nil {
		return err
	}

	n.updateNodeSpecAvailability(node, availability)

	if err := n.client.UpdateSwarmNode(ctx, id, node.Version, node.Spec); err != nil {
		return fmt.Errorf("failed to update node availability: %w", err)
	}

	n.log.Info(
		"Node availability updated successfully",
		"node_id",
		id,
		"availability",
		availability,
	)
	return nil
}

// RemoveNode removes a swarm node
func (n *NodeService) RemoveNode(ctx context.Context, id string, force bool) error {
	err := n.client.RemoveSwarmNode(ctx, id, force)
	if err != nil {
		return fmt.Errorf("failed to remove swarm node: %w", err)
	}

	n.log.Info("Node removed successfully", "node_id", id, "force", force)
	return nil
}

// GetActions returns the available actions for swarm nodes
func (n *NodeService) GetActions() map[rune]string {
	return map[rune]string{
		'i': "Inspect",
		'a': "Update Availability",
		'r': "Remove",
	}
}

// GetActionsString returns the available actions as a string
func (n *NodeService) GetActionsString() string {
	return "<i>: Inspect\n<a>: Update Availability\n<r>: Remove"
}

// GetNavigation returns the available navigation options for swarm nodes as a map
func (n *NodeService) GetNavigation() map[rune]string {
	return map[rune]string{
		'↑': "Up",
		'↓': "Down",
		':': "Command",
		'/': "Filter",
	}
}

// GetNavigationString returns the available navigation options for swarm nodes as a formatted string
func (n *NodeService) GetNavigationString() string {
	return "↑/↓: Navigate\n<:> Command mode\n/: Filter"
}

// convertToSharedNode converts a Docker swarm node to shared node
// nolint:gocritic // Docker API requires value parameter for compatibility
func (n *NodeService) convertToSharedNode(node swarm.Node) shared.SwarmNode {
	return shared.SwarmNode{
		ID:            node.ID,
		Hostname:      node.Description.Hostname,
		Role:          string(node.Spec.Role),
		Availability:  string(node.Spec.Availability),
		Status:        string(node.Status.State),
		ManagerStatus: getManagerStatus(node.ManagerStatus),
		EngineVersion: node.Description.Engine.EngineVersion,
		Address:       node.Status.Addr,
		CPUs:          node.Description.Resources.NanoCPUs,
		Memory:        node.Description.Resources.MemoryBytes,
		Labels:        node.Spec.Labels,
	}
}

// validateAvailability validates the availability string
func (n *NodeService) validateAvailability(availability string) error {
	switch availability {
	case "active", "pause", "drain":
		return nil
	default:
		return fmt.Errorf("invalid availability: %s", availability)
	}
}

// updateNodeSpecAvailability updates the node spec availability
func (n *NodeService) updateNodeSpecAvailability(node swarm.Node, availability string) {
	switch availability {
	case "active":
		node.Spec.Availability = swarm.NodeAvailabilityActive
	case "pause":
		node.Spec.Availability = swarm.NodeAvailabilityPause
	case "drain":
		node.Spec.Availability = swarm.NodeAvailabilityDrain
	}
}

// getManagerStatus returns the manager status as a string
func getManagerStatus(managerStatus *swarm.ManagerStatus) string {
	if managerStatus == nil {
		return ""
	}
	return string(managerStatus.Reachability)
}
