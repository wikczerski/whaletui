package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/wikczerski/whaletui/internal/config"
	"github.com/wikczerski/whaletui/internal/docker"
	"github.com/wikczerski/whaletui/internal/logger"
	"github.com/wikczerski/whaletui/internal/services"
	"github.com/wikczerski/whaletui/internal/ui/core"
	"github.com/wikczerski/whaletui/internal/ui/managers"
)

// App represents the main application instance
type App struct {
	cfg      *config.Config
	docker   *docker.Client
	services *services.ServiceFactory
	ui       *core.UI
	ctx      context.Context
	cancel   context.CancelFunc
	log      *slog.Logger
}

// New creates a new application instance
func New(cfg *config.Config) (*App, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	log := logger.GetLogger()

	client, err := docker.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("docker client creation failed: %w", err)
	}

	services := services.NewServiceFactory(client)
	ctx, cancel := context.WithCancel(context.Background())

	ui, err := core.New(services, cfg.Theme, nil, nil) // Managers will be created after UI creation
	if err != nil {
		cancel()
		if client != nil {
			client.Close()
		}
		return nil, fmt.Errorf("UI creation failed: %w", err)
	}

	// Create managers after UI creation
	headerManager := managers.NewHeaderManager(ui)
	modalManager := managers.NewModalManager(ui)

	// Update the managers with the real UI
	ui.SetHeaderManager(headerManager)
	ui.SetModalManager(modalManager)

	// Complete the UI initialization now that managers are set
	if err := ui.CompleteInitialization(); err != nil {
		cancel()
		if client != nil {
			client.Close()
		}
		return nil, fmt.Errorf("UI initialization failed: %w", err)
	}

	return &App{
		cfg:      cfg,
		docker:   client,
		services: services,
		ui:       ui,
		ctx:      ctx,
		cancel:   cancel,
		log:      log,
	}, nil
}

// Run starts the application
func (a *App) Run() error {
	ticker := time.NewTicker(time.Duration(a.cfg.RefreshInterval) * time.Second)
	defer ticker.Stop()

	go a.refreshLoop(ticker)

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	uiDone := make(chan error, 1)
	go func() { uiDone <- a.ui.Start() }()

	select {
	case err := <-uiDone:
		return err
	case <-a.ui.GetShutdownChan():
		return nil
	case sig := <-sigChan:
		a.log.Info("Received signal, shutting down gracefully", "signal", sig)
		a.Shutdown()
		return nil
	}
}

func (a *App) refreshLoop(ticker *time.Ticker) {
	lastRefresh := time.Now()
	minRefreshInterval := 5 * time.Second // Increased minimum time between refreshes to prevent UI issues

	for {
		select {
		case <-a.ctx.Done():
			return
		case <-ticker.C:
			// Prevent excessive refresh calls that might cause terminal issues
			if time.Since(lastRefresh) >= minRefreshInterval {
				a.ui.Refresh()
				lastRefresh = time.Now()
			}
		}
	}
}

// Shutdown gracefully shuts down the application
func (a *App) Shutdown() {
	a.log.Info("Application shutdown started")
	a.cancel()
	a.log.Info("UI stopping")
	a.ui.Stop()
	a.log.Info("UI stopped")

	if a.docker != nil {
		a.log.Info("Closing Docker client")
		if err := a.docker.Close(); err != nil {
			a.log.Error("Failed to close Docker client", "error", err)
		} else {
			a.log.Info("Docker client closed successfully")
		}
	} else {
		a.log.Warn("Docker client is nil during shutdown")
	}

	a.log.Info("Application shutdown completed")
}

// GetUI returns the UI instance
func (a *App) GetUI() *core.UI {
	return a.ui
}
