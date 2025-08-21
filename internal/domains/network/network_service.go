package network

import (
	"context"
	"fmt"

	"github.com/wikczerski/whaletui/internal/docker"
	"github.com/wikczerski/whaletui/internal/shared/interfaces"
)

// networkService implements NetworkService
type networkService struct {
	client *docker.Client
}

// NewNetworkService creates a new network service
func NewNetworkService(client *docker.Client) interfaces.NetworkService {
	return &networkService{
		client: client,
	}
}

// ListNetworks retrieves all networks
func (s *networkService) ListNetworks(ctx context.Context) ([]Network, error) {
	if s.client == nil {
		return nil, fmt.Errorf("docker client is not initialized")
	}

	networks, err := s.client.ListNetworks(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list networks: %w", err)
	}

	result := make([]Network, 0, len(networks))
	for _, net := range networks {
		result = append(result, Network{
			ID:      net.ID,
			Name:    net.Name,
			Driver:  net.Driver,
			Scope:   net.Scope,
			Created: net.Created,
		})
	}

	return result, nil
}

// RemoveNetwork removes a network
func (s *networkService) RemoveNetwork(ctx context.Context, id string) error {
	if s.client == nil {
		return fmt.Errorf("docker client is not initialized")
	}

	if id == "" {
		return fmt.Errorf("network ID cannot be empty")
	}

	return s.client.RemoveNetwork(ctx, id)
}

// InspectNetwork inspects a network
func (s *networkService) InspectNetwork(ctx context.Context, id string) (map[string]any, error) {
	if s.client == nil {
		return nil, fmt.Errorf("docker client is not initialized")
	}

	return s.client.InspectNetwork(ctx, id)
}

// GetActions returns the available actions for networks as a map
func (s *networkService) GetActions() map[rune]string {
	return map[rune]string{
		'r': "Remove",
		'h': "History",
		'f': "Filter",
		't': "Sort",
		'i': "Inspect",
	}
}

// GetActionsString returns the available actions for networks as a formatted string
func (s *networkService) GetActionsString() string {
	return "<r> Remove\n<h> History\n<f> Filter\n<t> Sort\n<i> Inspect\n<enter> Details\n<up/down> Navigate\n<?> Help\n<f5> Refresh\n<:> Command"
}
