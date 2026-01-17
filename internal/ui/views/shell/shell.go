package shell

import (
	"context"
	"strings"

	"github.com/rivo/tview"
	"github.com/wikczerski/whaletui/internal/ui/interfaces"
)

// View represents an interactive shell view for container interaction
type View struct {
	ui              interfaces.UIInterface
	onExit          func()
	view            *tview.Flex
	outputView      *tview.TextView
	inputField      *tview.InputField
	execFunc        func(context.Context, string, []string, bool) (string, error)
	containerName   string
	containerID     string
	currentInput    string
	commandHistory  []string
	multiLineBuffer []string
	historyIndex    int
	isMultiLine     bool
}

// NewView creates a new shell view for the specified container
func NewView(
	ui interfaces.UIInterface,
	containerID, containerName string,
	onExit func(),
	execFunc func(context.Context, string, []string, bool) (string, error),
) *View {
	sv := &View{
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
func (sv *View) GetView() tview.Primitive {
	return sv.view
}

// GetContainerID returns the container ID
func (sv *View) GetContainerID() string {
	return sv.containerID
}

// GetContainerName returns the container name
func (sv *View) GetContainerName() string {
	return sv.containerName
}

// GetInputField returns the input field for external access
func (sv *View) GetInputField() *tview.InputField {
	return sv.inputField
}

// exitShell exits the shell view
func (sv *View) exitShell() {
	if sv.onExit != nil {
		sv.onExit()
	}
}

// addOutput adds text to the output view
func (sv *View) addOutput(text string) {
	currentText := sv.outputView.GetText(true)
	sv.outputView.SetText(currentText + text)
	sv.outputView.ScrollToEnd()
}

// isMultiLineCommand checks if the command is incomplete (ends with backslash)
func (sv *View) isMultiLineCommand(command string) bool {
	return strings.HasSuffix(strings.TrimSpace(command), "\\")
}

// addToMultiLineBuffer adds a line to the multi-line buffer
func (sv *View) addToMultiLineBuffer(line string) {
	line = strings.TrimSuffix(strings.TrimSpace(line), "\\")
	sv.multiLineBuffer = append(sv.multiLineBuffer, line)
}

// getMultiLineCommand combines all lines in the buffer
func (sv *View) getMultiLineCommand() string {
	return strings.Join(sv.multiLineBuffer, " ")
}

// clearMultiLineBuffer clears the multi-line buffer
func (sv *View) clearMultiLineBuffer() {
	sv.multiLineBuffer = nil
	sv.isMultiLine = false
}
