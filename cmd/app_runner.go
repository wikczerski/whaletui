package cmd

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/wikczerski/whaletui/internal/app"
)

// runApplicationWithShutdown starts the application and waits for a shutdown signal
func runApplicationWithShutdown(application *app.App, log *slog.Logger) error {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	uiShutdownCh := application.GetUI().GetShutdownChan()

	go func() {
		if err := application.Run(); err != nil {
			log.Error("App run failed", "error", err)
		}
	}()

	waitForShutdownSignal(sigCh, uiShutdownCh, log)
	cleanupAndShutdown(application)

	return nil
}

// waitForShutdownSignal waits for shutdown signals from the OS or the UI
func waitForShutdownSignal(sigCh chan os.Signal, uiShutdownCh chan struct{}, log *slog.Logger) {
	select {
	case <-sigCh:
		log.Info("Received shutdown signal, shutting down gracefully...")
	case <-uiShutdownCh:
		log.Info("Received UI shutdown signal, shutting down gracefully...")
	}
}
