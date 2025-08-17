package shell

import (
	"github.com/gdamore/tcell/v2"
)

// handleInputCapture handles special key combinations
func (sv *View) handleInputCapture(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyEscape:
		sv.exitShell()
		return nil
	case tcell.KeyUp:
		if sv.historyIndex == len(sv.commandHistory) {
			sv.currentInput = sv.inputField.GetText()
		}
		sv.navigateHistory(1)
		return nil
	case tcell.KeyDown:
		sv.navigateHistory(-1)
		return nil
	case tcell.KeyTab:
		sv.handleTabCompletion()
		return nil
	}
	return event
}
