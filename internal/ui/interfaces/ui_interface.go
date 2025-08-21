package interfaces

import (
	"context"
	"time"

	"github.com/rivo/tview"
	"github.com/wikczerski/whaletui/internal/config"
)

// Type aliases to avoid import cycles
type Container = any
type Image = any
type Volume = any
type Network = any

// DockerInfoService defines the minimal interface needed by UI for Docker info
type DockerInfoService interface {
	GetDockerInfo(ctx context.Context) (DockerInfo, error)
}

// DockerInfo represents Docker system information for UI
type DockerInfo struct {
	Version         string
	Containers      int
	Images          int
	Volumes         int
	Networks        int
	OperatingSystem string
	Architecture    string
	Driver          string
	LoggingDriver   string
}

// ContainerService defines the minimal interface for container operations
type ContainerService interface {
	ListContainers(ctx context.Context) ([]Container, error)
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

// ImageService defines the minimal interface for image operations
type ImageService interface {
	ListImages(ctx context.Context) ([]Image, error)
	RemoveImage(ctx context.Context, id string, force bool) error
	InspectImage(ctx context.Context, id string) (map[string]any, error)
	GetActions() map[rune]string
	GetActionsString() string
}

// VolumeService defines the minimal interface for volume operations
type VolumeService interface {
	ListVolumes(ctx context.Context) ([]Volume, error)
	RemoveVolume(ctx context.Context, name string, force bool) error
	InspectVolume(ctx context.Context, name string) (map[string]any, error)
	GetActions() map[rune]string
	GetActionsString() string
}

// NetworkService defines the minimal interface for network operations
type NetworkService interface {
	ListNetworks(ctx context.Context) ([]Network, error)
	RemoveNetwork(ctx context.Context, id string) error
	InspectNetwork(ctx context.Context, id string) (map[string]any, error)
	GetActions() map[rune]string
	GetActionsString() string
}

// LogsService defines the minimal interface for logs operations
type LogsService interface {
	GetLogs(ctx context.Context, resourceType, resourceID string) (string, error)
	GetActions() map[rune]string
	GetActionsString() string
}

// ServiceWithActions defines the minimal interface for services that provide actions
type ServiceWithActions interface {
	GetActions() map[rune]string
	GetActionsString() string
}

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

// HeaderManagerInterface defines the interface for header management
type HeaderManagerInterface interface {
	CreateHeaderSection() *tview.Flex
	UpdateAll()
	UpdateDockerInfo()
	UpdateNavigation()
	UpdateActions()
}

// ModalManagerInterface defines the interface for modal management
type ModalManagerInterface interface {
	ShowHelp()
	ShowError(error)
	ShowConfirm(string, func(bool))
}

// UIInterface defines the interface that views need from the UI
type UIInterface interface {
	// Services
	GetServices() ServiceFactoryInterface

	// UI methods
	ShowError(error)
	ShowDetails(any)
	ShowCurrentView()
	ShowConfirm(string, func(bool))

	// App methods
	GetApp() any

	// State methods
	IsInLogsMode() bool
	IsInDetailsMode() bool
	IsModalActive() bool
	IsRefreshing() bool
	GetCurrentActions() map[rune]string
	GetCurrentViewActions() string
	GetViewRegistry() any

	// Additional methods needed by managers
	GetMainFlex() any
	SwitchView(string)
	ShowHelp()

	// Additional methods needed by handlers
	GetPages() any
	ShowLogs(string, string)
	ShowLogsForResource(string, string, string) // resourceType, resourceID, resourceName
	ShowShell(string, string)

	// Additional methods needed by modal manager
	GetViewContainer() any

	// Additional methods needed by views
	GetContainerService() any
	GetImageService() any
	GetVolumeService() any
	GetNetworkService() any

	// Theme management
	GetThemeManager() *config.ThemeManager

	// Shutdown management
	GetShutdownChan() chan struct{}
}
