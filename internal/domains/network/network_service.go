package network

import (
	"context"
	"errors"
	"fmt"

	"github.com/wikczerski/whaletui/internal/docker"
	"github.com/wikczerski/whaletui/internal/shared"
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
func (s *networkService) ListNetworks(ctx context.Context) ([]shared.Network, error) {
	if s.client == nil {
		return nil, errors.New("docker client is not initialized")
	}

	networks, err := s.client.ListNetworks(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list networks: %w", err)
	}

	result := make([]shared.Network, 0, len(networks))
	for _, net := range networks {
		result = append(result, shared.Network{
			ID:         net.ID,
			Name:       net.Name,
			Driver:     net.Driver,
			Scope:      net.Scope,
			Created:    net.Created,
			Internal:   net.Internal,
			Attachable: net.Attachable,
			Ingress:    net.Ingress,
			EnableIPv6: net.EnableIPv6,
			Labels:     net.Labels,
		})
	}

	return result, nil
}

// RemoveNetwork removes a network
func (s *networkService) RemoveNetwork(ctx context.Context, id string) error {
	if s.client == nil {
		return errors.New("docker client is not initialized")
	}

	if id == "" {
		return errors.New("network ID cannot be empty")
	}

	return s.client.RemoveNetwork(ctx, id)
}

// InspectNetwork inspects a network
func (s *networkService) InspectNetwork(ctx context.Context, id string) (map[string]any, error) {
	if s.client == nil {
		return nil, errors.New("docker client is not initialized")
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
	return "<r> Remove\n<h> History\n<f> Filter\n<t> Sort\n<i> Inspect\n<enter> Details\n" +
		"<up/down> Navigate\n<?> Help\n<f5> Refresh\n<:> Command"
}

// GetNavigation returns the available navigation options for networks as a map
func (s *networkService) GetNavigation() map[rune]string {
	return map[rune]string{
		'↑': "Up",
		'↓': "Down",
		':': "Command",
		'/': "Filter",
	}
}

// GetNavigationString returns the available navigation options for networks as a formatted string
func (s *networkService) GetNavigationString() string {
	return "↑/↓: Navigate\n<:> Command mode\n/: Filter"
}
