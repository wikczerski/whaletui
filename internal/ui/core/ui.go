package core

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/wikczerski/whaletui/internal/config"
	"github.com/wikczerski/whaletui/internal/domains/container"
	"github.com/wikczerski/whaletui/internal/domains/image"
	"github.com/wikczerski/whaletui/internal/domains/logs"
	"github.com/wikczerski/whaletui/internal/domains/network"
	swarmDomain "github.com/wikczerski/whaletui/internal/domains/swarm"
	"github.com/wikczerski/whaletui/internal/domains/volume"
	"github.com/wikczerski/whaletui/internal/logger"
	"github.com/wikczerski/whaletui/internal/shared"
	"github.com/wikczerski/whaletui/internal/ui/builders"
	"github.com/wikczerski/whaletui/internal/ui/constants"
	"github.com/wikczerski/whaletui/internal/ui/handlers"
	"github.com/wikczerski/whaletui/internal/ui/interfaces"
	"github.com/wikczerski/whaletui/internal/ui/managers"
	"github.com/wikczerski/whaletui/internal/ui/views/shell"
	"github.com/wikczerski/whaletui/internal/ui/views/swarm"
)

// UI represents the main UI
type UI struct {
	app            *tview.Application
	pages          *tview.Pages
	mainFlex       *tview.Flex
	statusBar      *tview.TextView
	viewContainer  *tview.Flex
	services       interfaces.ServiceFactoryInterface
	log            *slog.Logger
	shutdownChan   chan struct{}
	inDetailsMode  bool            // Track if we're in details view mode
	inLogsMode     bool            // Track if we're viewing container logs
	currentActions map[rune]string // Track current available actions

	// Theme management
	themeManager *config.ThemeManager

	// Abstracted managers
	viewRegistry   *ViewRegistry
	headerManager  interfaces.HeaderManagerInterface
	commandHandler *handlers.CommandHandler
	modalManager   interfaces.ModalManagerInterface

	// Individual views
	containersView *container.ContainersView
	imagesView     *image.ImagesView
	volumesView    *volume.VolumesView
	networksView   *network.NetworksView
	logsView       *logs.View
	shellView      *shell.View

	// Swarm views
	swarmServicesView *swarm.ServicesView
	swarmNodesView    *swarm.NodesView

	// Component builders
	componentBuilder *builders.ComponentBuilder
	viewBuilder      *builders.ViewBuilder
	tableBuilder     *builders.TableBuilder

	// Flags for refresh cycles
	isRefreshing bool

	// UI components
	headerSection *tview.Flex
	commandInput  *tview.InputField
}

// New creates a new UI
func New(
	serviceFactory interfaces.ServiceFactoryInterface,
	themePath string,
	headerManager interfaces.HeaderManagerInterface,
	modalManager interfaces.ModalManagerInterface,
	_ *config.Config,
) (*UI, error) {
	// TUI mode is already set globally in init()

	ui := &UI{
		services:       serviceFactory,
		app:            tview.NewApplication(),
		pages:          tview.NewPages(),
		log:            logger.GetLogger(),
		shutdownChan:   make(chan struct{}, 1), // Buffer channel to prevent deadlock
		currentActions: make(map[rune]string),
		headerManager:  headerManager,
		modalManager:   modalManager,
	}

	if e := ui.initializeManagers(themePath); e != nil {
		return nil, e
	}

	// Only initialize components if managers are provided
	if headerManager != nil && modalManager != nil {
		ui.initComponents()
		ui.setupKeyBindings()
	}

	return ui, nil
}

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

// IsContainerServiceAvailable checks if container service is available (for shared.UIInterface)
func (ui *UI) IsContainerServiceAvailable() bool {
	if ui.services != nil {
		return ui.services.IsContainerServiceAvailable()
	}
	return false
}

// GetContainerService returns the container service (for shared.UIInterface)
func (ui *UI) GetContainerService() any {
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

// IsRefreshing returns whether the UI is currently in a refresh cycle
func (ui *UI) IsRefreshing() bool {
	return ui.isRefreshing
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

// ShowLogsForResource shows logs for any resource type
func (ui *UI) ShowLogsForResource(resourceType, resourceID, resourceName string) {
	ui.showLogs(resourceType, resourceID, resourceName)
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

// GetImageService returns the image service
func (ui *UI) GetImageService() any {
	return ui.services.GetImageService()
}

// GetVolumeService returns the volume service
func (ui *UI) GetVolumeService() any {
	return ui.services.GetVolumeService()
}

// GetNetworkService returns the network service
func (ui *UI) GetNetworkService() any {
	return ui.services.GetNetworkService()
}

// GetThemeManager returns the theme manager
func (ui *UI) GetThemeManager() *config.ThemeManager {
	return ui.themeManager
}

// SetHeaderManager sets the header manager
func (ui *UI) SetHeaderManager(headerManager interfaces.HeaderManagerInterface) {
	ui.headerManager = headerManager
}

// SetModalManager sets the modal manager
func (ui *UI) SetModalManager(modalManager interfaces.ModalManagerInterface) {
	ui.modalManager = modalManager
}

// UpdateLegend updates the legend with current view information
func (ui *UI) UpdateLegend() {
	ui.updateLegend()
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

// initializeManagers initializes all the UI managers and builders
func (ui *UI) initializeManagers(themePath string) error {
	ui.themeManager = config.NewThemeManager(themePath)

	ui.componentBuilder = builders.NewComponentBuilderWithTheme(ui.themeManager)
	ui.viewBuilder = builders.NewViewBuilder()
	ui.tableBuilder = builders.NewTableBuilder()

	ui.viewRegistry = NewViewRegistry()
	// Managers are now passed as parameters to avoid circular imports
	ui.commandHandler = handlers.NewCommandHandler(ui)

	return nil
}

// initComponents initializes UI components
func (ui *UI) initComponents() {
	ui.mainFlex = ui.viewBuilder.CreateView()
	ui.setupHeaderAndCommandInput()
	ui.createAndRegisterViews()
	ui.createViewContainer()
	ui.createStatusBar()
	ui.setupMainLayout()
	ui.setupMainPages(ui.commandInput)
	ui.initializeUIState()
	ui.log.Info("UI components initialized")
}

// setupHeaderAndCommandInput sets up the header section and command input
func (ui *UI) setupHeaderAndCommandInput() {
	headerSection, ok := ui.headerManager.CreateHeaderSection().(*tview.Flex)
	if !ok {
		ui.log.Error("failed to create header section")
		return
	}
	ui.headerSection = headerSection
	// Force initial header update to populate content
	ui.headerManager.UpdateDockerInfo()
	ui.headerManager.UpdateNavigation()
	ui.headerManager.UpdateActions()
	ui.commandInput = ui.commandHandler.CreateCommandInput()
}

// setupMainLayout sets up the main layout with proper heights and direction
func (ui *UI) setupMainLayout() {
	// Ensure proper layout with fixed heights to prevent shifting
	ui.mainFlex.AddItem(ui.headerSection, constants.HeaderSectionHeight, 1, false)
	ui.mainFlex.AddItem(ui.viewContainer, 0, 1, true)
	ui.mainFlex.AddItem(ui.statusBar, constants.StatusBarHeight, 1, false)
	// Ensure layout stability
	ui.mainFlex.SetDirection(tview.FlexRow)
}

// setupMainPages sets up the main pages in the UI
func (ui *UI) setupMainPages(commandInput *tview.InputField) {
	ui.pages.AddPage("main", ui.mainFlex, true, true)
	ui.pages.AddPage("command", commandInput, true, false)

	ui.app.SetRoot(ui.pages, true)
}

// initializeUIState initializes the initial UI state
func (ui *UI) initializeUIState() {
	ui.log.Debug("Starting UI state initialization")

	if ui.services != nil {
		ui.log.Debug("Services available, updating headers")
		ui.headerManager.UpdateDockerInfo()
		ui.headerManager.UpdateNavigation()
		ui.headerManager.UpdateActions()
	} else {
		ui.log.Debug("No services available")
	}

	// Perform initial view refresh to populate data
	ui.log.Debug("Performing initial view refresh")
	ui.refreshCurrentView()
}

// createAndRegisterViews creates all views and registers them with the view registry
func (ui *UI) createAndRegisterViews() {
	ui.createResourceViews()
	ui.registerViewsWithActions()
	// Ensure views are fully registered before setting default view
	ui.setDefaultView()
}

// createResourceViews creates all the resource views
func (ui *UI) createResourceViews() {
	ui.containersView = container.NewContainersView(ui)
	ui.imagesView = image.NewImagesView(ui)
	ui.volumesView = volume.NewVolumesView(ui)
	ui.networksView = network.NewNetworksView(ui)

	// Create swarm views
	ui.swarmServicesView = swarm.NewServicesView(
		ui,
		ui.services.GetSwarmServiceService().(*swarmDomain.ServiceService),
		ui.modalManager.(*managers.ModalManager),
		ui.headerManager.(*managers.HeaderManager),
	)
	ui.swarmNodesView = swarm.NewNodesView(
		ui,
		ui.services.GetSwarmNodeService().(*swarmDomain.NodeService),
		ui.modalManager.(*managers.ModalManager),
		ui.headerManager.(*managers.HeaderManager),
	)
}

// registerViewsWithActions registers views with their metadata and actions
func (ui *UI) registerViewsWithActions() {
	if ui.services != nil {
		ui.registerContainerView()
		ui.registerResourceViews()
	} else {
		// Register views without actions when services are not available
		ui.registerViewsWithoutServices()
	}
}

// registerContainerView registers the containers view with its actions
func (ui *UI) registerContainerView() {
	containerActions := ""
	containerNavigation := ""
	if ui.services != nil && ui.services.GetContainerService() != nil {
		if actionService, ok := ui.services.GetContainerService().(interfaces.ServiceWithActions); ok {
			containerActions = actionService.GetActionsString()
		}
		if navigationService, ok := ui.services.GetContainerService().(interfaces.ServiceWithNavigation); ok {
			containerNavigation = navigationService.GetNavigationString()
		}
	}
	ui.viewRegistry.Register(
		"containers",
		"Containers",
		'c',
		ui.containersView.GetView(),
		ui.containersView.Refresh,
		containerActions,
		containerNavigation,
	)
}

// registerResourceViews registers the resource views with their actions
func (ui *UI) registerResourceViews() {
	actions := ui.collectServiceActions()
	ui.registerResourceViewsWithActions(actions)
}

// collectServiceActions collects actions from all available services
func (ui *UI) collectServiceActions() map[string]string {
	actions := make(map[string]string)

	if ui.services == nil {
		return actions
	}

	ui.collectImageActions(actions)
	ui.collectVolumeActions(actions)
	ui.collectNetworkActions(actions)
	ui.collectSwarmServiceActions(actions)
	ui.collectSwarmNodeActions(actions)

	return actions
}

// collectImageActions collects actions from the image service
func (ui *UI) collectImageActions(actions map[string]string) {
	if imageService := ui.services.GetImageService(); imageService != nil {
		if actionService, ok := imageService.(interfaces.ServiceWithActions); ok {
			actions["images"] = actionService.GetActionsString()
		}
	}
}

// collectVolumeActions collects actions from the volume service
func (ui *UI) collectVolumeActions(actions map[string]string) {
	if volumeService := ui.services.GetVolumeService(); volumeService != nil {
		if actionService, ok := volumeService.(interfaces.ServiceWithActions); ok {
			actions["volumes"] = actionService.GetActionsString()
		}
	}
}

// collectNetworkActions collects actions from the network service
func (ui *UI) collectNetworkActions(actions map[string]string) {
	if networkService := ui.services.GetNetworkService(); networkService != nil {
		if actionService, ok := networkService.(interfaces.ServiceWithActions); ok {
			actions["networks"] = actionService.GetActionsString()
		}
	}
}

// collectSwarmServiceActions collects actions from the swarm service service
func (ui *UI) collectSwarmServiceActions(actions map[string]string) {
	if swarmServiceService := ui.services.GetSwarmServiceService(); swarmServiceService != nil {
		if actionService, ok := swarmServiceService.(interfaces.ServiceWithActions); ok {
			actions["swarmServices"] = actionService.GetActionsString()
		}
	}
}

// collectSwarmNodeActions collects actions from the swarm node service
func (ui *UI) collectSwarmNodeActions(actions map[string]string) {
	if swarmNodeService := ui.services.GetSwarmNodeService(); swarmNodeService != nil {
		if actionService, ok := swarmNodeService.(interfaces.ServiceWithActions); ok {
			actions["swarmNodes"] = actionService.GetActionsString()
		}
	}
}

// registerResourceViewsWithActions registers resource views with their collected actions
func (ui *UI) registerResourceViewsWithActions(actions map[string]string) {
	ui.registerImagesView(actions)
	ui.registerVolumesView(actions)
	ui.registerNetworksView(actions)
	ui.registerSwarmServicesView(actions)
	ui.registerSwarmNodesView(actions)
}

// registerImagesView registers the images view
func (ui *UI) registerImagesView(actions map[string]string) {
	ui.viewRegistry.Register(
		"images",
		"Images",
		'i',
		ui.imagesView.GetView(),
		ui.imagesView.Refresh,
		actions["images"],
		"",
	)
}

// registerVolumesView registers the volumes view
func (ui *UI) registerVolumesView(actions map[string]string) {
	ui.viewRegistry.Register(
		"volumes",
		"Volumes",
		'v',
		ui.volumesView.GetView(),
		ui.volumesView.Refresh,
		actions["volumes"],
		"",
	)
}

// registerNetworksView registers the networks view
func (ui *UI) registerNetworksView(actions map[string]string) {
	ui.viewRegistry.Register(
		"networks",
		"Networks",
		'n',
		ui.networksView.GetView(),
		ui.networksView.Refresh,
		actions["networks"],
		"",
	)
}

// registerSwarmServicesView registers the swarm services view
func (ui *UI) registerSwarmServicesView(actions map[string]string) {
	ui.viewRegistry.Register(
		"swarmServices",
		"Swarm Services",
		's',
		ui.swarmServicesView.GetView(),
		ui.swarmServicesView.Refresh,
		actions["swarmServices"],
		"",
	)
}

// registerSwarmNodesView registers the swarm nodes view
func (ui *UI) registerSwarmNodesView(actions map[string]string) {
	ui.viewRegistry.Register(
		"swarmNodes",
		"Swarm Nodes",
		'w',
		ui.swarmNodesView.GetView(),
		ui.swarmNodesView.Refresh,
		actions["swarmNodes"],
		"",
	)
}

// registerViewsWithoutServices registers views without service actions
func (ui *UI) registerViewsWithoutServices() {
	ui.registerContainersViewWithoutActions()
	ui.registerImagesViewWithoutActions()
	ui.registerVolumesViewWithoutActions()
	ui.registerNetworksViewWithoutActions()
	ui.registerSwarmServicesViewWithoutActions()
	ui.registerSwarmNodesViewWithoutActions()
}

// registerContainersViewWithoutActions registers the containers view without actions
func (ui *UI) registerContainersViewWithoutActions() {
	ui.viewRegistry.Register(
		"containers",
		"Containers",
		'c',
		ui.containersView.GetView(),
		ui.containersView.Refresh,
		"",
		"",
	)
}

// registerImagesViewWithoutActions registers the images view without actions
func (ui *UI) registerImagesViewWithoutActions() {
	ui.viewRegistry.Register(
		"images",
		"Images",
		'i',
		ui.imagesView.GetView(),
		ui.imagesView.Refresh,
		"",
		"",
	)
}

// registerVolumesViewWithoutActions registers the volumes view without actions
func (ui *UI) registerVolumesViewWithoutActions() {
	ui.viewRegistry.Register(
		"volumes",
		"Volumes",
		'v',
		ui.volumesView.GetView(),
		ui.volumesView.Refresh,
		"",
		"",
	)
}

// registerNetworksViewWithoutActions registers the networks view without actions
func (ui *UI) registerNetworksViewWithoutActions() {
	ui.viewRegistry.Register(
		"networks",
		"Networks",
		'n',
		ui.networksView.GetView(),
		ui.networksView.Refresh,
		"",
		"",
	)
}

// registerSwarmServicesViewWithoutActions registers the swarm services view without actions
func (ui *UI) registerSwarmServicesViewWithoutActions() {
	ui.viewRegistry.Register(
		"swarmServices",
		"Swarm Services",
		's',
		ui.swarmServicesView.GetView(),
		ui.swarmServicesView.Refresh,
		"",
		"",
	)
}

// registerSwarmNodesViewWithoutActions registers the swarm nodes view without actions
func (ui *UI) registerSwarmNodesViewWithoutActions() {
	ui.viewRegistry.Register(
		"swarmNodes",
		"Swarm Nodes",
		'w',
		ui.swarmNodesView.GetView(),
		ui.swarmNodesView.Refresh,
		"",
		"",
	)
}

// setDefaultView sets the default view for the application
func (ui *UI) setDefaultView() {
	ui.viewRegistry.SetCurrent(constants.DefaultView)

	// Set the default service to container for initial navigation
	if ui.services != nil {
		ui.services.SetCurrentService("container")
	}
}

// createViewContainer creates the main view container
func (ui *UI) createViewContainer() {
	ui.viewContainer = ui.viewBuilder.CreateView()
	ui.viewContainer.SetBorder(true)
	ui.viewContainer.SetTitleColor(ui.themeManager.GetHeaderColor())
	ui.viewContainer.SetBorderColor(ui.themeManager.GetBorderColor())

	// Set initial view
	if currentView := ui.viewRegistry.GetCurrent(); currentView != nil {
		ui.viewContainer.SetTitle(fmt.Sprintf(" %s ", currentView.Title))
		ui.viewContainer.AddItem(currentView.View, 0, 1, true)
	}
}

// createStatusBar creates the status bar
func (ui *UI) createStatusBar() {
	ui.statusBar = ui.componentBuilder.CreateTextView(
		"",
		tview.AlignLeft,
		ui.themeManager.GetTextColor(),
	)
	ui.statusBar.SetBackgroundColor(ui.themeManager.GetBackgroundColor())

	// Ensure status bar has consistent height and doesn't expand
	ui.statusBar.SetDynamicColors(false)
	ui.statusBar.SetScrollable(false)
	ui.statusBar.SetWrap(false)

	ui.updateStatusBar()
}

// setupKeyBindings sets up global key bindings
func (ui *UI) setupKeyBindings() {
	ui.app.SetInputCapture(ui.handleGlobalKeyBindings)
}

// handleGlobalKeyBindings handles all global key bindings
func (ui *UI) handleGlobalKeyBindings(event *tcell.EventKey) *tcell.EventKey {
	if ui.shouldSkipGlobalKeyBindings(event) {
		return event
	}

	return ui.routeKeyBinding(event)
}

// shouldSkipGlobalKeyBindings checks if global key bindings should be skipped
func (ui *UI) shouldSkipGlobalKeyBindings(event *tcell.EventKey) bool {
	return ui.IsModalActive() || ui.isShellInputFieldFocused()
}

// routeKeyBinding routes the key binding to the appropriate handler
func (ui *UI) routeKeyBinding(event *tcell.EventKey) *tcell.EventKey {
	if ui.app.GetFocus() == ui.commandHandler.GetInput() {
		return ui.handleCommandModeKeyBindings(event)
	}

	if ui.isExecCommandInputActive() {
		return ui.handleExecCommandKeyBindings(event)
	}

	if ui.isShellViewActive() {
		return ui.handleShellViewKeyBindings(event)
	}

	return ui.handleNormalModeKeyBindings(event)
}

// handleCommandModeKeyBindings handles key bindings when in command mode
func (ui *UI) handleCommandModeKeyBindings(event *tcell.EventKey) *tcell.EventKey {
	// In command mode, only allow ESC to exit
	if event.Key() == tcell.KeyEscape {
		ui.commandHandler.HandleInput(tcell.KeyEscape)
		return nil
	}
	return event
}

// handleExecCommandKeyBindings handles key bindings when exec command input is active
func (ui *UI) handleExecCommandKeyBindings(event *tcell.EventKey) *tcell.EventKey {
	// In exec command mode, only allow ESC to exit
	if event.Key() != tcell.KeyEscape {
		return event
	}

	// Remove the exec input and return focus to view
	ui.removeExecInputAndRestoreFocus()
	return nil
}

// removeExecInputAndRestoreFocus removes the exec input and restores focus to the view container
func (ui *UI) removeExecInputAndRestoreFocus() {
	if ui.mainFlex == nil {
		return
	}

	focused := ui.app.GetFocus()
	if focused == nil {
		return
	}

	ui.removeInputFieldAndRestoreFocus(focused)
}

// removeInputFieldAndRestoreFocus removes the input field and restores focus
func (ui *UI) removeInputFieldAndRestoreFocus(focused tview.Primitive) {
	inputField, ok := focused.(*tview.InputField)
	if !ok {
		return
	}

	ui.mainFlex.RemoveItem(inputField)
	if ui.viewContainer != nil {
		ui.app.SetFocus(ui.viewContainer)
	}
}

// handleShellViewKeyBindings handles key bindings when shell view is active
func (ui *UI) handleShellViewKeyBindings(event *tcell.EventKey) *tcell.EventKey {
	// In shell mode, only allow ESC to exit (handled by shell view)
	if event.Key() == tcell.KeyEscape {
		return event // Let shell view handle ESC
	}
	// Block other global key bindings in shell mode
	return event
}

// handleNormalModeKeyBindings handles key bindings in normal mode
func (ui *UI) handleNormalModeKeyBindings(event *tcell.EventKey) *tcell.EventKey {
	// Log what component is focused
	if focused := ui.app.GetFocus(); focused != nil {
		ui.log.Info(
			"Normal mode key binding",
			"key",
			event.Key(),
			"focusedType",
			fmt.Sprintf("%T", focused),
		)
	}

	switch event.Key() {
	case tcell.KeyRune:
		// Only handle global rune keys, let others pass through to focused component
		return ui.handleGlobalRuneKeyBindings(event)
	case tcell.KeyCtrlC:
		return ui.handleCtrlCKeyBinding(event)
	}
	return event
}

// handleGlobalRuneKeyBindings handles only global rune key bindings
func (ui *UI) handleGlobalRuneKeyBindings(event *tcell.EventKey) *tcell.EventKey {
	ui.log.Info("Global rune key handler called", "key", string(event.Rune()))

	if ui.handleQuitKey(event) || ui.handleCommandModeKey(event) {
		return nil
	}

	ui.log.Info("Global handler passing through key", "key", string(event.Rune()))
	// Let all other keys pass through to the focused view (including action keys)
	return event
}

// handleQuitKey handles quit key presses
func (ui *UI) handleQuitKey(event *tcell.EventKey) bool {
	if event.Rune() == 'q' || event.Rune() == 'Q' {
		ui.log.Info("Quitting application...")
		// Send shutdown signal instead of direct exit to ensure cleanup
		select {
		case ui.shutdownChan <- struct{}{}:
		default:
		}
		return true
	}
	return false
}

// handleCommandModeKey handles command mode key presses
func (ui *UI) handleCommandModeKey(event *tcell.EventKey) bool {
	if event.Rune() == ':' {
		ui.log.Info("Entering command mode")
		ui.commandHandler.Enter()
		return true
	}
	return false
}

// handleCtrlCKeyBinding handles Ctrl+C key binding
func (ui *UI) handleCtrlCKeyBinding(_ *tcell.EventKey) *tcell.EventKey {
	ui.log.Info("Received Ctrl+C, shutting down...")
	select {
	case ui.shutdownChan <- struct{}{}:
	default:
	}
	return nil
}

// isExecCommandInputActive checks if exec command input is currently active
func (ui *UI) isExecCommandInputActive() bool {
	if focused := ui.app.GetFocus(); focused != nil {
		if inputField, ok := focused.(*tview.InputField); ok {
			return inputField.GetLabel() == " Exec Command: "
		}
	}
	return false
}

// isShellViewActive checks if shell view is currently active
func (ui *UI) isShellViewActive() bool {
	return ui.shellView != nil && ui.app.GetFocus() == ui.shellView.GetView()
}

// isShellInputFieldFocused checks if shell input field is currently focused
func (ui *UI) isShellInputFieldFocused() bool {
	if ui.shellView != nil {
		if focused := ui.app.GetFocus(); focused != nil {
			if inputField, ok := focused.(*tview.InputField); ok {
				// Check if this input field belongs to the shell view
				return inputField == ui.shellView.GetInputField()
			}
		}
	}
	return false
}

// cleanup performs terminal cleanup operations
func (ui *UI) cleanup() {
	// Skip terminal cleanup operations when in TUI mode to prevent interference
	if logger.IsTUIMode() {
		return
	}

	ui.clearScreen()
	ui.resetColors()
	ui.showCursor()
	ui.moveCursorToTop()
	ui.syncStdout()
}

// clearScreen clears the terminal screen
func (ui *UI) clearScreen() {
	if _, e := fmt.Fprint(os.Stdout, "\033[2J"); e != nil {
		ui.log.Warn("Failed to clear screen", "error", e)
	}
}

// resetColors resets terminal colors
func (ui *UI) resetColors() {
	if _, e := fmt.Fprint(os.Stdout, "\033[0m"); e != nil {
		ui.log.Warn("Failed to reset colors", "error", e)
	}
}

// showCursor shows the terminal cursor
func (ui *UI) showCursor() {
	if _, e := fmt.Fprint(os.Stdout, "\033[?25h"); e != nil {
		ui.log.Warn("Failed to show cursor", "error", e)
	}
}

// moveCursorToTop moves the cursor to the top of the terminal
func (ui *UI) moveCursorToTop() {
	if _, err := fmt.Fprint(os.Stdout, "\033[H"); err != nil {
		ui.log.Warn("Failed to move cursor", "error", err)
	}
}

// syncStdout synchronizes stdout
func (ui *UI) syncStdout() {
	if e := os.Stdout.Sync(); e != nil {
		ui.log.Debug("Failed to sync stdout", "error", e)
	}
}

// ensureStableLayout ensures the UI layout remains stable during refreshes
func (ui *UI) ensureStableLayout() {
	if ui.mainFlex != nil {
		// Ensure the main layout doesn't shift during refreshes
		ui.mainFlex.SetDirection(tview.FlexRow)
	}

	if ui.statusBar != nil {
		// Ensure status bar maintains its height
		ui.statusBar.SetDynamicColors(false)
		ui.statusBar.SetScrollable(false)
		ui.statusBar.SetWrap(false)
	}
}

// refreshCurrentView refreshes the currently active view
func (ui *UI) refreshCurrentView() {
	if currentView := ui.viewRegistry.GetCurrent(); currentView != nil &&
		currentView.Refresh != nil {
		ui.log.Debug("Refreshing current view", "view", currentView.Title)
		currentView.Refresh()
	} else {
		ui.log.Debug("No current view or refresh function available")
	}
}

// switchView switches to the specified view
func (ui *UI) switchView(view string) {
	ui.log.Debug("Switching to view", "view", view)

	if !ui.validateViewExists(view) {
		return
	}

	ui.performViewSwitch(view)
	ui.updateViewDisplay()
	ui.refreshViewAndFocus(view)

	ui.log.Debug("Switched to view", "view", view)
}

// validateViewExists checks if the specified view exists
func (ui *UI) validateViewExists(view string) bool {
	if !ui.viewRegistry.Exists(view) {
		ui.log.Warn("Unknown view", "view", view)
		return false
	}
	return true
}

// performViewSwitch performs the actual view switching logic
func (ui *UI) performViewSwitch(view string) {
	ui.viewRegistry.SetCurrent(view)

	// Set the current service based on the view to enable proper navigation
	if ui.services != nil {
		switch view {
		case "containers":
			ui.services.SetCurrentService("container")
		case "images":
			ui.services.SetCurrentService("image")
		case "volumes":
			ui.services.SetCurrentService("volume")
		case "networks":
			ui.services.SetCurrentService("network")
		}
	}
}

// updateViewDisplay updates the view container display
func (ui *UI) updateViewDisplay() {
	viewInfo := ui.viewRegistry.GetCurrent()
	ui.viewContainer.SetTitle(fmt.Sprintf(" %s ", viewInfo.Title))
	ui.viewContainer.Clear()
	ui.viewContainer.AddItem(viewInfo.View, 0, 1, true)
}

// refreshViewAndFocus refreshes the view and sets focus
func (ui *UI) refreshViewAndFocus(_ string) {
	viewInfo := ui.viewRegistry.GetCurrent()

	ui.updateStatusBar()
	ui.updateHeadersIfNeeded()
	ui.refreshViewInfo(viewInfo)
	ui.setFocusToTable(viewInfo)
}

// updateHeadersIfNeeded updates headers only if not in a refresh cycle
func (ui *UI) updateHeadersIfNeeded() {
	if !ui.isRefreshing {
		ui.headerManager.UpdateDockerInfo()
		ui.headerManager.UpdateNavigation()
		ui.headerManager.UpdateActions()
	}
}

// refreshViewInfo refreshes the current view if it has a refresh function
func (ui *UI) refreshViewInfo(viewInfo *ViewInfo) {
	if viewInfo.Refresh != nil {
		viewInfo.Refresh()
	}
}

// setFocusToTable sets focus to the table within the view
func (ui *UI) setFocusToTable(viewInfo *ViewInfo) {
	if view, ok := viewInfo.View.(*tview.Flex); ok {
		if table := ui.findTableInFlex(view); table != nil {
			ui.app.SetFocus(table)
			return
		}
	}

	// Fallback to setting focus on the view if no table is found
	ui.app.SetFocus(viewInfo.View)
}

// findTableInFlex finds a table within a Flex container
func (ui *UI) findTableInFlex(view *tview.Flex) *tview.Table {
	for i := 0; i < view.GetItemCount(); i++ {
		if item := view.GetItem(i); item != nil {
			if table, isTable := item.(*tview.Table); isTable {
				return table
			}
		}
	}
	return nil
}

// updateStatusBar updates the status bar with current information
func (ui *UI) updateStatusBar() {
	if ui.statusBar == nil {
		return
	}

	statusText := ui.buildStatusBarText()

	// Only update if the text has actually changed to prevent unnecessary redraws
	currentText := ui.statusBar.GetText(true)
	if currentText != statusText {
		ui.statusBar.SetText(statusText)
		ui.log.Debug("Status bar updated", "old", currentText, "new", statusText)
	}
}

// buildStatusBarText builds the status bar text with current information
func (ui *UI) buildStatusBarText() string {
	now := time.Now()
	timeStr := now.Format("15:04:05")

	// Ensure no newlines in status bar text to prevent terminal display issues
	statusText := fmt.Sprintf(constants.StatusBarTemplate, timeStr)
	return strings.TrimSpace(statusText)
}

// updateLegend updates the legend with view-specific shortcuts
func (ui *UI) updateLegend() {
	ui.log.Debug("updateLegend called", "isRefreshing", ui.isRefreshing)
	ui.headerManager.UpdateDockerInfo()
	ui.headerManager.UpdateNavigation()
	ui.headerManager.UpdateActions()
}

// showHelp shows the help modal
func (ui *UI) showHelp() {
	if ui.modalManager != nil {
		ui.modalManager.ShowHelp()
	}
}

// showError shows an error modal
func (ui *UI) showError(err error) {
	if ui.modalManager != nil {
		ui.modalManager.ShowError(err)
	}
}

// showInfo shows an info modal
func (ui *UI) showInfo(message string) {
	if ui.modalManager != nil {
		ui.modalManager.ShowInfo(message)
	}
}

// showContextualHelp shows context-sensitive help modal
func (ui *UI) showContextualHelp(context, operation string) {
	if ui.modalManager != nil {
		ui.modalManager.ShowContextualHelp(context, operation)
	}
}

// showRetryDialog shows retry dialog with automatic retry logic
func (ui *UI) showRetryDialog(
	operation string,
	err error,
	retryFunc func() error,
	onSuccess func(),
) {
	if ui.modalManager != nil {
		ui.modalManager.ShowRetryDialog(operation, err, retryFunc, onSuccess)
	}
}

// showFallbackDialog shows fallback operations dialog
func (ui *UI) showFallbackDialog(
	operation string,
	err error,
	fallbackOptions []string,
	onFallback func(string),
) {
	if ui.modalManager != nil {
		ui.modalManager.ShowFallbackDialog(operation, err, fallbackOptions, onFallback)
	}
}

// showConfirm shows a confirmation modal
func (ui *UI) showConfirm(text string, callback func(bool)) {
	if ui.modalManager != nil {
		ui.modalManager.ShowConfirm(text, callback)
	}
}

// showDetails shows a details view in the main view container
func (ui *UI) showDetails(detailsView tview.Primitive) {
	ui.viewContainer.Clear()
	ui.viewContainer.AddItem(detailsView, 0, 1, true)

	ui.app.SetFocus(detailsView)

	ui.inDetailsMode = true
	ui.updateLegend()
}

// showCurrentView returns to the current view's table
func (ui *UI) showCurrentView() {
	currentViewInfo := ui.viewRegistry.GetCurrent()
	if currentViewInfo == nil {
		return
	}

	ui.log.Debug("Returning to current view", "view", currentViewInfo.Name)

	ui.clearSpecialModes()
	ui.restoreCurrentView(currentViewInfo)
	ui.updateUIAfterViewRestore(currentViewInfo)
}

// clearSpecialModes clears special UI modes
func (ui *UI) clearSpecialModes() {
	ui.inDetailsMode = false
	ui.inLogsMode = false
	ui.currentActions = nil
}

// restoreCurrentView restores the current view in the container
func (ui *UI) restoreCurrentView(currentViewInfo *ViewInfo) {
	ui.viewContainer.Clear()
	ui.viewContainer.AddItem(currentViewInfo.View, 0, 1, true)
	ui.viewContainer.SetTitle(fmt.Sprintf(" %s ", currentViewInfo.Title))

	// Ensure the current service is set for proper navigation
	if ui.services != nil {
		switch currentViewInfo.Name {
		case "containers":
			ui.services.SetCurrentService("container")
		case "images":
			ui.services.SetCurrentService("image")
		case "volumes":
			ui.services.SetCurrentService("volume")
		case "networks":
			ui.services.SetCurrentService("network")
		}
	}
}

// updateUIAfterViewRestore updates the UI after restoring the view
func (ui *UI) updateUIAfterViewRestore(currentViewInfo *ViewInfo) {
	// Only update headers if not in a refresh cycle
	if !ui.isRefreshing {
		ui.headerManager.UpdateDockerInfo()
		ui.headerManager.UpdateNavigation()
		ui.headerManager.UpdateActions()
	}

	ui.app.SetFocus(currentViewInfo.View)

	if currentViewInfo.Refresh != nil {
		currentViewInfo.Refresh()
	}
}

// showLogs displays logs for any resource type in a dedicated view
func (ui *UI) showLogs(resourceType, resourceID, resourceName string) {
	ui.log.Debug(
		"Showing logs for resource",
		"type",
		resourceType,
		"id",
		resourceID,
		"name",
		resourceName,
	)

	ui.setLogsMode()
	ui.setupLogsActions()
	ui.updateLogsViewTitle(resourceType)
	ui.displayLogsView(resourceType, resourceID, resourceName)
	ui.updateLegend()
	ui.setLogsFocus()
}

// setLogsMode sets the UI to logs mode
func (ui *UI) setLogsMode() {
	ui.inLogsMode = true
	ui.inDetailsMode = false
}

// setupLogsActions sets up the available actions for logs view
func (ui *UI) setupLogsActions() {
	ui.currentActions = ui.services.GetLogsService().GetActions()
}

// updateLogsViewTitle updates the view container title for logs
func (ui *UI) updateLogsViewTitle(resourceType string) {
	ui.viewContainer.SetTitle(fmt.Sprintf(" %s<Logs> ", resourceType))
}

// displayLogsView displays the logs view in the container
func (ui *UI) displayLogsView(resourceType, resourceID, resourceName string) {
	logsView := ui.GetLogsView(resourceType, resourceID, resourceName)
	logsView.LoadLogs()

	ui.viewContainer.Clear()
	ui.viewContainer.AddItem(logsView.GetView(), 0, 1, true)

	// Store the logs view for focus setting
	ui.logsView = logsView
}

// setLogsFocus sets focus to the logs view
func (ui *UI) setLogsFocus() {
	if ui.logsView != nil {
		ui.app.SetFocus(ui.logsView.GetView())
	}
}

// createShellView creates a new shell view for the container
func (ui *UI) createShellView(containerID, containerName string) {
	containerService := ui.services.GetContainerService()
	if containerService != nil {
		ui.shellView = shell.NewView(
			ui,
			containerID,
			containerName,
			ui.handleShellExit,
			containerService.ExecContainer,
		)
	}
}

// displayShellView displays the shell view in the container
func (ui *UI) displayShellView(containerID, containerName string) {
	ui.viewContainer.Clear()
	ui.viewContainer.SetTitle(
		fmt.Sprintf(" Shell - %s (%s) ", containerName, shared.TruncName(containerID, 12)),
	)
	ui.viewContainer.AddItem(ui.shellView.GetView(), 0, 1, true)
	ui.app.SetFocus(ui.shellView.GetView())
}

// handleShellExit handles the shell exit callback
func (ui *UI) handleShellExit() {
	ui.switchView("containers")
}

// hasValidPages checks if the pages container is valid
func (ui *UI) hasValidPages() bool {
	return ui.pages != nil
}

// hasModalPages checks if any modal pages are currently shown
func (ui *UI) hasModalPages() bool {
	return ui.pages.HasPage("help_modal") ||
		ui.pages.HasPage("error_modal") ||
		ui.pages.HasPage("confirm_modal") ||
		ui.pages.HasPage("exec_output_modal")
}
