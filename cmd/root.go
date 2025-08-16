package cmd

import (
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
  • Remote Host Support: Connect to Docker hosts on different machines`,
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
	rootCmd.PersistentFlags().StringVar(&theme, "theme", "default", "UI theme")
}

// runApp runs the main application with the provided configuration
func runApp(cmd *cobra.Command, args []string) error {
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

	go func() {
		if err := application.Run(); err != nil {
			log.Error("App run failed: %v", err)
		}
	}()

	<-sigCh
	log.Info("Received shutdown signal, shutting down gracefully...")

	application.Shutdown()
	return nil
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

	if theme != "default" {
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
