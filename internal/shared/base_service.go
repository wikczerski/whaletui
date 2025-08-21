package shared

import (
	"context"
	"fmt"

	"github.com/wikczerski/whaletui/internal/docker"
)

// BaseService provides common functionality for all Docker resource services
type BaseService[T any] struct {
	client       *docker.Client
	resourceName string

	// Callbacks for specific behavior
	ListFunc    func(*docker.Client, context.Context) ([]T, error)
	RemoveFunc  func(*docker.Client, context.Context, string, bool) error
	InspectFunc func(*docker.Client, context.Context, string) (map[string]any, error)
}

// NewBaseService creates a new base service with common functionality
func NewBaseService[T any](client *docker.Client, resourceName string) *BaseService[T] {
	return &BaseService[T]{
		client:       client,
		resourceName: resourceName,
	}
}

// List retrieves all resources using the provided list function
func (bs *BaseService[T]) List(ctx context.Context) ([]T, error) {
	if bs.client == nil {
		return nil, fmt.Errorf("docker client is not initialized")
	}

	if bs.ListFunc == nil {
		return nil, fmt.Errorf("list function not implemented for %s", bs.resourceName)
	}

	items, err := bs.ListFunc(bs.client, ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list %s: %w", bs.resourceName, err)
	}

	return items, nil
}

// Remove removes a resource using the provided remove function
func (bs *BaseService[T]) Remove(ctx context.Context, id string, force bool) error {
	if bs.client == nil {
		return fmt.Errorf("docker client is not initialized")
	}

	if bs.RemoveFunc == nil {
		return fmt.Errorf("%s removal not yet implemented", bs.resourceName)
	}

	if err := bs.RemoveFunc(bs.client, ctx, id, force); err != nil {
		return fmt.Errorf("failed to remove %s %s: %w", bs.resourceName, id, err)
	}

	return nil
}

// Inspect inspects a resource using the provided inspect function
func (bs *BaseService[T]) Inspect(ctx context.Context, id string) (map[string]any, error) {
	if bs.client == nil {
		return nil, fmt.Errorf("docker client is not initialized")
	}

	if bs.InspectFunc == nil {
		return nil, fmt.Errorf("%s inspection not yet implemented", bs.resourceName)
	}

	result, err := bs.InspectFunc(bs.client, ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to inspect %s %s: %w", bs.resourceName, id, err)
	}

	return result, nil
}

// ValidateClient checks if the Docker client is initialized
func (bs *BaseService[T]) ValidateClient() error {
	if bs.client == nil {
		return fmt.Errorf("docker client is not initialized")
	}
	return nil
}
