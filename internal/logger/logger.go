package logger

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

// LogLevel represents the logging level
type LogLevel int

const (
	// DEBUG represents debug logging level
	DEBUG LogLevel = iota
	// INFO represents info logging level
	INFO
	// WARN represents warning logging level
	WARN
	// ERROR represents error logging level
	ERROR
	// FATAL represents fatal logging level
	FATAL
)

var levelNames = map[LogLevel]string{
	DEBUG: "DBG",
	INFO:  "INF",
	WARN:  "WRN",
	ERROR: "ERR",
	FATAL: "FAT",
}

// Logger represents a logger instance
type Logger struct {
	mu        sync.Mutex
	output    io.Writer
	minLevel  LogLevel
	logPrefix string
	enabled   bool
}

var (
	defaultLogger *Logger
	once          sync.Once
)

// GetLogger returns the default logger instance
func GetLogger() *Logger {
	once.Do(func() {
		defaultLogger = &Logger{
			output:    os.Stderr,
			minLevel:  INFO,
			logPrefix: "",
			enabled:   true,
		}
	})
	return defaultLogger
}

// SetOutput sets the output writer for the logger
func (l *Logger) SetOutput(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.output = w
}

// SetLevel sets the minimum logging level
func (l *Logger) SetLevel(level LogLevel) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.minLevel = level
}

// SetPrefix sets the log prefix
func (l *Logger) SetPrefix(prefix string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.logPrefix = prefix
}

// Enable enables or disables logging
func (l *Logger) Enable(enabled bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.enabled = enabled
}

// IsEnabled returns whether logging is enabled
func (l *Logger) IsEnabled() bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.enabled
}

func (l *Logger) log(level LogLevel, format string, args ...any) {
	if !l.enabled || level < l.minLevel {
		return
	}

	l.mu.Lock()

	prefix := ""
	if l.logPrefix != "" {
		prefix = "[" + l.logPrefix + "] "
	}

	message := fmt.Sprintf(format, args...)
	message = strings.TrimSuffix(message, "\n")

	if level == INFO {
		if _, err := fmt.Fprintf(l.output, "%s%s\n", prefix, message); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write log message: %v\n", err)
		}
	} else {
		timestamp := time.Now().Format("15:04:05")
		if _, err := fmt.Fprintf(l.output, "%s [%s] %s%s\n", timestamp, levelNames[level], prefix, message); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write log message: %v\n", err)
		}
	}

	if level == FATAL {
		l.mu.Unlock() // Unlock before exit
		os.Exit(1)
	}
	l.mu.Unlock()
}

// Debug logs a debug message
func (l *Logger) Debug(format string, args ...any) {
	l.log(DEBUG, format, args...)
}

// Info logs an info message
func (l *Logger) Info(format string, args ...any) {
	l.log(INFO, format, args...)
}

// Warn logs a warning message
func (l *Logger) Warn(format string, args ...any) {
	l.log(WARN, format, args...)
}

// Error logs an error message
func (l *Logger) Error(format string, args ...any) {
	l.log(ERROR, format, args...)
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(format string, args ...any) {
	l.log(FATAL, format, args...)
}

// Debug logs a debug message using the default logger
func Debug(format string, args ...any) {
	GetLogger().Debug(format, args...)
}

// Info logs an info message using the default logger
func Info(format string, args ...any) {
	GetLogger().Info(format, args...)
}

// Warn logs a warning message using the default logger
func Warn(format string, args ...any) {
	GetLogger().Warn(format, args...)
}

// Error logs an error message using the default logger
func Error(format string, args ...any) {
	GetLogger().Error(format, args...)
}

// Fatal logs a fatal message using the default logger and exits
func Fatal(format string, args ...any) {
	GetLogger().Fatal(format, args...)
}
