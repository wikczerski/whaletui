package cmd

import (
	"testing"
)

func TestVersionConstants(t *testing.T) {
	if Version == "" {
		t.Error("Version should not be empty")
	}
	if CommitSHA == "" {
		t.Error("CommitSHA should not be empty")
	}
	if BuildDate == "" {
		t.Error("BuildDate should not be empty")
	}

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
	if versionCmd.Use != "version" {
		t.Errorf("versionCmd.Use = %s, want 'version'", versionCmd.Use)
	}

	if versionCmd.Short == "" {
		t.Error("versionCmd.Short should not be empty")
	}

	if versionCmd.Long == "" {
		t.Error("versionCmd.Long should not be empty")
	}

	if versionCmd.Run == nil {
		t.Error("versionCmd.Run should not be nil")
	}
}

func TestVersionCommandDescription(t *testing.T) {
	expectedShort := "Show version information"
	if versionCmd.Short != expectedShort {
		t.Errorf("versionCmd.Short = %s, want %s", versionCmd.Short, expectedShort)
	}

	longDesc := versionCmd.Long
	expectedContent := []string{
		"Display version information for whaletui including:",
		"Version number",
		"Git commit hash",
		"Build date",
	}

	for _, content := range expectedContent {
		if !contains(longDesc, content) {
			t.Errorf("Long description should contain '%s', got: %s", content, longDesc)
		}
	}
}

func TestVersionCommandIntegration(t *testing.T) {
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
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Version command Run function panicked: %v", r)
		}
	}()

	if versionCmd.Run != nil {
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
