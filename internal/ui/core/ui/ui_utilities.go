package ui

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/wikczerski/whaletui/internal/logger"
	"github.com/wikczerski/whaletui/internal/ui/constants"
)

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
