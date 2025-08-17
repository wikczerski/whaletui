package shell

import (
	"context"
	"strings"

	"github.com/rivo/tview"
	"github.com/wikczerski/D5r/internal/ui/interfaces"
)

// ShellView represents an interactive shell view for container interaction
type ShellView struct {
	containerID   string
	containerName string
	ui            interfaces.UIInterface
	view          *tview.Flex
	outputView    *tview.TextView
	inputField    *tview.InputField
	onExit        func()
	execFunc      func(context.Context, string, []string, bool) (string, error)

	commandHistory []string
	historyIndex   int
	currentInput   string

	multiLineBuffer []string
	isMultiLine     bool
}

// NewShellView creates a new shell view for the specified container
func NewShellView(ui interfaces.UIInterface, containerID, containerName string, onExit func(), execFunc func(context.Context, string, []string, bool) (string, error)) *ShellView {
	sv := &ShellView{
		containerID:   containerID,
		containerName: containerName,
		ui:            ui,
		onExit:        onExit,
		execFunc:      execFunc,
	}

	sv.createView()
	return sv
}

// GetView returns the shell view primitive
func (sv *ShellView) GetView() tview.Primitive {
	return sv.view
}

// GetContainerID returns the container ID
func (sv *ShellView) GetContainerID() string {
	return sv.containerID
}

// GetContainerName returns the container name
func (sv *ShellView) GetContainerName() string {
	return sv.containerName
}

// GetInputField returns the input field for external access
func (sv *ShellView) GetInputField() *tview.InputField {
	return sv.inputField
}

// exitShell exits the shell view
func (sv *ShellView) exitShell() {
	if sv.onExit != nil {
		sv.onExit()
	}
}

// addOutput adds text to the output view
func (sv *ShellView) addOutput(text string) {
	currentText := sv.outputView.GetText(true)
	sv.outputView.SetText(currentText + text)
	sv.outputView.ScrollToEnd()
}

// isMultiLineCommand checks if the command is incomplete (ends with backslash)
func (sv *ShellView) isMultiLineCommand(command string) bool {
	return strings.HasSuffix(strings.TrimSpace(command), "\\")
}

// addToMultiLineBuffer adds a line to the multi-line buffer
func (sv *ShellView) addToMultiLineBuffer(line string) {
	line = strings.TrimSuffix(strings.TrimSpace(line), "\\")
	sv.multiLineBuffer = append(sv.multiLineBuffer, line)
}

// getMultiLineCommand combines all lines in the buffer
func (sv *ShellView) getMultiLineCommand() string {
	return strings.Join(sv.multiLineBuffer, " ")
}

// clearMultiLineBuffer clears the multi-line buffer
func (sv *ShellView) clearMultiLineBuffer() {
	sv.multiLineBuffer = nil
	sv.isMultiLine = false
}
