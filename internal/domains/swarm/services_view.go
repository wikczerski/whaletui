package swarm

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/wikczerski/whaletui/internal/logger"
	"github.com/wikczerski/whaletui/internal/shared"
	"github.com/wikczerski/whaletui/internal/ui/interfaces"
	"github.com/wikczerski/whaletui/internal/ui/utils"
)

// ServicesView represents the swarm services view
type ServicesView struct {
	*shared.BaseView[shared.SwarmService]
	presenter     *ServicesPresenter
	modalManager  interfaces.ModalManagerInterface
	headerManager interfaces.HeaderManagerInterface
	log           *slog.Logger
}

// NewServicesView creates a new swarm services view
func NewServicesView(
	ui shared.SharedUIInterface,
	serviceService *ServiceService,
	modalManager interfaces.ModalManagerInterface,
	headerManager interfaces.HeaderManagerInterface,
) *ServicesView {
	headers := []string{"ID", "Name", "Image", "Mode", "Replicas", "Status", "Created"}
	baseView := shared.NewBaseView[shared.SwarmService](ui, "swarm services", headers)

	presenter := NewServicesPresenter(serviceService, logger.GetLogger())

	view := &ServicesView{
		BaseView:      baseView,
		presenter:     presenter,
		modalManager:  modalManager,
		headerManager: headerManager,
		log:           logger.GetLogger(),
	}

	view.setupCallbacks()
	view.setupCharacterLimits(ui)

	return view
}

// handleInspect handles service inspection
func (v *ServicesView) handleInspect(ctx context.Context) (any, error) {
	selectedService := v.GetSelectedItem()
	if err := v.presenter.ValidateService(selectedService); err != nil {
		v.GetUI().ShowError(err)
		return v, err
	}

	serviceInfo, err := v.presenter.InspectService(ctx, selectedService.ID)
	if err != nil {
		errorMsg := v.presenter.FormatInspectError(selectedService.Name, err)
		v.GetUI().ShowError(fmt.Errorf("%s", errorMsg))
		return v, err
	}

	infoText := v.presenter.FormatServiceInfo(selectedService, serviceInfo)
	v.GetUI().ShowInfo(infoText)
	return v, nil
}

// handleScale handles service scaling
func (v *ServicesView) handleScale(ctx context.Context) (any, error) {
	selectedService := v.GetSelectedItem()
	if err := v.presenter.ValidateService(selectedService); err != nil {
		v.GetUI().ShowError(err)
		return v, err
	}

	currentReplicas, err := v.presenter.GetCurrentReplicas(ctx, selectedService.ID)
	if err != nil {
		v.GetUI().ShowError(fmt.Errorf("failed to get service info for '%s': %v",
			selectedService.Name, err))
		currentReplicas = 1
	}

	v.showScaleModal(ctx, selectedService, currentReplicas)
	return v, nil
}

// showScaleModal shows the scale modal for the service
func (v *ServicesView) showScaleModal(
	ctx context.Context,
	selectedService *shared.SwarmService,
	currentReplicas uint64,
) {
	v.GetUI().ShowServiceScaleModal(selectedService.Name, currentReplicas, func(newReplicas int) {
		v.handleServiceScaling(ctx, selectedService, newReplicas)
	})
}

// handleServiceScaling handles the actual service scaling
func (v *ServicesView) handleServiceScaling(
	ctx context.Context,
	selectedService *shared.SwarmService,
	newReplicas int,
) {
	if newReplicas < 0 {
		v.GetUI().ShowError(fmt.Errorf("invalid replica count: %d", newReplicas))
		return
	}

	err := v.presenter.ScaleService(ctx, selectedService.ID, uint64(newReplicas))
	if err != nil {
		v.handleScaleError(ctx, selectedService, newReplicas, err)
	} else {
		v.Refresh()
	}
}

// handleScaleError handles scaling errors with advanced recovery options
func (v *ServicesView) handleScaleError(
	ctx context.Context,
	service *shared.SwarmService,
	newReplicas int,
	err error,
) {
	if v.presenter.IsRetryableError(err) {
		v.showRetryDialog(ctx, service, newReplicas, err)
	} else {
		v.showFallbackDialog(ctx, service, newReplicas, err)
	}
}

// showRetryDialog shows the retry dialog for retryable errors
func (v *ServicesView) showRetryDialog(
	ctx context.Context,
	service *shared.SwarmService,
	newReplicas int,
	err error,
) {
	v.GetUI().ShowRetryDialog(
		fmt.Sprintf("scale service '%s' to %d replicas", service.Name, newReplicas),
		err,
		func() error {
			if newReplicas < 0 {
				return fmt.Errorf("invalid replica count: %d", newReplicas)
			}
			return v.presenter.ScaleService(ctx, service.ID, uint64(newReplicas))
		},
		func() {
			v.GetUI().ShowInfo(fmt.Sprintf("Service '%s' successfully scaled to %d replicas",
				service.Name, newReplicas))
			v.Refresh()
		},
	)
}

// showFallbackDialog shows the fallback dialog for non-retryable errors
func (v *ServicesView) showFallbackDialog(
	ctx context.Context,
	service *shared.SwarmService,
	newReplicas int,
	err error,
) {
	fallbackOptions := v.presenter.GetScaleFallbackOptions()

	v.GetUI().ShowFallbackDialog(
		fmt.Sprintf("scale service '%s' to %d replicas", service.Name, newReplicas),
		err,
		fallbackOptions,
		func(fallbackOption string) {
			v.executeScaleFallback(ctx, service, fallbackOption, newReplicas)
		},
	)
}

// executeScaleFallback executes fallback operations for scaling failures
func (v *ServicesView) executeScaleFallback(
	ctx context.Context,
	service *shared.SwarmService,
	fallbackOption string,
	newReplicas int,
) {
	if fallbackOption == "Try Different Replica Count" {
		v.handleTryDifferentReplicaCount(ctx, service)
		return
	}

	result, err := v.presenter.ExecuteScaleFallback(ctx, service, fallbackOption)
	if err != nil {
		v.GetUI().ShowError(err)
		return
	}

	v.GetUI().ShowInfo(result)
}

// handleTryDifferentReplicaCount handles the "Try Different Replica Count" fallback option
func (v *ServicesView) handleTryDifferentReplicaCount(
	ctx context.Context,
	service *shared.SwarmService,
) {
	currentReplicas, err := v.presenter.GetCurrentReplicas(ctx, service.ID)
	if err != nil {
		v.GetUI().ShowError(fmt.Errorf("failed to get service info: %v", err))
		return
	}

	v.GetUI().ShowServiceScaleModal(service.Name, currentReplicas, func(differentReplicas int) {
		v.handleDifferentReplicaCount(ctx, service, differentReplicas)
	})
}

// handleDifferentReplicaCount handles scaling to a different replica count
func (v *ServicesView) handleDifferentReplicaCount(
	ctx context.Context,
	service *shared.SwarmService,
	differentReplicas int,
) {
	if differentReplicas < 0 {
		v.GetUI().ShowError(fmt.Errorf("invalid replica count: %d", differentReplicas))
		return
	}

	err := v.presenter.ScaleService(ctx, service.ID, uint64(differentReplicas))
	if err != nil {
		v.GetUI().
			ShowError(fmt.Errorf("failed to scale service to %d replicas: %v", differentReplicas, err))
	} else {
		v.GetUI().ShowInfo(fmt.Sprintf("Service '%s' successfully scaled to %d replicas",
			service.Name, differentReplicas))
		v.Refresh()
	}
}

// handleRemove handles service removal
func (v *ServicesView) handleRemove(ctx context.Context) (any, error) {
	selectedService := v.GetSelectedItem()
	if err := v.presenter.ValidateService(selectedService); err != nil {
		v.GetUI().ShowError(err)
		return v, err
	}

	v.showRemoveConfirmation(ctx, selectedService)
	return v, nil
}

// showRemoveConfirmation shows the remove confirmation dialog
func (v *ServicesView) showRemoveConfirmation(
	ctx context.Context,
	selectedService *shared.SwarmService,
) {
	message := v.presenter.BuildRemoveConfirmationMessage(selectedService)
	v.GetUI().ShowConfirm(message, func(confirmed bool) {
		if confirmed {
			v.executeServiceRemoval(ctx, selectedService)
		}
	})
}

// executeServiceRemoval executes the actual service removal
func (v *ServicesView) executeServiceRemoval(
	ctx context.Context,
	selectedService *shared.SwarmService,
) {
	err := v.presenter.RemoveService(ctx, selectedService.ID)
	if err != nil {
		errorMsg := v.presenter.FormatRemoveError(selectedService.Name, err)
		v.GetUI().ShowError(fmt.Errorf("%s", errorMsg))
	} else {
		v.GetUI().ShowInfo(fmt.Sprintf("Service '%s' successfully removed", selectedService.Name))
		v.Refresh()
	}
}

// handleLogs handles service logs viewing
func (v *ServicesView) handleLogs(ctx context.Context) (any, error) {
	selectedService := v.GetSelectedItem()
	if err := v.presenter.ValidateService(selectedService); err != nil {
		v.GetUI().ShowError(err)
		return v, err
	}

	logs, err := v.presenter.GetServiceLogs(ctx, selectedService.ID)
	if err != nil {
		errorMsg := v.presenter.FormatLogsError(selectedService.Name, err)
		v.GetUI().ShowError(fmt.Errorf("%s", errorMsg))
		return v, err
	}

	logText := v.presenter.FormatServiceLogs(selectedService.Name, logs)
	v.GetUI().ShowInfo(logText)
	return v, nil
}

// setupCallbacks sets up the callbacks for the base view
func (v *ServicesView) setupCallbacks() {
	v.setupBasicCallbacks()
	v.setupActionCallbacks()
}

// setupBasicCallbacks sets up the basic view callbacks
func (v *ServicesView) setupBasicCallbacks() {
	v.ListItems = v.listServices
	v.FormatRow = func(s shared.SwarmService) []string { return v.formatServiceRow(&s) }
	v.GetItemID = func(s shared.SwarmService) string { return v.getServiceID(&s) }
	v.GetItemName = func(s shared.SwarmService) string { return v.getServiceName(&s) }
}

// setupActionCallbacks sets up the action-related callbacks
func (v *ServicesView) setupActionCallbacks() {
	v.HandleKeyPress = func(key rune, s shared.SwarmService) { v.handleAction(key, &s) }
	v.GetActions = v.getActions
}

// handleAction handles action key presses for swarm services
func (v *ServicesView) handleAction(key rune, service *shared.SwarmService) {
	ctx := context.Background()

	switch key {
	case 'i':
		if _, err := v.handleInspect(ctx); err != nil {
			v.log.Error("Failed to handle inspect action", "error", err)
		}
	case 's':
		if _, err := v.handleScale(ctx); err != nil {
			v.log.Error("Failed to handle scale action", "error", err)
		}
	case 'r':
		if _, err := v.handleRemove(ctx); err != nil {
			v.log.Error("Failed to handle remove action", "error", err)
		}
	case 'l':
		if _, err := v.handleLogs(ctx); err != nil {
			v.log.Error("Failed to handle logs action", "error", err)
		}
	case 'f':
		v.Refresh()
	default:
		v.log.Warn("Unknown action key", "key", string(key))
	}
}

// listServices lists all swarm services
func (v *ServicesView) listServices(ctx context.Context) ([]shared.SwarmService, error) {
	// Access service through presenter's serviceService field
	// Since presenter is internal, we need to expose this through presenter
	// For now, we'll need to add a method to presenter
	serviceService := v.GetUI().GetSwarmServiceService()
	if serviceService == nil {
		return nil, errors.New("swarm service service is not available")
	}

	if swarmService, ok := serviceService.(*ServiceService); ok {
		return swarmService.ListServices(ctx)
	}

	return nil, errors.New("swarm service service is not properly configured")
}

// formatServiceRow formats a service row for display
func (v *ServicesView) formatServiceRow(service *shared.SwarmService) []string {
	return []string{
		shared.TruncName(service.ID, 12),
		service.Name,
		service.Image,
		service.Mode,
		service.Replicas,
		service.Status,
		service.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

// getServiceID returns the service ID
func (v *ServicesView) getServiceID(service *shared.SwarmService) string {
	return service.ID
}

// getServiceName returns the service name
func (v *ServicesView) getServiceName(service *shared.SwarmService) string {
	return service.Name
}

// getActions returns the available actions for swarm services
func (v *ServicesView) getActions() map[rune]string {
	return map[rune]string{
		'i': "Inspect",
		's': "Scale",
		'r': "Remove",
		'l': "Logs",
		'f': "Refresh",
	}
}

// setupCharacterLimits sets up character limits for table columns
func (v *ServicesView) setupCharacterLimits(ui shared.SharedUIInterface) {
	// Define column types for swarm services table: ID, Name, Image, Mode, Replicas, Status, Created
	columnTypes := []string{"id", "name", "image", "mode", "replicas", "status", "created"}
	v.SetColumnTypes(columnTypes)

	// Create formatter from theme manager
	formatter := utils.NewTableFormatterFromTheme(ui.GetThemeManager())
	v.SetFormatter(formatter)
}
