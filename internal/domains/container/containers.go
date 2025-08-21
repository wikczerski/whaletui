package container

import (
	"context"

	"github.com/gdamore/tcell/v2"
	"github.com/wikczerski/whaletui/internal/shared"
	"github.com/wikczerski/whaletui/internal/ui/builders"
	"github.com/wikczerski/whaletui/internal/ui/constants"
	"github.com/wikczerski/whaletui/internal/ui/handlers"
	"github.com/wikczerski/whaletui/internal/ui/interfaces"
)

// ContainersView displays and manages Docker containers
type ContainersView struct {
	*shared.BaseView[Container]
	executor *handlers.OperationExecutor
	handlers *handlers.ActionHandlers
}

// NewContainersView creates a new containers view
func NewContainersView(ui interfaces.UIInterface) *ContainersView {
	headers := []string{"ID", "Name", "Image", "Status", "State", "Ports", "Created"}
	baseView := shared.NewBaseView[Container](ui, "containers", headers)

	cv := &ContainersView{
		BaseView: baseView,
		executor: handlers.NewOperationExecutor(ui),
		handlers: handlers.NewActionHandlers(ui),
	}

	// Set up callbacks
	cv.ListItems = cv.listContainers
	cv.FormatRow = func(c Container) []string { return cv.formatContainerRow(&c) }
	cv.GetRowColor = func(c Container) tcell.Color { return cv.getStateColor(&c) }
	cv.GetItemID = func(c Container) string { return c.ID }
	cv.GetItemName = func(c Container) string { return c.Name }
	cv.HandleKeyPress = func(key rune, c Container) { cv.handleContainerKey(key, &c) }
	cv.ShowDetails = func(c Container) { cv.showContainerDetails(&c) }
	cv.GetActions = cv.getContainerActions

	return cv
}

func (cv *ContainersView) listContainers(ctx context.Context) ([]Container, error) {
	services := cv.GetUI().GetServices()
	if !services.IsContainerServiceAvailable() {
		return []Container{}, nil
	}

	containers, err := services.GetContainerService().ListContainers(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]Container, len(containers))
	for i, container := range containers {
		if c, ok := container.(Container); ok {
			result[i] = c
		}
	}

	return result, nil
}

func (cv *ContainersView) formatContainerRow(container *Container) []string {
	return []string{
		container.ID,
		container.Name,
		container.Image,
		container.Status,
		container.State,
		container.Ports,
		builders.FormatTime(container.Created),
	}
}

func (cv *ContainersView) getContainerActions() map[rune]string {
	return cv.GetUI().GetServices().GetContainerService().GetActions()
}

func (cv *ContainersView) handleContainerKey(key rune, container *Container) {
	services := cv.GetUI().GetServices()
	cv.handlers.HandleContainerAction(key, container.ID, container.Name, services.GetContainerService(), func() { cv.Refresh() })
}

func (cv *ContainersView) getStateColor(container *Container) tcell.Color {
	switch container.State {
	case "running":
		return constants.TableSuccessColor
	case "exited":
		return constants.TableErrorColor
	case "created":
		return constants.TableWarningColor
	default:
		return constants.TableDefaultRowColor
	}
}

func (cv *ContainersView) showContainerDetails(container *Container) {
	ctx := context.Background()
	services := cv.GetUI().GetServices()
	inspectData, err := services.GetContainerService().InspectContainer(ctx, container.ID)
	cv.ShowItemDetails(*container, inspectData, err)
}
