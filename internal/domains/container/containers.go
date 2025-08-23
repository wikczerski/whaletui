package container

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/wikczerski/whaletui/internal/logger"
	"github.com/wikczerski/whaletui/internal/shared"
	"github.com/wikczerski/whaletui/internal/ui/builders"
	"github.com/wikczerski/whaletui/internal/ui/constants"
	"github.com/wikczerski/whaletui/internal/ui/handlers"
	"github.com/wikczerski/whaletui/internal/ui/interfaces"
)

// ContainersView displays and manages Docker containers
type ContainersView struct {
	*shared.BaseView[shared.Container]
	executor *handlers.OperationExecutor
	handlers *handlers.ActionHandlers
	log      *slog.Logger
}

// NewContainersView creates a new containers view
func NewContainersView(ui interfaces.UIInterface) *ContainersView {
	headers := []string{"ID", "Name", "Image", "Status", "State", "Ports", "Created"}
	baseView := shared.NewBaseView[shared.Container](ui, "containers", headers)

	cv := &ContainersView{
		BaseView: baseView,
		executor: handlers.NewOperationExecutor(ui),
		handlers: handlers.NewActionHandlers(ui),
		log:      logger.GetLogger(),
	}

	// Set up callbacks
	cv.ListItems = cv.listContainers
	cv.FormatRow = func(c shared.Container) []string { return cv.formatContainerRow(&c) }
	cv.GetRowColor = func(c shared.Container) tcell.Color { return cv.getStateColor(&c) }
	cv.GetItemID = func(c shared.Container) string { return c.ID }
	cv.GetItemName = func(c shared.Container) string { return c.Name }
	cv.HandleKeyPress = func(key rune, c shared.Container) { cv.handleContainerKey(key, &c) }
	cv.ShowDetails = func(c shared.Container) { cv.showContainerDetails(&c) }
	cv.GetActions = cv.getContainerActions

	return cv
}

func (cv *ContainersView) listContainers(ctx context.Context) ([]shared.Container, error) {
	cv.log.Info("listContainers called")

	ui := cv.GetUI()
	cv.log.Debug("UI instance", "ui_nil", ui == nil, "ui_type", fmt.Sprintf("%T", ui))

	services := ui.GetServicesAny()
	if services == nil {
		cv.log.Debug("services is nil")
		return []shared.Container{}, nil
	}

	cv.log.Debug("services type", "type", fmt.Sprintf("%T", services))

	// Type assertion to get the service factory
	if serviceFactory, ok := services.(interfaces.ServiceFactoryInterface); ok {
		cv.log.Debug("serviceFactory type assertion successful")
		if containerService := serviceFactory.GetContainerService(); containerService != nil {
			cv.log.Debug("containerService type", "type", fmt.Sprintf("%T", containerService))
			// Type assertion to get the ListContainers method
			if containerService != nil {
				cv.log.Debug("listService type assertion successful")
				containers, err := containerService.ListContainers(ctx)
				if err != nil {
					cv.log.Debug("ListContainers error", "error", err)
					return nil, err
				}
				cv.log.Debug("ListContainers returned containers", "count", len(containers))
				return containers, nil
			}
			cv.log.Debug("listService type assertion failed")
		} else {
			cv.log.Debug("containerService is nil")
		}
	} else {
		cv.log.Debug("serviceFactory type assertion failed")
	}

	return []shared.Container{}, nil
}

func (cv *ContainersView) formatContainerRow(container *shared.Container) []string {
	ports := strings.Join(container.Ports, ", ")
	return []string{
		container.ID,
		container.Name,
		container.Image,
		container.Status,
		container.State,
		ports,
		builders.FormatTime(container.Created),
	}
}

func (cv *ContainersView) getContainerActions() map[rune]string {
	services := cv.GetUI().GetServicesAny()
	cv.log.Debug("getContainerActions called", "services_nil", services == nil, "services_type", fmt.Sprintf("%T", services))

	// Type assertion to get the service factory
	if serviceFactory, ok := services.(interfaces.ServiceFactoryInterface); ok {
		cv.log.Debug("serviceFactory type assertion successful")
		if containerService := serviceFactory.GetContainerService(); containerService != nil {
			cv.log.Debug("containerService retrieved", "containerService_type", fmt.Sprintf("%T", containerService))
			// Type assertion to get the GetActions method
			if actionService, ok := containerService.(interface{ GetActions() map[rune]string }); ok {
				cv.log.Debug("actionService type assertion successful")
				actions := actionService.GetActions()
				cv.log.Debug("actions retrieved", "actions", actions)
				return actions
			}
			cv.log.Debug("actionService type assertion failed")
		} else {
			cv.log.Debug("containerService is nil")
		}
	} else {
		cv.log.Debug("serviceFactory type assertion failed")
	}

	cv.log.Debug("returning empty actions map")
	return make(map[rune]string)
}

func (cv *ContainersView) handleContainerKey(key rune, container *shared.Container) {
	services := cv.GetUI().GetServicesAny()

	// Type assertion to get the service factory
	if serviceFactory, ok := services.(interfaces.ServiceFactoryInterface); ok {
		if containerService := serviceFactory.GetContainerService(); containerService != nil {
			cv.handlers.HandleContainerAction(key, container.ID, container.Name, containerService, func() { cv.Refresh() })
		}
	}
}

func (cv *ContainersView) getStateColor(container *shared.Container) tcell.Color {
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

func (cv *ContainersView) showContainerDetails(container *shared.Container) {
	ctx := context.Background()
	services := cv.GetUI().GetServicesAny()

	// Type assertion to get the service factory
	if serviceFactory, ok := services.(interfaces.ServiceFactoryInterface); ok {
		if containerService := serviceFactory.GetContainerService(); containerService != nil {
			// Type assertion to get the InspectContainer method
			if inspectService, ok := containerService.(interface {
				InspectContainer(context.Context, string) (map[string]any, error)
			}); ok {
				inspectData, err := inspectService.InspectContainer(ctx, container.ID)
				cv.ShowItemDetails(*container, inspectData, err)
			}
		}
	}
}
