package managers

import (
	"errors"
	"testing"

	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
	mocks "github.com/wikczerski/whaletui/internal/mocks/ui"
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

func TestModalManager_ShowError_DoesNotPanic(t *testing.T) {
	manager := newModalManagerWithUI(t)
	assert.NotPanics(t, func() { manager.ShowError(errors.New("test error")) })
}

func TestModalManager_ShowConfirm_DoesNotPanic(t *testing.T) {
	manager := newModalManagerWithUI(t)
	assert.NotPanics(t, func() { manager.ShowConfirm("confirm?", func(bool) {}) })
}
