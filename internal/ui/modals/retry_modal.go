package modals

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/wikczerski/whaletui/internal/ui/builders"
)

// NewRetryModal creates a new retry modal
func NewRetryModal(
	operation string,
	err error,
	onRetry func(),
	onAutoRetry func(),
	onCancel func(),
) *tview.Modal {
	content := fmt.Sprintf(
		"🔄 Operation Failed: %s\n\nError: %v\n\nThis may be a temporary issue. Would you like to retry?",
		operation,
		err,
	)

	modal := builders.NewModalBuilder().
		SetText(content).
		AddButtons([]string{"Retry", "Retry (Auto)", "Cancel"}).
		Build()

	modal.SetDoneFunc(func(_ int, buttonLabel string) {
		switch buttonLabel {
		case "Retry":
			onRetry()
		case "Retry (Auto)":
			onAutoRetry()
		case "Cancel":
			onCancel()
		}
	})

	modal.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			onCancel()
			return nil
		}
		return event
	})

	return modal
}

// NewProgressModal creates a new progress modal for automatic retries
func NewProgressModal(
	operation string,
	onCancel func(),
) *tview.Modal {
	content := fmt.Sprintf(
		"🔄 Retrying: %s\n\nPlease wait while we attempt to recover...",
		operation,
	)

	modal := builders.NewModalBuilder().
		SetText(content).
		AddButtons([]string{"Cancel"}).
		Build()

	modal.SetDoneFunc(func(_ int, _ string) {
		onCancel()
	})

	return modal
}
