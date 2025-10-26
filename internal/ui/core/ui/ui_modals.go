package ui

import (
	"github.com/rivo/tview"
)

// showError shows an error modal
func (ui *UI) showError(err error) {
	if ui.modalManager != nil {
		ui.modalManager.ShowError(err)
	}
}

// showInfo shows an info modal
func (ui *UI) showInfo(message string) {
	if ui.modalManager != nil {
		ui.modalManager.ShowInfo(message)
	}
}

// showRetryDialog shows retry dialog with automatic retry logic
func (ui *UI) showRetryDialog(
	operation string,
	err error,
	retryFunc func() error,
	onSuccess func(),
) {
	if ui.modalManager != nil {
		ui.modalManager.ShowRetryDialog(operation, err, retryFunc, onSuccess)
	}
}

// showFallbackDialog shows fallback operations dialog
func (ui *UI) showFallbackDialog(
	operation string,
	err error,
	fallbackOptions []string,
	onFallback func(string),
) {
	if ui.modalManager != nil {
		ui.modalManager.ShowFallbackDialog(operation, err, fallbackOptions, onFallback)
	}
}

// showConfirm shows a confirmation modal
func (ui *UI) showConfirm(text string, callback func(bool)) {
	if ui.modalManager != nil {
		ui.modalManager.ShowConfirm(text, callback)
	}
}

// showDetails shows a details view in the main view container
func (ui *UI) showDetails(detailsView tview.Primitive) {
	ui.viewContainer.Clear()
	ui.viewContainer.AddItem(detailsView, 0, 1, true)

	ui.app.SetFocus(detailsView)

	ui.inDetailsMode = true
	ui.updateLegend()
}

// hasValidPages checks if the pages container is valid
func (ui *UI) hasValidPages() bool {
	return ui.pages != nil
}

// hasModalPages checks if any modal pages are currently shown
func (ui *UI) hasModalPages() bool {
	return ui.pages.HasPage("help_modal") ||
		ui.pages.HasPage("error_modal") ||
		ui.pages.HasPage("confirm_modal") ||
		ui.pages.HasPage("exec_output_modal")
}
