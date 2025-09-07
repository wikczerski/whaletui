package ui

import (
	"log/slog"

	"github.com/rivo/tview"
	"github.com/wikczerski/whaletui/internal/config"
	"github.com/wikczerski/whaletui/internal/domains/container"
	"github.com/wikczerski/whaletui/internal/domains/image"
	"github.com/wikczerski/whaletui/internal/domains/logs"
	"github.com/wikczerski/whaletui/internal/domains/network"
	"github.com/wikczerski/whaletui/internal/domains/swarm"
	"github.com/wikczerski/whaletui/internal/domains/volume"
	"github.com/wikczerski/whaletui/internal/logger"
	"github.com/wikczerski/whaletui/internal/ui/builders"
	"github.com/wikczerski/whaletui/internal/ui/core"
	"github.com/wikczerski/whaletui/internal/ui/handlers"
	"github.com/wikczerski/whaletui/internal/ui/interfaces"
	"github.com/wikczerski/whaletui/internal/ui/views/shell"
)

// UI represents the main UI
type UI struct {
	app            *tview.Application
	pages          *tview.Pages
	mainFlex       *tview.Flex
	statusBar      *tview.TextView
	viewContainer  *tview.Flex
	services       interfaces.ServiceFactoryInterface
	log            *slog.Logger
	shutdownChan   chan struct{}
	inDetailsMode  bool            // Track if we're in details view mode
	inLogsMode     bool            // Track if we're viewing container logs
	currentActions map[rune]string // Track current available actions

	// Theme management
	themeManager *config.ThemeManager

	// Abstracted managers
	viewRegistry   *core.ViewRegistry
	headerManager  interfaces.HeaderManagerInterface
	commandHandler *handlers.CommandHandler
	searchHandler  *handlers.SearchHandler
	modalManager   interfaces.ModalManagerInterface

	// Individual views
	containersView *container.ContainersView
	imagesView     *image.ImagesView
	volumesView    *volume.VolumesView
	networksView   *network.NetworksView
	logsView       *logs.View
	shellView      *shell.View

	// Swarm views
	swarmServicesView *swarm.ServicesView
	swarmNodesView    *swarm.NodesView

	// Component builders
	componentBuilder *builders.ComponentBuilder
	viewBuilder      *builders.ViewBuilder
	tableBuilder     *builders.TableBuilder

	// Flags for refresh cycles
	isRefreshing bool

	// UI components
	headerSection *tview.Flex
	commandInput  *tview.InputField
	searchInput   *tview.InputField
}

// New creates a new UI
func New(
	serviceFactory interfaces.ServiceFactoryInterface,
	themePath string,
	headerManager interfaces.HeaderManagerInterface,
	modalManager interfaces.ModalManagerInterface,
	_ *config.Config,
) (*UI, error) {
	// TUI mode is already set globally in init()

	ui := &UI{
		services:       serviceFactory,
		app:            tview.NewApplication(),
		pages:          tview.NewPages(),
		log:            logger.GetLogger(),
		shutdownChan:   make(chan struct{}, 1), // Buffer channel to prevent deadlock
		currentActions: make(map[rune]string),
		headerManager:  headerManager,
		modalManager:   modalManager,
	}

	if e := ui.initializeManagers(themePath); e != nil {
		return nil, e
	}

	// Only initialize components if managers are provided
	if headerManager != nil && modalManager != nil {
		ui.initComponents()
		ui.setupKeyBindings()
	}

	return ui, nil
}
