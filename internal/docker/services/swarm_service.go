package services

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
	"github.com/wikczerski/whaletui/internal/docker/utils"
)

// SwarmService handles swarm-related operations
type SwarmService struct {
	cli *client.Client
	log *slog.Logger
}

// NewSwarmService creates a new SwarmService
func NewSwarmService(cli *client.Client, log *slog.Logger) *SwarmService {
	return &SwarmService{
		cli: cli,
		log: log,
	}
}

// ListSwarmServices lists all swarm services
func (s *SwarmService) ListSwarmServices(ctx context.Context) ([]swarm.Service, error) {
	if err := s.validateClient(); err != nil {
		return nil, err
	}

	services, err := s.cli.ServiceList(ctx, swarm.ServiceListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list swarm services: %w", err)
	}

	return services, nil
}

// InspectSwarmService inspects a swarm service
func (s *SwarmService) InspectSwarmService(ctx context.Context, id string) (swarm.Service, error) {
	if err := s.validateClient(); err != nil {
		return swarm.Service{}, err
	}

	if err := utils.ValidateID(id, "service ID"); err != nil {
		return swarm.Service{}, err
	}

	service, _, err := s.cli.ServiceInspectWithRaw(ctx, id, swarm.ServiceInspectOptions{})
	if err != nil {
		return swarm.Service{}, fmt.Errorf("failed to inspect swarm service %s: %w", id, err)
	}

	return service, nil
}

// UpdateSwarmService updates a swarm service
func (s *SwarmService) UpdateSwarmService(
	ctx context.Context,
	id string,
	version swarm.Version,
	spec swarm.ServiceSpec,
) error {
	if err := s.validateClient(); err != nil {
		return err
	}

	if err := utils.ValidateID(id, "service ID"); err != nil {
		return err
	}

	response, err := s.cli.ServiceUpdate(ctx, id, version, spec, swarm.ServiceUpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to update swarm service %s: %w", id, err)
	}

	if len(response.Warnings) > 0 {
		s.log.Warn("Service update warnings", "service_id", id, "warnings", response.Warnings)
	}

	return nil
}

// RemoveSwarmService removes a swarm service
func (s *SwarmService) RemoveSwarmService(ctx context.Context, id string) error {
	if err := s.validateClient(); err != nil {
		return err
	}

	if err := utils.ValidateID(id, "service ID"); err != nil {
		return err
	}

	if err := s.cli.ServiceRemove(ctx, id); err != nil {
		return fmt.Errorf("failed to remove swarm service %s: %w", id, err)
	}

	return nil
}

// GetSwarmServiceLogs gets logs for a swarm service
func (s *SwarmService) GetSwarmServiceLogs(ctx context.Context, id string) (string, error) {
	if err := s.validateServiceLogsRequest(id); err != nil {
		return "", err
	}

	response, err := s.getServiceLogsResponse(ctx, id)
	if err != nil {
		return "", err
	}
	defer s.closeServiceLogsResponse(response)

	return s.readServiceLogs(response), nil
}

// ListSwarmNodes lists all swarm nodes
func (s *SwarmService) ListSwarmNodes(ctx context.Context) ([]swarm.Node, error) {
	if err := s.validateClient(); err != nil {
		return nil, err
	}

	nodes, err := s.cli.NodeList(ctx, swarm.NodeListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list swarm nodes: %w", err)
	}

	return nodes, nil
}

// InspectSwarmNode inspects a swarm node
func (s *SwarmService) InspectSwarmNode(ctx context.Context, id string) (swarm.Node, error) {
	if err := s.validateClient(); err != nil {
		return swarm.Node{}, err
	}

	if err := utils.ValidateID(id, "node ID"); err != nil {
		return swarm.Node{}, err
	}

	node, _, err := s.cli.NodeInspectWithRaw(ctx, id)
	if err != nil {
		return swarm.Node{}, fmt.Errorf("failed to inspect swarm node %s: %w", id, err)
	}

	return node, nil
}

// UpdateSwarmNode updates a swarm node
func (s *SwarmService) UpdateSwarmNode(
	ctx context.Context,
	id string,
	version swarm.Version,
	spec swarm.NodeSpec,
) error {
	if err := s.validateClient(); err != nil {
		return err
	}

	if err := utils.ValidateID(id, "node ID"); err != nil {
		return err
	}

	if err := s.cli.NodeUpdate(ctx, id, version, spec); err != nil {
		return fmt.Errorf("failed to update swarm node %s: %w", id, err)
	}

	return nil
}

// RemoveSwarmNode removes a swarm node
func (s *SwarmService) RemoveSwarmNode(ctx context.Context, id string, force bool) error {
	if err := s.validateClient(); err != nil {
		return err
	}

	if err := utils.ValidateID(id, "node ID"); err != nil {
		return err
	}

	options := swarm.NodeRemoveOptions{
		Force: force,
	}

	if err := s.cli.NodeRemove(ctx, id, options); err != nil {
		return fmt.Errorf("failed to remove swarm node %s: %w", id, err)
	}

	return nil
}

// Helper methods

func (s *SwarmService) validateClient() error {
	if s.cli == nil {
		return errors.New("docker client is not initialized")
	}
	return nil
}

func (s *SwarmService) validateServiceLogsRequest(id string) error {
	if err := s.validateClient(); err != nil {
		return err
	}
	return utils.ValidateID(id, "service ID")
}

func (s *SwarmService) getServiceLogsResponse(
	ctx context.Context,
	id string,
) (io.ReadCloser, error) {
	options := container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Timestamps: true,
		Follow:     false,
	}

	response, err := s.cli.ServiceLogs(ctx, id, options)
	if err != nil {
		return nil, fmt.Errorf("failed to get logs for swarm service %s: %w", id, err)
	}
	return response, nil
}

func (s *SwarmService) closeServiceLogsResponse(response io.ReadCloser) {
	if err := response.Close(); err != nil {
		s.log.Warn("Failed to close response", "error", err)
	}
}

func (s *SwarmService) readServiceLogs(response io.ReadCloser) string {
	output := &strings.Builder{}
	buffer := make([]byte, 1024)

	for {
		n, err := response.Read(buffer)
		if n > 0 {
			output.Write(buffer[:n])
		}
		if err != nil {
			break
		}
	}

	return output.String()
}
