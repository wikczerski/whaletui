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

	sm.inputField = sm.createInputField(currentReplicas)
	sm.form = sm.createForm(onConfirm, onCancel)
	sm.modal = sm.createModal(serviceName, currentReplicas, onCancel)
	sm.flex = sm.createLayout(onCancel)

	return sm.flex
}

func (sm *ScaleModal) createInputField(currentReplicas uint64) *tview.InputField {
	return tview.NewInputField().
		SetLabel("Replicas: ").
		SetText(fmt.Sprintf("%d", currentReplicas)).
		SetFieldWidth(10).
		SetAcceptanceFunc(tview.InputFieldInteger)
}

func (sm *ScaleModal) createForm(onConfirm func(int), onCancel func()) *tview.Form {
	return tview.NewForm().
		AddFormItem(sm.inputField).
		AddButton("Scale", func() {
			replicasStr := sm.inputField.GetText()
			replicas, err := strconv.Atoi(replicasStr)
			if err != nil || replicas < 0 {
				return
			}
			onConfirm(replicas)
		}).
		AddButton("Cancel", onCancel)
}

func (sm *ScaleModal) createModal(serviceName string, currentReplicas uint64, onCancel func()) *tview.Modal {
	return tview.NewModal().
		SetText(fmt.Sprintf("Scale Service: %s\nCurrent Replicas: %d", serviceName, currentReplicas)).
		SetDoneFunc(func(_ int, _ string) {
			onCancel()
		})
}

func (sm *ScaleModal) createLayout(onCancel func()) *tview.Flex {
	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(sm.modal, 0, 1, false).
		AddItem(sm.form, 0, 1, true)

	flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			onCancel()
			return nil
		}
		return event
	})

	return flex
}
