package interfaces

import (
	"github.com/wikczerski/whaletui/internal/config"
)

// UIInterface defines the interface that views need from the UI
type UIInterface interface {
	// Services
	GetServices() ServiceFactoryInterface

	// UI methods
	ShowError(error)
	ShowInfo(string)
	ShowDetails(any)
	ShowCurrentView()
	ShowConfirm(string, func(bool))
	ShowServiceScaleModal(string, uint64, func(int))
	ShowNodeAvailabilityModal(string, string, func(string))
	ShowContextualHelp(string, string)
	ShowRetryDialog(string, error, func() error, func())
	ShowFallbackDialog(string, error, []string, func(string))

	// App methods
	GetApp() any

	// State methods
	IsInLogsMode() bool
	IsInDetailsMode() bool
	IsModalActive() bool
	IsRefreshing() bool
	GetCurrentActions() map[rune]string
	GetCurrentViewActions() string
	GetCurrentViewNavigation() string
	GetViewRegistry() any

	// Additional methods needed by managers
	GetMainFlex() any
	SwitchView(string)
	ShowHelp()

	// Additional methods needed by handlers
	GetPages() any
	ShowLogs(string, string)
	ShowLogsForResource(string, string, string)
	ShowShell(string, string)

	// Additional methods needed by modal manager
	GetViewContainer() any

	// Additional methods needed by views
	GetContainerService() any
	GetImageService() any
	GetVolumeService() any
	GetNetworkService() any
	GetServicesAny() any
	GetSwarmServiceService() any
	GetSwarmNodeService() any
	IsContainerServiceAvailable() bool

	// Theme management
	GetThemeManager() *config.ThemeManager

	// Shutdown management
	GetShutdownChan() chan struct{}
}
