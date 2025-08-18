package core

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/wikczerski/D5r/internal/config"
	"github.com/wikczerski/D5r/internal/logger"
	"github.com/wikczerski/D5r/internal/services"
	"github.com/wikczerski/D5r/internal/ui/builders"
	"github.com/wikczerski/D5r/internal/ui/constants"
	"github.com/wikczerski/D5r/internal/ui/handlers"
	"github.com/wikczerski/D5r/internal/ui/managers"
	"github.com/wikczerski/D5r/internal/ui/views"
	"github.com/wikczerski/D5r/internal/ui/views/shell"
)

// UI represents the main UI
type UI struct {
	app            *tview.Application
	pages          *tview.Pages
	mainFlex       *tview.Flex
	statusBar      *tview.TextView
	viewContainer  *tview.Flex
	services       services.ServiceFactoryInterface
	log            *logger.Logger
	shutdownChan   chan struct{}
	inDetailsMode  bool            // Track if we're in details view mode
	inLogsMode     bool            // Track if we're viewing container logs
	currentActions map[rune]string // Track current available actions

	// Theme management
	themeManager *config.ThemeManager

	// Abstracted managers
	viewRegistry   *ViewRegistry
	headerManager  *managers.HeaderManager
	commandHandler *handlers.CommandHandler
	modalManager   *managers.ModalManager

	// Individual views
	containersView *views.ContainersView
	imagesView     *views.ImagesView
	volumesView    *views.VolumesView
	networksView   *views.NetworksView
	logsView       *views.LogsView
	shellView      *shell.View

	// Component builders
	componentBuilder *builders.ComponentBuilder
	viewBuilder      *builders.ViewBuilder
	tableBuilder     *builders.TableBuilder
}

// New creates a new UI
func New(serviceFactory *services.ServiceFactory, themePath string) (*UI, error) {
	ui := &UI{
		services:       serviceFactory,
		app:            tview.NewApplication(),
		pages:          tview.NewPages(),
		log:            logger.GetLogger(),
		shutdownChan:   make(chan struct{}, 1), // Buffer channel to prevent deadlock
		currentActions: make(map[rune]string),
	}

	ui.log.SetPrefix("UI")

	if err := ui.initializeManagers(themePath); err != nil {
		return nil, err
	}

	ui.initComponents()
	ui.setupKeyBindings()

	return ui, nil
}

// initializeManagers initializes all the UI managers and builders
func (ui *UI) initializeManagers(themePath string) error {
	ui.themeManager = config.NewThemeManager(themePath)

	ui.componentBuilder = builders.NewComponentBuilderWithTheme(ui.themeManager)
	ui.viewBuilder = builders.NewViewBuilder()
	ui.tableBuilder = builders.NewTableBuilder()

	ui.viewRegistry = NewViewRegistry()
	ui.headerManager = managers.NewHeaderManager(ui)
	ui.commandHandler = handlers.NewCommandHandler(ui)
	ui.modalManager = managers.NewModalManager(ui)

	return nil
}

// initComponents initializes UI components
func (ui *UI) initComponents() {
	ui.mainFlex = ui.viewBuilder.CreateView()

	headerSection := ui.headerManager.CreateHeaderSection()
	commandInput := ui.commandHandler.CreateCommandInput()

	ui.createAndRegisterViews()
	ui.createViewContainer()
	ui.createStatusBar()

	ui.mainFlex.AddItem(headerSection, constants.HeaderSectionHeight, 1, false)
	ui.mainFlex.AddItem(ui.viewContainer, 0, 1, true)
	ui.mainFlex.AddItem(ui.statusBar, constants.StatusBarHeight, 1, false)

	ui.setupMainPages(commandInput)
	ui.initializeUIState()

	ui.log.Info("UI components initialized")
}

// setupMainPages sets up the main pages in the UI
func (ui *UI) setupMainPages(commandInput *tview.InputField) {
	ui.pages.AddPage("main", ui.mainFlex, true, true)
	ui.pages.AddPage("command", commandInput, true, false)

	ui.app.SetRoot(ui.pages, true)
}

// initializeUIState initializes the initial UI state
func (ui *UI) initializeUIState() {
	ui.headerManager.UpdateAll()
	ui.refreshCurrentView()
}

// createAndRegisterViews creates all views and registers them with the view registry
func (ui *UI) createAndRegisterViews() {
	ui.createResourceViews()
	ui.registerViewsWithActions()
	ui.setDefaultView()
}

// createResourceViews creates all the resource views
func (ui *UI) createResourceViews() {
	ui.containersView = views.NewContainersView(ui)
	ui.imagesView = views.NewImagesView(ui)
	ui.volumesView = views.NewVolumesView(ui)
	ui.networksView = views.NewNetworksView(ui)
}

// registerViewsWithActions registers views with their metadata and actions
func (ui *UI) registerViewsWithActions() {
	ui.registerContainerView()
	ui.registerResourceViews()
}

// registerContainerView registers the containers view with its actions
func (ui *UI) registerContainerView() {
	containerActions := ui.services.GetContainerService().GetActionsString()
	ui.viewRegistry.Register("containers", "Containers", 'c', ui.containersView.GetView(), ui.containersView.Refresh, containerActions)
}

// registerResourceViews registers the resource views with their actions
func (ui *UI) registerResourceViews() {
	ui.viewRegistry.Register("images", "Images", 'i', ui.imagesView.GetView(), ui.imagesView.Refresh, ui.services.GetImageService().GetActionsString())
	ui.viewRegistry.Register("volumes", "Volumes", 'v', ui.volumesView.GetView(), ui.volumesView.Refresh, ui.services.GetVolumeService().GetActionsString())
	ui.viewRegistry.Register("networks", "Networks", 'n', ui.networksView.GetView(), ui.networksView.Refresh, ui.services.GetNetworkService().GetActionsString())
}

// setDefaultView sets the default view for the application
func (ui *UI) setDefaultView() {
	ui.viewRegistry.SetCurrent(constants.DefaultView)
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
	ui.statusBar = ui.componentBuilder.CreateTextView("", tview.AlignLeft, ui.themeManager.GetTextColor())
	ui.statusBar.SetBackgroundColor(ui.themeManager.GetBackgroundColor())
	ui.updateStatusBar()
}

// setupKeyBindings sets up global key bindings
func (ui *UI) setupKeyBindings() {
	ui.app.SetInputCapture(ui.handleGlobalKeyBindings)
}

// handleGlobalKeyBindings handles all global key bindings
func (ui *UI) handleGlobalKeyBindings(event *tcell.EventKey) *tcell.EventKey {
	// Check if a modal is currently active - don't interfere with modal key handling
	if ui.IsModalActive() {
		return event // Let the modal handle its own key events
	}

	// Check if command mode is active
	if ui.app.GetFocus() == ui.commandHandler.GetInput() {
		return ui.handleCommandModeKeyBindings(event)
	}

	// Check if exec command input is active
	if ui.isExecCommandInputActive() {
		return ui.handleExecCommandKeyBindings(event)
	}

	// Check if shell view is active
	if ui.isShellViewActive() {
		return ui.handleShellViewKeyBindings(event)
	}

	// Check if shell input field is focused
	if ui.isShellInputFieldFocused() {
		return event // Block global key bindings in shell input mode
	}

	// Normal mode key bindings
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
	if event.Key() == tcell.KeyEscape {
		// Remove the exec input and return focus to view
		if mainFlex := ui.mainFlex; mainFlex != nil {
			if focused := ui.app.GetFocus(); focused != nil {
				if inputField, ok := focused.(*tview.InputField); ok {
					mainFlex.RemoveItem(inputField)
					if viewContainer := ui.viewContainer; viewContainer != nil {
						ui.app.SetFocus(viewContainer)
					}
				}
			}
		}
		return nil
	}
	return event
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
	switch event.Key() {
	case tcell.KeyRune:
		return ui.handleRuneKeyBindings(event)
	case tcell.KeyCtrlC:
		return ui.handleCtrlCKeyBinding(event)
	}
	return event
}

// handleRuneKeyBindings handles rune key bindings in normal mode
func (ui *UI) handleRuneKeyBindings(event *tcell.EventKey) *tcell.EventKey {
	switch event.Rune() {
	case 'q', 'Q':
		ui.log.Info("Quitting application...")
		// Send shutdown signal instead of direct exit to ensure cleanup
		select {
		case ui.shutdownChan <- struct{}{}:
		default:
		}
		return nil
	case ':':
		ui.commandHandler.Enter()
		return nil
	}
	return event
}

// handleCtrlCKeyBinding handles Ctrl+C key binding
func (ui *UI) handleCtrlCKeyBinding(event *tcell.EventKey) *tcell.EventKey {
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

// cleanup performs terminal cleanup operations
func (ui *UI) cleanup() {
	ui.clearScreen()
	ui.resetColors()
	ui.showCursor()
	ui.moveCursorToTop()
	ui.syncStdout()
}

// clearScreen clears the terminal screen
func (ui *UI) clearScreen() {
	if _, err := fmt.Fprint(os.Stdout, "\033[2J"); err != nil {
		ui.log.Warn("Failed to clear screen: %v", err)
	}
}

// resetColors resets terminal colors
func (ui *UI) resetColors() {
	if _, err := fmt.Fprint(os.Stdout, "\033[0m"); err != nil {
		ui.log.Warn("Failed to reset colors: %v", err)
	}
}

// showCursor shows the terminal cursor
func (ui *UI) showCursor() {
	if _, err := fmt.Fprint(os.Stdout, "\033[?25h"); err != nil {
		ui.log.Warn("Failed to show cursor: %v", err)
	}
}

// moveCursorToTop moves the cursor to the top of the terminal
func (ui *UI) moveCursorToTop() {
	if _, err := fmt.Fprint(os.Stdout, "\033[H"); err != nil {
		ui.log.Warn("Failed to move cursor: %v", err)
	}
}

// syncStdout synchronizes stdout
func (ui *UI) syncStdout() {
	if err := os.Stdout.Sync(); err != nil {
		ui.log.Warn("Failed to sync stdout: %v", err)
	}
}

// GetShutdownChan returns the shutdown channel
func (ui *UI) GetShutdownChan() chan struct{} {
	return ui.shutdownChan
}

// GetServices returns the service factory
func (ui *UI) GetServices() services.ServiceFactoryInterface {
	return ui.services
}

// GetApp returns the tview application
func (ui *UI) GetApp() any {
	return ui.app
}

// ShowError displays an error message
func (ui *UI) ShowError(err error) {
	ui.showError(err)
}

// ShowSuccess displays a success message in the status bar
func (ui *UI) ShowSuccess(message string) {
	ui.UpdateStatusBar("✓ " + message)
}

// ShowDetails displays a details view
func (ui *UI) ShowDetails(details any) {
	if detailsView, ok := details.(tview.Primitive); ok {
		ui.showDetails(detailsView)
	} else {
		ui.log.Warn("ShowDetails called with non-Primitive type: %T", details)
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

// UpdateStatusBar updates the status bar with the given message
func (ui *UI) UpdateStatusBar(_ string) {
	ui.updateStatusBar()
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

// GetViewRegistry returns the view registry
func (ui *UI) GetViewRegistry() any {
	return ui.viewRegistry
}

// GetMainFlex returns the main flex container
func (ui *UI) GetMainFlex() any {
	return ui.mainFlex
}

// GetLog returns the logger
func (ui *UI) GetLog() any {
	return ui.log
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
	ui.showLogs(containerID, containerName)
}

// ShowShell shows shell view for a container
func (ui *UI) ShowShell(containerID, containerName string) {
	ui.createShellView(containerID, containerName)
	ui.displayShellView(containerID, containerName)
}

// createShellView creates a new shell view for the container
func (ui *UI) createShellView(containerID, containerName string) {
	containerService := ui.services.GetContainerService()
	ui.shellView = shell.NewView(ui, containerID, containerName, ui.handleShellExit, containerService.ExecContainer)
}

// displayShellView displays the shell view in the container
func (ui *UI) displayShellView(containerID, containerName string) {
	ui.viewContainer.Clear()
	ui.viewContainer.SetTitle(fmt.Sprintf(" Shell - %s (%s) ", containerName, containerID[:12]))
	ui.viewContainer.AddItem(ui.shellView.GetView(), 0, 1, true)
	ui.app.SetFocus(ui.shellView.GetView())
}

// handleShellExit handles the shell exit callback
func (ui *UI) handleShellExit() {
	ui.switchView("containers")
}

// GetLogsView returns the logs view for a container
func (ui *UI) GetLogsView(containerID, containerName string) *views.LogsView {
	if ui.logsView == nil || ui.logsView.ContainerID != containerID {
		ui.logsView = views.NewLogsView(ui, containerID, containerName)
	}
	return ui.logsView
}

// GetViewContainer returns the view container
func (ui *UI) GetViewContainer() any {
	return ui.viewContainer
}

// GetContainerService returns the container service
func (ui *UI) GetContainerService() any {
	return ui.services.GetContainerService()
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

// UpdateLegend updates the legend with current view information
func (ui *UI) UpdateLegend() {
	ui.updateLegend()
}

// Refresh refreshes the UI
func (ui *UI) Refresh() {
	ui.log.Debug("Refreshing UI")
	ui.updateStatusBar()
	ui.headerManager.UpdateAll()
	ui.refreshCurrentView()
}

// refreshCurrentView refreshes the currently active view
func (ui *UI) refreshCurrentView() {
	if currentView := ui.viewRegistry.GetCurrent(); currentView != nil && currentView.Refresh != nil {
		currentView.Refresh()
	}
}

// switchView switches to the specified view
func (ui *UI) switchView(view string) {
	ui.log.Debug("Switching to view: %s", view)

	if !ui.validateViewExists(view) {
		return
	}

	ui.performViewSwitch(view)
	ui.updateViewDisplay()
	ui.refreshViewAndFocus(view)

	ui.log.Debug("Switched to view: %s", view)
}

// validateViewExists checks if the specified view exists
func (ui *UI) validateViewExists(view string) bool {
	if !ui.viewRegistry.Exists(view) {
		ui.log.Warn("Unknown view: %s", view)
		return false
	}
	return true
}

// performViewSwitch performs the actual view switching logic
func (ui *UI) performViewSwitch(view string) {
	ui.viewRegistry.SetCurrent(view)
}

// updateViewDisplay updates the view container display
func (ui *UI) updateViewDisplay() {
	viewInfo := ui.viewRegistry.GetCurrent()
	ui.viewContainer.SetTitle(fmt.Sprintf(" %s ", viewInfo.Title))
	ui.viewContainer.Clear()
	ui.viewContainer.AddItem(viewInfo.View, 0, 1, true)
}

// refreshViewAndFocus refreshes the view and sets focus
func (ui *UI) refreshViewAndFocus(view string) {
	viewInfo := ui.viewRegistry.GetCurrent()

	ui.updateStatusBar()
	ui.headerManager.UpdateAll()

	if viewInfo.Refresh != nil {
		viewInfo.Refresh()
	}
	ui.app.SetFocus(viewInfo.View)
}

// updateStatusBar updates the status bar with current information
func (ui *UI) updateStatusBar() {
	if ui.statusBar == nil {
		return
	}

	statusText := ui.buildStatusBarText()
	ui.statusBar.SetText(statusText)
}

// buildStatusBarText builds the status bar text with current information
func (ui *UI) buildStatusBarText() string {
	now := time.Now()
	timeStr := now.Format("15:04:05")

	return fmt.Sprintf(constants.StatusBarTemplate, timeStr)
}

// updateLegend updates the legend with view-specific shortcuts
func (ui *UI) updateLegend() {
	ui.headerManager.UpdateAll()
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

	ui.log.Debug("Returning to current view: %s", currentViewInfo.Name)

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
}

// updateUIAfterViewRestore updates the UI after restoring the view
func (ui *UI) updateUIAfterViewRestore(currentViewInfo *ViewInfo) {
	ui.headerManager.UpdateAll()
	ui.app.SetFocus(currentViewInfo.View)

	if currentViewInfo.Refresh != nil {
		currentViewInfo.Refresh()
	}
}

// showLogs displays container logs in a dedicated view
func (ui *UI) showLogs(containerID, containerName string) {
	ui.log.Debug("Showing logs for container: %s (%s)", containerID, containerName)

	ui.setLogsMode()
	ui.setupLogsActions()
	ui.updateLogsViewTitle()
	ui.displayLogsView(containerID, containerName)
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

// parseActionsString parses an action string into a map of rune to action description
func (ui *UI) parseActionsString(actionsString string) map[rune]string {
	actions := make(map[rune]string)
	lines := strings.Split(actionsString, "\n")

	for _, line := range lines {
		if strings.Contains(line, "<") && strings.Contains(line, ">") {
			// Extract the key and description
			start := strings.Index(line, "<") + 1
			end := strings.Index(line, ">")
			if start < end && end < len(line) {
				key := line[start:end]
				description := strings.TrimSpace(line[end+1:])
				if len(key) == 1 {
					actions[rune(key[0])] = description
				}
			}
		}
	}

	return actions
}

// updateLogsViewTitle updates the view container title for logs
func (ui *UI) updateLogsViewTitle() {
	ui.viewContainer.SetTitle(" Containers<Logs> ")
}

// displayLogsView displays the logs view in the container
func (ui *UI) displayLogsView(containerID, containerName string) {
	logsView := ui.GetLogsView(containerID, containerName)
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
