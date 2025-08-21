package views

import (
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
	"github.com/wikczerski/whaletui/internal/config"
	"github.com/wikczerski/whaletui/internal/services"
	"github.com/wikczerski/whaletui/internal/ui/interfaces/mocks"
)

func TestNewLogsView(t *testing.T) {
	mockUI := mocks.NewMockUIInterface(t)
	mockThemeManager := config.NewThemeManager("")

	mockUI.On("GetThemeManager").Return(mockThemeManager)

	logsView := NewLogsView(mockUI, "container", "test-id", "test-name")

	assert.NotNil(t, logsView)
	assert.Equal(t, "container", logsView.ResourceType)
	assert.Equal(t, "test-id", logsView.ResourceID)
	assert.Equal(t, "test-name", logsView.ResourceName)
	assert.Equal(t, mockUI, logsView.ui)
	assert.Equal(t, mockThemeManager, logsView.themeManager)
}

func TestLogsView_GetView(t *testing.T) {
	mockUI := mocks.NewMockUIInterface(t)
	mockThemeManager := config.NewThemeManager("")

	mockUI.On("GetThemeManager").Return(mockThemeManager)

	logsView := NewLogsView(mockUI, "container", "test-id", "test-name")
	view := logsView.GetView()

	assert.NotNil(t, view)
	assert.IsType(t, &tview.Flex{}, view)
}

func TestLogsView_LoadLogs_Success(t *testing.T) {
	mockUI := mocks.NewMockUIInterface(t)
	mockThemeManager := config.NewThemeManager("")
	mockServices := &services.ServiceFactory{}

	mockUI.On("GetThemeManager").Return(mockThemeManager)
	mockUI.On("GetServices").Return(mockServices)

	logsView := NewLogsView(mockUI, "container", "test-id", "test-name")

	// Test that LoadLogs doesn't panic
	assert.NotPanics(t, func() {
		logsView.LoadLogs()
	})
}

func TestLogsView_LoadLogs_ServiceNotAvailable(t *testing.T) {
	mockUI := mocks.NewMockUIInterface(t)
	mockThemeManager := config.NewThemeManager("")
	mockServices := &services.ServiceFactory{}

	mockUI.On("GetThemeManager").Return(mockThemeManager)
	mockUI.On("GetServices").Return(mockServices)

	logsView := NewLogsView(mockUI, "container", "test-id", "test-name")

	// Test that LoadLogs handles service unavailability gracefully
	assert.NotPanics(t, func() {
		logsView.LoadLogs()
	})
}

func TestLogsView_GetActions(t *testing.T) {
	mockUI := mocks.NewMockUIInterface(t)
	mockThemeManager := config.NewThemeManager("")

	// Create a real ServiceFactory with a real LogsService
	realServices := services.NewServiceFactory(nil)

	mockUI.On("GetThemeManager").Return(mockThemeManager)
	mockUI.On("GetServices").Return(realServices)

	logsView := NewLogsView(mockUI, "container", "test-id", "test-name")
	actions := logsView.GetActions()

	// Since the ServiceFactory is created with nil client, it won't have a LogsService
	// So we expect an empty map
	assert.Equal(t, map[rune]string{}, actions)
}

func TestLogsView_KeyBindings_Escape(t *testing.T) {
	mockUI := mocks.NewMockUIInterface(t)
	mockThemeManager := config.NewThemeManager("")

	mockUI.On("GetThemeManager").Return(mockThemeManager)
	mockUI.On("ShowCurrentView").Return()

	logsView := NewLogsView(mockUI, "container", "test-id", "test-name")
	view := logsView.GetView()

	// Test escape key
	escapeEvent := tcell.NewEventKey(tcell.KeyEscape, 0, tcell.ModNone)
	result := view.(*tview.Flex).GetInputCapture()(escapeEvent)

	assert.Nil(t, result)
}

func TestLogsView_KeyBindings_Enter(t *testing.T) {
	mockUI := mocks.NewMockUIInterface(t)
	mockThemeManager := config.NewThemeManager("")

	mockUI.On("GetThemeManager").Return(mockThemeManager)
	mockUI.On("ShowCurrentView").Return()

	logsView := NewLogsView(mockUI, "container", "test-id", "test-name")
	view := logsView.GetView()

	// Test enter key
	enterEvent := tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone)
	result := view.(*tview.Flex).GetInputCapture()(enterEvent)

	assert.Nil(t, result)
}

func TestLogsView_KeyBindings_ScrollingKeys(t *testing.T) {
	mockUI := mocks.NewMockUIInterface(t)
	mockThemeManager := config.NewThemeManager("")

	mockUI.On("GetThemeManager").Return(mockThemeManager)

	logsView := NewLogsView(mockUI, "container", "test-id", "test-name")
	view := logsView.GetView()

	// Test scrolling keys
	scrollingKeys := []tcell.Key{tcell.KeyUp, tcell.KeyDown, tcell.KeyPgUp, tcell.KeyPgDn, tcell.KeyHome, tcell.KeyEnd}

	for _, key := range scrollingKeys {
		event := tcell.NewEventKey(key, 0, tcell.ModNone)
		result := view.(*tview.Flex).GetInputCapture()(event)

		assert.Equal(t, event, result)
	}
}

func TestLogsView_ResourceTypeDisplay(t *testing.T) {
	mockUI := mocks.NewMockUIInterface(t)
	mockThemeManager := config.NewThemeManager("")

	mockUI.On("GetThemeManager").Return(mockThemeManager)

	// Test with different resource types
	testCases := []struct {
		resourceType string
		resourceID   string
		resourceName string
	}{
		{"container", "test-container-id", "test-container"},
		{"service", "test-service-id", "test-service"},
		{"task", "test-task-id", "test-task"},
	}

	for _, tc := range testCases {
		t.Run(tc.resourceType, func(t *testing.T) {
			logsView := NewLogsView(mockUI, tc.resourceType, tc.resourceID, tc.resourceName)

			assert.Equal(t, tc.resourceType, logsView.ResourceType)
			assert.Equal(t, tc.resourceID, logsView.ResourceID)
			assert.Equal(t, tc.resourceName, logsView.ResourceName)
		})
	}
}

func TestLogsView_ResourceIDTruncation(t *testing.T) {
	mockUI := mocks.NewMockUIInterface(t)
	mockThemeManager := config.NewThemeManager("")

	mockUI.On("GetThemeManager").Return(mockThemeManager)

	// Test with long ID
	longID := "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
	logsView := NewLogsView(mockUI, "container", longID, "test-name")

	assert.Equal(t, longID, logsView.ResourceID)
}

func TestLogsView_EmptyResourceName(t *testing.T) {
	mockUI := mocks.NewMockUIInterface(t)
	mockThemeManager := config.NewThemeManager("")

	mockUI.On("GetThemeManager").Return(mockThemeManager)

	// Test with empty resource name
	logsView := NewLogsView(mockUI, "container", "test-id", "")

	assert.Equal(t, "", logsView.ResourceName)
}
