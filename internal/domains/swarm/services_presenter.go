package swarm

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/wikczerski/whaletui/internal/shared"
)

// ServicesPresenter handles business logic for swarm services
type ServicesPresenter struct {
	serviceService *ServiceService
	log            *slog.Logger
}

// NewServicesPresenter creates a new services presenter
func NewServicesPresenter(serviceService *ServiceService, log *slog.Logger) *ServicesPresenter {
	return &ServicesPresenter{
		serviceService: serviceService,
		log:            log,
	}
}

// ValidateService validates that a service is not nil
func (p *ServicesPresenter) ValidateService(service *shared.SwarmService) error {
	if service == nil {
		return errors.New("no service selected")
	}
	return nil
}

// InspectService inspects a service and returns its information
func (p *ServicesPresenter) InspectService(
	ctx context.Context,
	serviceID string,
) (map[string]any, error) {
	serviceInfo, err := p.serviceService.InspectService(ctx, serviceID)
	if err != nil {
		return nil, fmt.Errorf("failed to inspect service: %w", err)
	}
	return serviceInfo, nil
}

// FormatServiceInfo formats service information for display
func (p *ServicesPresenter) FormatServiceInfo(
	service *shared.SwarmService,
	serviceInfo map[string]any,
) string {
	if serviceInfo == nil {
		return fmt.Sprintf("Service '%s' has no detailed information available", service.Name)
	}

	return fmt.Sprintf("Service Details: %s\n\n"+
		"ID: %s\n"+
		"Image: %s\n"+
		"Mode: %s\n"+
		"Replicas: %s\n"+
		"Status: %s\n"+
		"Created: %s\n"+
		"Updated: %s",
		service.Name,
		shared.TruncName(service.ID, 12),
		service.Image,
		service.Mode,
		service.Replicas,
		service.Status,
		service.CreatedAt.Format("2006-01-02 15:04:05"),
		service.UpdatedAt.Format("2006-01-02 15:04:05"))
}

// GetCurrentReplicas gets the current replica count for a service
func (p *ServicesPresenter) GetCurrentReplicas(
	ctx context.Context,
	serviceID string,
) (uint64, error) {
	serviceInfo, err := p.serviceService.InspectService(ctx, serviceID)
	if err != nil {
		return 1, fmt.Errorf("failed to get service info: %w", err)
	}

	var currentReplicas uint64 = 1
	if replicas, ok := serviceInfo["Replicas"].(uint64); ok {
		currentReplicas = replicas
	}
	return currentReplicas, nil
}

// ScaleService scales a service to the specified number of replicas
func (p *ServicesPresenter) ScaleService(
	ctx context.Context,
	serviceID string,
	replicas uint64,
) error {
	return p.serviceService.ScaleService(ctx, serviceID, replicas)
}

// IsRetryableError determines if an error is retryable
func (p *ServicesPresenter) IsRetryableError(err error) bool {
	if err == nil {
		return false
	}

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

// GetScaleFallbackOptions returns the available fallback options for scaling failures
func (p *ServicesPresenter) GetScaleFallbackOptions() []string {
	return []string{
		"Check Service Status",
		"View Service Logs",
		"Try Different Replica Count",
		"Show Service Details",
	}
}

// ExecuteScaleFallback executes a fallback operation for scaling failures
func (p *ServicesPresenter) ExecuteScaleFallback(
	ctx context.Context,
	service *shared.SwarmService,
	fallbackOption string,
) (string, error) {
	switch fallbackOption {
	case "Check Service Status":
		return p.checkServiceStatus(ctx, service)
	case "View Service Logs":
		return p.getServiceLogsPreview(ctx, service)
	case "Show Service Details":
		return p.getServiceDetails(ctx, service)
	default:
		return "", fmt.Errorf("unknown fallback option: %s", fallbackOption)
	}
}

// RemoveService removes a service
func (p *ServicesPresenter) RemoveService(ctx context.Context, serviceID string) error {
	return p.serviceService.RemoveService(ctx, serviceID)
}

// BuildRemoveConfirmationMessage builds the remove confirmation message
func (p *ServicesPresenter) BuildRemoveConfirmationMessage(service *shared.SwarmService) string {
	return fmt.Sprintf("⚠️  Remove Service Confirmation\n\n"+
		"Service: %s\n"+
		"ID: %s\n"+
		"Image: %s\n\n"+
		"This action will:\n"+
		"• Stop all running tasks\n"+
		"• Remove the service definition\n"+
		"• Cannot be undone\n\n"+
		"Are you sure you want to continue?",
		service.Name, shared.TruncName(service.ID, 12), service.Image)
}

// GetServiceLogs retrieves service logs
func (p *ServicesPresenter) GetServiceLogs(ctx context.Context, serviceID string) (string, error) {
	logs, err := p.serviceService.GetServiceLogs(ctx, serviceID)
	if err != nil {
		return "", fmt.Errorf("failed to get service logs: %w", err)
	}
	return logs, nil
}

// FormatServiceLogs formats service logs for display
func (p *ServicesPresenter) FormatServiceLogs(serviceName, logs string) string {
	if logs == "" {
		return fmt.Sprintf("Service '%s' has no logs available", serviceName)
	}

	logPreview := logs
	if len(logs) > 500 {
		logPreview = logs[:500] + "\n\n... (truncated, logs too long for preview)"
	}
	return fmt.Sprintf("Service '%s' Logs:\n\n%s", serviceName, logPreview)
}

// FormatInspectError formats an inspection error message
func (p *ServicesPresenter) FormatInspectError(serviceName string, err error) string {
	return fmt.Sprintf(
		"failed to inspect service '%s': %v\n\nPlease check:\n"+
			"• Service is accessible\n"+
			"• You have sufficient permissions\n"+
			"• Docker daemon is running",
		serviceName,
		err,
	)
}

// FormatLogsError formats a logs retrieval error message
func (p *ServicesPresenter) FormatLogsError(serviceName string, err error) string {
	return fmt.Sprintf(
		"failed to get logs for service '%s': %v\n\nPlease check:\n"+
			"• Service is running\n"+
			"• You have sufficient permissions\n"+
			"• Docker daemon is accessible",
		serviceName,
		err,
	)
}

// FormatRemoveError formats a removal error message
func (p *ServicesPresenter) FormatRemoveError(serviceName string, err error) string {
	return fmt.Sprintf(
		"failed to remove service '%s': %v\n\nPlease check:\n"+
			"• Docker daemon is running\n"+
			"• You have sufficient permissions\n"+
			"• Service is not in a critical state",
		serviceName,
		err,
	)
}

// checkServiceStatus checks the service status
func (p *ServicesPresenter) checkServiceStatus(
	ctx context.Context,
	service *shared.SwarmService,
) (string, error) {
	serviceInfo, err := p.serviceService.InspectService(ctx, service.ID)
	if err != nil {
		return "", fmt.Errorf("failed to check service status: %w", err)
	}

	return fmt.Sprintf("Service '%s' Status:\n\n"+
		"Current Replicas: %v\n"+
		"Status: %v\n"+
		"Last Error: %v",
		service.Name, serviceInfo["Replicas"], serviceInfo["Status"], serviceInfo["LastError"]), nil
}

// getServiceLogsPreview gets a preview of service logs
func (p *ServicesPresenter) getServiceLogsPreview(
	ctx context.Context,
	service *shared.SwarmService,
) (string, error) {
	logs, err := p.serviceService.GetServiceLogs(ctx, service.ID)
	if err != nil {
		return "", fmt.Errorf("failed to get service logs: %w", err)
	}

	logPreview := logs
	if len(logs) > 500 {
		logPreview = logs[:500] + "\n\n... (truncated, logs too long for preview)"
	}
	return fmt.Sprintf("Service '%s' Logs:\n\n%s", service.Name, logPreview), nil
}

// getServiceDetails gets detailed service information
func (p *ServicesPresenter) getServiceDetails(
	ctx context.Context,
	service *shared.SwarmService,
) (string, error) {
	serviceInfo, err := p.serviceService.InspectService(ctx, service.ID)
	if err != nil {
		return "", fmt.Errorf("failed to get service details: %w", err)
	}

	return fmt.Sprintf("Service '%s' Details:\n\n"+
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
		service.UpdatedAt.Format("2006-01-02 15:04:05")), nil
}
