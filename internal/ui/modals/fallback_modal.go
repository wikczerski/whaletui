package modals

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/wikczerski/whaletui/internal/ui/builders"
)

// NewFallbackModal creates a new fallback modal
func NewFallbackModal(
	operation string,
	err error,
	options []string,
	onFallback func(string),
	onCancel func(),
) *tview.Modal {
	content := fmt.Sprintf(
		"⚠️  Operation Failed: %s\n\nError: %v\n\nAlternative operations are available:",
		operation,
		err,
	)

	// Create buttons list with Cancel at the end
	buttons := make([]string, len(options)+1)
	copy(buttons, options)
	buttons[len(options)] = "Cancel"

	modal := builders.NewModalBuilder().
		SetText(content).
		AddButtons(buttons).
		Build()

	modal.SetDoneFunc(func(_ int, buttonLabel string) {
		if buttonLabel == "Cancel" {
			onCancel()
		} else {
			onFallback(buttonLabel)
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
