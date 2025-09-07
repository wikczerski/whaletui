package ui

import (
	"fmt"

	"github.com/rivo/tview"
	"github.com/wikczerski/whaletui/internal/config"
	"github.com/wikczerski/whaletui/internal/ui/builders"
	"github.com/wikczerski/whaletui/internal/ui/constants"
	"github.com/wikczerski/whaletui/internal/ui/core"
	"github.com/wikczerski/whaletui/internal/ui/handlers"
)

// initializeManagers initializes all the UI managers and builders
func (ui *UI) initializeManagers(themePath string) error {
	ui.themeManager = config.NewThemeManager(themePath)

	ui.componentBuilder = builders.NewComponentBuilderWithTheme(ui.themeManager)
	ui.viewBuilder = builders.NewViewBuilder()
	ui.tableBuilder = builders.NewTableBuilder()

	ui.viewRegistry = core.NewViewRegistry()
	// Managers are now passed as parameters to avoid circular imports
	ui.commandHandler = handlers.NewCommandHandler(ui)
	ui.searchHandler = handlers.NewSearchHandler(ui)

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
	ui.searchInput = ui.searchHandler.CreateSearchInput()
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
