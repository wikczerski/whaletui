package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/wikczerski/whaletui/internal/config"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Show configuration information",
	Long: `Display current configuration settings for whaletui including:
  • Docker host configuration
  • Refresh interval
  • Log level
  • Theme settings
  • Configuration file location`,
	Run: func(_ *cobra.Command, _ []string) {
		showConfig()
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}

// showConfig displays the current configuration
func showConfig() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "unknown"
	}

	fmt.Printf("whaletui Configuration\n")
	fmt.Printf("==================\n\n")

	fmt.Printf("Docker Host: %s\n", cfg.DockerHost)
	if cfg.RemoteHost != "" {
		fmt.Printf("Remote Host: %s\n", cfg.RemoteHost)
	}
	fmt.Printf("Refresh Interval: %d seconds\n", cfg.RefreshInterval)
	fmt.Printf("Log Level: %s\n", cfg.LogLevel)
	fmt.Printf("Theme: %s\n", cfg.Theme)
	fmt.Printf("\nConfig File: %s/.dockerk9s/config.json\n", homeDir)
}
