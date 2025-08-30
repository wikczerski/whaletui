// Package logger provides structured logging functionality for the WhaleTUI application.
package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var (
	instance *slog.Logger
	logFile  *os.File
	logPath  string
	mu       sync.RWMutex
	tuiMode  bool // Flag to indicate if TUI mode is active
)

// GetLogger returns the singleton logger instance
func GetLogger() *slog.Logger {
	if instance := getExistingInstance(); instance != nil {
		return instance
	}

	return createNewInstance()
}

// getExistingInstance checks if an instance already exists
func getExistingInstance() *slog.Logger {
	mu.RLock()
	defer mu.RUnlock()
	return instance
}

// createNewInstance creates a new logger instance
func createNewInstance() *slog.Logger {
	mu.Lock()
	defer mu.Unlock()

	// Double-check after acquiring write lock
	if instance != nil {
		return instance
	}

	handler := createDefaultHandler()
	instance = slog.New(handler)
	slog.SetDefault(instance)
	return instance
}

// createDefaultHandler creates the default handler based on TUI mode
func createDefaultHandler() slog.Handler {
	if tuiMode {
		// In TUI mode, use file-only handler to prevent stderr interference
		return createFileOnlyHandler(slog.LevelInfo, "")
	}
	// Not in TUI mode, use stderr as usual
	return slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
}

// SetLevel sets the logging level and configures multistream if DEBUG
func SetLevel(level string) {
	SetLevelWithPath(level, "")
}

// SetLevelWithPath sets the logging level and configures multistream if DEBUG with a specific log file path
func SetLevelWithPath(level, logFilePath string) {
	slogLevel := parseLogLevel(level)

	mu.Lock()
	defer mu.Unlock()

	closeExistingLogFile()
	handler := createHandlerForLevel(level, slogLevel, logFilePath)
	updateLoggerInstance(handler)
}

// parseLogLevel converts string level to slog.Level
func parseLogLevel(level string) slog.Level {
	switch level {
	case "DEBUG":
		return slog.LevelDebug
	case "WARN":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// closeExistingLogFile closes the existing log file if any
func closeExistingLogFile() {
	if logFile != nil {
		if err := logFile.Close(); err != nil {
			// Log the error but continue since this is cleanup
			fmt.Fprintf(os.Stderr, "Failed to close log file: %v\n", err)
		}
		logFile = nil
	}
}

// createHandlerForLevel creates the appropriate handler for the given level
func createHandlerForLevel(level string, slogLevel slog.Level, logFilePath string) slog.Handler {
	if level == "DEBUG" {
		// For DEBUG level, always use multistream but with appropriate console output
		// In TUI mode, console output goes to a separate file to prevent interference
		return createMultistreamHandler(slogLevel, logFilePath)
	}

	// For other levels, check TUI mode
	if tuiMode {
		// In TUI mode, only log to file to prevent interference
		return createFileOnlyHandler(slogLevel, logFilePath)
	}

	// Not in TUI mode, use console only
	return slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slogLevel,
	})
}

// updateLoggerInstance updates the logger instance and default
func updateLoggerInstance(handler slog.Handler) {
	instance = slog.New(handler)
	slog.SetDefault(instance)
}

// createMultistreamHandler creates a handler that writes to both file and console
func createMultistreamHandler(level slog.Level, logFilePath string) slog.Handler {
	// Get the actual path to use (default if empty, otherwise use provided)
	actualPath := getDefaultLogPath(logFilePath)

	if !canCreateLogFile(actualPath) {
		return createDiscardHandler(level)
	}

	file, err := openLogFile(actualPath)
	if err != nil {
		return createDiscardHandler(level)
	}

	setupLogFile(file, actualPath)

	if tuiMode {
		return createFileOnlyHandler(level, actualPath)
	}

	return createConsoleAndFileHandler(level, file)
}

// canCreateLogFile checks if a log file can be created at the given path
func canCreateLogFile(logFilePath string) bool {
	if !isValidLogPath(logFilePath) {
		return false
	}

	return createLogsDirectory(logFilePath) == nil
}

// getDefaultLogPath returns the default log path if none is provided
func getDefaultLogPath(logFilePath string) string {
	if logFilePath == "" {
		return "./logs/whaletui.log"
	}
	return logFilePath
}

// setupLogFile sets up the global log file variables
func setupLogFile(file *os.File, logFilePath string) {
	logFile = file
	logPath = logFilePath
}

// createConsoleAndFileHandler creates a handler that writes to both console and file
func createConsoleAndFileHandler(level slog.Level, file *os.File) slog.Handler {
	multiWriter := io.MultiWriter(os.Stderr, file)
	return slog.NewTextHandler(multiWriter, &slog.HandlerOptions{
		Level: level,
	})
}

// isValidLogPath validates the log file path to prevent directory traversal
func isValidLogPath(logFilePath string) bool {
	// Allow relative paths that don't contain directory traversal
	cleanPath := filepath.Clean(logFilePath)

	// Only check for directory traversal, not for backslashes (which are normal on Windows)
	return !strings.Contains(cleanPath, "..")
}

// createDiscardHandler creates a handler that discards all logs
func createDiscardHandler(level slog.Level) slog.Handler {
	return slog.NewTextHandler(io.Discard, &slog.HandlerOptions{
		Level: level,
	})
}

// createLogsDirectory creates the logs directory if it doesn't exist
func createLogsDirectory(logFilePath string) error {
	logsDir := filepath.Dir(logFilePath)
	return os.MkdirAll(logsDir, 0o750)
}

// openLogFile opens the log file for writing
func openLogFile(logFilePath string) (*os.File, error) {
	const fileMode = 0o600
	const defaultLogPath = "./logs/whaletui.log"

	// Use a constant path to avoid gosec warning
	if logFilePath == defaultLogPath {
		return os.OpenFile(defaultLogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, fileMode)
	}

	// For non-default paths, validate the path to prevent directory traversal
	if !isValidLogPath(logFilePath) {
		return nil, fmt.Errorf("invalid log file path: %s", logFilePath)
	}

	// nolint:gosec // Path is validated by isValidLogPath
	return os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, fileMode)
}

// createFileOnlyHandler creates a handler that writes only to a file (for TUI mode)
func createFileOnlyHandler(level slog.Level, logFilePath string) slog.Handler {
	// Get the actual path to use (default if empty, otherwise use provided)
	actualPath := getDefaultLogPath(logFilePath)

	if !isValidLogPath(actualPath) {
		return createDiscardHandler(level)
	}

	if err := createLogsDirectory(actualPath); err != nil {
		return createDiscardHandler(level)
	}

	file, err := openLogFile(actualPath)
	if err != nil {
		return createDiscardHandler(level)
	}

	setupLogFile(file, actualPath)

	return slog.NewTextHandler(file, &slog.HandlerOptions{
		Level: level,
	})
}

// CloseLogFile closes the log file if it's open
func CloseLogFile() {
	if logFile != nil {
		if err := logFile.Close(); err != nil {
			// Log the error but continue since this is cleanup
			fmt.Fprintf(os.Stderr, "Failed to close log file: %v\n", err)
		}
		logFile = nil
	}
}

// IsDebugMode returns true if the logger is currently in DEBUG mode
func IsDebugMode() bool {
	mu.RLock()
	defer mu.RUnlock()
	return instance != nil && instance.Handler().Enabled(context.Background(), slog.LevelDebug)
}

// GetLogFilePath returns the current log file path if any
func GetLogFilePath() string {
	mu.RLock()
	defer mu.RUnlock()
	return logPath
}

// Convenience functions that use the singleton logger

// Debug logs a debug message using the singleton logger
func Debug(msg string, args ...any) {
	GetLogger().Debug(msg, args...)
}

// Info logs an info message using the singleton logger
func Info(msg string, args ...any) {
	GetLogger().Info(msg, args...)
}

// Warn logs a warning message using the singleton logger
func Warn(msg string, args ...any) {
	GetLogger().Warn(msg, args...)
}

// Error logs an error message using the singleton logger
func Error(msg string, args ...any) {
	GetLogger().Error(msg, args...)
}

// Fatal logs a fatal error message and exits the application
func Fatal(msg string, args ...any) {
	GetLogger().Error(msg, args...)
	os.Exit(1)
}

// SetTUIMode sets whether the logger should avoid stderr output (for TUI mode)
func SetTUIMode(enabled bool) {
	mu.Lock()
	defer mu.Unlock()
	tuiMode = enabled

	// If we're currently in DEBUG mode and have a logger instance,
	// recreate the handler to respect the new TUI mode setting
	// Check directly without calling IsDebugMode() to avoid deadlock
	if instance != nil && instance.Handler().Enabled(context.Background(), slog.LevelDebug) {
		var handler slog.Handler
		if tuiMode {
			// In TUI mode, only log to file
			handler = createFileOnlyHandler(slog.LevelDebug, logPath)
		} else {
			// Not in TUI mode, use multistream
			handler = createMultistreamHandler(slog.LevelDebug, logPath)
		}
		instance = slog.New(handler)
		slog.SetDefault(instance)
	}
}

// IsTUIMode returns whether TUI mode is currently active
func IsTUIMode() bool {
	mu.RLock()
	defer mu.RUnlock()
	return tuiMode
}
