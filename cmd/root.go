package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/wikczerski/D5r/internal/app"
	"github.com/wikczerski/D5r/internal/config"
	"github.com/wikczerski/D5r/internal/docker"
	"github.com/wikczerski/D5r/internal/logger"
)

var (
	refresh  int
	logLevel string
	theme    string
)

// connectCmd represents the connect command for remote Docker hosts
var connectCmd = &cobra.Command{
	Use:   "connect [flags]",
	Short: "Connect to a remote Docker host via SSH",
	Long: `Connect to a remote Docker host using SSH fallback when direct connection fails.
This command establishes an SSH connection to the remote host and sets up a Docker proxy.

Example:
  d5r connect --host 192.168.1.100 --user admin --port 2375`,
	RunE: runConnectCommand,
	// Disable automatic help display on errors
	SilenceUsage:  true,
	SilenceErrors: true,
}

// themeCmd represents the theme command
var themeCmd = &cobra.Command{
	Use:   "theme",
	Short: "Manage theme configuration",
	Long:  `Manage theme configuration for the D5r UI.`,
	RunE:  runThemeCommand,
	// Disable automatic help display on errors
	SilenceUsage:  true,
	SilenceErrors: true,
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
  • Remote Host Support: Connect to Docker hosts on different machines with SSH fallback
  • SSH Fallback: Automatic SSH connection when direct Docker connection fails
  • Configurable Ports: Customize SSH fallback proxy ports to avoid conflicts
  • Theme Support: Customizable color schemes via YAML/JSON configuration

Commands:
  connect  - Connect to a remote Docker host via SSH
  theme    - Manage theme configuration

Examples:
  d5r                    - Start with local Docker instance
  d5r connect --host 192.168.1.100 --user admin  - Connect to remote host
  d5r theme             - Manage theme configuration`,
	RunE: runApp,
	// Disable automatic help display on errors
	SilenceUsage:  true,
	SilenceErrors: true,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	// Set up the command to not show help on errors
	rootCmd.SilenceUsage = true
	rootCmd.SilenceErrors = true

	err := rootCmd.Execute()
	if err != nil {
		// Log the error and exit without showing help
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().IntVar(&refresh, "refresh", 5, "Refresh interval in seconds")
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "INFO", "Log level (DEBUG, INFO, WARN, ERROR)")
	rootCmd.PersistentFlags().StringVar(&theme, "theme", "", "Path to theme configuration file (YAML/JSON)")

	// Connect command flags
	connectCmd.Flags().String("host", "", "Remote Docker host (e.g., 192.168.1.100 or tcp://192.168.1.100)")
	connectCmd.Flags().String("user", "", "SSH username for remote host connection")
	connectCmd.Flags().Int("port", 2375, "Port for SSH fallback Docker proxy (default: 2375)")
	connectCmd.Flags().Bool("diagnose", false, "Run SSH connection diagnostics before connecting")
	_ = connectCmd.MarkFlagRequired("host")
	_ = connectCmd.MarkFlagRequired("user")

	// Add subcommands
	rootCmd.AddCommand(connectCmd)
	rootCmd.AddCommand(themeCmd)
}

// runConnectCommand handles the connect subcommand for remote Docker hosts
func runConnectCommand(cmd *cobra.Command, _ []string) error {
	log := logger.GetLogger()
	log.SetPrefix("Connect")

	// Get flags from the connect command
	host, _ := cmd.Flags().GetString("host")
	user, _ := cmd.Flags().GetString("user")
	port, _ := cmd.Flags().GetInt("port")
	diagnose, _ := cmd.Flags().GetBool("diagnose")

	cfg, err := config.Load()
	if err != nil {
		log.Error("Config load failed: %v", err)
		return err
	}

	// Set remote connection parameters
	cfg.RemoteHost = host
	cfg.RemoteUser = user
	cfg.RemotePort = port
	cfg.DockerHost = host

	setLogLevel(log, cfg.LogLevel)

	// Run diagnostics if requested
	if diagnose {
		log.Info("Running SSH connection diagnostics...")
		if err := runSSHDiagnostics(host, user, port); err != nil {
			log.Error("SSH diagnostics failed: %v", err)
			log.Info("Continuing with connection attempt...")
		} else {
			log.Info("SSH diagnostics passed successfully")
		}
	}

	log.Info("Connecting to remote Docker host: %s as user: %s", host, user)

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

// runSSHDiagnostics runs SSH connection diagnostics
func runSSHDiagnostics(host, user string, _ int) error {
	// Create SSH client for diagnostics
	sshHost := fmt.Sprintf("%s@%s", user, host)
	sshClient, err := docker.NewSSHClient(sshHost, 22) // Default SSH port
	if err != nil {
		return fmt.Errorf("failed to create SSH client for diagnostics: %w", err)
	}

	// Run diagnostics
	if err := sshClient.DiagnoseConnection(); err != nil {
		return fmt.Errorf("SSH diagnostics failed: %w", err)
	}

	return nil
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

	// Validate that --user is provided when --host is specified
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
	cleanupTerminalClearScreen()
	cleanupTerminalResetColors()
	cleanupTerminalShowCursor()
	cleanupTerminalMoveCursorToTop()
	cleanupTerminalSyncStdout()
}

// cleanupTerminalClearScreen clears the terminal screen
func cleanupTerminalClearScreen() {
	if _, err := fmt.Fprint(os.Stdout, "\033[2J"); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to clear screen: %v\n", err)
	}
}

// cleanupTerminalResetColors resets terminal colors
func cleanupTerminalResetColors() {
	if _, err := fmt.Fprint(os.Stdout, "\033[0m"); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to reset colors: %v\n", err)
	}
}

// cleanupTerminalShowCursor shows the terminal cursor
func cleanupTerminalShowCursor() {
	if _, err := fmt.Fprint(os.Stdout, "\033[?25h"); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to show cursor: %v\n", err)
	}
}

// cleanupTerminalMoveCursorToTop moves the cursor to the top of the terminal
func cleanupTerminalMoveCursorToTop() {
	if _, err := fmt.Fprint(os.Stdout, "\033[H"); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to move cursor: %v\n", err)
	}
}

// cleanupTerminalSyncStdout synchronizes stdout
func cleanupTerminalSyncStdout() {
	if err := os.Stdout.Sync(); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to sync stdout: %v\n", err)
	}
}

// applyFlagOverrides applies command line flag values to the configuration
func applyFlagOverrides(cfg *config.Config) {
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
