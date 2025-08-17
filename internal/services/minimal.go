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

// ListContainers lists all containers
func (sf *ServiceFactory) ListContainers(ctx context.Context) ([]models.Container, error) {
	return sf.ContainerService.ListContainers(ctx)
}

// ListImages lists all images
func (sf *ServiceFactory) ListImages(ctx context.Context) ([]models.Image, error) {
	return sf.ImageService.ListImages(ctx)
}

// ListVolumes lists all volumes
func (sf *ServiceFactory) ListVolumes(ctx context.Context) ([]models.Volume, error) {
	return sf.VolumeService.ListVolumes(ctx)
}

// ListNetworks lists all networks
func (sf *ServiceFactory) ListNetworks(ctx context.Context) ([]models.Network, error) {
	return sf.NetworkService.ListNetworks(ctx)
}
