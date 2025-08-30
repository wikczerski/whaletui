// Package core provides core UI components and functionality for WhaleTUI.
package core

import (
	"context"
	"time"

	"github.com/wikczerski/whaletui/internal/docker"
	"github.com/wikczerski/whaletui/internal/domains/container"
	"github.com/wikczerski/whaletui/internal/domains/image"
	"github.com/wikczerski/whaletui/internal/domains/logs"
	"github.com/wikczerski/whaletui/internal/domains/network"
	"github.com/wikczerski/whaletui/internal/domains/swarm"
	"github.com/wikczerski/whaletui/internal/domains/volume"
	"github.com/wikczerski/whaletui/internal/shared"
	sharedInterfaces "github.com/wikczerski/whaletui/internal/shared/interfaces"
	"github.com/wikczerski/whaletui/internal/ui/constants"
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
	ContainerService    interfaces.ContainerService
	ImageService        interfaces.ImageService
	VolumeService       interfaces.VolumeService
	NetworkService      interfaces.NetworkService
	DockerInfoService   interfaces.DockerInfoService
	LogsService         interfaces.LogsService
	SwarmServiceService sharedInterfaces.SwarmServiceService
	SwarmNodeService    sharedInterfaces.SwarmNodeService
	currentService      string // Track which service is currently active
}

// NewServiceFactory creates a new service factory
func NewServiceFactory(client *docker.Client) *ServiceFactory {
	if client == nil {
		return createEmptyServiceFactory()
	}

	containerService := container.NewContainerService(client)
	sharedDockerInfoService := shared.NewDockerInfoService(client)

	// Create adapter services that bridge between shared interfaces and UI interfaces
	uiServices := createUIServices(client, containerService)
	swarmServices := createSwarmServices(client)

	// Create a wrapper that converts between shared.DockerInfo and interfaces.DockerInfo
	dockerInfoService := &dockerInfoServiceWrapper{service: sharedDockerInfoService}

	return &ServiceFactory{
		ContainerService:    uiServices.container,
		ImageService:        uiServices.image,
		VolumeService:       uiServices.volume,
		NetworkService:      uiServices.network,
		DockerInfoService:   dockerInfoService,
		LogsService:         logs.NewLogsService(containerService),
		SwarmServiceService: swarmServices.service,
		SwarmNodeService:    swarmServices.node,
		currentService:      constants.ViewContainers, // Default to container service
	}
}

// createEmptyServiceFactory creates a service factory with nil services
func createEmptyServiceFactory() *ServiceFactory {
	return &ServiceFactory{
		ContainerService:  nil,
		ImageService:      nil,
		VolumeService:     nil,
		NetworkService:    nil,
		DockerInfoService: nil,
		LogsService:       nil,
		currentService:    constants.ViewContainers, // Default to container service
	}
}

// uiServices holds the UI service adapters
type uiServices struct {
	container *containerServiceAdapter
	image     *imageServiceAdapter
	volume    *volumeServiceAdapter
	network   *networkServiceAdapter
}

// createUIServices creates the UI service adapters
func createUIServices(
	client *docker.Client,
	containerService sharedInterfaces.ContainerService,
) uiServices {
	return uiServices{
		container: &containerServiceAdapter{service: containerService},
		image:     &imageServiceAdapter{service: image.NewImageService(client)},
		volume:    &volumeServiceAdapter{service: volume.NewVolumeService(client)},
		network:   &networkServiceAdapter{service: network.NewNetworkService(client)},
	}
}

// swarmServices holds the swarm services
type swarmServices struct {
	service sharedInterfaces.SwarmServiceService
	node    sharedInterfaces.SwarmNodeService
}

// createSwarmServices creates the swarm services
func createSwarmServices(client *docker.Client) swarmServices {
	return swarmServices{
		service: swarm.NewServiceService(client),
		node:    swarm.NewNodeService(client),
	}
}

// Adapter services that bridge between shared interfaces and UI interfaces

// containerServiceAdapter adapts shared.ContainerService to interfaces.ContainerService
type containerServiceAdapter struct {
	service sharedInterfaces.ContainerService
}

func (a *containerServiceAdapter) ListContainers(ctx context.Context) ([]shared.Container, error) {
	containers, err := a.service.ListContainers(ctx)
	if err != nil {
		return nil, err
	}
	return containers, nil
}

func (a *containerServiceAdapter) StartContainer(ctx context.Context, id string) error {
	return a.service.StartContainer(ctx, id)
}

func (a *containerServiceAdapter) StopContainer(
	ctx context.Context,
	id string,
	timeout *time.Duration,
) error {
	return a.service.StopContainer(ctx, id, timeout)
}

func (a *containerServiceAdapter) RestartContainer(
	ctx context.Context,
	id string,
	timeout *time.Duration,
) error {
	return a.service.RestartContainer(ctx, id, timeout)
}

func (a *containerServiceAdapter) RemoveContainer(
	ctx context.Context,
	id string,
	force bool,
) error {
	return a.service.RemoveContainer(ctx, id, force)
}

func (a *containerServiceAdapter) GetContainerLogs(ctx context.Context, id string) (string, error) {
	return a.service.GetContainerLogs(ctx, id)
}

func (a *containerServiceAdapter) InspectContainer(
	ctx context.Context,
	id string,
) (map[string]any, error) {
	return a.service.InspectContainer(ctx, id)
}

func (a *containerServiceAdapter) ExecContainer(
	ctx context.Context,
	id string,
	command []string,
	tty bool,
) (string, error) {
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

func (a *containerServiceAdapter) GetNavigation() map[rune]string {
	return a.service.GetNavigation()
}

func (a *containerServiceAdapter) GetNavigationString() string {
	return a.service.GetNavigationString()
}

// imageServiceAdapter adapts shared.ImageService to interfaces.ImageService
type imageServiceAdapter struct {
	service sharedInterfaces.ImageService
}

func (a *imageServiceAdapter) ListImages(ctx context.Context) ([]shared.Image, error) {
	images, err := a.service.ListImages(ctx)
	if err != nil {
		return nil, err
	}
	return images, nil
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

func (a *volumeServiceAdapter) ListVolumes(ctx context.Context) ([]shared.Volume, error) {
	volumes, err := a.service.ListVolumes(ctx)
	if err != nil {
		return nil, err
	}
	return volumes, nil
}

func (a *volumeServiceAdapter) RemoveVolume(ctx context.Context, name string, force bool) error {
	return a.service.RemoveVolume(ctx, name, force)
}

func (a *volumeServiceAdapter) InspectVolume(
	ctx context.Context,
	name string,
) (map[string]any, error) {
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

func (a *networkServiceAdapter) ListNetworks(ctx context.Context) ([]shared.Network, error) {
	networks, err := a.service.ListNetworks(ctx)
	if err != nil {
		return nil, err
	}
	return networks, nil
}

func (a *networkServiceAdapter) RemoveNetwork(ctx context.Context, id string) error {
	return a.service.RemoveNetwork(ctx, id)
}

func (a *networkServiceAdapter) InspectNetwork(
	ctx context.Context,
	id string,
) (map[string]any, error) {
	return a.service.InspectNetwork(ctx, id)
}

func (a *networkServiceAdapter) GetActions() map[rune]string {
	return a.service.GetActions()
}

func (a *networkServiceAdapter) GetActionsString() string {
	return a.service.GetActionsString()
}

// dockerInfoImpl implements interfaces.DockerInfo
type dockerInfoImpl struct {
	version         string
	operatingSystem string
	architecture    string
	driver          string
	loggingDriver   string
}

func (d *dockerInfoImpl) GetVersion() string {
	return d.version
}

func (d *dockerInfoImpl) GetOperatingSystem() string {
	return d.operatingSystem
}

func (d *dockerInfoImpl) GetLoggingDriver() string {
	return d.loggingDriver
}

// dockerInfoServiceWrapper wraps shared.DockerInfoService to implement interfaces.DockerInfoService
type dockerInfoServiceWrapper struct {
	service shared.DockerInfoService
}

// nolint:gocritic // Interface design requires pointer return type
func (w *dockerInfoServiceWrapper) GetDockerInfo(
	ctx context.Context,
) (*interfaces.DockerInfo, error) {
	sharedInfo, err := w.service.GetDockerInfo(ctx)
	if err != nil {
		return nil, err
	}

	// Create a concrete implementation of interfaces.DockerInfo
	info := &dockerInfoImpl{
		version:         sharedInfo.Version,
		operatingSystem: sharedInfo.OperatingSystem,
		architecture:    sharedInfo.Architecture,
		driver:          sharedInfo.Driver,
		loggingDriver:   sharedInfo.LoggingDriver,
	}

	// Convert to interface pointer
	var dockerInfo interfaces.DockerInfo = info
	return &dockerInfo, nil
}

func (w *dockerInfoServiceWrapper) GetActions() map[rune]string {
	// Return default actions for docker info
	return map[rune]string{
		'r': "Refresh",
		'h': "Help",
	}
}

func (w *dockerInfoServiceWrapper) GetActionsString() string {
	return "<r>:Refresh\n<h>:Help"
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

// GetSwarmServiceService returns the swarm service service
func (sf *ServiceFactory) GetSwarmServiceService() any {
	return sf.SwarmServiceService
}

// GetSwarmNodeService returns the swarm node service
func (sf *ServiceFactory) GetSwarmNodeService() any {
	return sf.SwarmNodeService
}

// IsServiceAvailable checks if a specific service is available
func (sf *ServiceFactory) IsServiceAvailable(serviceName string) bool {
	if sf == nil {
		return false
	}

	return sf.checkServiceAvailability(serviceName)
}

// IsContainerServiceAvailable checks if the container service is available
func (sf *ServiceFactory) IsContainerServiceAvailable() bool {
	return sf != nil && sf.ContainerService != nil
}

// SetCurrentService sets the currently active service
func (sf *ServiceFactory) SetCurrentService(serviceName string) {
	if sf != nil {
		sf.currentService = serviceName
	}
}

// getServiceByType returns the service based on the service type
func (sf *ServiceFactory) getServiceByType(serviceType string) any {
	switch serviceType {
	case constants.ViewContainers:
		return sf.ContainerService
	case constants.ViewImages:
		return sf.ImageService
	case constants.ViewVolumes:
		return sf.VolumeService
	case constants.ViewNetworks:
		return sf.NetworkService
	case constants.ViewDockerInfo:
		return sf.DockerInfoService
	case constants.ViewLogs:
		return sf.LogsService
	case constants.ViewSwarmServices:
		return sf.SwarmServiceService
	case constants.ViewSwarmNodes:
		return sf.SwarmNodeService
	default:
		return nil
	}
}

// checkServiceAvailability checks if a specific service is available
func (sf *ServiceFactory) checkServiceAvailability(serviceName string) bool {
	return sf.getServiceByType(serviceName) != nil
}
