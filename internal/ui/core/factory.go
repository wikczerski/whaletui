package core

import (
	"context"
	"time"

	"github.com/wikczerski/whaletui/internal/docker"
	"github.com/wikczerski/whaletui/internal/domains/container"
	"github.com/wikczerski/whaletui/internal/domains/image"
	"github.com/wikczerski/whaletui/internal/domains/logs"
	"github.com/wikczerski/whaletui/internal/domains/network"
	"github.com/wikczerski/whaletui/internal/domains/volume"
	"github.com/wikczerski/whaletui/internal/shared"
	sharedInterfaces "github.com/wikczerski/whaletui/internal/shared/interfaces"
	"github.com/wikczerski/whaletui/internal/ui/interfaces"
)

// ServiceFactoryInterface defines the interface for service factory operations
type ServiceFactoryInterface interface {
	GetContainerService() interfaces.ContainerService
	GetImageService() interfaces.ImageService
	GetVolumeService() interfaces.VolumeService
	GetNetworkService() interfaces.NetworkService
	GetDockerInfoService() interfaces.DockerInfoService
	GetLogsService() interfaces.LogsService
	IsServiceAvailable(serviceName string) bool
	IsContainerServiceAvailable() bool
}

// ServiceFactory creates and manages all services
type ServiceFactory struct {
	ContainerService  interfaces.ContainerService
	ImageService      interfaces.ImageService
	VolumeService     interfaces.VolumeService
	NetworkService    interfaces.NetworkService
	DockerInfoService interfaces.DockerInfoService
	LogsService       interfaces.LogsService
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

	containerService := container.NewContainerService(client)
	sharedDockerInfoService := shared.NewDockerInfoService(client)

	// Create adapter services that bridge between shared interfaces and UI interfaces
	uiContainerService := &containerServiceAdapter{service: containerService}
	uiImageService := &imageServiceAdapter{service: image.NewImageService(client)}
	uiVolumeService := &volumeServiceAdapter{service: volume.NewVolumeService(client)}
	uiNetworkService := &networkServiceAdapter{service: network.NewNetworkService(client)}

	// Create a wrapper that converts between shared.DockerInfo and interfaces.DockerInfo
	dockerInfoService := &dockerInfoServiceWrapper{service: sharedDockerInfoService}

	return &ServiceFactory{
		ContainerService:  uiContainerService,
		ImageService:      uiImageService,
		VolumeService:     uiVolumeService,
		NetworkService:    uiNetworkService,
		DockerInfoService: dockerInfoService,
		LogsService:       logs.NewLogsService(containerService),
	}
}

// Adapter services that bridge between shared interfaces and UI interfaces

// containerServiceAdapter adapts shared.ContainerService to interfaces.ContainerService
type containerServiceAdapter struct {
	service sharedInterfaces.ContainerService
}

func (a *containerServiceAdapter) ListContainers(ctx context.Context) ([]interfaces.Container, error) {
	containers, err := a.service.ListContainers(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]interfaces.Container, len(containers))
	for i, container := range containers {
		result[i] = container
	}
	return result, nil
}

func (a *containerServiceAdapter) StartContainer(ctx context.Context, id string) error {
	return a.service.StartContainer(ctx, id)
}

func (a *containerServiceAdapter) StopContainer(ctx context.Context, id string, timeout *time.Duration) error {
	return a.service.StopContainer(ctx, id, timeout)
}

func (a *containerServiceAdapter) RestartContainer(ctx context.Context, id string, timeout *time.Duration) error {
	return a.service.RestartContainer(ctx, id, timeout)
}

func (a *containerServiceAdapter) RemoveContainer(ctx context.Context, id string, force bool) error {
	return a.service.RemoveContainer(ctx, id, force)
}

func (a *containerServiceAdapter) GetContainerLogs(ctx context.Context, id string) (string, error) {
	return a.service.GetContainerLogs(ctx, id)
}

func (a *containerServiceAdapter) InspectContainer(ctx context.Context, id string) (map[string]any, error) {
	return a.service.InspectContainer(ctx, id)
}

func (a *containerServiceAdapter) ExecContainer(ctx context.Context, id string, command []string, tty bool) (string, error) {
	return a.service.ExecContainer(ctx, id, command, tty)
}

func (a *containerServiceAdapter) AttachContainer(ctx context.Context, id string) (any, error) {
	return a.service.AttachContainer(ctx, id)
}

func (a *containerServiceAdapter) GetActions() map[rune]string {
	return a.service.GetActions()
}

func (a *containerServiceAdapter) GetActionsString() string {
	return a.service.GetActionsString()
}

// imageServiceAdapter adapts shared.ImageService to interfaces.ImageService
type imageServiceAdapter struct {
	service sharedInterfaces.ImageService
}

func (a *imageServiceAdapter) ListImages(ctx context.Context) ([]interfaces.Image, error) {
	images, err := a.service.ListImages(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]interfaces.Image, len(images))
	for i, image := range images {
		result[i] = image
	}
	return result, nil
}

func (a *imageServiceAdapter) RemoveImage(ctx context.Context, id string, force bool) error {
	return a.service.RemoveImage(ctx, id, force)
}

func (a *imageServiceAdapter) InspectImage(ctx context.Context, id string) (map[string]any, error) {
	return a.service.InspectImage(ctx, id)
}

func (a *imageServiceAdapter) GetActions() map[rune]string {
	return a.service.GetActions()
}

func (a *imageServiceAdapter) GetActionsString() string {
	return a.service.GetActionsString()
}

// volumeServiceAdapter adapts shared.VolumeService to interfaces.VolumeService
type volumeServiceAdapter struct {
	service sharedInterfaces.VolumeService
}

func (a *volumeServiceAdapter) ListVolumes(ctx context.Context) ([]interfaces.Volume, error) {
	volumes, err := a.service.ListVolumes(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]interfaces.Volume, len(volumes))
	for i, volume := range volumes {
		result[i] = volume
	}
	return result, nil
}

func (a *volumeServiceAdapter) RemoveVolume(ctx context.Context, name string, force bool) error {
	return a.service.RemoveVolume(ctx, name, force)
}

func (a *volumeServiceAdapter) InspectVolume(ctx context.Context, name string) (map[string]any, error) {
	return a.service.InspectVolume(ctx, name)
}

func (a *volumeServiceAdapter) GetActions() map[rune]string {
	return a.service.GetActions()
}

func (a *volumeServiceAdapter) GetActionsString() string {
	return a.service.GetActionsString()
}

// networkServiceAdapter adapts shared.NetworkService to interfaces.NetworkService
type networkServiceAdapter struct {
	service sharedInterfaces.NetworkService
}

func (a *networkServiceAdapter) ListNetworks(ctx context.Context) ([]interfaces.Network, error) {
	networks, err := a.service.ListNetworks(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]interfaces.Network, len(networks))
	for i, network := range networks {
		result[i] = network
	}
	return result, nil
}

func (a *networkServiceAdapter) RemoveNetwork(ctx context.Context, id string) error {
	return a.service.RemoveNetwork(ctx, id)
}

func (a *networkServiceAdapter) InspectNetwork(ctx context.Context, id string) (map[string]any, error) {
	return a.service.InspectNetwork(ctx, id)
}

func (a *networkServiceAdapter) GetActions() map[rune]string {
	return a.service.GetActions()
}

func (a *networkServiceAdapter) GetActionsString() string {
	return a.service.GetActionsString()
}

// dockerInfoServiceWrapper wraps shared.DockerInfoService to implement interfaces.DockerInfoService
type dockerInfoServiceWrapper struct {
	service shared.DockerInfoService
}

func (w *dockerInfoServiceWrapper) GetDockerInfo(ctx context.Context) (interfaces.DockerInfo, error) {
	sharedInfo, err := w.service.GetDockerInfo(ctx)
	if err != nil {
		return interfaces.DockerInfo{}, err
	}

	// Convert shared.DockerInfo to interfaces.DockerInfo
	info := interfaces.DockerInfo{
		Version:         sharedInfo.Version,
		Containers:      sharedInfo.Containers,
		Images:          sharedInfo.Images,
		Volumes:         sharedInfo.Volumes,
		Networks:        sharedInfo.Networks,
		OperatingSystem: sharedInfo.OperatingSystem,
		Architecture:    sharedInfo.Architecture,
		Driver:          sharedInfo.Driver,
		LoggingDriver:   sharedInfo.LoggingDriver,
	}

	return info, nil
}

// GetContainerService returns the container service
func (sf *ServiceFactory) GetContainerService() interfaces.ContainerService {
	return sf.ContainerService
}

// GetImageService returns the image service
func (sf *ServiceFactory) GetImageService() interfaces.ImageService {
	return sf.ImageService
}

// GetVolumeService returns the volume service
func (sf *ServiceFactory) GetVolumeService() interfaces.VolumeService {
	return sf.VolumeService
}

// GetNetworkService returns the network service
func (sf *ServiceFactory) GetNetworkService() interfaces.NetworkService {
	return sf.NetworkService
}

// GetDockerInfoService returns the docker info service
func (sf *ServiceFactory) GetDockerInfoService() interfaces.DockerInfoService {
	return sf.DockerInfoService
}

// GetLogsService returns the logs service
func (sf *ServiceFactory) GetLogsService() interfaces.LogsService {
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
