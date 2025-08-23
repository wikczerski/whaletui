package logger

import (
	"os"
	"testing"
)

func TestMultistreamLogging(t *testing.T) {
	// Clean up any existing log files
	testLogPath := "./test_logs/test.log"
	defer os.RemoveAll("./test_logs")

	// Test that DEBUG level creates multistream handler
	SetLevelWithPath("DEBUG", testLogPath)

	logger := GetLogger()
	if logger == nil {
		t.Fatal("Logger should not be nil")
	}

	// Test that INFO level creates console-only handler
	SetLevelWithPath("INFO", testLogPath)

	logger = GetLogger()
	if logger == nil {
		t.Fatal("Logger should not be nil")
	}

	// Test that log file is created for DEBUG level
	SetLevelWithPath("DEBUG", testLogPath)

	// Verify log file exists
	if _, err := os.Stat(testLogPath); os.IsNotExist(err) {
		t.Fatal("Log file should be created for DEBUG level")
	}

	// Test log file cleanup
	CloseLogFile()

	// Verify log file is closed
	if logFile != nil {
		t.Fatal("Log file should be closed after CloseLogFile()")
	}
}

func TestLogLevelParsing(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"DEBUG", "DEBUG"},
		{"INFO", "INFO"},
		{"WARN", "WARN"},
		{"ERROR", "ERROR"},
		{"UNKNOWN", "INFO"}, // Default case
	}

	for _, tc := range testCases {
		SetLevel(tc.input)
		logger := GetLogger()
		if logger == nil {
			t.Fatalf("Logger should not be nil for level %s", tc.input)
		}
	}
}

func TestDefaultLogPath(t *testing.T) {
	// Test that default path is used when empty string is provided
	SetLevelWithPath("DEBUG", "")

	// Verify default log file exists
	defaultPath := "./logs/whaletui.log"
	if _, err := os.Stat(defaultPath); os.IsNotExist(err) {
		t.Fatal("Default log file should be created")
	}

	// Clean up
	CloseLogFile()
	os.RemoveAll("./logs")
}

func TestTUIModeToggle(t *testing.T) {
	// Test initial state
	initialState := IsTUIMode()

	// Test enabling TUI mode
	SetTUIMode(true)
	if !IsTUIMode() {
		t.Fatal("TUI mode should be enabled")
	}

	// Test disabling TUI mode
	SetTUIMode(false)
	if IsTUIMode() {
		t.Fatal("TUI mode should be disabled")
	}

	// Restore initial state
	SetTUIMode(initialState)
}

func TestTUIModeLogging(t *testing.T) {
	// Clean up any existing log files
	testLogPath := "./test_logs/tui_test.log"
	defer os.RemoveAll("./test_logs")

	// Test TUI mode uses file-only handler
	SetTUIMode(true)
	SetLevelWithPath("DEBUG", testLogPath)

	logger := GetLogger()
	if logger == nil {
		t.Fatal("Logger should not be nil in TUI mode")
	}

	// Test normal mode uses multistream handler
	SetTUIMode(false)
	SetLevelWithPath("DEBUG", testLogPath)

	logger = GetLogger()
	if logger == nil {
		t.Fatal("Logger should not be nil in normal mode")
	}

	// Clean up
	CloseLogFile()
	SetTUIMode(false) // Reset to default
}
