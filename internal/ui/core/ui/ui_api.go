package ui

import (
	"errors"
	"fmt"

	"github.com/rivo/tview"
	"github.com/wikczerski/whaletui/internal/config"
	"github.com/wikczerski/whaletui/internal/domains/logs"
	"github.com/wikczerski/whaletui/internal/ui/interfaces"
)

// Start starts the UI
func (ui *UI) Start() error {
	ui.log.Info("Starting TUI...")
	return ui.app.Run()
}

// Stop stops the UI
func (ui *UI) Stop() {
	ui.cleanup()
	ui.app.Stop()
}

// GetShutdownChan returns the shutdown channel
func (ui *UI) GetShutdownChan() chan struct{} {
	return ui.shutdownChan
}

// GetServices returns the service factory (for interfaces.UIInterface)
func (ui *UI) GetServices() interfaces.ServiceFactoryInterface {
	ui.log.Debug("GetServices called", "services_nil", ui.services == nil)
	return ui.services
}

// GetServicesAny returns the service factory as any (for shared.UIInterface)
func (ui *UI) GetServicesAny() any {
	ui.log.Debug("GetServices called (any)", "services_nil", ui.services == nil)
	return ui.services
}

// GetSwarmServiceService returns the swarm service service (for shared.UIInterface)
func (ui *UI) GetSwarmServiceService() any {
	if ui.services != nil {
		return ui.services.GetSwarmServiceService()
	}
	return nil
}

// GetSwarmNodeService returns the swarm node service (for shared.UIInterface)
func (ui *UI) GetSwarmNodeService() any {
	if ui.services != nil {
		return ui.services.GetSwarmNodeService()
	}
	return nil
}

// GetContainerService returns the container service (for shared.UIInterface)
func (ui *UI) GetContainerService() interfaces.ContainerService {
	if ui.services != nil {
		return ui.services.GetContainerService()
	}
	return nil
}

// GetApp returns the tview application
func (ui *UI) GetApp() any {
	return ui.app
}

// ShowError displays an error message
func (ui *UI) ShowError(err error) {
	ui.showError(err)
}

// ShowInfo displays an info message
func (ui *UI) ShowInfo(message string) {
	ui.showInfo(message)
}

// ShowContextualHelp displays context-sensitive help based on current operation
func (ui *UI) ShowContextualHelp(context, operation string) {
	ui.showContextualHelp(context, operation)
}

// ShowRetryDialog displays a retry dialog with automatic retry logic
func (ui *UI) ShowRetryDialog(
	operation string,
	err error,
	retryFunc func() error,
	onSuccess func(),
) {
	ui.showRetryDialog(operation, err, retryFunc, onSuccess)
}

// ShowFallbackDialog displays a fallback operations dialog
func (ui *UI) ShowFallbackDialog(
	operation string,
	err error,
	fallbackOptions []string,
	onFallback func(string),
) {
	ui.showFallbackDialog(operation, err, fallbackOptions, onFallback)
}

// ShowDetails displays a details view
func (ui *UI) ShowDetails(details any) {
	if detailsView, ok := details.(tview.Primitive); ok {
		ui.showDetails(detailsView)
	} else {
		ui.log.Warn("ShowDetails called with non-Primitive type", "type", fmt.Sprintf("%T", details))
	}
}

// ShowCurrentView returns to the current view
func (ui *UI) ShowCurrentView() {
	ui.showCurrentView()
}

// ShowConfirm shows a confirmation dialog
func (ui *UI) ShowConfirm(message string, onConfirm func(bool)) {
	ui.showConfirm(message, onConfirm)
}

// ShowServiceScaleModal shows a modal for scaling swarm services
func (ui *UI) ShowServiceScaleModal(
	serviceName string,
	currentReplicas uint64,
	callback func(int),
) {
	ui.modalManager.ShowServiceScaleModal(serviceName, currentReplicas, callback)
}

// ShowNodeAvailabilityModal shows a modal for updating node availability
func (ui *UI) ShowNodeAvailabilityModal(
	nodeName, currentAvailability string,
	callback func(string),
) {
	ui.modalManager.ShowNodeAvailabilityModal(nodeName, currentAvailability, callback)
}

// IsInLogsMode returns whether the UI is currently in logs mode
func (ui *UI) IsInLogsMode() bool {
	return ui.inLogsMode
}

// IsInDetailsMode returns whether the UI is currently in details mode
func (ui *UI) IsInDetailsMode() bool {
	return ui.inDetailsMode
}

// IsModalActive returns whether a modal is currently active
func (ui *UI) IsModalActive() bool {
	if !ui.hasValidPages() {
		return false
	}

	return ui.hasModalPages()
}

// GetCurrentActions returns the current available actions
func (ui *UI) GetCurrentActions() map[rune]string {
	return ui.currentActions
}

// GetCurrentViewActions returns the actions string from the current view
func (ui *UI) GetCurrentViewActions() string {
	if ui.viewRegistry != nil {
		return ui.viewRegistry.GetCurrentActionsString()
	}
	return ""
}

// GetCurrentViewNavigation returns the navigation string from the current view
func (ui *UI) GetCurrentViewNavigation() string {
	if ui.viewRegistry != nil {
		return ui.viewRegistry.GetCurrentNavigationString()
	}
	return ""
}

// GetViewRegistry returns the view registry
func (ui *UI) GetViewRegistry() any {
	return ui.viewRegistry
}

// GetMainFlex returns the main flex container
func (ui *UI) GetMainFlex() any {
	return ui.mainFlex
}

// SwitchView switches to the specified view
func (ui *UI) SwitchView(view string) {
	ui.switchView(view)
}

// ShowHelp shows the help dialog
func (ui *UI) ShowHelp() {
	ui.showHelp()
}

// GetPages returns the pages container
func (ui *UI) GetPages() any {
	return ui.pages
}

// ShowLogs shows logs for a container
func (ui *UI) ShowLogs(containerID, containerName string) {
	ui.showLogs("container", containerID, containerName)
}

// ShowShell shows shell view for a container
func (ui *UI) ShowShell(containerID, containerName string) {
	ui.createShellView(containerID, containerName)
	ui.displayShellView(containerID, containerName)
}

// GetLogsView returns the logs view for any resource type
func (ui *UI) GetLogsView(resourceType, resourceID, resourceName string) *logs.View {
	if ui.logsView == nil || ui.logsView.ResourceID != resourceID ||
		ui.logsView.ResourceType != resourceType {
		ui.logsView = logs.NewView(ui, resourceType, resourceID, resourceName)
	}
	return ui.logsView
}

// GetViewContainer returns the view container
func (ui *UI) GetViewContainer() any {
	return ui.viewContainer
}

// GetThemeManager returns the theme manager
func (ui *UI) GetThemeManager() *config.ThemeManager {
	return ui.themeManager
}

// ReloadTheme reloads the theme configuration and refreshes all views
func (ui *UI) ReloadTheme() error {
	if ui.themeManager == nil {
		return errors.New("theme manager not initialized")
	}

	// Reload the theme from file
	err := ui.themeManager.ReloadTheme()
	if err != nil {
		return fmt.Errorf("failed to reload theme: %w", err)
	}

	// Refresh all views to apply the new character limits
	ui.refreshAllViews()

	ui.log.Info("Theme reloaded successfully")
	return nil
}

// refreshAllViews refreshes all registered views
func (ui *UI) refreshAllViews() {
	// Refresh the current view
	currentView := ui.viewRegistry.GetCurrent()
	if currentView != nil && currentView.Refresh != nil {
		currentView.Refresh()
	}

	// Update headers to reflect any theme changes
	ui.headerManager.UpdateDockerInfo()
	ui.headerManager.UpdateNavigation()
	ui.headerManager.UpdateActions()
}

// SetHeaderManager sets the header manager
func (ui *UI) SetHeaderManager(headerManager interfaces.HeaderManagerInterface) {
	ui.headerManager = headerManager
}

// SetModalManager sets the modal manager
func (ui *UI) SetModalManager(modalManager interfaces.ModalManagerInterface) {
	ui.modalManager = modalManager
}

// Refresh refreshes the UI
func (ui *UI) Refresh() {
	ui.log.Debug("Refreshing UI")

	// Set a flag to prevent header updates during refresh cycles
	ui.isRefreshing = true
	defer func() {
		ui.isRefreshing = false
		ui.log.Debug("Refresh completed, isRefreshing set to false")
	}()

	ui.log.Debug("Starting refresh cycle", "isRefreshing", ui.isRefreshing)

	// Ensure layout stability before refreshing
	ui.ensureStableLayout()

	// Only update components that actually need refreshing
	// This prevents unnecessary terminal redraws that might cause empty lines
	ui.updateStatusBar()

	// Skip header updates during refresh cycles to prevent newlines from causing empty spaces
	// Headers are only updated when switching views or showing details
	// if ui.services != nil {
	// 	ui.headerManager.UpdateAll()
	// }

	// Only refresh current view if it exists and has a refresh function
	ui.refreshCurrentView()
}

// CompleteInitialization completes the UI initialization after managers are set
func (ui *UI) CompleteInitialization() error {
	if ui.headerManager == nil || ui.modalManager == nil {
		return errors.New("managers must be set before completing initialization")
	}

	ui.initComponents()
	ui.setupKeyBindings()

	return nil
}
