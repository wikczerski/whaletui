package logger

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLogger(t *testing.T) {
	logger1 := GetLogger()
	logger2 := GetLogger()

	// Should return the same instance
	if logger1 != logger2 {
		t.Errorf("Expected same logger instance")
	}
	if logger1 == nil {
		t.Errorf("Expected non-nil logger")
	}
}

func TestLogger_SetOutput(t *testing.T) {
	logger := GetLogger()
	buf := &bytes.Buffer{}

	logger.SetOutput(buf)
	logger.Info("test message")

	assert.Contains(t, buf.String(), "test message")
}

func TestLogger_SetLevel(t *testing.T) {
	logger := GetLogger()
	buf := &bytes.Buffer{}
	logger.SetOutput(buf)

	// Test DEBUG level
	logger.SetLevel(DEBUG)
	logger.Debug("debug message")
	assert.Contains(t, buf.String(), "debug message")

	// Test INFO level
	buf.Reset()
	logger.SetLevel(INFO)
	logger.Debug("debug message")
	logger.Info("info message")
	assert.NotContains(t, buf.String(), "debug message")
	assert.Contains(t, buf.String(), "info message")

	// Test WARN level
	buf.Reset()
	logger.SetLevel(WARN)
	logger.Info("info message")
	logger.Warn("warn message")
	assert.NotContains(t, buf.String(), "info message")
	assert.Contains(t, buf.String(), "warn message")

	// Test ERROR level
	buf.Reset()
	logger.SetLevel(ERROR)
	logger.Warn("warn message")
	logger.Error("error message")
	assert.NotContains(t, buf.String(), "warn message")
	assert.Contains(t, buf.String(), "error message")
}

func TestLogger_SetPrefix(t *testing.T) {
	logger := GetLogger()
	buf := &bytes.Buffer{}
	logger.SetOutput(buf)

	// Reset logger state
	logger.SetPrefix("")
	logger.SetLevel(INFO)

	logger.SetPrefix("TEST")
	logger.Info("message")

	output := buf.String()
	assert.Contains(t, output, "[TEST]")
	assert.Contains(t, output, "message")
}

func TestLogger_Enable(t *testing.T) {
	logger := GetLogger()
	buf := &bytes.Buffer{}
	logger.SetOutput(buf)

	// Reset logger state
	logger.SetPrefix("")
	logger.SetLevel(INFO)

	// Disable logging
	logger.Enable(false)
	logger.Info("should not appear")

	assert.Empty(t, buf.String())

	// Re-enable logging
	logger.Enable(true)
	logger.Info("should appear")

	assert.Contains(t, buf.String(), "should appear")
}

func TestLogger_IsEnabled(t *testing.T) {
	logger := GetLogger()

	// Should be enabled by default
	assert.True(t, logger.IsEnabled())

	// Disable and check
	logger.Enable(false)
	assert.False(t, logger.IsEnabled())

	// Re-enable and check
	logger.Enable(true)
	assert.True(t, logger.IsEnabled())
}

func TestLogger_LogLevels(t *testing.T) {
	logger := GetLogger()
	buf := &bytes.Buffer{}
	logger.SetOutput(buf)
	logger.SetLevel(DEBUG)

	// Test all log levels
	logger.Debug("debug message")
	logger.Info("info message")
	logger.Warn("warn message")
	logger.Error("error message")

	output := buf.String()
	assert.Contains(t, output, "debug message")
	assert.Contains(t, output, "info message")
	assert.Contains(t, output, "warn message")
	assert.Contains(t, output, "error message")
}

func TestLogger_TimestampFormat(t *testing.T) {
	logger := GetLogger()
	buf := &bytes.Buffer{}
	logger.SetOutput(buf)

	// Reset logger state
	logger.SetPrefix("")
	logger.SetLevel(DEBUG)

	// Test that non-INFO messages include timestamp
	logger.Debug("debug message")
	logger.Warn("warn message")

	output := buf.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")

	// Check timestamp format (HH:MM:SS)
	for _, line := range lines {
		if strings.Contains(line, "debug message") || strings.Contains(line, "warn message") {
			// Should contain timestamp in format [HH:MM:SS]
			if !strings.Contains(line, "[") || !strings.Contains(line, "]") {
				t.Errorf("Expected timestamp format, got: %s", line)
			}
		}
	}
}

func TestLogger_InfoNoTimestamp(t *testing.T) {
	logger := GetLogger()
	buf := &bytes.Buffer{}
	logger.SetOutput(buf)

	logger.Info("info message")

	output := buf.String()
	// INFO messages should not have timestamp
	assert.NotRegexp(t, `\[\d{2}:\d{2}:\d{2}\]`, output)
	assert.Contains(t, output, "info message")
}

func TestLogger_ConcurrentAccess(t *testing.T) {
	logger := GetLogger()
	buf := &bytes.Buffer{}
	logger.SetOutput(buf)

	// Test concurrent access with simple messages
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func(_ int) {
			logger.Info("message from goroutine")
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	// Should have 10 messages
	output := buf.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")
	assert.Len(t, lines, 10)
}

func TestLogger_Formatting(t *testing.T) {
	logger := GetLogger()
	buf := &bytes.Buffer{}
	logger.SetOutput(buf)

	// Test formatting with arguments
	logger.Info("formatted message: string, 42")

	output := buf.String()
	assert.Contains(t, output, "formatted message: string, 42")
}

func TestLogger_NewlineHandling(t *testing.T) {
	logger := GetLogger()
	buf := &bytes.Buffer{}
	logger.SetOutput(buf)

	// Test message with trailing newline
	logger.Info("message with newline\n")

	output := buf.String()
	// Should not have double newlines
	assert.NotContains(t, output, "\n\n")
	assert.Contains(t, output, "message with newline")
}

func TestGlobalFunctions(t *testing.T) {
	// Test global convenience functions
	buf := &bytes.Buffer{}
	originalOutput := GetLogger().output
	defer GetLogger().SetOutput(originalOutput)

	GetLogger().SetOutput(buf)
	GetLogger().SetLevel(DEBUG)

	Debug("global debug")
	Info("global info")
	Warn("global warn")
	Error("global error")

	output := buf.String()
	assert.Contains(t, output, "global debug")
	assert.Contains(t, output, "global info")
	assert.Contains(t, output, "global warn")
	assert.Contains(t, output, "global error")
}

func TestLogger_LevelNames(t *testing.T) {
	// Test level name mapping
	assert.Equal(t, "DBG", levelNames[DEBUG])
	assert.Equal(t, "INF", levelNames[INFO])
	assert.Equal(t, "WRN", levelNames[WARN])
	assert.Equal(t, "ERR", levelNames[ERROR])
	assert.Equal(t, "FAT", levelNames[FATAL])
}

func TestLogger_EdgeCases(t *testing.T) {
	logger := GetLogger()
	buf := &bytes.Buffer{}
	logger.SetOutput(buf)

	// Reset logger state
	logger.SetPrefix("")
	logger.SetLevel(INFO)

	// Test empty message
	logger.Info("")
	assert.Empty(t, strings.TrimSpace(buf.String()))

	// Test very long message
	logger.Info("a")
	assert.Contains(t, buf.String(), "a")

	// Test special characters
	logger.Info("message with special chars: !@#$&*()")
	assert.Contains(t, buf.String(), "message with special chars")
}
