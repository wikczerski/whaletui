package services

import (
	"context"
	"fmt"

	"github.com/wikczerski/D5r/internal/docker"
	"github.com/wikczerski/D5r/internal/models"
)

// networkService implements NetworkService
type networkService struct {
	client *docker.Client
}

// NewNetworkService creates a new network service
func NewNetworkService(client *docker.Client) NetworkService {
	return &networkService{
		client: client,
	}
}

// ListNetworks retrieves all networks
func (s *networkService) ListNetworks(ctx context.Context) ([]models.Network, error) {
	networks, err := s.client.ListNetworks(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list networks: %w", err)
	}

	var result []models.Network
	for _, net := range networks {
		result = append(result, models.Network{
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
	return s.client.InspectNetwork(ctx, id)
}
