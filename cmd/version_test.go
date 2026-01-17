package cmd

import (
	"strings"
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
	if versionCmd.Run == nil {
		t.Fatal("Version command Run function should not be nil")
	}
	// Run function is verified in integration tests to avoid stdout interference
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
