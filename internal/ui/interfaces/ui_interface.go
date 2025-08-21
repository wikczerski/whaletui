package interfaces

import (
	"github.com/rivo/tview"
	"github.com/wikczerski/whaletui/internal/config"
	"github.com/wikczerski/whaletui/internal/services"
)

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
	GetServices() services.ServiceFactoryInterface

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
