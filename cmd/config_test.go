package cmd

import (
	"bytes"
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestConfigCmd(t *testing.T) {
	// Test that config command is properly registered
	assert.NotNil(t, configCmd)
	assert.Equal(t, "config", configCmd.Use)
	assert.Equal(t, "Show configuration information", configCmd.Short)
}

func TestConfigCmdExecution(t *testing.T) {
	// Test that the config command can be executed without panicking
	cmd := &cobra.Command{}
	cmd.AddCommand(configCmd)

	// This is a basic test to ensure the command structure is valid
	assert.NotNil(t, cmd.Commands())
	assert.Len(t, cmd.Commands(), 1)
	assert.Equal(t, "config", cmd.Commands()[0].Use)
}

func TestConfigCmdLongDescription(t *testing.T) {
	// Test that the long description contains expected content
	assert.Contains(t, configCmd.Long, "Docker host configuration")
	assert.Contains(t, configCmd.Long, "Refresh interval")
	assert.Contains(t, configCmd.Long, "Log level")
	assert.Contains(t, configCmd.Long, "Theme settings")
	assert.Contains(t, configCmd.Long, "Configuration file location")
}

func TestShowConfig_PrintsConfiguration(t *testing.T) {
	// Isolate HOME to avoid touching the real user config
	tempHome := t.TempDir()
	origHome := os.Getenv("HOME")
	origUserProfile := os.Getenv("USERPROFILE")
	_ = os.Setenv("HOME", tempHome)
	_ = os.Setenv("USERPROFILE", tempHome)
	t.Cleanup(func() {
		_ = os.Setenv("HOME", origHome)
		_ = os.Setenv("USERPROFILE", origUserProfile)
	})

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	t.Cleanup(func() { os.Stdout = oldStdout })

	// Run
	showConfig()

	// Read output
	_ = w.Close()
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	out := buf.String()

	assert.Contains(t, out, "D5r Configuration")
	assert.Contains(t, out, "Docker Host:")
	assert.Contains(t, out, "Refresh Interval:")
	assert.Contains(t, out, "Log Level:")
	assert.Contains(t, out, "Theme:")
}
