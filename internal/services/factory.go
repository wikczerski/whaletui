package services

import "github.com/wikczerski/D5r/internal/docker"

// ServiceFactoryInterface defines the interface for service factory operations
type ServiceFactoryInterface interface {
	GetContainerService() ContainerService
	GetImageService() ImageService
	GetVolumeService() VolumeService
	GetNetworkService() NetworkService
	GetDockerInfoService() DockerInfoService
	GetLogsService() LogsService
}

// ServiceFactory creates and manages all services
type ServiceFactory struct {
	ContainerService  ContainerService
	ImageService      ImageService
	VolumeService     VolumeService
	NetworkService    NetworkService
	DockerInfoService DockerInfoService
	LogsService       LogsService
}

// NewServiceFactory creates a new service factory
func NewServiceFactory(client *docker.Client) *ServiceFactory {
	return &ServiceFactory{
		ContainerService:  NewContainerService(client),
		ImageService:      NewImageService(client),
		VolumeService:     NewVolumeService(client),
		NetworkService:    NewNetworkService(client),
		DockerInfoService: NewDockerInfoService(client),
		LogsService:       NewLogsService(),
	}
}

// GetContainerService returns the container service
func (sf *ServiceFactory) GetContainerService() ContainerService {
	return sf.ContainerService
}

// GetImageService returns the image service
func (sf *ServiceFactory) GetImageService() ImageService {
	return sf.ImageService
}

// GetVolumeService returns the volume service
func (sf *ServiceFactory) GetVolumeService() VolumeService {
	return sf.VolumeService
}

// GetNetworkService returns the network service
func (sf *ServiceFactory) GetNetworkService() NetworkService {
	return sf.NetworkService
}

// GetDockerInfoService returns the docker info service
func (sf *ServiceFactory) GetDockerInfoService() DockerInfoService {
	return sf.DockerInfoService
}

// GetLogsService returns the logs service
func (sf *ServiceFactory) GetLogsService() LogsService {
	return sf.LogsService
}
