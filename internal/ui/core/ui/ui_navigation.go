package ui

import (
	"fmt"

	"github.com/rivo/tview"
	"github.com/wikczerski/whaletui/internal/shared"
	"github.com/wikczerski/whaletui/internal/ui/constants"
	"github.com/wikczerski/whaletui/internal/ui/core"
	"github.com/wikczerski/whaletui/internal/ui/views/shell"
)

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
		ui.services.SetCurrentService(view)
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
func (ui *UI) refreshViewInfo(viewInfo *core.ViewInfo) {
	if viewInfo.Refresh != nil {
		viewInfo.Refresh()
	}
}

// setFocusToTable sets focus to the table within the view
func (ui *UI) setFocusToTable(viewInfo *core.ViewInfo) {
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
func (ui *UI) restoreCurrentView(currentViewInfo *core.ViewInfo) {
	ui.viewContainer.Clear()
	ui.viewContainer.AddItem(currentViewInfo.View, 0, 1, true)
	ui.viewContainer.SetTitle(fmt.Sprintf(" %s ", currentViewInfo.Title))

	// Ensure the current service is set for proper navigation
	if ui.services != nil {
		ui.services.SetCurrentService(currentViewInfo.Name)
	}
}

// updateUIAfterViewRestore updates the UI after restoring the view
func (ui *UI) updateUIAfterViewRestore(currentViewInfo *core.ViewInfo) {
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
	containerService := ui.GetContainerService()
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
	ui.switchView(constants.ViewContainers)
}
