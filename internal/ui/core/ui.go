package core

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/user/d5r/internal/logger"
	"github.com/user/d5r/internal/services"
	"github.com/user/d5r/internal/ui/builders"
	"github.com/user/d5r/internal/ui/constants"
	"github.com/user/d5r/internal/ui/managers"
	"github.com/user/d5r/internal/ui/views"
)

// UI represents the main UI
type UI struct {
	app            *tview.Application
	pages          *tview.Pages
	mainFlex       *tview.Flex
	statusBar      *tview.TextView
	viewContainer  *tview.Flex
	services       *services.ServiceFactory
	log            *logger.Logger
	shutdownChan   chan struct{}
	inDetailsMode  bool            // Track if we're in details view mode
	inLogsMode     bool            // Track if we're viewing container logs
	currentActions map[rune]string // Track current available actions

	// Abstracted managers
	viewRegistry   *ViewRegistry
	headerManager  *managers.HeaderManager
	commandHandler *managers.CommandHandler
	modalManager   *managers.ModalManager

	// Individual views
	containersView *views.ContainersView
	imagesView     *views.ImagesView
	volumesView    *views.VolumesView
	networksView   *views.NetworksView

	// Component builders
	componentBuilder *builders.ComponentBuilder
	viewBuilder      *builders.ViewBuilder
	tableBuilder     *builders.TableBuilder
}

// New creates a new UI
func New(serviceFactory *services.ServiceFactory) (*UI, error) {
	app := tview.NewApplication()

	ui := &UI{
		services:       serviceFactory,
		app:            app,
		pages:          tview.NewPages(),
		log:            logger.GetLogger(),
		shutdownChan:   make(chan struct{}, 1), // Buffer channel to prevent deadlock
		currentActions: make(map[rune]string),
	}

	ui.log.SetPrefix("UI")

	ui.componentBuilder = builders.NewComponentBuilder()
	ui.viewBuilder = builders.NewViewBuilder()
	ui.tableBuilder = builders.NewTableBuilder()

	ui.viewRegistry = NewViewRegistry()

	ui.headerManager = managers.NewHeaderManager(ui)

	ui.commandHandler = managers.NewCommandHandler(ui)
	ui.modalManager = managers.NewModalManager(ui)

	ui.initComponents()

	ui.setupKeyBindings()

	return ui, nil
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

	ui.pages.AddPage("main", ui.mainFlex, true, true)
	ui.pages.AddPage("command", commandInput, true, false)

	ui.app.SetRoot(ui.pages, true)

	ui.headerManager.UpdateAll()

	ui.refreshCurrentView()

	ui.log.Info("UI components initialized")
}

// createAndRegisterViews creates all views and registers them with the view registry
func (ui *UI) createAndRegisterViews() {
	ui.containersView = views.NewContainersView(ui)
	ui.imagesView = views.NewImagesView(ui)
	ui.volumesView = views.NewVolumesView(ui)
	ui.networksView = views.NewNetworksView(ui)

	// Register views with their metadata
	containerActions := "<s> Start\n<S> Stop\n<r> Restart\n<d> Delete\n<a> Attach\n<l> Logs\n<i> Inspect\n<n> New\n<e> Exec\n<f> Filter\n<t> Sort\n<h> History\n<enter> Details\n<:> Command"
	resourceActions := "<r> Remove\n<h> History\n<f> Filter\n<t> Sort\n<i> Inspect\n<enter> Details\n<up/down> Navigate\n<?> Help\n<f5> Refresh\n<:> Command"

	ui.viewRegistry.Register("containers", "Containers", 'c', ui.containersView.GetView(), ui.containersView.Refresh, containerActions)
	ui.viewRegistry.Register("images", "Images", 'i', ui.imagesView.GetView(), ui.imagesView.Refresh, resourceActions)
	ui.viewRegistry.Register("volumes", "Volumes", 'v', ui.volumesView.GetView(), ui.volumesView.Refresh, resourceActions)
	ui.viewRegistry.Register("networks", "Networks", 'n', ui.networksView.GetView(), ui.networksView.Refresh, resourceActions)

	ui.viewRegistry.SetCurrent(constants.DefaultView)
}

// createViewContainer creates the main view container
func (ui *UI) createViewContainer() {
	ui.viewContainer = ui.viewBuilder.CreateView()
	ui.viewContainer.SetBorder(true)
	ui.viewContainer.SetTitleColor(constants.HeaderColor)
	ui.viewContainer.SetBorderColor(constants.BorderColor)

	// Set initial view
	if currentView := ui.viewRegistry.GetCurrent(); currentView != nil {
		ui.viewContainer.SetTitle(fmt.Sprintf(" %s ", currentView.Title))
		ui.viewContainer.AddItem(currentView.View, 0, 1, true)
	}
}

// createStatusBar creates the status bar
func (ui *UI) createStatusBar() {
	ui.statusBar = ui.componentBuilder.CreateTextView("", tview.AlignLeft, constants.TextColor)
	ui.statusBar.SetBackgroundColor(constants.BackgroundColor)
	ui.updateStatusBar()
}

// setupKeyBindings sets up global key bindings
func (ui *UI) setupKeyBindings() {
	ui.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// Check if command mode is active
		if ui.app.GetFocus() == ui.commandHandler.GetInput() {
			// In command mode, only allow ESC to exit
			if event.Key() == tcell.KeyEscape {
				ui.commandHandler.HandleInput(tcell.KeyEscape)
				return nil
			}
			return event
		}

		// Normal mode key bindings
		switch event.Key() {
		case tcell.KeyRune:
			switch event.Rune() {
			case 'q', 'Q':
				ui.log.Info("Quitting application...")
				os.Exit(0)
			case ':':
				ui.commandHandler.Enter()
				return nil
			}
		case tcell.KeyCtrlC:
			ui.log.Info("Received Ctrl+C, shutting down...")
			select {
			case ui.shutdownChan <- struct{}{}:
			default:
			}
			return nil
		}
		return event
	})
}

// Start starts the UI
func (ui *UI) Start() error {
	ui.log.Info("Starting TUI...")
	return ui.app.Run()
}

// Stop stops the UI
func (ui *UI) Stop() {
	ui.app.Stop()
}

// GetShutdownChan returns the shutdown channel
func (ui *UI) GetShutdownChan() <-chan struct{} {
	return ui.shutdownChan
}

// GetServices returns the service factory
func (ui *UI) GetServices() *services.ServiceFactory {
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
func (ui *UI) UpdateStatusBar(message string) {
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

// GetCurrentActions returns the current available actions
func (ui *UI) GetCurrentActions() map[rune]string {
	return ui.currentActions
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

// GetViewContainer returns the view container
func (ui *UI) GetViewContainer() any {
	return ui.viewContainer
}

// GetContainerService returns the container service
func (ui *UI) GetContainerService() any {
	return ui.services.ContainerService
}

// GetImageService returns the image service
func (ui *UI) GetImageService() any {
	return ui.services.ImageService
}

// GetVolumeService returns the volume service
func (ui *UI) GetVolumeService() any {
	return ui.services.VolumeService
}

// GetNetworkService returns the network service
func (ui *UI) GetNetworkService() any {
	return ui.services.NetworkService
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

	if !ui.viewRegistry.Exists(view) {
		ui.log.Warn("Unknown view: %s", view)
		return
	}

	ui.viewRegistry.SetCurrent(view)
	viewInfo := ui.viewRegistry.GetCurrent()

	ui.viewContainer.SetTitle(fmt.Sprintf(" %s ", viewInfo.Title))
	ui.viewContainer.Clear()
	ui.viewContainer.AddItem(viewInfo.View, 0, 1, true)

	ui.updateStatusBar()
	ui.headerManager.UpdateAll()

	if viewInfo.Refresh != nil {
		viewInfo.Refresh()
	}
	ui.app.SetFocus(viewInfo.View)

	ui.log.Debug("Switched to view: %s", view)
}

// updateStatusBar updates the status bar with current information
func (ui *UI) updateStatusBar() {
	if ui.statusBar == nil {
		return
	}

	now := time.Now()
	timeStr := now.Format("15:04:05")

	statusText := fmt.Sprintf(constants.StatusBarTemplate,
		timeStr,
		"Enter",
		"Q")

	ui.statusBar.SetText(statusText)
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

	// Clear special modes
	ui.inDetailsMode = false
	ui.inLogsMode = false
	ui.currentActions = nil

	// Restore the view
	ui.viewContainer.Clear()
	ui.viewContainer.AddItem(currentViewInfo.View, 0, 1, true)
	ui.viewContainer.SetTitle(fmt.Sprintf(" %s ", currentViewInfo.Title))

	ui.headerManager.UpdateAll()
	ui.app.SetFocus(currentViewInfo.View)

	if currentViewInfo.Refresh != nil {
		currentViewInfo.Refresh()
	}
}

// showLogs displays container logs in a dedicated view
func (ui *UI) showLogs(containerID, containerName string) {
	ui.log.Debug("Showing logs for container: %s (%s)", containerID, containerName)

	ui.inLogsMode = true
	ui.inDetailsMode = false
	ui.currentActions = map[rune]string{
		'f': "Follow logs",
		't': "Tail logs",
		's': "Save logs",
		'c': "Clear logs",
		'w': "Wrap text",
	}

	ui.viewContainer.SetTitle(" Containers<Logs> ")

	logsView := ui.createLogsView(containerID, containerName)

	ui.viewContainer.Clear()
	ui.viewContainer.AddItem(logsView, 0, 1, true)

	ui.updateLegend()

	ui.app.SetFocus(logsView)
}

// createLogsView creates a view for displaying container logs
func (ui *UI) createLogsView(containerID, containerName string) *tview.Flex {
	logsFlex := ui.viewBuilder.CreateView()

	bottomTitleView := ui.componentBuilder.CreateBorderedTextView(
		fmt.Sprintf(" %s<%s> ", containerName, containerID[:12]),
		"",
		constants.HeaderColor,
	)
	bottomTitleView.SetTextAlign(tview.AlignCenter)

	logsText := ui.componentBuilder.CreateTextView("Loading logs...", tview.AlignLeft, constants.TextColor)
	logsText.SetDynamicColors(true)
	logsText.SetScrollable(true)
	logsText.SetBorder(true)
	logsText.SetBorderColor(constants.BorderColor)

	go ui.loadContainerLogs(containerID, logsText)

	backButton := ui.componentBuilder.CreateButton("Back to Table", func() {
		ui.showCurrentView()
	})

	logsFlex.AddItem(bottomTitleView, constants.TitleViewHeight, 0, false)
	logsFlex.AddItem(logsText, 0, 1, true)
	logsFlex.AddItem(backButton, constants.BackButtonHeight, 0, false)

	// Set up key bindings
	logsFlex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape, tcell.KeyEnter:
			ui.showCurrentView()
			return nil
		}

		// Let the text view handle scrolling keys
		if event.Key() == tcell.KeyUp || event.Key() == tcell.KeyDown ||
			event.Key() == tcell.KeyPgUp || event.Key() == tcell.KeyPgDn ||
			event.Key() == tcell.KeyHome || event.Key() == tcell.KeyEnd {
			return event
		}

		return event
	})

	return logsFlex
}

// loadContainerLogs loads container logs from Docker
func (ui *UI) loadContainerLogs(containerID string, logsText *tview.TextView) {
	ctx := context.Background()
	logs, err := ui.services.ContainerService.GetContainerLogs(ctx, containerID)
	if err != nil {
		logsText.SetText(fmt.Sprintf("Error loading logs: %v", err))
		return
	}

	logsText.SetText(logs)
}
