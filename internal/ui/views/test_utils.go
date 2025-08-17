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

// Implement UIInterface methods
func (m *MockUI) GetServices() *services.ServiceFactory { return m.services }
func (m *MockUI) GetApp() any                           { return m.app }
func (m *MockUI) GetPages() any                         { return m.pages }
func (m *MockUI) GetMainFlex() any                      { return m.mainFlex }
func (m *MockUI) GetLog() any                           { return m.log }
func (m *MockUI) GetViewContainer() any                 { return m.viewContainer }
func (m *MockUI) GetContainerService() any {
	if m.services != nil {
		return m.services.ContainerService
	}
	return nil
}
func (m *MockUI) GetImageService() any {
	if m.services != nil {
		return m.services.ImageService
	}
	return nil
}
func (m *MockUI) GetVolumeService() any {
	if m.services != nil {
		return m.services.VolumeService
	}
	return nil
}
func (m *MockUI) GetNetworkService() any {
	if m.services != nil {
		return m.services.NetworkService
	}
	return nil
}
func (m *MockUI) IsInLogsMode() bool                               { return false }
func (m *MockUI) IsInDetailsMode() bool                            { return false }
func (m *MockUI) GetCurrentActions() map[rune]string               { return nil }
func (m *MockUI) GetCurrentViewActions() string                    { return "" }
func (m *MockUI) GetViewRegistry() any                             { return nil }
func (m *MockUI) SwitchView(view string)                           {}
func (m *MockUI) ShowHelp()                                        {}
func (m *MockUI) ShowLogs(containerID, containerName string)       {}
func (m *MockUI) ShowShell(containerID, containerName string)      {}
func (m *MockUI) ShowError(err error)                              {}
func (m *MockUI) ShowDetails(details any)                          {}
func (m *MockUI) ShowCurrentView()                                 {}
func (m *MockUI) ShowConfirm(message string, onConfirm func(bool)) {}
func (m *MockUI) ShowSuccess(message string)                       {}
func (m *MockUI) GetThemeManager() *config.ThemeManager            { return nil }
func (m *MockUI) UpdateStatusBar(message string)                   {}
func (m *MockUI) GetInputField() any                               { return nil }

// GetShutdownChan returns a mock shutdown channel for testing
func (m *MockUI) GetShutdownChan() chan struct{} {
	return make(chan struct{})
}
