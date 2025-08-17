package views

import (
	"github.com/rivo/tview"
	"github.com/wikczerski/D5r/internal/config"
	"github.com/wikczerski/D5r/internal/logger"
	"github.com/wikczerski/D5r/internal/services"
)

// MockUI is a mock implementation for testing
type MockUI struct {
	app           *tview.Application
	pages         *tview.Pages
	viewContainer *tview.Flex
	mainFlex      *tview.Flex
	statusBar     *tview.TextView
	log           *logger.Logger
	services      *services.ServiceFactory
}

// NewMockUI creates a new mock UI for testing
func NewMockUI() *MockUI {
	return &MockUI{
		app:           tview.NewApplication(),
		pages:         tview.NewPages(),
		viewContainer: tview.NewFlex(),
		mainFlex:      tview.NewFlex(),
		statusBar:     tview.NewTextView(),
		log:           logger.GetLogger(),
	}
}

// SetServices sets the services for the mock UI
func (m *MockUI) SetServices(services *services.ServiceFactory) {
	m.services = services
}

// GetServices returns the service factory
func (m *MockUI) GetServices() *services.ServiceFactory { return m.services }

// GetApp returns the tview application
func (m *MockUI) GetApp() any { return m.app }

// GetPages returns the tview pages
func (m *MockUI) GetPages() any { return m.pages }

// GetMainFlex returns the main flex container
func (m *MockUI) GetMainFlex() any { return m.mainFlex }

// GetLog returns the logger
func (m *MockUI) GetLog() any { return m.log }

// GetViewContainer returns the view container
func (m *MockUI) GetViewContainer() any { return m.viewContainer }

// GetContainerService returns the container service
func (m *MockUI) GetContainerService() any {
	if m.services != nil {
		return m.services.ContainerService
	}
	return nil
}

// GetImageService returns the image service
func (m *MockUI) GetImageService() any {
	if m.services != nil {
		return m.services.ImageService
	}
	return nil
}

// GetVolumeService returns the volume service
func (m *MockUI) GetVolumeService() any {
	if m.services != nil {
		return m.services.VolumeService
	}
	return nil
}

// GetNetworkService returns the network service
func (m *MockUI) GetNetworkService() any {
	if m.services != nil {
		return m.services.NetworkService
	}
	return nil
}

// IsInLogsMode returns whether the UI is in logs mode
func (m *MockUI) IsInLogsMode() bool { return false }

// IsInDetailsMode returns whether the UI is in details mode
func (m *MockUI) IsInDetailsMode() bool { return false }

// GetCurrentActions returns the current actions
func (m *MockUI) GetCurrentActions() map[rune]string { return nil }

// GetCurrentViewActions returns the current view actions
func (m *MockUI) GetCurrentViewActions() string { return "" }

// GetViewRegistry returns the view registry
func (m *MockUI) GetViewRegistry() any { return nil }

// SwitchView switches to a different view
func (m *MockUI) SwitchView(_ string) {}

// ShowHelp shows the help information
func (m *MockUI) ShowHelp() {}

// ShowLogs shows the logs for a container
func (m *MockUI) ShowLogs(_, _ string) {}

// ShowShell shows the shell for a container
func (m *MockUI) ShowShell(_, _ string) {}

// ShowError shows an error message
func (m *MockUI) ShowError(_ error) {}

// ShowDetails shows detailed information
func (m *MockUI) ShowDetails(_ any) {}

// ShowCurrentView shows the current view
func (m *MockUI) ShowCurrentView() {}

// ShowConfirm shows a confirmation dialog
func (m *MockUI) ShowConfirm(_ string, _ func(bool)) {}

// ShowSuccess shows a success message
func (m *MockUI) ShowSuccess(_ string) {}

// GetThemeManager returns the theme manager
func (m *MockUI) GetThemeManager() *config.ThemeManager { return nil }

// UpdateStatusBar updates the status bar
func (m *MockUI) UpdateStatusBar(_ string) {}

// GetInputField returns the input field
func (m *MockUI) GetInputField() any { return nil }

// GetShutdownChan returns a mock shutdown channel for testing
func (m *MockUI) GetShutdownChan() chan struct{} {
	return make(chan struct{})
}
