package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/wikczerski/whaletui/internal/logger"
)

func (h *dockerErrorHandler) showLogOption() {
	logFilePath := logger.GetLogFilePath()
	if h.hasLogFile(logFilePath) {
		fmt.Printf("Recent logs are available at: %s\n", logFilePath)
		if h.interaction.askYesNo("Would you like to view recent logs?") {
			h.showRecentLogs(logFilePath)
		}
	}
}

func (h *dockerErrorHandler) hasLogFile(logFilePath string) bool {
	return logFilePath != ""
}

func (h *dockerErrorHandler) showRecentLogs(logFilePath string) {
	h.showRecentLogsHeader()
	h.validateAndReadLogFile(logFilePath)
}

// showRecentLogsHeader shows the recent logs header
func (h *dockerErrorHandler) showRecentLogsHeader() {
	fmt.Println()
	fmt.Println("Recent logs:")
	fmt.Println("============")
}

// validateAndReadLogFile validates and reads the log file
func (h *dockerErrorHandler) validateAndReadLogFile(logFilePath string) {
	if !h.isValidLogFilePath(logFilePath) {
		h.showInvalidLogPathMessage()
		return
	}

	h.readAndDisplayLogFile(logFilePath)
	fmt.Println()
}

// isValidLogFilePath checks if the log file path is valid
func (h *dockerErrorHandler) isValidLogFilePath(logFilePath string) bool {
	// Clean the path to remove any directory traversal attempts
	cleanPath := filepath.Clean(logFilePath)

	// Ensure it's an absolute path
	if !filepath.IsAbs(cleanPath) {
		return false
	}

	// Additional security: check for suspicious patterns
	if strings.Contains(cleanPath, "..") || strings.Contains(cleanPath, "~") {
		return false
	}

	// Ensure the cleaned path matches the original (after cleaning)
	return cleanPath == filepath.Clean(logFilePath)
}

// showInvalidLogPathMessage shows the invalid log path message
func (h *dockerErrorHandler) showInvalidLogPathMessage() {
	fmt.Println("Invalid log file path")
	fmt.Println()
}

// readAndDisplayLogFile reads and displays the log file
func (h *dockerErrorHandler) readAndDisplayLogFile(logFilePath string) {
	// nolint:gosec // Path is validated by isValidLogFilePath before this function is called
	logContent, readErr := os.ReadFile(logFilePath)
	if h.canReadLogFile(readErr) {
		h.displayLastLogLines(string(logContent))
	} else {
		fmt.Printf("Could not read log file: %v\n", readErr)
	}
}

func (h *dockerErrorHandler) canReadLogFile(err error) bool {
	return err == nil
}

func (h *dockerErrorHandler) displayLastLogLines(content string) {
	lines := strings.Split(content, "\n")
	start := h.calculateLogStartIndex(len(lines))

	for i := start; i < len(lines); i++ {
		if h.isValidLogLine(lines[i]) {
			fmt.Println(lines[i])
		}
	}
}

func (h *dockerErrorHandler) calculateLogStartIndex(totalLines int) int {
	const maxLogLines = 20
	start := totalLines - maxLogLines
	if h.isValidStartIndex(start) {
		return start
	}
	return 0
}

func (h *dockerErrorHandler) isValidStartIndex(start int) bool {
	return start >= 0
}

func (h *dockerErrorHandler) isValidLogLine(line string) bool {
	return line != ""
}
