package services

import (
	"context"
	"fmt"
	"time"

	"github.com/wikczerski/D5r/internal/docker"
)

// CommonOperations provides reusable Docker operations
type CommonOperations struct {
	client *docker.Client
}

// NewCommonOperations creates a new common operations helper
func NewCommonOperations(client *docker.Client) *CommonOperations {
	return &CommonOperations{client: client}
}

// StartContainer is a reusable start operation
func (co *CommonOperations) StartContainer(ctx context.Context, id string) error {
	if co.client == nil {
		return fmt.Errorf("docker client is not initialized")
	}
	return co.client.StartContainer(ctx, id)
}

// StopContainer is a reusable stop operation
func (co *CommonOperations) StopContainer(ctx context.Context, id string, timeout *time.Duration) error {
	if co.client == nil {
		return fmt.Errorf("docker client is not initialized")
	}
	return co.client.StopContainer(ctx, id, timeout)
}

// RestartContainer is a reusable restart operation
func (co *CommonOperations) RestartContainer(ctx context.Context, id string, timeout *time.Duration) error {
	if co.client == nil {
		return fmt.Errorf("docker client is not initialized")
	}
	return co.client.RestartContainer(ctx, id, timeout)
}

// RemoveContainer is a reusable remove operation
func (co *CommonOperations) RemoveContainer(ctx context.Context, id string, force bool) error {
	if co.client == nil {
		return fmt.Errorf("docker client is not initialized")
	}
	return co.client.RemoveContainer(ctx, id, force)
}

// GetContainerLogs is a reusable logs operation
func (co *CommonOperations) GetContainerLogs(ctx context.Context, id string) (string, error) {
	if co.client == nil {
		return "", fmt.Errorf("docker client is not initialized")
	}
	return co.client.GetContainerLogs(ctx, id)
}

// InspectResource provides a generic inspect operation
func (co *CommonOperations) InspectResource(ctx context.Context, resourceType, id string) (map[string]any, error) {
	if co.client == nil {
		return nil, fmt.Errorf("docker client is not initialized")
	}

	switch resourceType {
	case "container":
		return co.client.InspectContainer(ctx, id)
	case "image":
		return co.client.InspectImage(ctx, id)
	case "volume":
		return co.client.InspectVolume(ctx, id)
	case "network":
		return co.client.InspectNetwork(ctx, id)
	default:
		return nil, fmt.Errorf("unsupported resource type: %s", resourceType)
	}
}
