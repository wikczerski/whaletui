package modals

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/wikczerski/whaletui/internal/ui/builders"
)

// NewInfoModal creates a new info modal
func NewInfoModal(text string, onDone func()) *tview.Modal {
	modal := builders.NewModalBuilder().
		SetText(text).
		AddButtons([]string{"OK"}).
		Build()

	modal.SetDoneFunc(func(_ int, _ string) {
		onDone()
	})

	modal.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			onDone()
			return nil
		}
		return event
	})

	return modal
}
