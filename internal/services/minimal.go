package services

import (
	"context"

	"github.com/wikczerski/D5r/internal/models"
)

// MinimalDockerService provides a minimal interface for basic Docker operations
type MinimalDockerService interface {
	ListContainers(ctx context.Context) ([]models.Container, error)
	ListImages(ctx context.Context) ([]models.Image, error)
	ListVolumes(ctx context.Context) ([]models.Volume, error)
	ListNetworks(ctx context.Context) ([]models.Network, error)
}

// Ensure ServiceFactory implements MinimalDockerService
var _ MinimalDockerService = (*ServiceFactory)(nil)

func (sf *ServiceFactory) ListContainers(ctx context.Context) ([]models.Container, error) {
	return sf.ContainerService.ListContainers(ctx)
}

func (sf *ServiceFactory) ListImages(ctx context.Context) ([]models.Image, error) {
	return sf.ImageService.ListImages(ctx)
}

func (sf *ServiceFactory) ListVolumes(ctx context.Context) ([]models.Volume, error) {
	return sf.VolumeService.ListVolumes(ctx)
}

func (sf *ServiceFactory) ListNetworks(ctx context.Context) ([]models.Network, error) {
	return sf.NetworkService.ListNetworks(ctx)
}
