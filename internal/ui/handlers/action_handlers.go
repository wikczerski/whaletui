package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/wikczerski/whaletui/internal/config"
	"github.com/wikczerski/whaletui/internal/services"
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

// setupInspectCloseButton configures the close button for the inspect view
func (ah *ActionHandlers) setupInspectCloseButton(inspectFlex *tview.Flex) {
	inspectFlex.GetItem(1).(*tview.Button).SetSelectedFunc(func() {
		pages := ah.ui.GetPages().(*tview.Pages)
		pages.RemovePage("inspect")
	})
}

// addInspectPage adds the inspect page to the UI
func (ah *ActionHandlers) addInspectPage(inspectFlex *tview.Flex) {
	pages := ah.ui.GetPages().(*tview.Pages)
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
		app := ah.ui.GetApp().(*tview.Application)
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

// HandleContainerAction provides container-specific action handling
func (ah *ActionHandlers) HandleContainerAction(
	action rune,
	containerID, containerName string,
	containerService services.ContainerService,
	onRefresh func(),
) {
	if ah.handleContainerLifecycleAction(action, containerID, containerService, onRefresh) {
		return
	}

	if ah.handleContainerManagementAction(action, containerID, containerName, containerService, onRefresh) {
		return
	}

	ah.handleContainerAccessAction(action, containerID, containerName, containerService)
}

// handleContainerLifecycleAction handles container lifecycle operations (start, stop, restart)
func (ah *ActionHandlers) handleContainerLifecycleAction(
	action rune,
	containerID string,
	containerService services.ContainerService,
	onRefresh func(),
) bool {
	switch action {
	case 's':
		ah.executor.StartOperation("container", containerID, containerService.StartContainer, onRefresh)
		return true
	case 'S':
		ah.executor.StopOperation("container", containerID, containerService.StopContainer, onRefresh)
		return true
	case 'r':
		ah.executor.RestartOperation("container", containerID, containerService.RestartContainer, onRefresh)
		return true
	}
	return false
}

// handleContainerManagementAction handles container management operations (delete, inspect)
func (ah *ActionHandlers) handleContainerManagementAction(
	action rune,
	containerID, containerName string,
	containerService services.ContainerService,
	onRefresh func(),
) bool {
	switch action {
	case 'd':
		ah.HandleDeleteAction("container", containerID, containerName, containerService.RemoveContainer, onRefresh)
		return true
	case 'i':
		ah.HandleInspectAction("container", containerID, containerService.InspectContainer)
		return true
	}
	return false
}

// handleContainerAccessAction handles container access operations (attach, logs, exec)
func (ah *ActionHandlers) handleContainerAccessAction(
	action rune,
	containerID, containerName string,
	containerService services.ContainerService,
) bool {
	switch action {
	case 'a':
		ah.HandleAttachAction(containerID, containerName)
		return true
	case 'l':
		ah.ui.ShowLogs(containerID, containerName)
		return true
	case 'e':
		ah.HandleExecAction(containerID, containerName, containerService.ExecContainer)
		return true
	}
	return false
}

// HandleContainerLifecycleAction handles container lifecycle operations (start, stop, restart)
func (ah *ActionHandlers) HandleContainerLifecycleAction(
	action rune,
	containerID string,
	containerService services.ContainerService,
	onRefresh func(),
) {
	switch action {
	case 's':
		ah.executor.StartOperation("container", containerID, containerService.StartContainer, onRefresh)
	case 'S':
		ah.executor.StopOperation("container", containerID, containerService.StopContainer, onRefresh)
	case 'r':
		ah.executor.RestartOperation("container", containerID, containerService.RestartContainer, onRefresh)
	}
}

// HandleContainerManagementAction handles container management operations (delete, inspect)
func (ah *ActionHandlers) HandleContainerManagementAction(
	action rune,
	containerID, containerName string,
	containerService services.ContainerService,
	onRefresh func(),
) {
	switch action {
	case 'd':
		ah.HandleDeleteAction("container", containerID, containerName, containerService.RemoveContainer, onRefresh)
	case 'i':
		ah.HandleInspectAction("container", containerID, containerService.InspectContainer)
	}
}

// HandleContainerAccessAction handles container access operations (attach, logs, exec)
func (ah *ActionHandlers) HandleContainerAccessAction(
	action rune,
	containerID, containerName string,
	containerService services.ContainerService,
) {
	switch action {
	case 'a':
		ah.HandleAttachAction(containerID, containerName)
	case 'l':
		ah.ui.ShowLogs(containerID, containerName)
	case 'e':
		ah.HandleExecAction(containerID, containerName, containerService.ExecContainer)
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
	mainFlex := ah.ui.GetMainFlex().(*tview.Flex)
	execInput := ah.createExecInput(containerName)

	ah.setupExecInputHandlers(execInput, mainFlex, containerID, containerName, execFunc)
	ah.addExecInputToUI(execInput, mainFlex)
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
func (ah *ActionHandlers) configureExecInputStyling(execInput *tview.InputField, themeManager *config.ThemeManager, containerName string) {
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
func (ah *ActionHandlers) configureExecInputBehavior(execInput *tview.InputField, themeManager *config.ThemeManager) {
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
		ah.ui.ShowError(fmt.Errorf("command cannot be empty"))
		return
	}

	args := ah.parseExecCommand(command)
	if len(args) == 0 {
		ah.ui.ShowError(fmt.Errorf("invalid command"))
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

	app := ah.ui.GetApp().(*tview.Application)
	app.SetFocus(execInput)
}

// showExecOutput displays the command output in a modal
func (ah *ActionHandlers) showExecOutput(containerName, command, output string) {
	outputModal := ah.createExecOutputModal(containerName, command, output)

	ah.setupExecOutputModalHandlers(outputModal)
	ah.addExecOutputModalToUI(outputModal)
}

// createExecOutputModal creates the exec output modal with content
func (ah *ActionHandlers) createExecOutputModal(containerName, command, output string) *tview.Modal {
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
	pages := ah.ui.GetPages().(*tview.Pages)
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

	app := ah.ui.GetApp().(*tview.Application)
	app.SetFocus(vc)
}

// addExecOutputModalToUI adds the exec output modal to the UI and sets focus
func (ah *ActionHandlers) addExecOutputModalToUI(outputModal *tview.Modal) {
	pages := ah.ui.GetPages().(*tview.Pages)
	pages.AddPage("exec_output_modal", outputModal, true, true)

	app := ah.ui.GetApp().(*tview.Application)
	app.SetFocus(outputModal)
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
