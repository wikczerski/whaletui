package handlers

import (
	"github.com/rivo/tview"
	"github.com/wikczerski/whaletui/internal/ui/interfaces"
)

// BaseInputHandler provides common functionality for input handlers
type BaseInputHandler struct {
	ui           interfaces.UIInterface
	commandInput *tview.InputField
	isActive     bool
}

// NewBaseInputHandler creates a new base input handler
func NewBaseInputHandler(ui interfaces.UIInterface) *BaseInputHandler {
	return &BaseInputHandler{ui: ui}
}

// GetInput returns the input widget
func (bh *BaseInputHandler) GetInput() *tview.InputField {
	return bh.commandInput
}

// IsActive returns whether the handler is active
func (bh *BaseInputHandler) IsActive() bool {
	return bh.isActive
}

// SetActive sets the active state
func (bh *BaseInputHandler) SetActive(active bool) {
	bh.isActive = active
}

// SetInput sets the input field
func (bh *BaseInputHandler) SetInput(input *tview.InputField) {
	bh.commandInput = input
}

// ShowInput shows the input field in the main flex
func (bh *BaseInputHandler) ShowInput() {
	mainFlex, ok := bh.ui.GetMainFlex().(*tview.Flex)
	if !ok {
		return
	}
	mainFlex.AddItem(bh.commandInput, 3, 1, true)
	app, ok := bh.ui.GetApp().(*tview.Application)
	if !ok {
		return
	}
	app.SetFocus(bh.commandInput)
}

// HideInput hides the input field from the main flex
func (bh *BaseInputHandler) HideInput() {
	mainFlex, ok := bh.ui.GetMainFlex().(*tview.Flex)
	if !ok {
		return
	}
	mainFlex.RemoveItem(bh.commandInput)
	app, ok := bh.ui.GetApp().(*tview.Application)
	if !ok {
		return
	}
	app.SetFocus(mainFlex)
}

// ClearError clears any error messages
func (bh *BaseInputHandler) ClearError() {
	// Implementation depends on how errors are displayed in the UI
	// This is a placeholder for now
}
