package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/wikczerski/D5r/internal/app"
	"github.com/wikczerski/D5r/internal/config"
	"github.com/wikczerski/D5r/internal/logger"
)

var (
	remoteHost string
	refresh    int
	logLevel   string
	theme      string
)

// themeCmd represents the theme command
var themeCmd = &cobra.Command{
	Use:   "theme",
	Short: "Manage theme configuration",
	Long:  `Manage theme configuration for the D5r UI.`,
	RunE:  runThemeCommand,
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "d5r",
	Short: "D5r - Docker CLI Dashboard",
	Long: `D5r is a terminal-based Docker management tool inspired by k9s,
providing an intuitive and powerful interface for managing Docker containers,
images, volumes, and networks with a modern, responsive TUI.

Features:
  • Container Management: View, start, stop, restart, delete, and manage containers
  • Image Management: Browse, inspect, and remove Docker images
  • Volume Management: Manage Docker volumes with ease
  • Network Management: View and manage Docker networks
  • Remote Host Support: Connect to Docker hosts on different machines
  • Theme Support: Customizable color schemes via YAML/JSON configuration`,
	RunE: runApp,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&remoteHost, "host", "", "Remote Docker host (e.g., tcp://192.168.1.100:2375)")
	rootCmd.PersistentFlags().IntVar(&refresh, "refresh", 5, "Refresh interval in seconds")
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "INFO", "Log level (DEBUG, INFO, WARN, ERROR)")
	rootCmd.PersistentFlags().StringVar(&theme, "theme", "", "Path to theme configuration file (YAML/JSON)")

	// Add theme subcommand
	rootCmd.AddCommand(themeCmd)
}

// runApp runs the main application with the provided configuration
func runApp(_ *cobra.Command, _ []string) error {
	log := logger.GetLogger()
	log.SetPrefix("D5r")

	cfg, err := config.Load()
	if err != nil {
		log.Error("Config load failed: %v", err)
		return err
	}

	applyFlagOverrides(cfg)
	setLogLevel(log, cfg.LogLevel)

	if cfg.RemoteHost != "" {
		log.Info("Connecting to remote Docker host: %s", cfg.RemoteHost)
	} else {
		log.Info("Connecting to local Docker instance")
	}

	application, err := app.New(cfg)
	if err != nil {
		log.Error("App init failed: %v", err)
		return err
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	uiShutdownCh := application.GetUI().GetShutdownChan()

	go func() {
		if err := application.Run(); err != nil {
			log.Error("App run failed: %v", err)
		}
	}()

	select {
	case <-sigCh:
		log.Info("Received shutdown signal, shutting down gracefully...")
	case <-uiShutdownCh:
		log.Info("Received UI shutdown signal, shutting down gracefully...")
	}

	defer cleanupTerminal()

	application.Shutdown()
	return nil
}

// cleanupTerminal performs additional terminal cleanup operations
func cleanupTerminal() {
	if _, err := fmt.Fprint(os.Stdout, "\033[2J"); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to clear screen: %v\n", err)
	}
	if _, err := fmt.Fprint(os.Stdout, "\033[0m"); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to reset colors: %v\n", err)
	}
	if _, err := fmt.Fprint(os.Stdout, "\033[?25h"); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to show cursor: %v\n", err)
	}
	if _, err := fmt.Fprint(os.Stdout, "\033[H"); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to move cursor: %v\n", err)
	}
	if err := os.Stdout.Sync(); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to sync stdout: %v\n", err)
	}
}

// applyFlagOverrides applies command line flag values to the configuration
func applyFlagOverrides(cfg *config.Config) {
	if remoteHost != "" {
		cfg.RemoteHost = remoteHost
		cfg.DockerHost = remoteHost
	}

	if refresh != 5 {
		cfg.RefreshInterval = refresh
	}

	if logLevel != "INFO" {
		cfg.LogLevel = logLevel
	}

	// Handle theme configuration
	if theme != "" {
		cfg.Theme = theme
	}
}

// setLogLevel sets the log level based on the configuration
func setLogLevel(log *logger.Logger, level string) {
	switch level {
	case "DEBUG":
		log.SetLevel(logger.DEBUG)
	case "WARN":
		log.SetLevel(logger.WARN)
	case "ERROR":
		log.SetLevel(logger.ERROR)
	default:
		log.SetLevel(logger.INFO)
	}
}

// runThemeCommand handles theme-related commands
func runThemeCommand(_ *cobra.Command, _ []string) error {
	log := logger.GetLogger()
	log.SetPrefix("Theme")

	// Create a default theme manager
	themeManager := config.NewThemeManager("")

	// Save the current theme to the default location
	defaultPath := "./config/theme.yaml"
	if err := themeManager.SaveTheme(defaultPath); err != nil {
		log.Error("Failed to save theme: %v", err)
		return err
	}

	log.Info("Theme configuration saved to: %s", defaultPath)
	log.Info("You can now customize the colors in this file and restart D5r with --theme %s", defaultPath)

	return nil
}
