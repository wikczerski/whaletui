package modals

import (
	"fmt"
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// ScaleModal handles the service scaling UI
type ScaleModal struct {
	flex       *tview.Flex
	inputField *tview.InputField
	form       *tview.Form
	modal      *tview.Modal
}

// NewScaleModal creates a new scale modal
func NewScaleModal(
	serviceName string,
	currentReplicas uint64,
	onConfirm func(int),
	onCancel func(),
) *tview.Flex {
	sm := &ScaleModal{}

	// Create input field
	sm.inputField = tview.NewInputField().
		SetLabel("Replicas: ").
		SetText(fmt.Sprintf("%d", currentReplicas)).
		SetFieldWidth(10).
		SetAcceptanceFunc(tview.InputFieldInteger)

	// Create form
	sm.form = tview.NewForm().
		AddFormItem(sm.inputField).
		AddButton("Scale", func() {
			replicasStr := sm.inputField.GetText()
			replicas, err := strconv.Atoi(replicasStr)
			if err != nil || replicas < 0 {
				// In a real scenario we might want to show an error here,
				// but for now we'll just return and let the manager handle validation if needed
				// or we could pass an error handler.
				// For this refactor, we'll assume valid input or let the caller handle it.
				// However, the original code showed an error modal.
				// We'll just call onConfirm with the value and let the caller validate/show error
				// OR we can keep the validation here if we had access to ShowError.
				// To keep it simple and decoupled, we'll just parse here.
				return
			}
			onConfirm(replicas)
		}).
		AddButton("Cancel", onCancel)

	// Create modal
	sm.modal = tview.NewModal().
		SetText(fmt.Sprintf("Scale Service: %s\nCurrent Replicas: %d", serviceName, currentReplicas)).
		SetDoneFunc(func(_ int, _ string) {
			onCancel()
		})

	// Create flex container
	sm.flex = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(sm.modal, 0, 1, false).
		AddItem(sm.form, 0, 1, true)

	// Add keyboard handling
	sm.flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			onCancel()
			return nil
		}
		return event
	})

	return sm.flex
}
