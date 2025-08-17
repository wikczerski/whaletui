package interfaces

import (
	"github.com/wikczerski/D5r/internal/config"
	"github.com/wikczerski/D5r/internal/services"
)

// UIInterface defines the interface that views need from the UI
type UIInterface interface {
	// Services
	GetServices() *services.ServiceFactory

	// UI methods
	ShowError(error)
	ShowSuccess(string)
	ShowDetails(any)
	ShowCurrentView()
	ShowConfirm(string, func(bool))

	// App methods
	GetApp() any

	// Status methods
	UpdateStatusBar(string)

	// State methods
	IsInLogsMode() bool
	IsInDetailsMode() bool
	GetCurrentActions() map[rune]string
	GetCurrentViewActions() string
	GetViewRegistry() any

	// Additional methods needed by managers
	GetMainFlex() any
	GetLog() any
	SwitchView(string)
	ShowHelp()

	// Additional methods needed by handlers
	GetPages() any
	ShowLogs(string, string)
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
