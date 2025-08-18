package managers

import (
	"testing"

	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
	"github.com/wikczerski/D5r/internal/config"
	"github.com/wikczerski/D5r/internal/ui/interfaces/mocks"
)

func newHeaderManagerWithTheme(t *testing.T) *HeaderManager {
	mockUI := mocks.NewMockUIInterface(t)
	// Return a real ThemeManager from the mock
	tm := config.NewThemeManager("")
	mockUI.On("GetThemeManager").Return(tm).Maybe()
	return NewHeaderManager(mockUI)
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
	assert.NotNil(t, manager.GetDockerInfoCol())
}
