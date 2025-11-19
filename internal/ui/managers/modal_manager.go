package managers

import (
	"fmt"

	"github.com/rivo/tview"
	"github.com/wikczerski/whaletui/internal/ui/interfaces"
	"github.com/wikczerski/whaletui/internal/ui/modals"
)

// ModalManager handles various modal dialogs
type ModalManager struct {
	ui interfaces.UIInterface
}

// NewModalManager creates a new modal manager
func NewModalManager(ui interfaces.UIInterface) *ModalManager {
	return &ModalManager{ui: ui}
}

// ShowError displays an error modal
func (mm *ModalManager) ShowError(err error) {
	errorText := fmt.Sprintf("Error: %v", err)
	modal := modals.NewInfoModal(errorText, func() {
		mm.closeModalAndRestoreFocus("error_modal")
	})

	mm.addModalToPages("error_modal", modal)
}

// ShowInfo displays an info modal
func (mm *ModalManager) ShowInfo(message string) {
	modal := modals.NewInfoModal(message, func() {
		mm.closeModalAndRestoreFocus("info_modal")
	})

	mm.addModalToPages("info_modal", modal)
}

// ShowConfirm displays a confirmation modal with Yes/No buttons
func (mm *ModalManager) ShowConfirm(text string, callback func(bool)) {
	modal := modals.NewConfirmModal(text, func(confirmed bool) {
		mm.closeModalAndRestoreFocus("confirm_modal")
		callback(confirmed)
	})

	mm.addModalToPages("confirm_modal", modal)
}

// ShowServiceScaleModal displays a modal for scaling a service
func (mm *ModalManager) ShowServiceScaleModal(
	serviceName string,
	currentReplicas uint64,
	onConfirm func(int),
) {
	flex := modals.NewScaleModal(
		serviceName,
		currentReplicas,
		func(replicas int) {
			mm.closeModalAndRestoreFocus("scale_modal")
			onConfirm(replicas)
		},
		func() {
			mm.closeModalAndRestoreFocus("scale_modal")
		},
	)

	mm.addModalToPages("scale_modal", flex)
}

// ShowNodeAvailabilityModal displays a modal for updating node availability
func (mm *ModalManager) ShowNodeAvailabilityModal(
	nodeName, currentAvailability string,
	onConfirm func(string),
) {
	modal := modals.NewNodeAvailabilityModal(
		nodeName,
		currentAvailability,
		func(availability string) {
			mm.closeModalAndRestoreFocus("availability_modal")
			onConfirm(availability)
		},
		func() {
			mm.closeModalAndRestoreFocus("availability_modal")
		},
	)

	mm.addModalToPages("availability_modal", modal)
}

// ShowRetryDialog displays a retry dialog with automatic retry logic
func (mm *ModalManager) ShowRetryDialog(
	operation string,
	err error,
	retryFunc func() error,
	onSuccess func(),
) {
	modal := modals.NewRetryModal(
		operation,
		err,
		func() {
			// Manual retry
			mm.closeModalAndRestoreFocus("retry_modal")
			// The caller is responsible for triggering the retry logic again if manual retry is chosen?
			// Wait, the original code just closed the dialog for manual retry.
			// "Manual retry - close dialog and let user retry"
			// So this is correct.
		},
		func() {
			// Automatic retry
			mm.performAutomaticRetry(operation, retryFunc, onSuccess)
		},
		func() {
			// Cancel
			mm.closeModalAndRestoreFocus("retry_modal")
		},
	)

	mm.addModalToPages("retry_modal", modal)
}

// ShowFallbackDialog displays a fallback operations dialog
func (mm *ModalManager) ShowFallbackDialog(
	operation string,
	err error,
	fallbackOptions []string,
	onFallback func(string),
) {
	modal := modals.NewFallbackModal(
		operation,
		err,
		fallbackOptions,
		func(option string) {
			mm.closeModalAndRestoreFocus("fallback_modal")
			onFallback(option)
		},
		func() {
			mm.closeModalAndRestoreFocus("fallback_modal")
		},
	)

	mm.addModalToPages("fallback_modal", modal)
}

// performAutomaticRetry performs automatic retry with progress indication
func (mm *ModalManager) performAutomaticRetry(
	operation string,
	retryFunc func() error,
	onSuccess func(),
) {
	// Close the retry dialog
	mm.closeModal("retry_modal")

	// Show progress modal
	progressModal := modals.NewProgressModal(operation, func() {
		mm.closeModalAndRestoreFocus("retry_progress_modal")
		// Should we cancel the goroutine? The original code didn't seem to support cancellation of the goroutine itself easily.
		// It just closed the modal.
	})

	mm.addModalToPages("retry_progress_modal", progressModal)

	// Perform retry in a goroutine
	go mm.executeRetryOperation(operation, retryFunc, onSuccess)
}

// executeRetryOperation executes the retry operation in a goroutine
func (mm *ModalManager) executeRetryOperation(
	operation string,
	retryFunc func() error,
	onSuccess func(),
) {
	// Attempt retry
	err := retryFunc()

	// Close progress modal from main thread
	app, ok := mm.ui.GetApp().(*tview.Application)
	if !ok {
		return
	}
	app.QueueUpdateDraw(func() {
		mm.handleRetryResult(operation, err, onSuccess)
	})
}

// handleRetryResult handles the result of the retry operation
func (mm *ModalManager) handleRetryResult(operation string, err error, onSuccess func()) {
	mm.closeModal("retry_progress_modal")

	if err != nil {
		// Retry failed - show error
		mm.ShowError(fmt.Errorf("automatic retry failed for %s: %v", operation, err))
	} else {
		// Retry succeeded - show success and execute callback
		mm.ShowInfo(fmt.Sprintf("✅ Operation '%s' recovered successfully!", operation))
		if onSuccess != nil {
			onSuccess()
		}
	}

	// Restore focus to the main view
	mm.restoreFocusToMainView()
}

// Helper methods

func (mm *ModalManager) addModalToPages(name string, item tview.Primitive) {
	pages, ok := mm.ui.GetPages().(*tview.Pages)
	if !ok {
		return
	}
	pages.AddPage(name, item, true, true)

	app, ok := mm.ui.GetApp().(*tview.Application)
	if !ok {
		return
	}

	// If the item is a flex containing a form (like ScaleModal), we might want to focus the form?
	// The ScaleModal returns a Flex.
	// In original code: mm.setFocusToForm(form)
	// Here we just focus the item. If it's a Flex, tview usually focuses the first focusable item.
	// Let's ensure ScaleModal sets up focus correctly or we handle it here.
	// ScaleModal returns a Flex with Modal (non-focusable usually?) and Form.
	// We should probably focus the item.
	app.SetFocus(item)
}

func (mm *ModalManager) closeModal(name string) {
	pages, ok := mm.ui.GetPages().(*tview.Pages)
	if !ok {
		return
	}
	pages.RemovePage(name)
}

func (mm *ModalManager) closeModalAndRestoreFocus(name string) {
	mm.closeModal(name)
	mm.restoreFocusToMainView()
}

func (mm *ModalManager) restoreFocusToMainView() {
	if viewContainer := mm.ui.GetViewContainer(); viewContainer != nil {
		if vc, ok := viewContainer.(*tview.Flex); ok {
			app, ok := mm.ui.GetApp().(*tview.Application)
			if !ok {
				return
			}
			app.SetFocus(vc)
		}
	}
}
