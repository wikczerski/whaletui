package swarm

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/docker/docker/api/types/swarm"
	"github.com/wikczerski/whaletui/internal/docker"
	"github.com/wikczerski/whaletui/internal/logger"
	"github.com/wikczerski/whaletui/internal/shared"
	"github.com/wikczerski/whaletui/internal/shared/interfaces"
)

// ServiceService implements the SwarmServiceService interface
type ServiceService struct {
	client *docker.Client
	log    *slog.Logger
}

// NewServiceService creates a new swarm service service
func NewServiceService(client *docker.Client) interfaces.SwarmServiceService {
	return &ServiceService{
		client: client,
		log:    logger.GetLogger(),
	}
}

// ListServices lists all swarm services
func (s *ServiceService) ListServices(ctx context.Context) ([]shared.SwarmService, error) {
	dockerServices, err := s.client.ListSwarmServices(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list swarm services: %w", err)
	}

	services := make([]shared.SwarmService, len(dockerServices))
	for i := range dockerServices {
		services[i] = s.convertToSharedService(dockerServices[i])
	}

	return services, nil
}

// InspectService inspects a swarm service
func (s *ServiceService) InspectService(ctx context.Context, id string) (map[string]any, error) {
	service, err := s.client.InspectSwarmService(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to inspect swarm service: %w", err)
	}

	// Convert to map for inspection
	result := map[string]any{
		"ID":          service.ID,
		"Name":        service.Spec.Name,
		"Image":       service.Spec.TaskTemplate.ContainerSpec.Image,
		"Mode":        getServiceMode(service.Spec.Mode),
		"Replicas":    getServiceReplicas(service.Spec.Mode),
		"Ports":       getServicePorts(service.Spec.EndpointSpec),
		"CreatedAt":   service.CreatedAt,
		"UpdatedAt":   service.UpdatedAt,
		"Status":      getServiceStatus(service.UpdateStatus),
		"Labels":      service.Spec.Labels,
		"Env":         service.Spec.TaskTemplate.ContainerSpec.Env,
		"Command":     service.Spec.TaskTemplate.ContainerSpec.Command,
		"Args":        service.Spec.TaskTemplate.ContainerSpec.Args,
		"Constraints": service.Spec.TaskTemplate.Placement.Constraints,
	}

	return result, nil
}

// ScaleService scales a swarm service
func (s *ServiceService) ScaleService(ctx context.Context, id string, replicas uint64) error {
	service, err := s.client.InspectSwarmService(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to inspect service for scaling: %w", err)
	}

	if err := s.validateServiceForScaling(service); err != nil {
		return err
	}

	s.updateServiceReplicas(service, replicas)

	if err := s.client.UpdateSwarmService(ctx, id, service.Version, service.Spec); err != nil {
		return fmt.Errorf("failed to scale service: %w", err)
	}

	s.log.Info("Service scaled successfully", "service_id", id, "replicas", replicas)
	return nil
}

// UpdateService updates a swarm service
func (s *ServiceService) UpdateService(_ context.Context, _ string, _ any) error {
	// This is a placeholder - in a real implementation, you would need to
	// properly handle the spec conversion and validation
	return errors.New("update service not implemented yet")
}

// RemoveService removes a swarm service
func (s *ServiceService) RemoveService(ctx context.Context, id string) error {
	err := s.client.RemoveSwarmService(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to remove swarm service: %w", err)
	}

	s.log.Info("Service removed successfully", "service_id", id)
	return nil
}

// GetServiceLogs gets logs for a swarm service
func (s *ServiceService) GetServiceLogs(ctx context.Context, id string) (string, error) {
	logs, err := s.client.GetSwarmServiceLogs(ctx, id)
	if err != nil {
		return "", fmt.Errorf("failed to get service logs: %w", err)
	}

	return logs, nil
}

// GetActions returns the available actions for swarm services
func (s *ServiceService) GetActions() map[rune]string {
	return map[rune]string{
		'i': "Inspect",
		's': "Scale",
		'r': "Remove",
		'l': "Logs",
	}
}

// GetActionsString returns the available actions as a string
func (s *ServiceService) GetActionsString() string {
	return "<i>: Inspect\n<s>: Scale\n<r>: Remove\n<l>: Logs"
}

// GetNavigation returns the available navigation options for swarm services as a map
func (s *ServiceService) GetNavigation() map[rune]string {
	return map[rune]string{
		'↑': "Up",
		'↓': "Down",
		':': "Command",
		'/': "Filter",
	}
}

// GetNavigationString returns the available navigation options for swarm services as a formatted string
func (s *ServiceService) GetNavigationString() string {
	return "↑/↓: Navigate\n<:> Command mode\n/: Filter"
}

// convertToSharedService converts a Docker swarm service to shared service
func (s *ServiceService) convertToSharedService(dockerService any) shared.SwarmService {
	service, ok := dockerService.(swarm.Service)
	if !ok {
		s.log.Error("Failed to convert docker service to swarm.Service type")
		return shared.SwarmService{}
	}

	return shared.SwarmService{
		ID:        service.ID,
		Name:      service.Spec.Name,
		Image:     service.Spec.TaskTemplate.ContainerSpec.Image,
		Mode:      getServiceMode(service.Spec.Mode),
		Replicas:  getServiceReplicas(service.Spec.Mode),
		Ports:     getServicePorts(service.Spec.EndpointSpec),
		CreatedAt: service.CreatedAt,
		UpdatedAt: service.UpdatedAt,
		Status:    getServiceStatus(service.UpdateStatus),
	}
}

// validateServiceForScaling validates that the service can be scaled
func (s *ServiceService) validateServiceForScaling(service swarm.Service) error {
	if service.Spec.Mode.Replicated == nil {
		return errors.New("cannot scale global service")
	}
	return nil
}

// updateServiceReplicas updates the service replicas count
func (s *ServiceService) updateServiceReplicas(service swarm.Service, replicas uint64) {
	service.Spec.Mode.Replicated.Replicas = &replicas
}

// Helper functions for service inspection
func getServiceMode(mode any) string {
	if mode == nil {
		return "unknown"
	}

	switch mode.(type) {
	case *swarm.ReplicatedService:
		return "replicated"
	case *swarm.GlobalService:
		return "global"
	default:
		return "unknown"
	}
}

func getServiceReplicas(mode any) string {
	if mode == nil {
		return "0/0"
	}

	switch m := mode.(type) {
	case *swarm.ReplicatedService:
		if m.Replicas != nil {
			return fmt.Sprintf("%d", *m.Replicas)
		}
		return "0"
	case *swarm.GlobalService:
		return "global"
	default:
		return "0/0"
	}
}

func getServicePorts(endpointSpec *swarm.EndpointSpec) []string {
	if endpointSpec == nil || endpointSpec.Ports == nil {
		return []string{}
	}

	ports := make([]string, len(endpointSpec.Ports))
	for i, port := range endpointSpec.Ports {
		if port.PublishedPort != 0 {
			ports[i] = fmt.Sprintf("%d:%d/%s", port.PublishedPort, port.TargetPort, port.Protocol)
		} else {
			ports[i] = fmt.Sprintf("%d/%s", port.TargetPort, port.Protocol)
		}
	}

	return ports
}

func getServiceStatus(updateStatus *swarm.UpdateStatus) string {
	if updateStatus == nil {
		return "running"
	}

	// Only transform states that need it (most states map to themselves)
	if updateStatus.State == swarm.UpdateStateUpdating {
		return "updating"
	}

	// For all other states, return the state as-is (they're already strings)
	return string(updateStatus.State)
}
