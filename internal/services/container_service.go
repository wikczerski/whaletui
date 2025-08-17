package services

import (
	"context"
	"time"

	"github.com/wikczerski/D5r/internal/docker"
	"github.com/wikczerski/D5r/internal/models"
)

type containerService struct {
	*BaseService[models.Container]
	operations *CommonOperations
}

// NewContainerService creates a new container service
func NewContainerService(client *docker.Client) ContainerService {
	base := NewBaseService[models.Container](client, "container")
	ops := NewCommonOperations(client)

	base.ListFunc = func(client *docker.Client, ctx context.Context) ([]models.Container, error) {
		dockerContainers, err := client.ListContainers(ctx, true)
		if err != nil {
			return nil, err
		}

		result := make([]models.Container, len(dockerContainers))
		for i := range dockerContainers {
			result[i] = models.Container(dockerContainers[i])
		}
		return result, nil
	}

	base.RemoveFunc = func(client *docker.Client, ctx context.Context, id string, force bool) error {
		return client.RemoveContainer(ctx, id, force)
	}

	base.InspectFunc = func(client *docker.Client, ctx context.Context, id string) (map[string]any, error) {
		return client.InspectContainer(ctx, id)
	}

	return &containerService{
		BaseService: base,
		operations:  ops,
	}
}

func (s *containerService) ListContainers(ctx context.Context) ([]models.Container, error) {
	return s.List(ctx)
}

func (s *containerService) StartContainer(ctx context.Context, id string) error {
	return s.operations.StartContainer(ctx, id)
}

func (s *containerService) StopContainer(ctx context.Context, id string, timeout *time.Duration) error {
	return s.operations.StopContainer(ctx, id, timeout)
}

func (s *containerService) RestartContainer(ctx context.Context, id string, timeout *time.Duration) error {
	return s.operations.RestartContainer(ctx, id, timeout)
}

func (s *containerService) RemoveContainer(ctx context.Context, id string, force bool) error {
	return s.Remove(ctx, id, force)
}

func (s *containerService) GetContainerLogs(ctx context.Context, id string) (string, error) {
	return s.operations.GetContainerLogs(ctx, id)
}

func (s *containerService) InspectContainer(ctx context.Context, id string) (map[string]any, error) {
	return s.Inspect(ctx, id)
}

func (s *containerService) ExecContainer(ctx context.Context, id string, command []string, tty bool) (string, error) {
	return s.operations.ExecContainer(ctx, id, command, tty)
}

func (s *containerService) AttachContainer(ctx context.Context, id string) (any, error) {
	return s.operations.AttachContainer(ctx, id)
}
