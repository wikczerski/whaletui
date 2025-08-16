package logger

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

var levelNames = map[LogLevel]string{
	DEBUG: "DBG",
	INFO:  "INF",
	WARN:  "WRN",
	ERROR: "ERR",
	FATAL: "FAT",
}

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

func (l *Logger) SetOutput(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.output = w
}

func (l *Logger) SetLevel(level LogLevel) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.minLevel = level
}

func (l *Logger) SetPrefix(prefix string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.logPrefix = prefix
}

func (l *Logger) Enable(enabled bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.enabled = enabled
}

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
	defer l.mu.Unlock()

	prefix := ""
	if l.logPrefix != "" {
		prefix = "[" + l.logPrefix + "] "
	}

	message := fmt.Sprintf(format, args...)
	message = strings.TrimSuffix(message, "\n")

	if level == INFO {
		fmt.Fprintf(l.output, "%s%s\n", prefix, message)
	} else {
		timestamp := time.Now().Format("15:04:05")
		fmt.Fprintf(l.output, "%s [%s] %s%s\n", timestamp, levelNames[level], prefix, message)
	}

	if level == FATAL {
		os.Exit(1)
	}
}

func (l *Logger) Debug(format string, args ...any) {
	l.log(DEBUG, format, args...)
}

func (l *Logger) Info(format string, args ...any) {
	l.log(INFO, format, args...)
}

func (l *Logger) Warn(format string, args ...any) {
	l.log(WARN, format, args...)
}

func (l *Logger) Error(format string, args ...any) {
	l.log(ERROR, format, args...)
}

func (l *Logger) Fatal(format string, args ...any) {
	l.log(FATAL, format, args...)
}

func Debug(format string, args ...any) {
	GetLogger().Debug(format, args...)
}

func Info(format string, args ...any) {
	GetLogger().Info(format, args...)
}

func Warn(format string, args ...any) {
	GetLogger().Warn(format, args...)
}

func Error(format string, args ...any) {
	GetLogger().Error(format, args...)
}

func Fatal(format string, args ...any) {
	GetLogger().Fatal(format, args...)
}
