package ui

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// handleGlobalKeyBindings handles all global key bindings
func (ui *UI) handleGlobalKeyBindings(event *tcell.EventKey) *tcell.EventKey {
	if ui.shouldSkipGlobalKeyBindings() {
		return event
	}

	return ui.routeKeyBinding(event)
}

// shouldSkipGlobalKeyBindings checks if global key bindings should be skipped
func (ui *UI) shouldSkipGlobalKeyBindings() bool {
	return ui.IsModalActive() || ui.isShellInputFieldFocused()
}

// routeKeyBinding routes the key binding to the appropriate handler
func (ui *UI) routeKeyBinding(event *tcell.EventKey) *tcell.EventKey {
	currentFocus := ui.app.GetFocus()
	commandInput := ui.commandHandler.GetInput()
	searchInput := ui.searchHandler.GetInput()

	if currentFocus == commandInput {
		return ui.handleCommandModeKeyBindings(event)
	}

	if currentFocus == searchInput {
		return ui.handleSearchModeKeyBindings(event)
	}

	if ui.isExecCommandInputActive() {
		return ui.handleExecCommandKeyBindings(event)
	}

	if ui.isShellViewActive() {
		return ui.handleShellViewKeyBindings(event)
	}

	return ui.handleNormalModeKeyBindings(event)
}

// handleSearchModeKeyBindings handles key bindings when in search mode
func (ui *UI) handleSearchModeKeyBindings(event *tcell.EventKey) *tcell.EventKey {
	// Handle special keys
	switch event.Key() {
	case tcell.KeyEnter:
		ui.searchHandler.HandleInput(event.Key())
		return nil
	case tcell.KeyEscape:
		ui.searchHandler.HandleInput(event.Key())
		return nil
	case tcell.KeyRune, tcell.KeyBackspace, tcell.KeyBackspace2:
		// For typing and backspace, let the input field handle it first
		// Then manually trigger search with the updated text
		result := event
		// Let the input field process the key first
		if ui.searchHandler.GetInput() != nil {
			// The input field will update its text, then we trigger search
			go func() {
				// Small delay to let the input field update
				time.Sleep(1 * time.Millisecond)
				text := ui.searchHandler.GetInput().GetText()
				ui.log.Info("Key binding goroutine calling ProcessSearch", "text", text)
				ui.searchHandler.ProcessSearch(text)
			}()
		}
		return result
	default:
		// Let the input field handle other keys
		return event
	}
}

// handleCommandModeKeyBindings handles key bindings when in command mode
func (ui *UI) handleCommandModeKeyBindings(event *tcell.EventKey) *tcell.EventKey {
	// In command mode, only allow ESC to exit
	if event.Key() == tcell.KeyEscape {
		ui.commandHandler.HandleInput(tcell.KeyEscape)
		return nil
	}
	return event
}

// handleExecCommandKeyBindings handles key bindings when exec command input is active
func (ui *UI) handleExecCommandKeyBindings(event *tcell.EventKey) *tcell.EventKey {
	// In exec command mode, only allow ESC to exit
	if event.Key() != tcell.KeyEscape {
		return event
	}

	// Remove the exec input and return focus to view
	ui.removeExecInputAndRestoreFocus()
	return nil
}

// removeExecInputAndRestoreFocus removes the exec input and restores focus to the view container
func (ui *UI) removeExecInputAndRestoreFocus() {
	if ui.mainFlex == nil {
		return
	}

	focused := ui.app.GetFocus()
	if focused == nil {
		return
	}

	ui.removeInputFieldAndRestoreFocus(focused)
}

// removeInputFieldAndRestoreFocus removes the input field and restores focus
func (ui *UI) removeInputFieldAndRestoreFocus(focused tview.Primitive) {
	inputField, ok := focused.(*tview.InputField)
	if !ok {
		return
	}

	ui.mainFlex.RemoveItem(inputField)
	if ui.viewContainer != nil {
		ui.app.SetFocus(ui.viewContainer)
	}
}

// handleShellViewKeyBindings handles key bindings when shell view is active
func (ui *UI) handleShellViewKeyBindings(event *tcell.EventKey) *tcell.EventKey {
	// In shell mode, only allow ESC to exit (handled by shell view)
	if event.Key() == tcell.KeyEscape {
		return event // Let shell view handle ESC
	}
	// Block other global key bindings in shell mode
	return event
}

// handleNormalModeKeyBindings handles key bindings in normal mode
func (ui *UI) handleNormalModeKeyBindings(event *tcell.EventKey) *tcell.EventKey {
	// Log what component is focused
	if focused := ui.app.GetFocus(); focused != nil {
		ui.log.Info(
			"Normal mode key binding",
			"key",
			event.Key(),
			"focusedType",
			fmt.Sprintf("%T", focused),
		)
	}

	switch event.Key() {
	case tcell.KeyRune:
		// Only handle global rune keys, let others pass through to focused component
		return ui.handleGlobalRuneKeyBindings(event)
	case tcell.KeyCtrlC:
		return ui.handleCtrlCKeyBinding(event)
	case tcell.KeyBackspace:
		// Handle Backspace to go back from subviews
		return ui.handleBackspaceKeyBinding(event)
	}
	return event
}

// handleGlobalRuneKeyBindings handles only global rune key bindings
func (ui *UI) handleGlobalRuneKeyBindings(event *tcell.EventKey) *tcell.EventKey {
	ui.log.Info("Global rune key handler called", "key", string(event.Rune()))

	if ui.handleQuitKey(event) || ui.handleCommandModeKey(event) {
		return nil
	}

	ui.log.Info("Global handler passing through key", "key", string(event.Rune()))
	// Let all other keys pass through to the focused view (including action keys)
	return event
}

// handleQuitKey handles quit key presses
func (ui *UI) handleQuitKey(event *tcell.EventKey) bool {
	if event.Rune() == 'q' || event.Rune() == 'Q' {
		ui.log.Info("Quitting application...")
		// Send shutdown signal instead of direct exit to ensure cleanup
		select {
		case ui.shutdownChan <- struct{}{}:
		default:
		}
		return true
	}
	return false
}

// handleCommandModeKey handles command mode key presses
func (ui *UI) handleCommandModeKey(event *tcell.EventKey) bool {
	if event.Rune() == ':' {
		ui.log.Info("Entering command mode")
		ui.commandHandler.Enter()
		return true
	}
	if event.Rune() == '/' {
		ui.log.Info("Entering search mode")
		ui.searchHandler.Enter()
		return true
	}
	return false
}

// handleCtrlCKeyBinding handles Ctrl+C key binding
func (ui *UI) handleCtrlCKeyBinding(_ *tcell.EventKey) *tcell.EventKey {
	ui.log.Info("Received Ctrl+C, shutting down...")
	select {
	case ui.shutdownChan <- struct{}{}:
	default:
	}
	return nil
}

// handleBackspaceKeyBinding handles Backspace key binding for subview navigation
func (ui *UI) handleBackspaceKeyBinding(_ *tcell.EventKey) *tcell.EventKey {
	// Only handle Backspace when in details mode or logs mode, but NOT in shell mode
	if (ui.inDetailsMode || ui.inLogsMode) && !ui.isShellViewActive() {
		ui.log.Info("Backspace pressed in subview, returning to main view")
		ui.ShowCurrentView()
		return nil // Consume the event
	}
	// If not in a subview or in shell mode, let the event pass through
	return nil
}

// isExecCommandInputActive checks if exec command input is currently active
func (ui *UI) isExecCommandInputActive() bool {
	if focused := ui.app.GetFocus(); focused != nil {
		if inputField, ok := focused.(*tview.InputField); ok {
			return inputField.GetLabel() == " Exec Command: "
		}
	}
	return false
}

// isShellViewActive checks if shell view is currently active
func (ui *UI) isShellViewActive() bool {
	return ui.shellView != nil && ui.app.GetFocus() == ui.shellView.GetView()
}

// isShellInputFieldFocused checks if shell input field is currently focused
func (ui *UI) isShellInputFieldFocused() bool {
	if ui.shellView != nil {
		if focused := ui.app.GetFocus(); focused != nil {
			if inputField, ok := focused.(*tview.InputField); ok {
				// Check if this input field belongs to the shell view
				return inputField == ui.shellView.GetInputField()
			}
		}
	}
	return false
}
