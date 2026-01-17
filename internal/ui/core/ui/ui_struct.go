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
	services          interfaces.ServiceFactoryInterface
	modalManager      interfaces.ModalManagerInterface
	headerManager     interfaces.HeaderManagerInterface
	swarmNodesView    *swarm.NodesView
	shellView         *shell.View
	statusBar         *tview.TextView
	log               *slog.Logger
	shutdownChan      chan struct{}
	searchInput       *tview.InputField
	commandInput      *tview.InputField
	currentActions    map[rune]string
	themeManager      *config.ThemeManager
	viewRegistry      *core.ViewRegistry
	mainFlex          *tview.Flex
	commandHandler    *handlers.CommandHandler
	searchHandler     *handlers.SearchHandler
	pages             *tview.Pages
	viewContainer     *tview.Flex
	volumesView       *volume.VolumesView
	containersView    *container.ContainersView
	networksView      *network.NetworksView
	logsView          *logs.View
	imagesView        *image.ImagesView
	swarmServicesView *swarm.ServicesView
	app               *tview.Application
	componentBuilder  *builders.ComponentBuilder
	viewBuilder       *builders.ViewBuilder
	tableBuilder      *builders.TableBuilder
	headerSection     *tview.Flex
	isRefreshing      bool
	inLogsMode        bool
	inDetailsMode     bool
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
