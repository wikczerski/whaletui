package cmd

import (
	"fmt"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/wikczerski/whaletui/internal/config"
)

func TestRootCmd(t *testing.T) {
	// Test that root command is properly configured
	assert.NotNil(t, rootCmd)
	assert.Equal(t, "whaletui", rootCmd.Use)
	assert.Equal(t, "whaletui - Docker CLI Dashboard", rootCmd.Short)
	assert.Contains(t, rootCmd.Long, "Docker containers")
	assert.Contains(t, rootCmd.Long, "images, volumes, and networks")
}

func TestThemeCmd(t *testing.T) {
	// Test that theme command is properly configured
	assert.NotNil(t, themeCmd)
	assert.Equal(t, "theme", themeCmd.Use)
	assert.Equal(t, "Manage theme configuration", themeCmd.Short)
}

func TestRootCmdFlags(t *testing.T) {
	// Test that all expected global flags are registered
	flags := rootCmd.PersistentFlags()

	refreshFlag := flags.Lookup("refresh")
	assert.NotNil(t, refreshFlag)
	assert.Equal(t, "refresh", refreshFlag.Name)

	logLevelFlag := flags.Lookup("log-level")
	assert.NotNil(t, logLevelFlag)
	assert.Equal(t, "log-level", logLevelFlag.Name)

	themeFlag := flags.Lookup("theme")
	assert.NotNil(t, themeFlag)
	assert.Equal(t, "theme", themeFlag.Name)
}

func TestApplyFlagOverrides(t *testing.T) {
	// Test flag override functionality
	cfg := &config.Config{
		DockerHost:      "unix:///var/run/docker.sock",
		RemoteHost:      "",
		RefreshInterval: 5,
		LogLevel:        "INFO",
		Theme:           "default",
	}

	// Test refresh interval override
	refresh = 15
	applyFlagOverrides(cfg)
	assert.Equal(t, 15, cfg.RefreshInterval)

	// Test log level override
	logLevel = "DEBUG"
	applyFlagOverrides(cfg)
	assert.Equal(t, "DEBUG", cfg.LogLevel)

	// Test theme override
	theme = "dark"
	applyFlagOverrides(cfg)
	assert.Equal(t, "dark", cfg.Theme)

	// Reset global variables
	refresh = 5
	logLevel = "INFO"
	theme = ""
}

func TestSetLogLevel(t *testing.T) {
	assert.NotPanics(t, func() {
		setLogLevel("DEBUG")
	})

	assert.NotPanics(t, func() {
		setLogLevel("WARN")
	})

	assert.NotPanics(t, func() {
		setLogLevel("ERROR")
	})

	assert.NotPanics(t, func() {
		setLogLevel("INFO")
	})

	assert.NotPanics(t, func() {
		setLogLevel("INVALID")
	})
}

func TestRootCmdSubcommands(t *testing.T) {
	// Test that all expected subcommands are added
	subcommands := rootCmd.Commands()

	// Check that expected commands are present (don't check exact count as cobra may add built-in commands)
	commandNames := make([]string, len(subcommands))
	for i, cmd := range subcommands {
		commandNames[i] = cmd.Use
	}

	assert.Contains(t, commandNames, "connect [flags]")
	assert.Contains(t, commandNames, "theme")
}

func TestConnectCommand(t *testing.T) {
	// Test connect command flags
	connectFlags := connectCmd.Flags()

	// Check that required flags exist
	hostFlag := connectFlags.Lookup("host")
	assert.NotNil(t, hostFlag)
	assert.Equal(t, "host", hostFlag.Name)

	userFlag := connectFlags.Lookup("user")
	assert.NotNil(t, userFlag)
	assert.Equal(t, "user", userFlag.Name)

	portFlag := connectFlags.Lookup("port")
	assert.NotNil(t, portFlag)
	assert.Equal(t, "port", portFlag.Name)
	assert.Equal(t, "2375", portFlag.DefValue)
}

func TestRootCmdIntegration(t *testing.T) {
	// Test that root command can be created and configured without panicking
	cmd := &cobra.Command{}
	cmd.AddCommand(rootCmd)

	// This is a basic test to ensure the command structure is valid
	assert.NotNil(t, cmd.Commands())
	assert.Len(t, cmd.Commands(), 1)
	assert.Equal(t, "whaletui", cmd.Commands()[0].Use)
}

func TestIsDockerConnectionError(t *testing.T) {
	// Test various Docker connection error patterns
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "docker client creation failed",
			err:      assert.AnError,
			expected: false, // assert.AnError doesn't contain our patterns
		},
		{
			name:     "connection refused error",
			err:      fmt.Errorf("connection refused"),
			expected: true,
		},
		{
			name:     "permission denied error",
			err:      fmt.Errorf("permission denied"),
			expected: true,
		},
		{
			name:     "timeout error",
			err:      fmt.Errorf("timeout"),
			expected: true,
		},
		{
			name:     "docker client creation failed error",
			err:      fmt.Errorf("docker client creation failed"),
			expected: true,
		},
		{
			name:     "non-docker error",
			err:      fmt.Errorf("some other error"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isDockerConnectionError(tt.err)
			assert.Equal(t, tt.expected, result)
		})
	}
}
