package app

import (
	"context"
	"fmt"
	"time"

	"github.com/wikczerski/D5r/internal/config"
	"github.com/wikczerski/D5r/internal/docker"
	"github.com/wikczerski/D5r/internal/logger"
	"github.com/wikczerski/D5r/internal/services"
	"github.com/wikczerski/D5r/internal/ui/core"
)

// App represents the main application instance
type App struct {
	cfg      *config.Config
	docker   *docker.Client
	services *services.ServiceFactory
	ui       *core.UI
	ctx      context.Context
	cancel   context.CancelFunc
	log      *logger.Logger
}

// New creates a new application instance
func New(cfg *config.Config) (*App, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	log := logger.GetLogger()
	log.SetPrefix("App")

	client, err := docker.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("docker client creation failed: %w", err)
	}

	services := services.NewServiceFactory(client)
	ctx, cancel := context.WithCancel(context.Background())

	ui, err := core.New(services, cfg.Theme)
	if err != nil {
		cancel()
		if client != nil {
			client.Close()
		}
		return nil, fmt.Errorf("UI creation failed: %w", err)
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

	uiDone := make(chan error, 1)
	go func() { uiDone <- a.ui.Start() }()

	select {
	case err := <-uiDone:
		return err
	case <-a.ui.GetShutdownChan():
		return nil
	}
}

func (a *App) refreshLoop(ticker *time.Ticker) {
	for {
		select {
		case <-a.ctx.Done():
			return
		case <-ticker.C:
			a.ui.Refresh()
		}
	}
}

// Shutdown gracefully shuts down the application
func (a *App) Shutdown() {
	a.cancel()
	a.ui.Stop()
	if a.docker != nil {
		if err := a.docker.Close(); err != nil {
			a.log.Error("Failed to close Docker client: %v", err)
		}
	}
}

// GetUI returns the UI instance
func (a *App) GetUI() *core.UI {
	return a.ui
}
