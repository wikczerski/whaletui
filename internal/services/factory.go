package services

import "github.com/wikczerski/whaletui/internal/docker"

// ServiceFactoryInterface defines the interface for service factory operations
type ServiceFactoryInterface interface {
	GetContainerService() ContainerService
	GetImageService() ImageService
	GetVolumeService() VolumeService
	GetNetworkService() NetworkService
	GetDockerInfoService() DockerInfoService
	GetLogsService() LogsService
	IsServiceAvailable(serviceName string) bool
	IsContainerServiceAvailable() bool
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
	if client == nil {
		return &ServiceFactory{
			ContainerService:  nil,
			ImageService:      nil,
			VolumeService:     nil,
			NetworkService:    nil,
			DockerInfoService: nil,
			LogsService:       nil,
		}
	}

	containerService := NewContainerService(client)

	return &ServiceFactory{
		ContainerService:  containerService,
		ImageService:      NewImageService(client),
		VolumeService:     NewVolumeService(client),
		NetworkService:    NewNetworkService(client),
		DockerInfoService: NewDockerInfoService(client),
		LogsService:       NewLogsService(containerService),
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

// IsServiceAvailable checks if a specific service is available
func (sf *ServiceFactory) IsServiceAvailable(serviceName string) bool {
	if sf == nil {
		return false
	}

	switch serviceName {
	case "container":
		return sf.ContainerService != nil
	case "image":
		return sf.ImageService != nil
	case "volume":
		return sf.VolumeService != nil
	case "network":
		return sf.NetworkService != nil
	case "dockerInfo":
		return sf.DockerInfoService != nil
	case "logs":
		return sf.LogsService != nil
	default:
		return false
	}
}

// IsContainerServiceAvailable checks if the container service is available
func (sf *ServiceFactory) IsContainerServiceAvailable() bool {
	return sf != nil && sf.ContainerService != nil
}
