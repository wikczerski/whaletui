package swarm

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/wikczerski/whaletui/internal/domains/swarm"
	"github.com/wikczerski/whaletui/internal/logger"
	"github.com/wikczerski/whaletui/internal/shared"
	"github.com/wikczerski/whaletui/internal/ui/interfaces"
)

// ServicesView represents the swarm services view
type ServicesView struct {
	*shared.BaseView[shared.SwarmService]
	serviceService *swarm.ServiceService
	modalManager   interfaces.ModalManagerInterface
	headerManager  interfaces.HeaderManagerInterface
	log            *slog.Logger
}

// NewServicesView creates a new swarm services view
func NewServicesView(
	ui interfaces.UIInterface,
	serviceService *swarm.ServiceService,
	modalManager interfaces.ModalManagerInterface,
	headerManager interfaces.HeaderManagerInterface,
) *ServicesView {
	headers := []string{"ID", "Name", "Image", "Mode", "Replicas", "Status", "Created"}

	view := &ServicesView{
		BaseView:       shared.NewBaseView[shared.SwarmService](ui, "Swarm Services", headers),
		serviceService: serviceService,
		modalManager:   modalManager,
		headerManager:  headerManager,
		log:            logger.GetLogger(),
	}

	view.setupCallbacks()

	return view
}

// Render renders the swarm services view
func (v *ServicesView) Render(_ context.Context) error {
	// The base view handles rendering automatically through the callbacks
	// Just refresh the data
	v.Refresh()
	return nil
}

// HandleInput handles user input for the services view
func (v *ServicesView) HandleInput(ctx context.Context, input rune) (interface{}, error) {
	switch input {
	case 'i':
		return v.handleInspect(ctx)
	case 's':
		return v.handleScale(ctx)
	case 'r':
		return v.handleRemove(ctx)
	case 'l':
		return v.handleLogs(ctx)
	case 'f':
		return v, nil // Refresh current view
	case 'n':
		return v.handleNavigateToNodes(ctx)
	case 'q':
		return v.handleBackToMain(ctx)
	case 'h':
		v.handleHelp()
		return v, nil
	default:
		return v, nil
	}
}

// handleInspect handles service inspection
func (v *ServicesView) handleInspect(ctx context.Context) (interface{}, error) {
	// Get the currently selected service
	selectedService := v.GetSelectedItem()
	if selectedService == nil {
		v.GetUI().ShowError(fmt.Errorf("please select a service first"))
		return v, fmt.Errorf("no service selected")
	}

	// Get the service service to access inspection functionality
	serviceService := v.GetUI().GetSwarmServiceService()
	if serviceService == nil {
		v.GetUI().ShowError(fmt.Errorf("swarm service service is not available - please check your Docker connection"))
		return v, fmt.Errorf("swarm service service not available")
	}

	// Cast to the correct type
	if swarmService, ok := serviceService.(*swarm.ServiceService); ok {
		// Get detailed service information
		serviceInfo, err := swarmService.InspectService(ctx, selectedService.ID)
		if err != nil {
			errorMsg := fmt.Sprintf("failed to inspect service '%s': %v\n\nPlease check:\n• Service is accessible\n• You have sufficient permissions\n• Docker daemon is running", selectedService.Name, err)
			v.GetUI().ShowError(fmt.Errorf("%s", errorMsg))
			return v, fmt.Errorf("failed to inspect service: %w", err)
		}

		// Show service information in a modal
		if serviceInfo == nil {
			v.GetUI().ShowInfo(fmt.Sprintf("Service '%s' has no detailed information available", selectedService.Name))
		} else {
			// Format service information for display
			infoText := fmt.Sprintf("Service Details: %s\n\n"+
				"ID: %s\n"+
				"Image: %s\n"+
				"Mode: %s\n"+
				"Replicas: %s\n"+
				"Status: %s\n"+
				"Created: %s\n"+
				"Updated: %s",
				selectedService.Name, shared.TruncName(selectedService.ID, 12), selectedService.Image, selectedService.Mode,
				selectedService.Replicas, selectedService.Status, selectedService.CreatedAt.Format("2006-01-02 15:04:05"),
				selectedService.UpdatedAt.Format("2006-01-02 15:04:05"))

			v.GetUI().ShowInfo(infoText)
		}

		return v, nil
	}

	v.GetUI().ShowError(fmt.Errorf("swarm service service is not properly configured"))
	return v, fmt.Errorf("swarm service service not available")
}

// handleScale handles service scaling
func (v *ServicesView) handleScale(ctx context.Context) (interface{}, error) {
	// Get the currently selected service
	selectedService := v.GetSelectedItem()
	if selectedService == nil {
		v.GetUI().ShowError(fmt.Errorf("please select a service first"))
		return v, fmt.Errorf("no service selected")
	}

	// Get the service service to access scaling functionality
	serviceService := v.GetUI().GetSwarmServiceService()
	if serviceService == nil {
		v.GetUI().ShowError(fmt.Errorf("swarm service service is not available - please check your Docker connection"))
		return v, fmt.Errorf("swarm service service not available")
	}

	// Cast to the correct type
	if swarmService, ok := serviceService.(*swarm.ServiceService); ok {
		// Get current service info to get replicas
		serviceInfo, err := swarmService.InspectService(ctx, selectedService.ID)
		if err != nil {
			errorMsg := fmt.Sprintf("failed to get service info for '%s': %v", selectedService.Name, err)
			v.GetUI().ShowError(fmt.Errorf("%s", errorMsg))
			return v, fmt.Errorf("failed to get service info: %w", err)
		}

		// Extract current replicas from service info
		var currentReplicas uint64 = 1
		if replicas, ok := serviceInfo["Replicas"].(uint64); ok {
			currentReplicas = replicas
		}

		// Show scale modal
		v.GetUI().ShowServiceScaleModal(selectedService.Name, currentReplicas, func(newReplicas int) {
			// Callback when user confirms scaling
			if newReplicas < 0 {
				v.GetUI().ShowError(fmt.Errorf("invalid replica count: %d", newReplicas))
				return
			}
			err := swarmService.ScaleService(ctx, selectedService.ID, uint64(newReplicas))
			if err != nil {
				// Enhanced error handling with retry and fallback options
				v.handleScaleError(ctx, selectedService, newReplicas, err, swarmService)
			} else {
				// Show success feedback and refresh
				v.Refresh()
			}
		})

		return v, nil
	}

	v.GetUI().ShowError(fmt.Errorf("swarm service service is not properly configured"))
	return v, fmt.Errorf("swarm service service not available")
}

// handleScaleError handles scaling errors with advanced recovery options
func (v *ServicesView) handleScaleError(ctx context.Context, service *shared.SwarmService, newReplicas int, err error, swarmService *swarm.ServiceService) {
	// Check if this is a retryable error
	if v.isRetryableError(err) {
		// Show retry dialog with automatic retry option
		v.GetUI().ShowRetryDialog(
			fmt.Sprintf("scale service '%s' to %d replicas", service.Name, newReplicas),
			err,
			func() error {
				// Retry function
				if newReplicas < 0 {
					return fmt.Errorf("invalid replica count: %d", newReplicas)
				}
				return swarmService.ScaleService(ctx, service.ID, uint64(newReplicas))
			},
			func() {
				// Success callback
				v.GetUI().ShowInfo(fmt.Sprintf("Service '%s' successfully scaled to %d replicas", service.Name, newReplicas))
				v.Refresh()
			},
		)
	} else {
		// Show fallback options for non-retryable errors
		fallbackOptions := []string{
			"Check Service Status",
			"View Service Logs",
			"Try Different Replica Count",
			"Show Service Details",
		}

		v.GetUI().ShowFallbackDialog(
			fmt.Sprintf("scale service '%s' to %d replicas", service.Name, newReplicas),
			err,
			fallbackOptions,
			func(fallbackOption string) {
				v.executeScaleFallback(ctx, service, fallbackOption, newReplicas, swarmService)
			},
		)
	}
}

// isRetryableError determines if an error is retryable
func (v *ServicesView) isRetryableError(err error) bool {
	errStr := err.Error()

	// Common retryable errors
	retryablePatterns := []string{
		"connection refused",
		"timeout",
		"temporary failure",
		"network unreachable",
		"service unavailable",
		"too many requests",
		"rate limit exceeded",
	}

	for _, pattern := range retryablePatterns {
		if strings.Contains(strings.ToLower(errStr), pattern) {
			return true
		}
	}

	return false
}

// executeScaleFallback executes fallback operations for scaling failures
func (v *ServicesView) executeScaleFallback(ctx context.Context, service *shared.SwarmService, fallbackOption string, _ int, swarmService *swarm.ServiceService) {
	switch fallbackOption {
	case "Check Service Status":
		// Show current service status
		serviceInfo, err := swarmService.InspectService(ctx, service.ID)
		if err != nil {
			v.GetUI().ShowError(fmt.Errorf("failed to check service status: %v", err))
		} else {
			statusText := fmt.Sprintf("Service '%s' Status:\n\n"+
				"Current Replicas: %v\n"+
				"Status: %v\n"+
				"Last Error: %v",
				service.Name, serviceInfo["Replicas"], serviceInfo["Status"], serviceInfo["LastError"])
			v.GetUI().ShowInfo(statusText)
		}

	case "View Service Logs":
		// Show service logs for troubleshooting
		logs, err := swarmService.GetServiceLogs(ctx, service.ID)
		if err != nil {
			v.GetUI().ShowError(fmt.Errorf("failed to get service logs: %v", err))
		} else {
			logPreview := logs
			if len(logs) > 500 {
				logPreview = logs[:500] + "\n\n... (truncated, logs too long for preview)"
			}
			v.GetUI().ShowInfo(fmt.Sprintf("Service '%s' Logs:\n\n%s", service.Name, logPreview))
		}

	case "Try Different Replica Count":
		// Show scale modal again with current replicas
		serviceInfo, err := swarmService.InspectService(ctx, service.ID)
		if err != nil {
			v.GetUI().ShowError(fmt.Errorf("failed to get service info: %v", err))
			return
		}

		var currentReplicas uint64 = 1
		if replicas, ok := serviceInfo["Replicas"].(uint64); ok {
			currentReplicas = replicas
		}

		v.GetUI().ShowServiceScaleModal(service.Name, currentReplicas, func(differentReplicas int) {
			if differentReplicas < 0 {
				v.GetUI().ShowError(fmt.Errorf("invalid replica count: %d", differentReplicas))
				return
			}
			err := swarmService.ScaleService(ctx, service.ID, uint64(differentReplicas))
			if err != nil {
				v.GetUI().ShowError(fmt.Errorf("failed to scale service to %d replicas: %v", differentReplicas, err))
			} else {
				v.GetUI().ShowInfo(fmt.Sprintf("Service '%s' successfully scaled to %d replicas", service.Name, differentReplicas))
				v.Refresh()
			}
		})

	case "Show Service Details":
		// Show detailed service information
		serviceInfo, err := swarmService.InspectService(ctx, service.ID)
		if err != nil {
			v.GetUI().ShowError(fmt.Errorf("failed to get service details: %v", err))
		} else {
			detailsText := fmt.Sprintf("Service '%s' Details:\n\n"+
				"ID: %s\n"+
				"Image: %s\n"+
				"Mode: %s\n"+
				"Replicas: %v\n"+
				"Status: %v\n"+
				"Created: %s\n"+
				"Updated: %s",
				service.Name, shared.TruncName(service.ID, 12), service.Image, service.Mode,
				serviceInfo["Replicas"], serviceInfo["Status"],
				service.CreatedAt.Format("2006-01-02 15:04:05"),
				service.UpdatedAt.Format("2006-01-02 15:04:05"))
			v.GetUI().ShowInfo(detailsText)
		}
	}
}

// handleRemove handles service removal
func (v *ServicesView) handleRemove(ctx context.Context) (interface{}, error) {
	// Get the currently selected service
	selectedService := v.GetSelectedItem()
	if selectedService == nil {
		v.GetUI().ShowError(fmt.Errorf("please select a service first"))
		return v, fmt.Errorf("no service selected")
	}

	// Show confirmation dialog with more detailed information
	message := fmt.Sprintf("⚠️  Remove Service Confirmation\n\n"+
		"Service: %s\n"+
		"ID: %s\n"+
		"Image: %s\n\n"+
		"This action will:\n"+
		"• Stop all running tasks\n"+
		"• Remove the service definition\n"+
		"• Cannot be undone\n\n"+
		"Are you sure you want to continue?",
		selectedService.Name, shared.TruncName(selectedService.ID, 12), selectedService.Image)

	v.GetUI().ShowConfirm(message, func(confirmed bool) {
		if confirmed {
			// Get the service service to access removal functionality
			serviceService := v.GetUI().GetSwarmServiceService()
			if serviceService == nil {
				v.GetUI().ShowError(fmt.Errorf("swarm service service is not available - please check your Docker connection"))
				return
			}

			// Cast to the correct type
			if swarmService, ok := serviceService.(swarm.ServiceService); ok {
				// Remove the service
				err := swarmService.RemoveService(ctx, selectedService.ID)
				if err != nil {
					// Show detailed error modal with recovery suggestions
					errorMsg := fmt.Sprintf("failed to remove service '%s': %v\n\nPlease check:\n• Docker daemon is running\n• You have sufficient permissions\n• Service is not in a critical state", selectedService.Name, err)
					v.GetUI().ShowError(fmt.Errorf("%s", errorMsg))
				} else {
					// Show success feedback and refresh
					v.GetUI().ShowInfo(fmt.Sprintf("Service '%s' successfully removed", selectedService.Name))
					v.Refresh()
				}
			} else {
				v.GetUI().ShowError(fmt.Errorf("swarm service service is not properly configured"))
			}
		}
	})

	return v, nil
}

// handleLogs handles service logs display
func (v *ServicesView) handleLogs(ctx context.Context) (interface{}, error) {
	// Get the currently selected service
	selectedService := v.GetSelectedItem()
	if selectedService == nil {
		v.GetUI().ShowError(fmt.Errorf("please select a service first"))
		return v, fmt.Errorf("no service selected")
	}

	// Get the service service to access logs functionality
	serviceService := v.GetUI().GetSwarmServiceService()
	if serviceService == nil {
		v.GetUI().ShowError(fmt.Errorf("swarm service service is not available - please check your Docker connection"))
		return v, fmt.Errorf("swarm service service not available")
	}

	// Cast to the correct type
	if swarmService, ok := serviceService.(*swarm.ServiceService); ok {
		// Get service logs
		logs, err := swarmService.GetServiceLogs(ctx, selectedService.ID)
		if err != nil {
			errorMsg := fmt.Sprintf("failed to get logs for service '%s': %v\n\nPlease check:\n• Service is running\n• You have sufficient permissions\n• Docker daemon is accessible", selectedService.Name, err)
			v.GetUI().ShowError(fmt.Errorf("%s", errorMsg))
			return v, fmt.Errorf("failed to get service logs: %w", err)
		}

		// Show logs in a modal or dedicated view
		if logs == "" {
			v.GetUI().ShowInfo(fmt.Sprintf("Service '%s' has no logs available", selectedService.Name))
		} else {
			// For now, show logs in an info modal (in future, could be a dedicated logs view)
			logPreview := logs
			if len(logs) > 500 {
				logPreview = logs[:500] + "\n\n... (truncated, logs too long for preview)"
			}
			v.GetUI().ShowInfo(fmt.Sprintf("Service '%s' Logs:\n\n%s", selectedService.Name, logPreview))
		}

		return v, nil
	}

	v.GetUI().ShowError(fmt.Errorf("swarm service service is not properly configured"))
	return v, fmt.Errorf("swarm service service not available")
}

// handleNavigateToNodes handles navigation to swarm nodes view
func (v *ServicesView) handleNavigateToNodes(_ context.Context) (interface{}, error) {
	// This would return a nodes view - placeholder for now
	return v, fmt.Errorf("nodes view not implemented yet")
}

// handleBackToMain handles navigation back to main menu
func (v *ServicesView) handleBackToMain(_ context.Context) (interface{}, error) {
	// This would return the main menu view - placeholder for now
	return v, fmt.Errorf("main menu view not implemented yet")
}

// handleHelp shows contextual help for swarm services
func (v *ServicesView) handleHelp() {
	// Show general swarm services help
	v.GetUI().ShowContextualHelp("swarm_services", "")
}

// showOperationHelp shows contextual help for a specific operation
// This function is intentionally unused - it's a placeholder for future help functionality
// nolint:unused // Intentionally unused - placeholder for future help functionality
func (v *ServicesView) showOperationHelp(operation string) {
	v.GetUI().ShowContextualHelp("swarm_services", operation)
}

// setupCallbacks sets up the callbacks for the base view
func (v *ServicesView) setupCallbacks() {
	v.ListItems = v.listServices
	v.FormatRow = func(s shared.SwarmService) []string { return v.formatServiceRow(&s) }
	v.GetItemID = func(s shared.SwarmService) string { return v.getServiceID(&s) }
	v.GetItemName = func(s shared.SwarmService) string { return v.getServiceName(&s) }
	v.GetActions = v.getActions
}

// listServices lists all swarm services
func (v *ServicesView) listServices(ctx context.Context) ([]shared.SwarmService, error) {
	return v.serviceService.ListServices(ctx)
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
