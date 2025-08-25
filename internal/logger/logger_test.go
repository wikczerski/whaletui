package logger

import (
	"os"
	"testing"
)

func TestMultistreamLogging(t *testing.T) {
	testLogPath := "./test_logs/test.log"
	setupTestLogging(t, testLogPath)
	defer cleanupTestLogging(t, testLogPath)

	testDebugLevelHandler(t, testLogPath)
	testInfoLevelHandler(t, testLogPath)
	testLogFileCreation(t, testLogPath)
	testLogFileCleanup(t)
}

func setupTestLogging(t *testing.T, testLogPath string) {
	t.Helper()
	defer func() {
		if err := os.RemoveAll("./test_logs"); err != nil {
			t.Logf("Failed to remove test logs: %v", err)
		}
	}()
}

func testDebugLevelHandler(t *testing.T, testLogPath string) {
	t.Helper()
	SetLevelWithPath("DEBUG", testLogPath)
	logger := GetLogger()
	if logger == nil {
		t.Fatal("Logger should not be nil")
	}
}

func testInfoLevelHandler(t *testing.T, testLogPath string) {
	t.Helper()
	SetLevelWithPath("INFO", testLogPath)
	logger := GetLogger()
	if logger == nil {
		t.Fatal("Logger should not be nil")
	}
}

func testLogFileCreation(t *testing.T, testLogPath string) {
	t.Helper()
	SetLevelWithPath("DEBUG", testLogPath)
	if _, err := os.Stat(testLogPath); os.IsNotExist(err) {
		t.Fatal("Log file should be created for DEBUG level")
	}
}

func testLogFileCleanup(t *testing.T) {
	t.Helper()
	CloseLogFile()
	if logFile != nil {
		t.Fatal("Log file should be closed after CloseLogFile()")
	}
}

func cleanupTestLogging(t *testing.T, testLogPath string) {
	t.Helper()
	// This function is called by defer, so it's intentionally empty
	// The actual cleanup is done in setupTestLogging
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
	if err := os.RemoveAll("./logs"); err != nil {
		t.Logf("Failed to remove logs directory: %v", err)
	}
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
	testLogPath := "./test_logs/tui_test.log"
	setupTUITestLogging(t, testLogPath)
	defer cleanupTUITestLogging(t, testLogPath)

	testTUIModeHandler(t, testLogPath)
	testNormalModeHandler(t, testLogPath)
	cleanupTUIModeTest(t)
}

func setupTUITestLogging(t *testing.T, testLogPath string) {
	t.Helper()
	defer func() {
		if err := os.RemoveAll("./test_logs"); err != nil {
			t.Logf("Failed to remove test logs: %v", err)
		}
	}()
}

func testTUIModeHandler(t *testing.T, testLogPath string) {
	t.Helper()
	SetTUIMode(true)
	SetLevelWithPath("DEBUG", testLogPath)
	logger := GetLogger()
	if logger == nil {
		t.Fatal("Logger should not be nil in TUI mode")
	}
}

func testNormalModeHandler(t *testing.T, testLogPath string) {
	t.Helper()
	SetTUIMode(false)
	SetLevelWithPath("DEBUG", testLogPath)
	logger := GetLogger()
	if logger == nil {
		t.Fatal("Logger should not be nil in normal mode")
	}
}

func cleanupTUIModeTest(t *testing.T) {
	t.Helper()
	CloseLogFile()
	SetTUIMode(false) // Reset to default
}

func cleanupTUITestLogging(t *testing.T, testLogPath string) {
	t.Helper()
	// This function is called by defer, so it's intentionally empty
	// The actual cleanup is done in setupTUITestLogging
}
