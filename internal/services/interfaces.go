package services

import (
	"context"
	"time"

	"github.com/wikczerski/whaletui/internal/models"
)

// ContainerService defines the interface for container business operations
type ContainerService interface {
	ListContainers(ctx context.Context) ([]models.Container, error)
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
}

// ImageService defines the interface for image business operations
type ImageService interface {
	ListImages(ctx context.Context) ([]models.Image, error)
	RemoveImage(ctx context.Context, id string, force bool) error
	InspectImage(ctx context.Context, id string) (map[string]any, error)
	GetActions() map[rune]string
	GetActionsString() string
}

// VolumeService defines the interface for volume business operations
type VolumeService interface {
	ListVolumes(ctx context.Context) ([]models.Volume, error)
	RemoveVolume(ctx context.Context, name string, force bool) error
	InspectVolume(ctx context.Context, name string) (map[string]any, error)
	GetActions() map[rune]string
	GetActionsString() string
}

// NetworkService defines the interface for network business operations
type NetworkService interface {
	ListNetworks(ctx context.Context) ([]models.Network, error)
	RemoveNetwork(ctx context.Context, id string) error
	InspectNetwork(ctx context.Context, id string) (map[string]any, error)
	GetActions() map[rune]string
	GetActionsString() string
}

// DockerInfoService defines the interface for Docker system information
type DockerInfoService interface {
	GetDockerInfo(ctx context.Context) (*models.DockerInfo, error)
}

// LogsService defines the interface for logs operations
type LogsService interface {
	GetActions() map[rune]string
	GetActionsString() string
}
