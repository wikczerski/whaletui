package swarm

import (
	"context"
	"fmt"
	"log/slog"
	"time"

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

	// Update the service spec with new replicas
	if service.Spec.Mode.Replicated != nil {
		service.Spec.Mode.Replicated.Replicas = &replicas
	} else {
		return fmt.Errorf("cannot scale global service")
	}

	err = s.client.UpdateSwarmService(ctx, id, service.Version, service.Spec)
	if err != nil {
		return fmt.Errorf("failed to scale service: %w", err)
	}

	s.log.Info("Service scaled successfully", "service_id", id, "replicas", replicas)
	return nil
}

// UpdateService updates a swarm service
func (s *ServiceService) UpdateService(_ context.Context, _ string, _ any) error {
	// This is a placeholder - in a real implementation, you would need to
	// properly handle the spec conversion and validation
	return fmt.Errorf("update service not implemented yet")
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
	return "i: Inspect, s: Scale, r: Remove, l: Logs"
}

// convertToSharedService converts a Docker swarm service to shared service
func (s *ServiceService) convertToSharedService(_ interface{}) shared.SwarmService {
	// This is a placeholder - you would need to properly convert the Docker service
	// to the shared service type based on the actual Docker API types
	return shared.SwarmService{
		ID:        "placeholder",
		Name:      "placeholder",
		Image:     "placeholder",
		Mode:      "placeholder",
		Replicas:  "placeholder",
		Ports:     []string{},
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
		Status:    "placeholder",
	}
}

// Helper functions for service inspection
func getServiceMode(_ interface{}) string {
	// Placeholder implementation
	return "replicated"
}

func getServiceReplicas(_ interface{}) string {
	// Placeholder implementation
	return "1/1"
}

func getServicePorts(_ interface{}) []string {
	// Placeholder implementation
	return []string{}
}

func getServiceStatus(_ interface{}) string {
	// Placeholder implementation
	return "running"
}
