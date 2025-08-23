package interfaces

import (
	"context"
	"time"

	"github.com/wikczerski/whaletui/internal/shared"
)

// ContainerService defines the interface for container business operations
type ContainerService interface {
	ListContainers(ctx context.Context) ([]shared.Container, error)
	StartContainer(ctx context.Context, id string) error
	StopContainer(ctx context.Context, id string, timeout *time.Duration) error
	RestartContainer(ctx context.Context, id string, timeout *time.Duration) error
	RemoveContainer(ctx context.Context, id string, force bool) error
	GetContainerLogs(ctx context.Context, id string) (string, error)
	InspectContainer(ctx context.Context, id string) (map[string]any, error)
	ExecContainer(ctx context.Context, id string, command []string, tty bool) (string, error)
	AttachContainer(ctx context.Context, id string) (any, error)
	GetActions() map[rune]string
	GetActionsString() string
	GetNavigation() map[rune]string
	GetNavigationString() string
}

// ImageService defines the interface for image business operations
type ImageService interface {
	ListImages(ctx context.Context) ([]shared.Image, error)
	RemoveImage(ctx context.Context, id string, force bool) error
	InspectImage(ctx context.Context, id string) (map[string]any, error)
	GetActions() map[rune]string
	GetActionsString() string
}

// VolumeService defines the interface for volume business operations
type VolumeService interface {
	ListVolumes(ctx context.Context) ([]shared.Volume, error)
	RemoveVolume(ctx context.Context, name string, force bool) error
	InspectVolume(ctx context.Context, name string) (map[string]any, error)
	GetActions() map[rune]string
	GetActionsString() string
}

// NetworkService defines the interface for network business operations
type NetworkService interface {
	ListNetworks(ctx context.Context) ([]shared.Network, error)
	RemoveNetwork(ctx context.Context, id string) error
	InspectNetwork(ctx context.Context, id string) (map[string]any, error)
	GetActions() map[rune]string
	GetActionsString() string
}

// DockerInfoService defines the interface for Docker system information
type DockerInfoService interface {
	GetDockerInfo(ctx context.Context) (*shared.DockerInfo, error)
}

// LogsService defines the interface for logs operations
type LogsService interface {
	GetLogs(ctx context.Context, resourceType, resourceID string) (string, error)
	GetActions() map[rune]string
	GetActionsString() string
}

// SwarmServiceService defines the interface for swarm service operations
type SwarmServiceService interface {
	ListServices(ctx context.Context) ([]shared.SwarmService, error)
	InspectService(ctx context.Context, id string) (map[string]any, error)
	ScaleService(ctx context.Context, id string, replicas uint64) error
	UpdateService(ctx context.Context, id string, spec any) error
	RemoveService(ctx context.Context, id string) error
	GetServiceLogs(ctx context.Context, id string) (string, error)
	GetActions() map[rune]string
	GetActionsString() string
}

// SwarmNodeService defines the interface for swarm node operations
type SwarmNodeService interface {
	ListNodes(ctx context.Context) ([]shared.SwarmNode, error)
	InspectNode(ctx context.Context, id string) (map[string]any, error)
	UpdateNodeAvailability(ctx context.Context, id string, availability string) error
	RemoveNode(ctx context.Context, id string, force bool) error
	GetActions() map[rune]string
	GetActionsString() string
}
