package cmd

import (
	"fmt"
	"os"

	"github.com/wikczerski/whaletui/internal/app"
	"github.com/wikczerski/whaletui/internal/logger"
)

// cleanupAndShutdown performs cleanup and shutdown
func cleanupAndShutdown(application *app.App) {
	defer cleanupTerminal()
	defer logger.CloseLogFile()
	defer logger.SetTUIMode(false)

	application.Shutdown()
}

// cleanupTerminal performs additional terminal cleanup operations
func cleanupTerminal() {
	if logger.IsTUIMode() {
		return
	}

	cleanupOperations := []func(){
		cleanupTerminalClearScreen,
		cleanupTerminalResetColors,
		cleanupTerminalShowCursor,
		cleanupTerminalMoveCursorToTop,
		cleanupTerminalSyncStdout,
	}

	for _, operation := range cleanupOperations {
		operation()
	}
}

// cleanupTerminalClearScreen clears the terminal screen
func cleanupTerminalClearScreen() {
	if _, err := fmt.Fprint(os.Stdout, "\033[2J"); err != nil {
		logger.Warn("Failed to clear screen", "error", err)
	}
}

// cleanupTerminalResetColors resets terminal colors
func cleanupTerminalResetColors() {
	if _, err := fmt.Fprint(os.Stdout, "\033[0m"); err != nil {
		logger.Warn("Failed to reset colors", "error", err)
	}
}

// cleanupTerminalShowCursor shows the terminal cursor
func cleanupTerminalShowCursor() {
	if _, err := fmt.Fprint(os.Stdout, "\033[?25h"); err != nil {
		logger.Warn("Failed to show cursor", "error", err)
	}
}

// cleanupTerminalMoveCursorToTop moves the cursor to the top of the terminal
func cleanupTerminalMoveCursorToTop() {
	if _, err := fmt.Fprint(os.Stdout, "\033[H"); err != nil {
		logger.Warn("Failed to move cursor", "error", err)
	}
}

// cleanupTerminalSyncStdout synchronizes stdout
func cleanupTerminalSyncStdout() {
	if err := os.Stdout.Sync(); err != nil {
		logger.Debug("Failed to sync stdout", "error", err)
	}
}
