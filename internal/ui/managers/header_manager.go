package managers

import (
	"context"
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/wikczerski/D5r/internal/config"
	"github.com/wikczerski/D5r/internal/ui/constants"
	"github.com/wikczerski/D5r/internal/ui/interfaces"
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
		hm.dockerInfoCol.SetText("Context: docker\nCluster: local\nUser: docker\nD5r Rev: dev\nDocker Rev: --\nCPU: --\nMEM: --")
		return
	}

	ctx := context.Background()
	dockerInfo, err := services.DockerInfoService.GetDockerInfo(ctx)

	var infoText string
	if err != nil {
		// Show connection info even if Docker info fails
		infoText = "Context: docker\nCluster: local\nUser: docker\nD5r Rev: dev\nDocker Rev: --\nCPU: --\nMEM: --"
	} else {
		// Use the DockerInfoTemplate constant with comprehensive information
		infoText = fmt.Sprintf(constants.DockerInfoTemplate,
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
	}

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
		logsActions := map[rune]string{
			'f': "Follow logs",
			't': "Tail logs",
			's': "Save logs",
			'c': "Clear logs",
			'w': "Wrap text",
		}
		var actionsText string
		for key, action := range logsActions {
			actionsText += fmt.Sprintf("<%c> %s\n", key, action)
		}
		actionsText += "ESC/Enter: Back to table"
		hm.actionsCol.SetText(actionsText)
		return
	}

	// If in details mode, show current actions
	if hm.ui.IsInDetailsMode() {
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
	}

	// Get actions for current view from registry
	viewActions := hm.ui.GetCurrentViewActions()
	if viewActions != "" {
		hm.actionsCol.SetText(viewActions)
		return
	}

	// Fallback to default container actions
	defaultActions := "<s> Start\n<S> Stop\n<r> Restart\n<d> Delete\n<a> Attach\n<l> Logs\n<i> Inspect\n<n> New\n<e> Exec\n<f> Filter\n<t> Sort\n<h> History\n<enter> Details\n<:> Command"
	hm.actionsCol.SetText(defaultActions)
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
