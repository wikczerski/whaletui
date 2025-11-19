package services

import (
	"context"
	"fmt"
	"log/slog"
	"sort"

	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	domaintypes "github.com/wikczerski/whaletui/internal/docker/types"
	"github.com/wikczerski/whaletui/internal/docker/utils"
)

// NetworkService handles network-related operations
type NetworkService struct {
	cli *client.Client
	log *slog.Logger
}

// NewNetworkService creates a new NetworkService
func NewNetworkService(cli *client.Client, log *slog.Logger) *NetworkService {
	return &NetworkService{
		cli: cli,
		log: log,
	}
}

// ListNetworks lists all networks
func (s *NetworkService) ListNetworks(ctx context.Context) ([]domaintypes.Network, error) {
	networks, err := s.cli.NetworkList(ctx, network.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list networks: %w", err)
	}

	result := make([]domaintypes.Network, 0, len(networks))
	for i := range networks {
		net := &networks[i]
		result = append(result, domaintypes.Network{
			ID:         net.ID[:12],
			Name:       net.Name,
			Driver:     net.Driver,
			Scope:      net.Scope,
			Created:    net.Created,
			Internal:   net.Internal,
			Attachable: net.Attachable,
			Ingress:    net.Ingress,
			IPv6:       net.EnableIPv6,
			EnableIPv6: net.EnableIPv6,
			Labels:     net.Labels,
			Containers: len(net.Containers),
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})

	return result, nil
}

// InspectNetwork inspects a network
func (s *NetworkService) InspectNetwork(ctx context.Context, id string) (map[string]any, error) {
	networkInfo, err := s.cli.NetworkInspect(ctx, id, network.InspectOptions{})
	if err != nil {
		return nil, fmt.Errorf("network inspect failed %s: %w", id, err)
	}
	return utils.MarshalToMap(networkInfo)
}

// RemoveNetwork removes a network
func (s *NetworkService) RemoveNetwork(ctx context.Context, id string) error {
	if err := utils.ValidateID(id, "network ID"); err != nil {
		return err
	}

	if err := s.cli.NetworkRemove(ctx, id); err != nil {
		return fmt.Errorf("failed to remove network %s: %w", id, err)
	}
	return nil
}
