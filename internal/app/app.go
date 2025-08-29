// Package app provides the main application logic and coordination for WhaleTUI.
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
	"github.com/wikczerski/whaletui/internal/errors"
	"github.com/wikczerski/whaletui/internal/logger"
	"github.com/wikczerski/whaletui/internal/ui/core"
	"github.com/wikczerski/whaletui/internal/ui/core/ui"
	"github.com/wikczerski/whaletui/internal/ui/managers"
)

// App represents the main application instance
type App struct {
	cfg      *config.Config
	docker   *docker.Client
	services *core.ServiceFactory
	ui       *ui.UI
	ctx      context.Context
	cancel   context.CancelFunc
	log      *slog.Logger
}

// New creates a new application instance
func New(cfg *config.Config) (*App, error) {
	if cfg == nil {
		return nil, errors.NewConfigError("config cannot be nil")
	}

	log := logger.GetLogger()
	client, services, err := createDockerClientAndServices(cfg)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	ui, err := createAndSetupUI(services, cfg, client, cancel)
	if err != nil {
		return nil, err
	}

	log.Info("UI created", "ui", ui != nil)

	return createAppInstance(cfg, client, services, ui, ctx, cancel, log), nil
}

// createDockerClientAndServices creates a new Docker client and service factory
func createDockerClientAndServices(
	cfg *config.Config,
) (*docker.Client, *core.ServiceFactory, error) {
	client, err := docker.New(cfg)
	if err != nil {
		return nil, nil, errors.NewDockerError("docker client creation", err)
	}

	services := core.NewServiceFactory(client)
	logger.GetLogger().Info("Service factory created", "services", services != nil)

	return client, services, nil
}

// createAndSetupUI creates and sets up the UI
func createAndSetupUI(
	services *core.ServiceFactory,
	cfg *config.Config,
	client *docker.Client,
	cancel context.CancelFunc,
) (*ui.UI, error) {
	ui, err := createUI(services, cfg)
	if err != nil {
		cleanupOnError(client, cancel)
		return nil, err
	}

	setupManagers(ui)

	if err := ui.CompleteInitialization(); err != nil {
		cleanupOnError(client, cancel)
		return nil, errors.UIError("UI initialization", err)
	}

	return ui, nil
}

// createAppInstance creates the App instance
func createAppInstance(
	cfg *config.Config,
	client *docker.Client,
	services *core.ServiceFactory,
	ui *ui.UI,
	ctx context.Context,
	cancel context.CancelFunc,
	log *slog.Logger,
) *App {
	return &App{
		cfg:      cfg,
		docker:   client,
		services: services,
		ui:       ui,
		ctx:      ctx,
		cancel:   cancel,
		log:      log,
	}
}

func createUI(services *core.ServiceFactory, cfg *config.Config) (*ui.UI, error) {
	uiInstance, err := ui.New(services, cfg.Theme, nil, nil, cfg)
	if err != nil {
		return nil, errors.UIError("UI creation", err)
	}
	return uiInstance, nil
}

func setupManagers(uiInstance *ui.UI) {
	headerManager := managers.NewHeaderManager(uiInstance)
	modalManager := managers.NewModalManager(uiInstance)
	uiInstance.SetHeaderManager(headerManager)
	uiInstance.SetModalManager(modalManager)
}

func cleanupOnError(client *docker.Client, cancel context.CancelFunc) {
	cancel()
	if client != nil {
		if err := client.Close(); err != nil {
			// Log the error but continue since this is cleanup
			fmt.Fprintf(os.Stderr, "Failed to close Docker client: %v\n", err)
		}
	}
}

// Run starts the application
func (a *App) Run() error {
	ticker := time.NewTicker(time.Duration(a.cfg.RefreshInterval) * time.Second)
	defer ticker.Stop()

	go a.refreshLoop(ticker)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	uiDone := make(chan error, 1)
	go func() { uiDone <- a.ui.Start() }()

	return a.waitForShutdown(uiDone, sigChan)
}

// Shutdown gracefully shuts down the application
func (a *App) Shutdown() {
	a.log.Info("Application shutdown started")
	a.cancel()
	a.stopUI()
	a.closeDockerClient()
	a.log.Info("Application shutdown completed")
}

// GetUI returns the UI instance
func (a *App) GetUI() *ui.UI {
	return a.ui
}

func (a *App) waitForShutdown(uiDone chan error, sigChan chan os.Signal) error {
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
	const minRefreshInterval = 5 * time.Second
	lastRefresh := time.Now()

	for {
		select {
		case <-a.ctx.Done():
			return
		case <-ticker.C:
			if time.Since(lastRefresh) >= minRefreshInterval {
				a.ui.Refresh()
				lastRefresh = time.Now()
			}
		}
	}
}

func (a *App) stopUI() {
	a.log.Info("UI stopping")
	a.ui.Stop()
	a.log.Info("UI stopped")
}

func (a *App) closeDockerClient() {
	if a.docker == nil {
		a.log.Warn("Docker client is nil during shutdown")
		return
	}

	a.log.Info("Closing Docker client")
	if err := a.docker.Close(); err != nil {
		a.log.Error("Failed to close Docker client", "error", err)
		return
	}
	a.log.Info("Docker client closed successfully")
}
