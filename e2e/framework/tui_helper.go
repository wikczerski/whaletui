package framework

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
)

// TUIHelper provides TUI interaction utilities for testing.
type TUIHelper struct {
	fw *TestFramework
}

// NewTUIHelper creates a new TUI helper.
func NewTUIHelper(fw *TestFramework) *TUIHelper {
	return &TUIHelper{fw: fw}
}

// NavigateToView navigates to a specific view using command mode.
func (th *TUIHelper) NavigateToView(viewName string) {
	th.fw.t.Helper()

	// Enter command mode
	th.fw.SimulateKeyRune(':')
	th.fw.Sleep(200 * time.Millisecond)

	// Type view name
	for _, ch := range viewName {
		th.fw.SimulateKeyRune(ch)
		th.fw.Sleep(50 * time.Millisecond)
	}

	// Press Enter
	th.fw.SimulateKeyPress(tcell.KeyEnter)
	th.fw.Sleep(300 * time.Millisecond)
}

// PressKey simulates pressing a single key.
func (th *TUIHelper) PressKey(key rune) {
	th.fw.t.Helper()

	th.fw.SimulateKeyRune(key)
	th.fw.Sleep(200 * time.Millisecond)
}

// PressSpecialKey simulates pressing a special key.
func (th *TUIHelper) PressSpecialKey(key tcell.Key) {
	th.fw.t.Helper()

	th.fw.SimulateKeyPress(key)
	th.fw.Sleep(200 * time.Millisecond)
}

// EnterSearchMode enters search mode.
func (th *TUIHelper) EnterSearchMode() {
	th.fw.t.Helper()

	th.fw.SimulateKeyRune('/')
	th.fw.Sleep(200 * time.Millisecond)
}

// TypeText types a string of text.
func (th *TUIHelper) TypeText(text string) {
	th.fw.t.Helper()

	for _, ch := range text {
		th.fw.SimulateKeyRune(ch)
		th.fw.Sleep(50 * time.Millisecond)
	}
}

// PressEnter presses the Enter key.
func (th *TUIHelper) PressEnter() {
	th.fw.t.Helper()

	th.fw.SimulateKeyPress(tcell.KeyEnter)
	th.fw.Sleep(200 * time.Millisecond)
}

// PressEscape presses the Escape key.
func (th *TUIHelper) PressEscape() {
	th.fw.t.Helper()

	th.fw.SimulateKeyPress(tcell.KeyEscape)
	th.fw.Sleep(200 * time.Millisecond)
}

// PressBackspace presses the Backspace key.
func (th *TUIHelper) PressBackspace() {
	th.fw.t.Helper()

	th.fw.SimulateKeyPress(tcell.KeyBackspace)
	th.fw.Sleep(200 * time.Millisecond)
}

// SelectTableRow selects a row in a table by index.
func (th *TUIHelper) SelectTableRow(index int) {
	th.fw.t.Helper()

	// Navigate to row using arrow keys
	for i := 0; i < index; i++ {
		th.fw.SimulateKeyPress(tcell.KeyDown)
		th.fw.Sleep(100 * time.Millisecond)
	}
}

// NavigateTableDown moves down in the table.
func (th *TUIHelper) NavigateTableDown(count int) {
	th.fw.t.Helper()

	for i := 0; i < count; i++ {
		th.fw.SimulateKeyPress(tcell.KeyDown)
		th.fw.Sleep(100 * time.Millisecond)
	}
}

// NavigateTableUp moves up in the table.
func (th *TUIHelper) NavigateTableUp(count int) {
	th.fw.t.Helper()

	for i := 0; i < count; i++ {
		th.fw.SimulateKeyPress(tcell.KeyUp)
		th.fw.Sleep(100 * time.Millisecond)
	}
}

// ConfirmModal confirms a modal dialog by pressing Yes/OK.
func (th *TUIHelper) ConfirmModal() {
	th.fw.t.Helper()

	// Typically Enter or first button
	th.fw.SimulateKeyPress(tcell.KeyEnter)
	th.fw.Sleep(200 * time.Millisecond)
}

// CancelModal cancels a modal dialog by pressing No/Cancel.
func (th *TUIHelper) CancelModal() {
	th.fw.t.Helper()

	// Typically Tab then Enter, or Escape
	th.fw.SimulateKeyPress(tcell.KeyTab)
	th.fw.Sleep(100 * time.Millisecond)
	th.fw.SimulateKeyPress(tcell.KeyEnter)
	th.fw.Sleep(200 * time.Millisecond)
}

// CloseModalWithEscape closes a modal using Escape key.
func (th *TUIHelper) CloseModalWithEscape() {
	th.fw.t.Helper()

	th.fw.SimulateKeyPress(tcell.KeyEscape)
	th.fw.Sleep(200 * time.Millisecond)
}

// WaitForTableUpdate waits for table to update with new data.
func (th *TUIHelper) WaitForTableUpdate(timeout time.Duration) {
	th.fw.t.Helper()

	// Give time for async data loading
	th.fw.Sleep(timeout)
}

// GetFocusedPrimitive returns the currently focused primitive.
func (th *TUIHelper) GetFocusedPrimitive() tview.Primitive {
	if th.fw.tviewApp == nil {
		return nil
	}
	return th.fw.tviewApp.GetFocus()
}

// AssertFocusedType asserts the type of the focused primitive.
func (th *TUIHelper) AssertFocusedType(expectedType string) {
	th.fw.t.Helper()

	focused := th.GetFocusedPrimitive()
	assert.NotNil(th.fw.t, focused, "No primitive is focused")

	actualType := fmt.Sprintf("%T", focused)
	assert.Contains(th.fw.t, actualType, expectedType, "Focused primitive type mismatch")
}

// WaitForFocus waits for a specific primitive type to be focused.
func (th *TUIHelper) WaitForFocus(expectedType string, timeout time.Duration) {
	th.fw.t.Helper()

	th.fw.WaitForCondition(func() bool {
		focused := th.GetFocusedPrimitive()
		if focused == nil {
			return false
		}
		actualType := fmt.Sprintf("%T", focused)
		return actualType == expectedType
	}, timeout, fmt.Sprintf("focus to be on %s", expectedType))
}

// ExecuteContainerAction executes a container action (start, stop, restart, delete, etc.).
func (th *TUIHelper) ExecuteContainerAction(action rune, confirm bool) {
	th.fw.t.Helper()

	// Press action key
	th.PressKey(action)

	// Handle confirmation if needed
	if confirm {
		th.ConfirmModal()
	}
}

// SearchInView searches for text in the current view.
func (th *TUIHelper) SearchInView(searchText string) {
	th.fw.t.Helper()

	// Enter search mode
	th.EnterSearchMode()

	// Type search text
	th.TypeText(searchText)

	// Wait for search to filter
	th.fw.Sleep(300 * time.Millisecond)
}

// ClearSearch clears the current search.
func (th *TUIHelper) ClearSearch() {
	th.fw.t.Helper()

	// Press Escape to clear search
	th.PressEscape()
	th.fw.Sleep(200 * time.Millisecond)
}

// ExecuteCommand executes a command in command mode.
func (th *TUIHelper) ExecuteCommand(command string) {
	th.fw.t.Helper()

	// Enter command mode
	th.fw.SimulateKeyRune(':')
	th.fw.Sleep(200 * time.Millisecond)

	// Type command
	th.TypeText(command)

	// Press Enter
	th.PressEnter()
}

// ViewDetails opens the details view for the selected item.
func (th *TUIHelper) ViewDetails() {
	th.fw.t.Helper()

	th.PressEnter()
	th.fw.Sleep(300 * time.Millisecond)
}

// CloseDetails closes the details view.
func (th *TUIHelper) CloseDetails() {
	th.fw.t.Helper()

	th.PressBackspace()
	th.fw.Sleep(200 * time.Millisecond)
}

// ViewLogs opens the logs view for the selected container.
func (th *TUIHelper) ViewLogs() {
	th.fw.t.Helper()

	th.PressKey('l')
	th.fw.Sleep(300 * time.Millisecond)
}

// CloseLogs closes the logs view.
func (th *TUIHelper) CloseLogs() {
	th.fw.t.Helper()

	th.PressEscape()
	th.fw.Sleep(200 * time.Millisecond)
}

// AttachToShell attaches to container shell.
func (th *TUIHelper) AttachToShell() {
	th.fw.t.Helper()

	th.PressKey('a')
	th.fw.Sleep(500 * time.Millisecond)
}

// ExitShell exits the container shell.
func (th *TUIHelper) ExitShell() {
	th.fw.t.Helper()

	// Type exit command
	th.TypeText("exit")
	th.PressEnter()
	th.fw.Sleep(300 * time.Millisecond)
}

// ExecuteShellCommand executes a command in the shell.
func (th *TUIHelper) ExecuteShellCommand(command string) {
	th.fw.t.Helper()

	th.TypeText(command)
	th.PressEnter()
	th.fw.Sleep(500 * time.Millisecond)
}

// OpenExecInput opens the exec command input.
func (th *TUIHelper) OpenExecInput() {
	th.fw.t.Helper()

	th.PressKey('e')
	th.fw.Sleep(300 * time.Millisecond)
}

// ExecuteExecCommand executes a command via exec.
func (th *TUIHelper) ExecuteExecCommand(command string) {
	th.fw.t.Helper()

	th.OpenExecInput()
	th.TypeText(command)
	th.PressEnter()
	th.fw.Sleep(500 * time.Millisecond)
}

// CloseExecOutput closes the exec output modal.
func (th *TUIHelper) CloseExecOutput() {
	th.fw.t.Helper()

	th.PressEnter() // Close button
	th.fw.Sleep(200 * time.Millisecond)
}

// InspectItem inspects the selected item.
func (th *TUIHelper) InspectItem() {
	th.fw.t.Helper()

	th.PressKey('i')
	th.fw.Sleep(300 * time.Millisecond)
}

// CloseInspect closes the inspect view.
func (th *TUIHelper) CloseInspect() {
	th.fw.t.Helper()

	th.PressEscape()
	th.fw.Sleep(200 * time.Millisecond)
}

// ScaleService opens the scale modal for a service.
func (th *TUIHelper) ScaleService(replicas int) {
	th.fw.t.Helper()

	// Press scale key
	th.PressKey('s')
	th.fw.Sleep(300 * time.Millisecond)

	// Type replica count
	th.TypeText(fmt.Sprintf("%d", replicas))

	// Confirm
	th.PressEnter()
	th.fw.Sleep(500 * time.Millisecond)
}

// RefreshView refreshes the current view.
func (th *TUIHelper) RefreshView() {
	th.fw.t.Helper()

	th.PressKey('f')
	th.fw.Sleep(500 * time.Millisecond)
}
