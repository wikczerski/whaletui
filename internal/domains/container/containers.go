package container

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/wikczerski/whaletui/internal/logger"
	"github.com/wikczerski/whaletui/internal/shared"
	"github.com/wikczerski/whaletui/internal/ui/constants"
	"github.com/wikczerski/whaletui/internal/ui/formatters"
	"github.com/wikczerski/whaletui/internal/ui/handlers"
	"github.com/wikczerski/whaletui/internal/ui/interfaces"
	"github.com/wikczerski/whaletui/internal/ui/utils"
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

	cv.setupCallbacks()
	cv.setupCharacterLimits(ui)
	return cv
}

// setupCallbacks sets up all the callback functions for the containers view
func (cv *ContainersView) setupCallbacks() {
	cv.setupBasicCallbacks()
	cv.setupActionCallbacks()
}

// setupBasicCallbacks sets up the basic view callbacks
func (cv *ContainersView) setupBasicCallbacks() {
	cv.ListItems = cv.listContainers
	cv.FormatRow = func(c shared.Container) []string { return cv.formatContainerRow(&c) }
	cv.GetRowColor = func(c shared.Container) tcell.Color { return cv.getStateColor(&c) }
	cv.GetItemID = func(c shared.Container) string { return c.ID }
	cv.GetItemName = func(c shared.Container) string { return c.Name }
}

// setupActionCallbacks sets up the action-related callbacks
func (cv *ContainersView) setupActionCallbacks() {
	cv.HandleKeyPress = func(key rune, c shared.Container) { cv.handleContainerKey(key, &c) }
	cv.ShowDetailsCallback = func(c shared.Container) { cv.showContainerDetails(&c) }
	cv.GetActions = cv.getContainerActions
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

	return cv.getContainersFromService(ctx, services)
}

// getContainersFromService retrieves containers from the service factory
func (cv *ContainersView) getContainersFromService(
	ctx context.Context,
	services any,
) ([]shared.Container, error) {
	serviceFactory := cv.getServiceFactory(services)
	if serviceFactory == nil {
		return []shared.Container{}, nil
	}

	containerService := cv.getContainerService(serviceFactory)
	if containerService == nil {
		return []shared.Container{}, nil
	}

	return cv.executeContainerList(ctx, containerService)
}

// getServiceFactory gets the service factory from services
func (cv *ContainersView) getServiceFactory(services any) interfaces.ServiceFactoryInterface {
	serviceFactory, ok := services.(interfaces.ServiceFactoryInterface)
	if !ok {
		cv.log.Debug("serviceFactory type assertion failed")
		return nil
	}
	cv.log.Debug("serviceFactory type assertion successful")
	return serviceFactory
}

// getContainerService gets the container service from the service factory
func (cv *ContainersView) getContainerService(
	serviceFactory interfaces.ServiceFactoryInterface,
) any {
	containerService := serviceFactory.GetContainerService()
	if containerService == nil {
		cv.log.Debug("containerService is nil")
		return nil
	}
	cv.log.Debug("containerService type", "type", fmt.Sprintf("%T", containerService))
	return containerService
}

// executeContainerList executes the container list operation
func (cv *ContainersView) executeContainerList(
	ctx context.Context,
	containerService any,
) ([]shared.Container, error) {
	if containerService == nil {
		cv.log.Debug("listService type assertion failed")
		return []shared.Container{}, nil
	}

	cv.log.Debug("listService type assertion successful")
	containers, err := containerService.(interface {
		ListContainers(context.Context) ([]shared.Container, error)
	}).ListContainers(ctx)
	if err != nil {
		cv.log.Debug("ListContainers error", "error", err)
		return nil, err
	}

	cv.log.Debug("ListContainers returned containers", "count", len(containers))
	return containers, nil
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
		formatters.FormatTime(container.Created),
	}
}

func (cv *ContainersView) getContainerActions() map[rune]string {
	services := cv.GetUI().GetServicesAny()
	cv.log.Debug(
		"getContainerActions called",
		"services_nil",
		services == nil,
		"services_type",
		fmt.Sprintf("%T", services),
	)

	return cv.getActionsFromService(services)
}

// getActionsFromService retrieves actions from the service factory
func (cv *ContainersView) getActionsFromService(services any) map[rune]string {
	serviceFactory, ok := services.(interfaces.ServiceFactoryInterface)
	if !ok {
		cv.log.Debug("serviceFactory type assertion failed")
		return make(map[rune]string)
	}

	cv.log.Debug("serviceFactory type assertion successful")
	containerService := serviceFactory.GetContainerService()
	if containerService == nil {
		cv.log.Debug("containerService is nil")
		return make(map[rune]string)
	}

	return cv.extractActionsFromContainerService(containerService)
}

// extractActionsFromContainerService extracts actions from the container service
func (cv *ContainersView) extractActionsFromContainerService(containerService any) map[rune]string {
	cv.log.Debug(
		"containerService retrieved",
		"containerService_type",
		fmt.Sprintf("%T", containerService),
	)

	actionService, ok := containerService.(interface{ GetActions() map[rune]string })
	if !ok {
		cv.log.Debug("actionService type assertion failed")
		return make(map[rune]string)
	}

	cv.log.Debug("actionService type assertion successful")
	actions := actionService.GetActions()
	cv.log.Debug("actions retrieved", "actions", actions)
	return actions
}

func (cv *ContainersView) handleContainerKey(key rune, container *shared.Container) {
	services := cv.GetUI().GetServicesAny()

	// Type assertion to get the service factory
	if serviceFactory, ok := services.(interfaces.ServiceFactoryInterface); ok {
		if containerService := serviceFactory.GetContainerService(); containerService != nil {
			cv.handlers.HandleContainerAction(
				key,
				container.ID,
				container.Name,
				containerService,
				func() { cv.Refresh() },
			)
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

// setupCharacterLimits sets up character limits for table columns
func (cv *ContainersView) setupCharacterLimits(ui interfaces.UIInterface) {
	// Define column types for containers table
	columnTypes := []string{"id", "name", "image", "status", "state", "ports", "created"}
	cv.SetColumnTypes(columnTypes)

	// Create formatter from theme manager
	formatter := utils.NewTableFormatterFromTheme(ui.GetThemeManager())
	cv.SetFormatter(formatter)
}
