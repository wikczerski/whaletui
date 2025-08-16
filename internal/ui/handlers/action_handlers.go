package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/rivo/tview"
	"github.com/wikczerski/D5r/internal/ui/builders"
	"github.com/wikczerski/D5r/internal/ui/interfaces"
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
func (ah *ActionHandlers) HandleDeleteAction(resourceType, resourceID, resourceName string, deleteFunc func(context.Context, string, bool) error, onRefresh func()) {
	ah.executor.DeleteOperation(resourceType, resourceID, resourceName, deleteFunc, onRefresh)
}

// HandleInspectAction provides a common inspect action pattern
func (ah *ActionHandlers) HandleInspectAction(resourceType, resourceID string, inspectFunc func(context.Context, string) (map[string]any, error)) {
	inspectView, inspectFlex := builders.CreateInspectView(fmt.Sprintf("Inspect: %s", resourceID))

	inspectFlex.GetItem(1).(*tview.Button).SetSelectedFunc(func() {
		pages := ah.ui.GetPages().(*tview.Pages)
		pages.RemovePage("inspect")
	})

	pages := ah.ui.GetPages().(*tview.Pages)
	pages.AddPage("inspect", inspectFlex, true, true)

	// Load inspect data asynchronously
	go func() {
		inspectData, err := inspectFunc(context.Background(), resourceID)
		app := ah.ui.GetApp().(*tview.Application)
		app.QueueUpdateDraw(func() {
			if err != nil {
				inspectView.SetText(fmt.Sprintf("%s inspection failed: %v", resourceType, err))
			} else {
				data, jsonErr := json.MarshalIndent(inspectData, "", "  ")
				if jsonErr != nil {
					inspectView.SetText(fmt.Sprintf("Failed to format %s data: %v", resourceType, jsonErr))
				} else {
					inspectView.SetText(string(data))
				}
			}
		})
	}()
}

// HandleContainerAction provides container-specific action handling
func (ah *ActionHandlers) HandleContainerAction(action rune, containerID, containerName string, containerService interface {
	StartContainer(context.Context, string) error
	StopContainer(context.Context, string, *time.Duration) error
	RestartContainer(context.Context, string, *time.Duration) error
	RemoveContainer(context.Context, string, bool) error
	InspectContainer(context.Context, string) (map[string]any, error)
	GetContainerLogs(context.Context, string) (string, error)
}, onRefresh func()) {

	switch action {
	case 's':
		ah.executor.StartOperation("container", containerID, containerService.StartContainer, onRefresh)
	case 'S':
		ah.executor.StopOperation("container", containerID, containerService.StopContainer, onRefresh)
	case 'r':
		ah.executor.RestartOperation("container", containerID, containerService.RestartContainer, onRefresh)
	case 'd':
		ah.HandleDeleteAction("container", containerID, containerName, containerService.RemoveContainer, onRefresh)
	case 'l':
		ah.ui.ShowLogs(containerID, containerName)
	case 'i':
		ah.HandleInspectAction("container", containerID, containerService.InspectContainer)
	}
}

// HandleResourceAction provides generic resource action handling (for images, volumes, networks)
func (ah *ActionHandlers) HandleResourceAction(action rune, resourceType, resourceID, resourceName string, inspectFunc func(context.Context, string) (map[string]any, error), deleteFunc func(context.Context, string, bool) error, onRefresh func()) {
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
