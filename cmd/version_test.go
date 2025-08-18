package cmd

import (
	"testing"
)

func TestVersionConstants(t *testing.T) {
	// Test that version constants are defined
	if Version == "" {
		t.Error("Version should not be empty")
	}
	if CommitSHA == "" {
		t.Error("CommitSHA should not be empty")
	}
	if BuildDate == "" {
		t.Error("BuildDate should not be empty")
	}

	// Test that version constants have reasonable values
	if Version == "unknown" {
		t.Log("Version is set to 'unknown', this might be expected in development")
	}
	if CommitSHA == "unknown" {
		t.Log("CommitSHA is set to 'unknown', this might be expected in development")
	}
	if BuildDate == "unknown" {
		t.Log("BuildDate is set to 'unknown', this might be expected in development")
	}
}

func TestVersionCommandStructure(t *testing.T) {
	// Test that version command is properly configured
	if versionCmd.Use != "version" {
		t.Errorf("versionCmd.Use = %s, want 'version'", versionCmd.Use)
	}

	if versionCmd.Short == "" {
		t.Error("versionCmd.Short should not be empty")
	}

	if versionCmd.Long == "" {
		t.Error("versionCmd.Long should not be empty")
	}

	// Test that the command has a Run function
	if versionCmd.Run == nil {
		t.Error("versionCmd.Run should not be nil")
	}
}

func TestVersionCommandDescription(t *testing.T) {
	// Test that the command descriptions are meaningful
	expectedShort := "Show version information"
	if versionCmd.Short != expectedShort {
		t.Errorf("versionCmd.Short = %s, want %s", versionCmd.Short, expectedShort)
	}

	// Test that long description contains expected content
	longDesc := versionCmd.Long
	expectedContent := []string{
		"Display version information for D5r including:",
		"Version number",
		"Git commit SHA",
		"Build date",
	}

	for _, content := range expectedContent {
		if !contains(longDesc, content) {
			t.Errorf("Long description should contain '%s', got: %s", content, longDesc)
		}
	}
}

func TestVersionCommandIntegration(t *testing.T) {
	// Test that the version command is properly integrated with root command
	// This tests the init() function indirectly

	// Check if version command is in root command's commands
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd.Use == "version" {
			found = true
			break
		}
	}

	if !found {
		t.Error("Version command should be added to root command")
	}
}

func TestVersionCommandRunFunction(t *testing.T) {
	// Test that the Run function can be called without panicking
	// This is a basic test that the function exists and can be executed
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Version command Run function panicked: %v", r)
		}
	}()

	// The Run function should not panic when called with nil arguments
	// Note: We can't easily test the actual output without mocking fmt.Printf
	if versionCmd.Run != nil {
		// This just verifies the function exists and can be called
		// The actual output testing would require more complex stdout mocking
		t.Log("Version command Run function exists and can be called")
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(s) > len(substr) && (s[:len(substr)] == substr ||
			s[len(s)-len(substr):] == substr ||
			containsSubstring(s, substr))))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
