package services

import "github.com/wikczerski/D5r/internal/docker"

// ServiceFactory creates and manages all services
type ServiceFactory struct {
	ContainerService  ContainerService
	ImageService      ImageService
	VolumeService     VolumeService
	NetworkService    NetworkService
	DockerInfoService DockerInfoService
}

// NewServiceFactory creates a new service factory
func NewServiceFactory(client *docker.Client) *ServiceFactory {
	return &ServiceFactory{
		ContainerService:  NewContainerService(client),
		ImageService:      NewImageService(client),
		VolumeService:     NewVolumeService(client),
		NetworkService:    NewNetworkService(client),
		DockerInfoService: NewDockerInfoService(client),
	}
}
