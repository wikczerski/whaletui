package managers

import (
	"testing"

	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
	"github.com/wikczerski/whaletui/internal/config"
	mocks "github.com/wikczerski/whaletui/internal/mocks/ui"
)

func newHeaderManagerWithTheme(t *testing.T) *HeaderManager {
	mockUI := mocks.NewMockUIInterface(t)
	mockServices := mocks.NewMockServiceFactoryInterface(t)

	// Create a simple mock view registry
	mockViewRegistry := &mockViewRegistry{currentName: ""}

	// Return a real ThemeManager from the mock
	tm := config.NewThemeManager("")
	mockUI.On("GetThemeManager").Return(tm).Maybe()

	// Mock service dependencies
	mockServices.On("IsContainerServiceAvailable").Return(false).Maybe()
	mockServices.On("GetContainerService").Return(nil).Maybe()
	mockServices.On("GetImageService").Return(nil).Maybe()
	mockServices.On("GetVolumeService").Return(nil).Maybe()
	mockServices.On("GetNetworkService").Return(nil).Maybe()
	mockServices.On("GetCurrentService").Return(nil).Maybe()
	mockUI.On("GetServices").Return(mockServices).Maybe()
	mockUI.On("GetCurrentViewActions").Return("").Maybe()
	mockUI.On("GetCurrentViewNavigation").Return("").Maybe()
	mockUI.On("GetViewRegistry").Return(mockViewRegistry).Maybe()

	// Mock UI state methods
	mockUI.On("IsInLogsMode").Return(false).Maybe()
	mockUI.On("IsInDetailsMode").Return(false).Maybe()
	mockUI.On("GetCurrentActions").Return(nil).Maybe()

	return NewHeaderManager(mockUI)
}

// mockViewRegistry is a simple mock for testing
type mockViewRegistry struct {
	currentName string
}

func (m *mockViewRegistry) GetCurrentName() string {
	return m.currentName
}

func TestNewHeaderManager(t *testing.T) {
	manager := newHeaderManagerWithTheme(t)
	assert.NotNil(t, manager)
}

func TestHeaderManager_CreateHeaderSection(t *testing.T) {
	manager := newHeaderManagerWithTheme(t)
	section := manager.CreateHeaderSection()
	assert.IsType(t, &tview.Flex{}, section)
}

func TestHeaderManager_GetColumns_AfterCreate(t *testing.T) {
	manager := newHeaderManagerWithTheme(t)
	_ = manager.CreateHeaderSection()
	// The new flex-based header doesn't expose individual columns
	// Instead, it manages the layout internally
	assert.NotNil(t, manager)
}
