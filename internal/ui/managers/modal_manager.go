package managers

import (
	"fmt"
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/wikczerski/whaletui/internal/ui/interfaces"
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
	modal := mm.createModal(errorText, []string{"OK"})

	mm.setupErrorModalHandlers(modal)
	mm.showErrorModal(modal)
}

// ShowInfo displays an info modal
func (mm *ModalManager) ShowInfo(message string) {
	modal := mm.createModal(message, []string{"OK"})

	// Add done function to handle OK button click
	modal.SetDoneFunc(func(_ int, _ string) {
		mm.closeInfoModalAndRestoreFocus()
	})

	// Add keyboard handling for ESC key to close modal
	modal.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			mm.closeInfoModalAndRestoreFocus()
			return nil // Consume the event
		}
		return event
	})

	mm.showInfoModal(modal)
}

// ShowConfirm displays a confirmation modal with Yes/No buttons
func (mm *ModalManager) ShowConfirm(text string, callback func(bool)) {
	modal := mm.createConfirmModal(text, callback)
	mm.setupConfirmModalHandlers(modal, callback)
	mm.showConfirmModal(modal)
}

// ShowServiceScaleModal displays a modal for scaling a service
func (mm *ModalManager) ShowServiceScaleModal(
	serviceName string,
	currentReplicas uint64,
	onConfirm func(int),
) {
	// Create input field for replicas
	inputField := mm.createReplicasInputField(currentReplicas)

	// Create form with input and buttons including help
	form := mm.createScaleForm(inputField, onConfirm)

	// Create modal container
	modal := mm.createScaleModal(serviceName, currentReplicas)

	// Create a flex container to hold both modal and form
	flex := mm.createScaleFlexContainer(modal, form)

	// Add the modal to the pages
	mm.addScaleModalToPages(flex)

	// Set focus to the form
	mm.setFocusToForm(form)
}

// ShowNodeAvailabilityModal displays a modal for updating node availability
func (mm *ModalManager) ShowNodeAvailabilityModal(
	nodeName, currentAvailability string,
	onConfirm func(string),
) {
	// Create the modal content
	content := mm.createNodeAvailabilityContent(nodeName, currentAvailability)

	// Create modal with help button
	modal := mm.createModal(content, []string{"Active", "Pause", "Drain", "Cancel"})

	// Add done function to handle button clicks
	modal.SetDoneFunc(func(_ int, buttonLabel string) {
		mm.handleNodeAvailabilityButtonClick(buttonLabel, onConfirm)
	})

	// Add keyboard handling for ESC key to close modal
	modal.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			mm.closeNodeAvailabilityModalAndRestoreFocus()
			return nil // Consume the event
		}
		return event
	})

	// Add the modal to the pages
	mm.addNodeAvailabilityModalToPages(modal)

	// Set focus to the modal so it can receive keyboard input
	app, ok := mm.ui.GetApp().(*tview.Application)
	if !ok {
		return
	}
	app.SetFocus(modal)
}

// ShowRetryDialog displays a retry dialog with automatic retry logic
func (mm *ModalManager) ShowRetryDialog(
	operation string,
	err error,
	retryFunc func() error,
	onSuccess func(),
) {
	// Create retry dialog content
	content := mm.createRetryDialogContent(operation, err)

	// Create modal with retry options
	modal := mm.createModal(content, []string{"Retry", "Retry (Auto)", "Cancel"})

	// Add done function to handle button clicks
	modal.SetDoneFunc(func(_ int, buttonLabel string) {
		mm.handleRetryDialogButtonClick(buttonLabel, operation, retryFunc, onSuccess)
	})

	// Add keyboard handling for ESC key to close modal
	modal.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			mm.closeRetryDialogAndRestoreFocus()
			return nil // Consume the event
		}
		return event
	})

	// Add the modal to the pages
	mm.addRetryDialogToPages(modal)

	// Set focus to the modal so it can receive keyboard input
	app, ok := mm.ui.GetApp().(*tview.Application)
	if !ok {
		return
	}
	app.SetFocus(modal)
}

// ShowFallbackDialog displays a fallback operations dialog
func (mm *ModalManager) ShowFallbackDialog(
	operation string,
	err error,
	fallbackOptions []string,
	onFallback func(string),
) {
	content := mm.createFallbackContent(operation, err)
	buttons := mm.createFallbackButtons(fallbackOptions)
	modal := mm.createModal(content, buttons)
	mm.setupFallbackModalHandlers(modal, onFallback)
	mm.showFallbackModal(modal)
}

// createFallbackContent creates the content for the fallback dialog
func (mm *ModalManager) createFallbackContent(operation string, err error) string {
	return fmt.Sprintf(
		"⚠️  Operation Failed: %s\n\nError: %v\n\nAlternative operations are available:",
		operation,
		err,
	)
}

// setupErrorModalHandlers sets up the handlers for the error modal
func (mm *ModalManager) setupErrorModalHandlers(modal *tview.Modal) {
	// Add done function to handle OK button click
	modal.SetDoneFunc(func(_ int, _ string) {
		mm.closeErrorModalAndRestoreFocus()
	})

	// Add keyboard handling for ESC key to close modal
	modal.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			mm.closeErrorModalAndRestoreFocus()
			return nil // Consume the event
		}
		return event
	})
}

// closeErrorModalAndRestoreFocus closes the error modal and restores focus
func (mm *ModalManager) closeErrorModalAndRestoreFocus() {
	pages, ok := mm.ui.GetPages().(*tview.Pages)
	if !ok {
		return
	}
	pages.RemovePage("error_modal")
	mm.restoreFocusToMainView()
}

// showErrorModal shows the error modal
func (mm *ModalManager) showErrorModal(modal *tview.Modal) {
	pages, ok := mm.ui.GetPages().(*tview.Pages)
	if !ok {
		return
	}
	pages.AddPage("error_modal", modal, true, true)

	// Set focus to the modal so it can receive keyboard input
	app, ok := mm.ui.GetApp().(*tview.Application)
	if !ok {
		return
	}
	app.SetFocus(modal)
}

// createReplicasInputField creates the replicas input field
func (mm *ModalManager) createReplicasInputField(currentReplicas uint64) *tview.InputField {
	return tview.NewInputField().
		SetLabel("Replicas: ").
		SetText(fmt.Sprintf("%d", currentReplicas)).
		SetFieldWidth(10).
		SetAcceptanceFunc(tview.InputFieldInteger)
}

// createScaleForm creates the scale form with buttons
func (mm *ModalManager) createScaleForm(
	inputField *tview.InputField,
	onConfirm func(int),
) *tview.Form {
	return tview.NewForm().
		AddFormItem(inputField).
		AddButton("Scale", func() {
			// Parse replicas from input
			replicasStr := inputField.GetText()
			replicas, err := strconv.Atoi(replicasStr)
			if err != nil || replicas < 0 {
				mm.ShowError(fmt.Errorf("invalid replicas value: %s", replicasStr))
				return
			}

			// Close modal and call callback
			mm.closeScaleModalAndRestoreFocus()
			onConfirm(replicas)
		}).
		AddButton("Cancel", func() {
			// Close modal without action
			mm.closeScaleModalAndRestoreFocus()
		})
}

// createScaleModal creates the scale modal
func (mm *ModalManager) createScaleModal(serviceName string, currentReplicas uint64) *tview.Modal {
	return tview.NewModal().
		SetText(fmt.Sprintf("Scale Service: %s\nCurrent Replicas: %d", serviceName, currentReplicas)).
		SetDoneFunc(func(_ int, _ string) {
			mm.closeScaleModalAndRestoreFocus()
		})
}

// createScaleFlexContainer creates the flex container for the scale modal
func (mm *ModalManager) createScaleFlexContainer(modal *tview.Modal, form *tview.Form) *tview.Flex {
	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(modal, 0, 1, false).
		AddItem(form, 0, 1, true)

	// Add keyboard handling for ESC key
	flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			mm.closeScaleModalAndRestoreFocus()
			return nil
		}
		return event
	})

	return flex
}

// createNodeAvailabilityContent creates the content for the node availability modal
func (mm *ModalManager) createNodeAvailabilityContent(nodeName, currentAvailability string) string {
	return fmt.Sprintf(
		"Update Node Availability: %s\n\nCurrent Availability: %s\n\nSelect new availability:",
		nodeName,
		currentAvailability,
	)
}

// handleNodeAvailabilityButtonClick handles button clicks in the node availability modal
func (mm *ModalManager) handleNodeAvailabilityButtonClick(
	buttonLabel string,
	onConfirm func(string),
) {
	switch buttonLabel {
	case "Active":
		mm.closeNodeAvailabilityModal()
		onConfirm("active")
	case "Pause":
		mm.closeNodeAvailabilityModal()
		onConfirm("pause")
	case "Drain":
		mm.closeNodeAvailabilityModal()
		onConfirm("drain")
	case "Cancel":
		// Close the modal without action
		mm.closeNodeAvailabilityModal()
	}
}

// closeNodeAvailabilityModal closes the node availability modal
func (mm *ModalManager) closeNodeAvailabilityModal() {
	pages, ok := mm.ui.GetPages().(*tview.Pages)
	if !ok {
		return
	}
	pages.RemovePage("availability_modal")
}

// closeNodeAvailabilityModalAndRestoreFocus closes the node availability modal and restores focus
func (mm *ModalManager) closeNodeAvailabilityModalAndRestoreFocus() {
	pages, ok := mm.ui.GetPages().(*tview.Pages)
	if !ok {
		return
	}
	pages.RemovePage("availability_modal")
	mm.restoreFocusToMainView()
}

// addNodeAvailabilityModalToPages adds the node availability modal to the pages
func (mm *ModalManager) addNodeAvailabilityModalToPages(modal *tview.Modal) {
	pages, ok := mm.ui.GetPages().(*tview.Pages)
	if !ok {
		return
	}
	pages.AddPage("availability_modal", modal, true, true)
}

// createRetryDialogContent creates the content for the retry dialog
func (mm *ModalManager) createRetryDialogContent(operation string, err error) string {
	return fmt.Sprintf(
		"🔄 Operation Failed: %s\n\nError: %v\n\nThis may be a temporary issue. Would you like to retry?",
		operation,
		err,
	)
}

// handleRetryDialogButtonClick handles button clicks in the retry dialog
func (mm *ModalManager) handleRetryDialogButtonClick(
	buttonLabel, operation string,
	retryFunc func() error,
	onSuccess func(),
) {
	switch buttonLabel {
	case "Retry":
		// Manual retry - close dialog and let user retry
		mm.closeRetryDialogAndRestoreFocus()
	case "Retry (Auto)":
		// Automatic retry with progress indication
		mm.performAutomaticRetry(operation, retryFunc, onSuccess)
	case "Cancel":
		// Close dialog without retry
		mm.closeRetryDialogAndRestoreFocus()
	}
}

// closeRetryDialogAndRestoreFocus closes the retry dialog and restores focus
func (mm *ModalManager) closeRetryDialogAndRestoreFocus() {
	pages, ok := mm.ui.GetPages().(*tview.Pages)
	if !ok {
		return
	}
	pages.RemovePage("retry_modal")
	mm.restoreFocusToMainView()
}

// addRetryDialogToPages adds the retry dialog to the pages
func (mm *ModalManager) addRetryDialogToPages(modal *tview.Modal) {
	pages, ok := mm.ui.GetPages().(*tview.Pages)
	if !ok {
		return
	}
	pages.AddPage("retry_modal", modal, true, true)
}

// createFallbackButtons creates the buttons for the fallback dialog
func (mm *ModalManager) createFallbackButtons(fallbackOptions []string) []string {
	buttons := make([]string, len(fallbackOptions)+1)
	copy(buttons, fallbackOptions)
	buttons[len(fallbackOptions)] = "Cancel"
	return buttons
}

// setupFallbackModalHandlers sets up the handlers for the fallback modal
func (mm *ModalManager) setupFallbackModalHandlers(modal *tview.Modal, onFallback func(string)) {
	modal.SetDoneFunc(func(_ int, buttonLabel string) {
		mm.handleFallbackButtonClick(buttonLabel, onFallback)
	})

	modal.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			mm.handleFallbackEscape()
			return nil
		}
		return event
	})
}

// handleFallbackButtonClick handles button clicks in the fallback modal
func (mm *ModalManager) handleFallbackButtonClick(buttonLabel string, onFallback func(string)) {
	if buttonLabel == "Cancel" {
		mm.closeFallbackModal()
		mm.restoreFocusToMainView()
	} else {
		onFallback(buttonLabel)
		mm.closeFallbackModal()
		mm.restoreFocusToMainView()
	}
}

// handleFallbackEscape handles ESC key press in the fallback modal
func (mm *ModalManager) handleFallbackEscape() {
	mm.closeFallbackModal()
	mm.restoreFocusToMainView()
}

// closeFallbackModal closes the fallback modal
func (mm *ModalManager) closeFallbackModal() {
	pages, ok := mm.ui.GetPages().(*tview.Pages)
	if !ok {
		return
	}
	pages.RemovePage("fallback_modal")
}

// showFallbackModal displays the fallback modal
func (mm *ModalManager) showFallbackModal(modal *tview.Modal) {
	pages, ok := mm.ui.GetPages().(*tview.Pages)
	if !ok {
		return
	}
	pages.AddPage("fallback_modal", modal, true, true)

	// Set focus to the modal so it can receive keyboard input
	app, ok := mm.ui.GetApp().(*tview.Application)
	if !ok {
		return
	}
	app.SetFocus(modal)
}

// createModal creates a standard modal with consistent styling
func (mm *ModalManager) createModal(text string, buttons []string) *tview.Modal {
	return tview.NewModal().
		SetText(text).
		AddButtons(buttons)
}

// performAutomaticRetry performs automatic retry with progress indication
func (mm *ModalManager) performAutomaticRetry(
	operation string,
	retryFunc func() error,
	onSuccess func(),
) {
	// Close the retry dialog
	mm.closeRetryDialog()

	// Show progress modal
	progressModal := mm.createProgressModal(operation)

	// Add the progress modal to the pages
	mm.addProgressModalToPages(progressModal)

	// Perform retry in a goroutine to avoid blocking UI
	go mm.executeRetryOperation(operation, retryFunc, onSuccess)
}

// closeRetryDialog closes the retry dialog
func (mm *ModalManager) closeRetryDialog() {
	pages, ok := mm.ui.GetPages().(*tview.Pages)
	if !ok {
		return
	}
	pages.RemovePage("retry_modal")
}

// createProgressModal creates the progress modal for automatic retry
func (mm *ModalManager) createProgressModal(operation string) *tview.Modal {
	progressContent := fmt.Sprintf(
		"🔄 Retrying: %s\n\nPlease wait while we attempt to recover...",
		operation,
	)
	progressModal := mm.createModal(progressContent, []string{"Cancel"})

	// Add cancel functionality
	progressModal.SetDoneFunc(func(_ int, _ string) {
		mm.closeProgressModalAndRestoreFocus()
	})

	return progressModal
}

// closeProgressModalAndRestoreFocus closes the progress modal and restores focus
func (mm *ModalManager) closeProgressModalAndRestoreFocus() {
	pages, ok := mm.ui.GetPages().(*tview.Pages)
	if !ok {
		return
	}
	pages.RemovePage("retry_progress_modal")
	mm.restoreFocusToMainView()
}

// addProgressModalToPages adds the progress modal to the pages
func (mm *ModalManager) addProgressModalToPages(progressModal *tview.Modal) {
	pages, ok := mm.ui.GetPages().(*tview.Pages)
	if !ok {
		return
	}
	pages.AddPage("retry_progress_modal", progressModal, true, true)
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
	mm.closeProgressModal()

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

// closeProgressModal closes the progress modal
func (mm *ModalManager) closeProgressModal() {
	pages, ok := mm.ui.GetPages().(*tview.Pages)
	if !ok {
		return
	}
	pages.RemovePage("retry_progress_modal")
}

// createConfirmModal creates the confirmation modal
func (mm *ModalManager) createConfirmModal(text string, callback func(bool)) *tview.Modal {
	return tview.NewModal().
		SetText(text).
		AddButtons([]string{"Yes", "No"}).
		SetDoneFunc(func(buttonIndex int, _ string) {
			mm.handleConfirmButtonClick(buttonIndex, callback)
		})
}

// setupConfirmModalHandlers sets up the keyboard handlers for the confirm modal
func (mm *ModalManager) setupConfirmModalHandlers(modal *tview.Modal, callback func(bool)) {
	modal.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			mm.handleConfirmEscape(callback)
			return nil
		}
		return event
	})
}

// handleConfirmButtonClick handles button clicks in the confirm modal
func (mm *ModalManager) handleConfirmButtonClick(buttonIndex int, callback func(bool)) {
	mm.closeConfirmModal()
	callback(buttonIndex == 0)
	mm.restoreFocusToMainView()
}

// handleConfirmEscape handles ESC key press in the confirm modal
func (mm *ModalManager) handleConfirmEscape(callback func(bool)) {
	mm.closeConfirmModal()
	callback(false)
	mm.restoreFocusToMainView()
}

// closeConfirmModal closes the confirmation modal
func (mm *ModalManager) closeConfirmModal() {
	pages, ok := mm.ui.GetPages().(*tview.Pages)
	if !ok {
		return
	}
	pages.RemovePage("confirm_modal")
}

// showConfirmModal displays the confirmation modal
func (mm *ModalManager) showConfirmModal(modal *tview.Modal) {
	pages, ok := mm.ui.GetPages().(*tview.Pages)
	if !ok {
		return
	}
	pages.AddPage("confirm_modal", modal, true, true)

	app, ok := mm.ui.GetApp().(*tview.Application)
	if !ok {
		return
	}
	app.SetFocus(modal)
}

// closeInfoModalAndRestoreFocus closes the info modal and restores focus
func (mm *ModalManager) closeInfoModalAndRestoreFocus() {
	pages, ok := mm.ui.GetPages().(*tview.Pages)
	if !ok {
		return
	}
	pages.RemovePage("info_modal")
	mm.restoreFocusToMainView()
}

// showInfoModal shows the info modal and sets focus
func (mm *ModalManager) showInfoModal(modal *tview.Modal) {
	pages, ok := mm.ui.GetPages().(*tview.Pages)
	if !ok {
		return
	}

	pages.AddPage("info_modal", modal, true, true)

	// Set focus to the modal so it can receive keyboard input
	app, ok := mm.ui.GetApp().(*tview.Application)
	if !ok {
		return
	}

	app.SetFocus(modal)
}

// restoreFocusToMainView restores focus to the main view container
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

// closeScaleModalAndRestoreFocus closes the scale modal and restores focus
func (mm *ModalManager) closeScaleModalAndRestoreFocus() {
	pages, ok := mm.ui.GetPages().(*tview.Pages)
	if !ok {
		return
	}
	pages.RemovePage("scale_modal")
	mm.restoreFocusToMainView()
}

// addScaleModalToPages adds the scale modal to the pages
func (mm *ModalManager) addScaleModalToPages(flex *tview.Flex) {
	pages, ok := mm.ui.GetPages().(*tview.Pages)
	if !ok {
		return
	}
	pages.AddPage("scale_modal", flex, true, true)
}

// setFocusToForm sets focus to the form
func (mm *ModalManager) setFocusToForm(form *tview.Form) {
	app, ok := mm.ui.GetApp().(*tview.Application)
	if !ok {
		return
	}
	app.SetFocus(form)
}
