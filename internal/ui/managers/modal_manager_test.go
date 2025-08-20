package managers

import (
	"errors"
	"testing"

	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
	"github.com/wikczerski/whaletui/internal/ui/interfaces/mocks"
)

func newModalManagerWithUI(t *testing.T) *ModalManager {
	mockUI := mocks.NewMockUIInterface(t)
	pages := tview.NewPages()
	app := tview.NewApplication()
	viewContainer := tview.NewFlex()

	mockUI.On("GetPages").Return(pages).Maybe()
	mockUI.On("GetApp").Return(app).Maybe()
	mockUI.On("GetViewContainer").Return(viewContainer).Maybe()
	return NewModalManager(mockUI)
}

func TestNewModalManager(t *testing.T) {
	manager := newModalManagerWithUI(t)
	assert.NotNil(t, manager)
}

func TestModalManager_ShowHelp_DoesNotPanic(t *testing.T) {
	manager := newModalManagerWithUI(t)
	assert.NotPanics(t, func() { manager.ShowHelp() })
}

func TestModalManager_ShowError_DoesNotPanic(t *testing.T) {
	manager := newModalManagerWithUI(t)
	assert.NotPanics(t, func() { manager.ShowError(errors.New("test error")) })
}

func TestModalManager_ShowConfirm_DoesNotPanic(t *testing.T) {
	manager := newModalManagerWithUI(t)
	assert.NotPanics(t, func() { manager.ShowConfirm("confirm?", func(bool) {}) })
}

// Private helpers like createModal/buildHelpText are not tested directly.
