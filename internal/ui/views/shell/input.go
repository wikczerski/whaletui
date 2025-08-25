package shell

import (
	"github.com/gdamore/tcell/v2"
)

// handleInputCapture handles special key combinations
func (sv *View) handleInputCapture(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyEscape:
		return sv.handleEscapeKey()
	case tcell.KeyUp:
		return sv.handleUpKey()
	case tcell.KeyDown:
		return sv.handleDownKey()
	case tcell.KeyTab:
		return sv.handleTabKey()
	}
	return event
}

// handleEscapeKey handles the escape key press
func (sv *View) handleEscapeKey() *tcell.EventKey {
	sv.exitShell()
	return nil
}

// handleUpKey handles the up arrow key press
func (sv *View) handleUpKey() *tcell.EventKey {
	if sv.historyIndex == len(sv.commandHistory) {
		sv.currentInput = sv.inputField.GetText()
	}
	sv.navigateHistory(1)
	return nil
}

// handleDownKey handles the down arrow key press
func (sv *View) handleDownKey() *tcell.EventKey {
	sv.navigateHistory(-1)
	return nil
}

// handleTabKey handles the tab key press
func (sv *View) handleTabKey() *tcell.EventKey {
	sv.handleTabCompletion()
	return nil
}
