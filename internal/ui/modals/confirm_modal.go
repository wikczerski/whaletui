package modals

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/wikczerski/whaletui/internal/ui/builders"
)

// NewConfirmModal creates a new confirmation modal
func NewConfirmModal(text string, onConfirm func(bool)) *tview.Modal {
	modal := builders.NewModalBuilder().
		SetText(text).
		AddButtons([]string{"Yes", "No"}).
		Build()

	modal.SetDoneFunc(func(buttonIndex int, _ string) {
		onConfirm(buttonIndex == 0)
	})

	modal.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			onConfirm(false)
			return nil
		}
		return event
	})

	return modal
}
