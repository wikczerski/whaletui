package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
)

var (
	instance *slog.Logger
	logFile  *os.File
	logPath  string
	mu       sync.RWMutex
)

// GetLogger returns the singleton logger instance
func GetLogger() *slog.Logger {
	mu.RLock()
	if instance != nil {
		defer mu.RUnlock()
		return instance
	}
	mu.RUnlock()

	mu.Lock()
	defer mu.Unlock()

	// Double-check after acquiring write lock
	if instance != nil {
		return instance
	}

	// Initialize with default console handler
	instance = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(instance)
	return instance
}

// SetLevel sets the logging level and configures multistream if DEBUG
func SetLevel(level string) {
	SetLevelWithPath(level, "")
}

// SetLevelWithPath sets the logging level and configures multistream if DEBUG with a specific log file path
func SetLevelWithPath(level, logFilePath string) {
	var slogLevel slog.Level
	switch level {
	case "DEBUG":
		slogLevel = slog.LevelDebug
	case "WARN":
		slogLevel = slog.LevelWarn
	case "ERROR":
		slogLevel = slog.LevelError
	default:
		slogLevel = slog.LevelInfo
	}

	mu.Lock()
	defer mu.Unlock()

	// Close existing log file if any
	if logFile != nil {
		logFile.Close()
		logFile = nil
	}

	var handler slog.Handler

	if level == "DEBUG" {
		// For DEBUG level, use multistream to log to both file and console
		handler = createMultistreamHandler(slogLevel, logFilePath)
	} else {
		// For other levels, use console only
		handler = slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			Level: slogLevel,
		})
	}

	// Update instance and default
	instance = slog.New(handler)
	slog.SetDefault(instance)
}

// createMultistreamHandler creates a handler that writes to both file and console
func createMultistreamHandler(level slog.Level, logFilePath string) slog.Handler {
	// Use provided path or default
	if logFilePath == "" {
		logFilePath = "./logs/whaletui.log"
	}

	// Create logs directory if it doesn't exist
	logsDir := filepath.Dir(logFilePath)
	if err := os.MkdirAll(logsDir, 0o755); err != nil {
		// Fallback to console only if directory creation fails
		return slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			Level: level,
		})
	}

	// Open log file
	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		// Fallback to console only if file creation fails
		return slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			Level: level,
		})
	}

	logFile = file
	logPath = logFilePath

	// Create multistream writer
	multiWriter := io.MultiWriter(os.Stderr, file)

	// Create handler with multistream
	return slog.NewTextHandler(multiWriter, &slog.HandlerOptions{
		Level: level,
	})
}

// CloseLogFile closes the log file if it's open
func CloseLogFile() {
	if logFile != nil {
		logFile.Close()
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
func Debug(msg string, args ...any) {
	GetLogger().Debug(msg, args...)
}

func Info(msg string, args ...any) {
	GetLogger().Info(msg, args...)
}

func Warn(msg string, args ...any) {
	GetLogger().Warn(msg, args...)
}

func Error(msg string, args ...any) {
	GetLogger().Error(msg, args...)
}

func Fatal(msg string, args ...any) {
	GetLogger().Error(msg, args...)
	os.Exit(1)
}

// With returns a new logger with the given key-value pairs
func With(args ...any) *slog.Logger {
	return GetLogger().With(args...)
}

// WithGroup returns a new logger with the given group name
func WithGroup(name string) *slog.Logger {
	return GetLogger().WithGroup(name)
}
