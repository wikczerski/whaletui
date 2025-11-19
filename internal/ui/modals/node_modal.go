package modals

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/wikczerski/whaletui/internal/ui/builders"
)

// NewNodeAvailabilityModal creates a new node availability modal
func NewNodeAvailabilityModal(
	nodeName, currentAvailability string,
	onConfirm func(string),
	onCancel func(),
) *tview.Modal {
	content := fmt.Sprintf(
		"Update Node Availability: %s\n\nCurrent Availability: %s\n\nSelect new availability:",
		nodeName,
		currentAvailability,
	)

	modal := builders.NewModalBuilder().
		SetText(content).
		AddButtons([]string{"Active", "Pause", "Drain", "Cancel"}).
		Build()

	modal.SetDoneFunc(func(_ int, buttonLabel string) {
		switch buttonLabel {
		case "Active":
			onConfirm("active")
		case "Pause":
			onConfirm("pause")
		case "Drain":
			onConfirm("drain")
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
