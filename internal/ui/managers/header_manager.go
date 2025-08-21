package managers

import (
	"context"
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/wikczerski/whaletui/internal/config"
	"github.com/wikczerski/whaletui/internal/ui/constants"
	"github.com/wikczerski/whaletui/internal/ui/interfaces"
)

// HeaderManager manages the header columns and their updates
type HeaderManager struct {
	ui            interfaces.UIInterface
	dockerInfoCol *tview.TextView
	navCol        *tview.TextView
	actionsCol    *tview.TextView
	logoCol       *tview.TextView
	themeManager  *config.ThemeManager
}

// NewHeaderManager creates a new header manager
func NewHeaderManager(ui interfaces.UIInterface) *HeaderManager {
	return &HeaderManager{
		ui:           ui,
		themeManager: ui.GetThemeManager(),
	}
}

// CreateHeaderSection creates the complete header section with all columns
func (hm *HeaderManager) CreateHeaderSection() *tview.Flex {
	headerSection := tview.NewFlex().SetDirection(tview.FlexColumn)

	hm.dockerInfoCol = hm.createColumn(tview.AlignLeft, hm.themeManager.GetTextColor())
	hm.navCol = hm.createColumn(tview.AlignLeft, hm.themeManager.GetTextColor())
	hm.actionsCol = hm.createColumn(tview.AlignLeft, hm.themeManager.GetTextColor())
	hm.logoCol = hm.createLogoColumn()

	headerSection.AddItem(hm.dockerInfoCol, 0, 1, false)
	headerSection.AddItem(hm.navCol, 0, 1, false)
	headerSection.AddItem(hm.actionsCol, 0, 1, false)
	headerSection.AddItem(hm.logoCol, 0, 1, false)

	return headerSection
}

// createColumn creates a standard header column with consistent styling
func (hm *HeaderManager) createColumn(align int, color tcell.Color) *tview.TextView {
	col := tview.NewTextView()
	col.SetBorder(false)
	col.SetBackgroundColor(hm.themeManager.GetBackgroundColor())
	col.SetTextColor(color)
	col.SetTextAlign(align)
	col.SetWordWrap(false)
	col.SetDynamicColors(false)
	return col
}

// createLogoColumn creates the logo column with special styling
func (hm *HeaderManager) createLogoColumn() *tview.TextView {
	logoCol := hm.createColumn(tview.AlignRight, hm.themeManager.GetHeaderColor())
	logoCol.SetText(`  ____  ____
 |  _ \|  _ \
 | | | | | | |
 | |_| | |_| |
 |____/|____/
 `)
	return logoCol
}

// UpdateAll updates all header columns
func (hm *HeaderManager) UpdateAll() {
	// Skip header updates during refresh cycles to prevent empty spaces and newlines
	if hm.ui.IsRefreshing() {
		return
	}

	// Ensure columns are initialized before updating
	if hm.dockerInfoCol == nil || hm.navCol == nil || hm.actionsCol == nil {
		return
	}

	hm.UpdateDockerInfo()
	hm.UpdateNavigation()
	hm.UpdateActions()
}

// UpdateDockerInfo updates the Docker info column
func (hm *HeaderManager) UpdateDockerInfo() {
	if hm.dockerInfoCol == nil {
		return
	}

	services := hm.ui.GetServices()
	if services == nil {
		hm.setDefaultDockerInfo()
		return
	}

	hm.updateDockerInfoFromService(services)
}

// setDefaultDockerInfo sets the default Docker info when services are not available
func (hm *HeaderManager) setDefaultDockerInfo() {
	hm.dockerInfoCol.SetText("Context: docker\nCluster: local\nUser: docker\nwhaletui Rev: dev\nDocker Rev: --\nCPU: --\nMEM: --")
}

// updateDockerInfoFromService updates Docker info from the service
func (hm *HeaderManager) updateDockerInfoFromService(services interfaces.ServiceFactoryInterface) {
	if services == nil {
		hm.setDefaultDockerInfo()
		return
	}

	dockerInfoService := services.GetDockerInfoService()
	if dockerInfoService == nil {
		hm.setDefaultDockerInfo()
		return
	}

	ctx := context.Background()
	dockerInfo, err := dockerInfoService.GetDockerInfo(ctx)
	if err != nil {
		// Show connection info even if Docker info fails
		hm.setDefaultDockerInfo()
		return
	}

	hm.setDockerInfoFromData(&dockerInfo)
}

// setDockerInfoFromData sets Docker info from the retrieved data
func (hm *HeaderManager) setDockerInfoFromData(dockerInfo *interfaces.DockerInfo) {
	// Use the DockerInfoTemplate constant with comprehensive information
	infoText := fmt.Sprintf(constants.DockerInfoTemplate,
		dockerInfo.Version,
		dockerInfo.Containers,
		dockerInfo.Images,
		dockerInfo.Volumes,
		dockerInfo.Networks,
		dockerInfo.OperatingSystem,
		dockerInfo.Architecture,
		dockerInfo.Driver,
		dockerInfo.LoggingDriver,
	)

	hm.dockerInfoCol.SetText(infoText)
}

// UpdateNavigation updates the navigation column based on current mode
func (hm *HeaderManager) UpdateNavigation() {
	if hm.navCol == nil {
		return
	}

	var navText string
	switch {
	case hm.ui.IsInLogsMode():
		navText = "Logs Navigation:\n<up/down> Scroll line\n<pgup/pgdn> Page\n<home/end> Top/Bottom\n<space> Half page\n<esc> Back to table"
	case hm.ui.IsInDetailsMode():
		navText = "Navigation:\n<up/down> Scroll line\n<pgup/pgdn> Page\n<home/end> Top/Bottom\n<space> Half page\n<esc> Back to table"
	default:
		navText = "View Actions:\n<:> Command mode\n<enter> Inspect item\n<l> View logs\n<up/down> Navigate rows\n<q> Quit app\n<ctrl-c> Exit"
	}

	hm.navCol.SetText(navText)
}

// UpdateActions updates the actions column based on current view and mode
func (hm *HeaderManager) UpdateActions() {
	if hm.actionsCol == nil {
		return
	}

	// If in logs mode, show logs actions
	if hm.ui.IsInLogsMode() {
		hm.updateLogsActions()
		return
	}

	// If in details mode, show current actions
	if hm.ui.IsInDetailsMode() {
		hm.updateDetailsActions()
		return
	}

	// Get actions for current view from registry
	hm.updateViewActions()
}

// updateLogsActions updates the actions column with logs-specific actions
func (hm *HeaderManager) updateLogsActions() {
	services := hm.ui.GetServices()
	if !services.IsServiceAvailable("logs") {
		hm.actionsCol.SetText("ESC/Enter: Back to table")
		return
	}

	logsActions := services.GetLogsService().GetActions()

	var actionsText string
	for key, action := range logsActions {
		actionsText += fmt.Sprintf("<%c> %s\n", key, action)
	}
	actionsText += "ESC/Enter: Back to table"
	hm.actionsCol.SetText(actionsText)
}

// updateDetailsActions updates the actions column with details-specific actions
func (hm *HeaderManager) updateDetailsActions() {
	currentActions := hm.ui.GetCurrentActions()
	if currentActions != nil {
		var actionsText string
		for key, action := range currentActions {
			actionsText += fmt.Sprintf("<%c> %s\n", key, action)
		}
		actionsText += "ESC/Enter: Back\n<up/down> Scroll JSON\n<:> Command mode"
		hm.actionsCol.SetText(actionsText)
		return
	}
	// Fallback to default actions
	hm.actionsCol.SetText("ESC/Enter: Back\n<up/down> Scroll JSON\n<:> Command mode")
}

// updateViewActions updates the actions column with view-specific actions
func (hm *HeaderManager) updateViewActions() {
	// Get actions for current view from registry
	viewActions := hm.ui.GetCurrentViewActions()
	if viewActions != "" {
		hm.actionsCol.SetText(viewActions)
		return
	}

	// Fallback to default container actions
	hm.setDefaultContainerActions()
}

// setDefaultContainerActions sets the default container actions
func (hm *HeaderManager) setDefaultContainerActions() {
	services := hm.ui.GetServices()
	if !services.IsContainerServiceAvailable() {
		hm.actionsCol.SetText("")
		return
	}

	containerService := services.GetContainerService()
	if actionService, ok := containerService.(interfaces.ServiceWithActions); ok {
		defaultActions := actionService.GetActionsString()
		hm.actionsCol.SetText(defaultActions)
	} else {
		hm.actionsCol.SetText("")
	}
}

// GetDockerInfoCol returns the docker info column for external access
func (hm *HeaderManager) GetDockerInfoCol() *tview.TextView {
	return hm.dockerInfoCol
}

// GetNavCol returns the navigation column for external access
func (hm *HeaderManager) GetNavCol() *tview.TextView {
	return hm.navCol
}

// GetActionsCol returns the actions column for external access
func (hm *HeaderManager) GetActionsCol() *tview.TextView {
	return hm.actionsCol
}
