package views

import (
	"context"

	"github.com/gdamore/tcell/v2"
	"github.com/wikczerski/D5r/internal/models"
	"github.com/wikczerski/D5r/internal/ui/builders"
	"github.com/wikczerski/D5r/internal/ui/constants"
	"github.com/wikczerski/D5r/internal/ui/handlers"
	"github.com/wikczerski/D5r/internal/ui/interfaces"
)

// ContainersView displays and manages Docker containers
type ContainersView struct {
	*BaseView[models.Container]
	executor *handlers.OperationExecutor
	handlers *handlers.ActionHandlers
}

// NewContainersView creates a new containers view
func NewContainersView(ui interfaces.UIInterface) *ContainersView {
	headers := []string{"ID", "Name", "Image", "Status", "State", "Ports", "Created"}
	baseView := NewBaseView[models.Container](ui, "containers", headers)

	cv := &ContainersView{
		BaseView: baseView,
		executor: handlers.NewOperationExecutor(ui),
		handlers: handlers.NewActionHandlers(ui),
	}

	// Set up callbacks
	cv.ListItems = cv.listContainers
	cv.FormatRow = func(c models.Container) []string { return cv.formatContainerRow(&c) }
	cv.GetRowColor = func(c models.Container) tcell.Color { return cv.getStateColor(&c) }
	cv.GetItemID = func(c models.Container) string { return c.ID }
	cv.GetItemName = func(c models.Container) string { return c.Name }
	cv.HandleKeyPress = func(key rune, c models.Container) { cv.handleContainerKey(key, &c) }
	cv.ShowDetails = func(c models.Container) { cv.showContainerDetails(&c) }
	cv.GetActions = cv.getContainerActions

	return cv
}

func (cv *ContainersView) listContainers(ctx context.Context) ([]models.Container, error) {
	services := cv.ui.GetServices()
	if services == nil || services.ContainerService == nil {
		return []models.Container{}, nil
	}
	return services.ContainerService.ListContainers(ctx)
}

func (cv *ContainersView) formatContainerRow(container *models.Container) []string {
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
	return map[rune]string{
		's': "Start",
		'S': "Stop",
		'r': "Restart",
		'd': "Delete",
		'a': "Attach",
		'l': "View Logs",
		'i': "Inspect",
		'e': "Exec",
	}
}

func (cv *ContainersView) handleContainerKey(key rune, container *models.Container) {
	services := cv.ui.GetServices()
	cv.handlers.HandleContainerAction(key, container.ID, container.Name, services.ContainerService, func() { cv.Refresh() })
}

func (cv *ContainersView) getStateColor(container *models.Container) tcell.Color {
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

func (cv *ContainersView) showContainerDetails(container *models.Container) {
	ctx := context.Background()
	services := cv.ui.GetServices()
	inspectData, err := services.ContainerService.InspectContainer(ctx, container.ID)
	cv.ShowItemDetails(*container, inspectData, err)
}
