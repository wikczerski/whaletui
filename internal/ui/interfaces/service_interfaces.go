//nolint:max-public-structs
package interfaces

import (
	"context"
	"time"

	"github.com/wikczerski/whaletui/internal/shared"
)

// ServiceWithActions defines the minimal interface for services that provide actions
type ServiceWithActions interface {
	GetActions() map[rune]string
	GetActionsString() string
}

// ServiceWithNavigation defines the minimal interface for services that provide navigation
type ServiceWithNavigation interface {
	GetNavigation() map[rune]string
	GetNavigationString() string
}

// ServiceFactoryInterface defines the interface for service factory operations
type ServiceFactoryInterface interface {
	GetContainerService() ContainerService
	GetImageService() ImageService
	GetVolumeService() VolumeService
	GetNetworkService() NetworkService
	GetDockerInfoService() DockerInfoService
	GetLogsService() LogsService
	GetSwarmServiceService() any
	GetSwarmNodeService() any
	GetCurrentService() any
	SetCurrentService(serviceName string)
	IsServiceAvailable(serviceName string) bool
	IsContainerServiceAvailable() bool
}

// ContainerService defines the interface for container operations
type ContainerService interface {
	ListContainers(ctx context.Context) ([]shared.Container, error)
	StartContainer(ctx context.Context, id string) error
	StopContainer(ctx context.Context, id string, timeout *time.Duration) error
	RestartContainer(ctx context.Context, id string, timeout *time.Duration) error
	RemoveContainer(ctx context.Context, id string, force bool) error
	InspectContainer(ctx context.Context, id string) (map[string]any, error)
	ExecContainer(ctx context.Context, id string, command []string, tty bool) (string, error)
}

// ImageService defines the interface for image operations
type ImageService interface {
	ListImages(ctx context.Context) ([]shared.Image, error)
	RemoveImage(ctx context.Context, id string, force bool) error
	InspectImage(ctx context.Context, id string) (map[string]any, error)
	GetActions() map[rune]string
	GetActionsString() string
}

// VolumeService defines the interface for volume operations
type VolumeService interface {
	ListVolumes(ctx context.Context) ([]shared.Volume, error)
	RemoveVolume(ctx context.Context, id string, force bool) error
	InspectVolume(ctx context.Context, id string) (map[string]any, error)
	GetActions() map[rune]string
	GetActionsString() string
}

// NetworkService defines the interface for network operations
type NetworkService interface {
	ListNetworks(ctx context.Context) ([]shared.Network, error)
	RemoveNetwork(ctx context.Context, id string) error
	InspectNetwork(ctx context.Context, id string) (map[string]any, error)
	GetActions() map[rune]string
	GetActionsString() string
}

// DockerInfoService defines the interface for Docker info operations
type DockerInfoService interface {
	GetDockerInfo(ctx context.Context) (*DockerInfo, error)
	GetActions() map[rune]string
	GetActionsString() string
}

// LogsService defines the interface for logs operations
type LogsService interface {
	GetLogs(ctx context.Context, resourceType, resourceID string) (string, error)
	GetActions() map[rune]string
	GetActionsString() string
}

// Domain types are now imported from internal/docker package

// DockerInfo represents Docker system information
type DockerInfo interface {
	GetVersion() string
	GetOperatingSystem() string
	GetLoggingDriver() string
	GetContainers() int
	GetImages() int
	GetVolumes() int
	GetNetworks() int
}
