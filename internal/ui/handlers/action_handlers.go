// Package handlers provides UI event handlers for WhaleTUI.
package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/wikczerski/whaletui/internal/config"
	"github.com/wikczerski/whaletui/internal/ui/builders"
	"github.com/wikczerski/whaletui/internal/ui/interfaces"
)

// ActionHandlers provides common action handling patterns for different resource types
type ActionHandlers struct {
	ui       interfaces.UIInterface
	executor *OperationExecutor
}

// NewActionHandlers creates a new action handlers helper
func NewActionHandlers(ui interfaces.UIInterface) *ActionHandlers {
	return &ActionHandlers{
		ui:       ui,
		executor: NewOperationExecutor(ui),
	}
}

// HandleDeleteAction provides a common delete action pattern
func (ah *ActionHandlers) HandleDeleteAction(
	resourceType, resourceID, resourceName string,
	deleteFunc func(context.Context, string, bool) error,
	onRefresh func(),
) {
	ah.executor.DeleteOperation(resourceType, resourceID, resourceName, deleteFunc, onRefresh)
}

// HandleInspectAction provides a common inspect action pattern
func (ah *ActionHandlers) HandleInspectAction(
	resourceType, resourceID string,
	inspectFunc func(context.Context, string) (map[string]any, error),
) {
	inspectView, inspectFlex := builders.CreateInspectView(fmt.Sprintf("Inspect: %s", resourceID))

	ah.setupInspectCloseButton(inspectFlex)
	ah.addInspectPage(inspectFlex)
	ah.loadInspectDataAsync(resourceType, resourceID, inspectView, inspectFunc)
}

// HandleResourceAction provides generic resource action handling (for images, volumes, networks)
func (ah *ActionHandlers) HandleResourceAction(
	action rune,
	resourceType, resourceID, resourceName string,
	inspectFunc func(context.Context, string) (map[string]any, error),
	deleteFunc func(context.Context, string, bool) error,
	onRefresh func(),
) {
	switch action {
	case 'd':
		if deleteFunc != nil {
			ah.HandleDeleteAction(resourceType, resourceID, resourceName, deleteFunc, onRefresh)
		}
	case 'i':
		if inspectFunc != nil {
			ah.HandleInspectAction(resourceType, resourceID, inspectFunc)
		}
	}
}

// HandleContainerAction provides container-specific action handling
func (ah *ActionHandlers) HandleContainerAction(
	action rune,
	containerID, containerName string,
	containerService any,
	onRefresh func(),
) {
	if ah.handleContainerLifecycleAction(action, containerID, containerService, onRefresh) {
		return
	}

	if ah.handleContainerManagementAction(
		action,
		containerID,
		containerName,
		containerService,
		onRefresh,
	) {
		return
	}

	ah.handleContainerAccessAction(action, containerID, containerName, containerService)
}

// HandleContainerLifecycleAction handles container lifecycle operations (start, stop, restart)
func (ah *ActionHandlers) HandleContainerLifecycleAction(
	action rune,
	containerID string,
	containerService any,
	onRefresh func(),
) {
	switch action {
	case 's':
		if cs, ok := containerService.(interfaces.ContainerService); ok {
			ah.executor.StartOperation("container", containerID, cs.StartContainer, onRefresh)
		}
	case 'S':
		if cs, ok := containerService.(interfaces.ContainerService); ok {
			ah.executor.StopOperation("container", containerID, cs.StopContainer, onRefresh)
		}
	case 'r':
		if cs, ok := containerService.(interfaces.ContainerService); ok {
			ah.executor.RestartOperation("container", containerID, cs.RestartContainer, onRefresh)
		}
	}
}

// HandleContainerManagementAction handles container management operations (delete, inspect)
func (ah *ActionHandlers) HandleContainerManagementAction(
	action rune,
	containerID, containerName string,
	containerService any,
	onRefresh func(),
) {
	switch action {
	case 'd':
		if cs, ok := containerService.(interfaces.ContainerService); ok {
			ah.HandleDeleteAction(
				"container",
				containerID,
				containerName,
				cs.RemoveContainer,
				onRefresh,
			)
		}
	case 'i':
		if cs, ok := containerService.(interfaces.ContainerService); ok {
			ah.HandleInspectAction("container", containerID, cs.InspectContainer)
		}
	}
}

// HandleContainerAccessAction handles container access operations (attach, logs, exec)
func (ah *ActionHandlers) HandleContainerAccessAction(
	action rune,
	containerID, containerName string,
	containerService any,
) {
	switch action {
	case 'a':
		ah.HandleAttachAction(containerID, containerName)
	case 'l':
		ah.ui.ShowLogs(containerID, containerName)
	case 'e':
		if cs, ok := containerService.(interfaces.ContainerService); ok {
			ah.HandleExecAction(containerID, containerName, cs.ExecContainer)
		}
	}
}

// HandleAttachAction handles container attach action
func (ah *ActionHandlers) HandleAttachAction(containerID, containerName string) {
	ah.ui.ShowShell(containerID, containerName)
}

// HandleExecAction handles container exec action
func (ah *ActionHandlers) HandleExecAction(
	containerID, containerName string,
	execFunc func(context.Context, string, []string, bool) (string, error),
) {
	mainFlex, ok := ah.ui.GetMainFlex().(*tview.Flex)
	if !ok {
		ah.ui.ShowError(errors.New("failed to get main flex container"))
		return
	}
	execInput := ah.createExecInput(containerName)

	ah.setupExecInputHandlers(execInput, mainFlex, containerID, containerName, execFunc)
	ah.addExecInputToUI(execInput, mainFlex)
}

// setupInspectCloseButton configures the close button for the inspect view
func (ah *ActionHandlers) setupInspectCloseButton(inspectFlex *tview.Flex) {
	inspectFlex.GetItem(1).(*tview.Button).SetSelectedFunc(func() {
		pages, ok := ah.ui.GetPages().(*tview.Pages)
		if !ok {
			ah.ui.ShowError(errors.New("failed to get pages container"))
			return
		}
		pages.RemovePage("inspect")
	})
}

// addInspectPage adds the inspect page to the UI
func (ah *ActionHandlers) addInspectPage(inspectFlex *tview.Flex) {
	pages, ok := ah.ui.GetPages().(*tview.Pages)
	if !ok {
		ah.ui.ShowError(errors.New("failed to get pages container"))
		return
	}
	pages.AddPage("inspect", inspectFlex, true, true)
}

// loadInspectDataAsync loads the inspect data asynchronously and updates the UI
func (ah *ActionHandlers) loadInspectDataAsync(
	resourceType, resourceID string,
	inspectView *tview.TextView,
	inspectFunc func(context.Context, string) (map[string]any, error),
) {
	go func() {
		inspectData, err := inspectFunc(context.Background(), resourceID)
		app, ok := ah.ui.GetApp().(*tview.Application)
		if !ok {
			ah.ui.ShowError(errors.New("failed to get application instance"))
			return
		}
		app.QueueUpdateDraw(func() {
			ah.updateInspectView(resourceType, inspectView, inspectData, err)
		})
	}()
}

// updateInspectView updates the inspect view with the loaded data or error
func (ah *ActionHandlers) updateInspectView(
	resourceType string,
	inspectView *tview.TextView,
	inspectData map[string]any,
	err error,
) {
	if err != nil {
		inspectView.SetText(fmt.Sprintf("%s inspection failed: %v", resourceType, err))
		return
	}

	data, jsonErr := json.MarshalIndent(inspectData, "", "  ")
	if jsonErr != nil {
		inspectView.SetText(fmt.Sprintf("Failed to format %s data: %v", resourceType, jsonErr))
		return
	}

	inspectView.SetText(string(data))
}

// handleContainerLifecycleAction handles container lifecycle operations (start, stop, restart)
func (ah *ActionHandlers) handleContainerLifecycleAction(
	action rune,
	containerID string,
	containerService any,
	onRefresh func(),
) bool {
	switch action {
	case 's':
		return ah.handleStartContainer(containerID, containerService, onRefresh)
	case 'S':
		return ah.handleStopContainer(containerID, containerService, onRefresh)
	case 'r':
		return ah.handleRestartContainer(containerID, containerService, onRefresh)
	}
	return false
}

// handleStartContainer handles starting a container
func (ah *ActionHandlers) handleStartContainer(
	containerID string,
	containerService any,
	onRefresh func(),
) bool {
	if cs, ok := containerService.(interfaces.ContainerService); ok {
		ah.executor.StartOperation("container", containerID, cs.StartContainer, onRefresh)
	}
	return true
}

// handleStopContainer handles stopping a container
func (ah *ActionHandlers) handleStopContainer(
	containerID string,
	containerService any,
	onRefresh func(),
) bool {
	if cs, ok := containerService.(interfaces.ContainerService); ok {
		ah.executor.StopOperation("container", containerID, cs.StopContainer, onRefresh)
	}
	return true
}

// handleRestartContainer handles restarting a container
func (ah *ActionHandlers) handleRestartContainer(
	containerID string,
	containerService any,
	onRefresh func(),
) bool {
	if cs, ok := containerService.(interfaces.ContainerService); ok {
		ah.executor.RestartOperation("container", containerID, cs.RestartContainer, onRefresh)
	}
	return true
}

// handleContainerManagementAction handles container management operations (delete, inspect)
func (ah *ActionHandlers) handleContainerManagementAction(
	action rune,
	containerID, containerName string,
	containerService any,
	onRefresh func(),
) bool {
	switch action {
	case 'd':
		if cs, ok := containerService.(interfaces.ContainerService); ok {
			ah.HandleDeleteAction(
				"container",
				containerID,
				containerName,
				cs.RemoveContainer,
				onRefresh,
			)
		}
		return true
	case 'i':
		if cs, ok := containerService.(interfaces.ContainerService); ok {
			ah.HandleInspectAction("container", containerID, cs.InspectContainer)
		}
		return true
	}
	return false
}

// handleContainerAccessAction handles container access operations (attach, logs, exec)
func (ah *ActionHandlers) handleContainerAccessAction(
	action rune,
	containerID, containerName string,
	containerService any,
) bool {
	switch action {
	case 'a':
		ah.HandleAttachAction(containerID, containerName)
		return true
	case 'l':
		ah.ui.ShowLogs(containerID, containerName)
		return true
	case 'e':
		if cs, ok := containerService.(interfaces.ContainerService); ok {
			ah.HandleExecAction(containerID, containerName, cs.ExecContainer)
		}
		return true
	}
	return false
}

// createExecInput creates and configures the exec command input field
func (ah *ActionHandlers) createExecInput(containerName string) *tview.InputField {
	themeManager := ah.ui.GetThemeManager()

	execInput := tview.NewInputField()
	ah.configureExecInputStyling(execInput, themeManager, containerName)
	ah.configureExecInputBehavior(execInput, themeManager)

	return execInput
}

// configureExecInputStyling configures the visual styling of the exec input
func (ah *ActionHandlers) configureExecInputStyling(
	execInput *tview.InputField,
	themeManager *config.ThemeManager,
	containerName string,
) {
	execInput.SetLabel(" Exec Command: ")
	execInput.SetLabelColor(themeManager.GetContainerExecLabelColor())
	execInput.SetFieldTextColor(themeManager.GetContainerExecTextColor())
	execInput.SetBorder(true)
	execInput.SetBorderColor(themeManager.GetContainerExecBorderColor())
	execInput.SetTitle(fmt.Sprintf(" Execute in %s ", containerName))
	execInput.SetTitleColor(themeManager.GetContainerExecTitleColor())
	execInput.SetBackgroundColor(themeManager.GetContainerExecBackgroundColor())
}

// configureExecInputBehavior configures the behavior of the exec input
func (ah *ActionHandlers) configureExecInputBehavior(
	execInput *tview.InputField,
	themeManager *config.ThemeManager,
) {
	execInput.SetPlaceholder("Type command to execute (e.g., ls -la, pwd, whoami)")
	execInput.SetPlaceholderTextColor(themeManager.GetContainerExecPlaceholderColor())
	execInput.SetFieldWidth(80)
}

// setupExecInputHandlers configures the event handlers for the exec input
func (ah *ActionHandlers) setupExecInputHandlers(
	execInput *tview.InputField,
	mainFlex *tview.Flex,
	containerID, containerName string,
	execFunc func(context.Context, string, []string, bool) (string, error),
) {
	execInput.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			ah.handleExecCommand(execInput, containerID, containerName, execFunc)
		}
		mainFlex.RemoveItem(execInput)
	})
}

// handleExecCommand processes the exec command when Enter is pressed
func (ah *ActionHandlers) handleExecCommand(
	execInput *tview.InputField,
	containerID, containerName string,
	execFunc func(context.Context, string, []string, bool) (string, error),
) {
	command := execInput.GetText()
	if command == "" {
		ah.ui.ShowError(errors.New("command cannot be empty"))
		return
	}

	args := ah.parseExecCommand(command)
	if len(args) == 0 {
		ah.ui.ShowError(errors.New("invalid command"))
		return
	}

	ah.executeCommand(containerID, containerName, command, args, execFunc)
}

// parseExecCommand parses the command string into arguments
func (ah *ActionHandlers) parseExecCommand(command string) []string {
	if ah.isComplexCommand(command) {
		return ah.parseComplexCommand(command)
	}
	return strings.Fields(command)
}

// isComplexCommand checks if the command contains shell operators
func (ah *ActionHandlers) isComplexCommand(command string) bool {
	return strings.Contains(command, "|") || strings.Contains(command, ">") ||
		strings.Contains(command, "<") || strings.Contains(command, "&&") ||
		strings.Contains(command, "||")
}

// parseComplexCommand parses a complex command with shell operators
func (ah *ActionHandlers) parseComplexCommand(command string) []string {
	return []string{"/bin/sh", "-c", command}
}

// executeCommand executes the command and handles the result
func (ah *ActionHandlers) executeCommand(
	containerID, containerName, command string,
	args []string,
	execFunc func(context.Context, string, []string, bool) (string, error),
) {
	ctx := context.Background()
	output, err := execFunc(ctx, containerID, args, false)
	if err != nil {
		ah.ui.ShowError(fmt.Errorf("exec failed: %w", err))
		return
	}

	ah.showExecOutput(containerName, command, output)
}

// addExecInputToUI adds the exec input to the UI and sets focus
func (ah *ActionHandlers) addExecInputToUI(execInput *tview.InputField, mainFlex *tview.Flex) {
	mainFlex.AddItem(execInput, 3, 0, true)

	app, ok := ah.ui.GetApp().(*tview.Application)
	if !ok {
		ah.ui.ShowError(errors.New("failed to get application instance"))
		return
	}
	app.SetFocus(execInput)
}

// showExecOutput displays the command output in a modal
func (ah *ActionHandlers) showExecOutput(containerName, command, output string) {
	outputModal := ah.createExecOutputModal(containerName, command, output)

	ah.setupExecOutputModalHandlers(outputModal)
	ah.addExecOutputModalToUI(outputModal)
}

// createExecOutputModal creates the exec output modal with content
func (ah *ActionHandlers) createExecOutputModal(
	containerName, command, output string,
) *tview.Modal {
	return tview.NewModal().
		SetText(fmt.Sprintf("Command Output: %s\n\nContainer: %s\nCommand: %s\n\nOutput:\n%s",
			containerName, containerName, command, output)).
		AddButtons([]string{"Close"})
}

// setupExecOutputModalHandlers configures the event handlers for the exec output modal
func (ah *ActionHandlers) setupExecOutputModalHandlers(outputModal *tview.Modal) {
	outputModal.SetDoneFunc(func(_ int, _ string) {
		ah.closeExecOutputModal()
		ah.returnFocusToViewContainer()
	})
}

// closeExecOutputModal removes the exec output modal from the UI
func (ah *ActionHandlers) closeExecOutputModal() {
	pages, ok := ah.ui.GetPages().(*tview.Pages)
	if !ok {
		ah.ui.ShowError(errors.New("failed to get pages container"))
		return
	}
	pages.RemovePage("exec_output_modal")
}

// returnFocusToViewContainer returns focus to the view container after closing the modal
func (ah *ActionHandlers) returnFocusToViewContainer() {
	viewContainer := ah.ui.GetViewContainer()
	if viewContainer == nil {
		return
	}

	vc, ok := viewContainer.(*tview.Flex)
	if !ok {
		return
	}

	ah.setFocusToViewContainer(vc)
}

// setFocusToViewContainer sets focus to the specified view container
func (ah *ActionHandlers) setFocusToViewContainer(vc *tview.Flex) {
	app, ok := ah.ui.GetApp().(*tview.Application)
	if !ok {
		ah.ui.ShowError(errors.New("failed to get application instance"))
		return
	}
	app.SetFocus(vc)
}

// addExecOutputModalToUI adds the exec output modal to the UI and sets focus
func (ah *ActionHandlers) addExecOutputModalToUI(outputModal *tview.Modal) {
	pages, ok := ah.ui.GetPages().(*tview.Pages)
	if !ok {
		ah.ui.ShowError(errors.New("failed to get pages container"))
		return
	}
	pages.AddPage("exec_output_modal", outputModal, true, true)

	app, ok := ah.ui.GetApp().(*tview.Application)
	if !ok {
		ah.ui.ShowError(errors.New("failed to get application instance"))
		return
	}
	app.SetFocus(outputModal)
}
