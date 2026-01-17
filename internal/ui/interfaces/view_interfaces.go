package interfaces

import (
	"github.com/rivo/tview"
	"github.com/wikczerski/whaletui/internal/config"
)

// SharedUIInterface defines the interface that views need from the UI
type SharedUIInterface interface {
	// Basic UI methods
	ShowError(error)
	ShowInfo(string)
	ShowDetails(any)
	ShowCurrentView()
	ShowConfirm(string, func(bool))

	// Advanced UI methods
	ShowServiceScaleModal(string, uint64, func(int))
	ShowNodeAvailabilityModal(string, string, func(string))
	ShowRetryDialog(string, error, func() error, func())
	ShowFallbackDialog(string, error, []string, func(string))

	// Service methods
	GetServicesAny() any
	GetSwarmServiceService() any
	GetSwarmNodeService() any

	// Theme management
	GetThemeManager() *config.ThemeManager
}

// BaseViewInterface defines common functionality for all Docker resource views
type BaseViewInterface interface {
	GetView() tview.Primitive
	GetUI() SharedUIInterface
	Refresh()
	Search(searchTerm string)
	ClearSearch()
	IsSearchActive() bool
	GetSearchTerm() string
}
