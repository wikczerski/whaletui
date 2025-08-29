// Package managers provides UI management components for WhaleTUI.
package managers

import (
	"context"
	"fmt"
	"strings"

	"github.com/rivo/tview"
	"github.com/wikczerski/whaletui/internal/ui/constants"
	"github.com/wikczerski/whaletui/internal/ui/interfaces"
)

// HeaderManager manages the header section of the UI
type HeaderManager struct {
	ui         interfaces.UIInterface
	headerFlex *tview.Flex
}

// NewHeaderManager creates a new header manager
func NewHeaderManager(ui interfaces.UIInterface) *HeaderManager {
	return &HeaderManager{ui: ui}
}

// CreateHeaderSection creates the header section with Docker info, navigation, and actions
func (hm *HeaderManager) CreateHeaderSection() tview.Primitive {
	headerFlex := tview.NewFlex()
	headerFlex.SetDirection(tview.FlexColumn)

	// Add header views with proportions
	dockerInfoView := hm.createHeaderView("Docker Info", hm.getDockerInfoText(), tview.AlignLeft)
	navigationView := hm.createHeaderView("Navigation", hm.getNavigationText(), tview.AlignLeft)
	actionsView := hm.createHeaderView("Actions", hm.getActionsText(), tview.AlignLeft)
	logoView := hm.createHeaderView("WhaleTui", constants.WhaleTuiLogo, tview.AlignCenter)

	headerFlex.AddItem(dockerInfoView, 0, 1, false)
	headerFlex.AddItem(navigationView, 0, 1, false)
	headerFlex.AddItem(actionsView, 0, 3, false)
	headerFlex.AddItem(logoView, 0, 2, false)

	hm.headerFlex = headerFlex
	return headerFlex
}

// UpdateAll updates all header content
func (hm *HeaderManager) UpdateAll() {
	if hm.headerFlex == nil || hm.headerFlex.GetItemCount() < 4 {
		return
	}

	// Update each view directly
	hm.updateView(0, hm.getDockerInfoText())
	hm.updateView(1, hm.getNavigationText())
	hm.updateView(2, hm.getActionsText())
	hm.updateView(3, constants.WhaleTuiLogo)
}

// UpdateDockerInfo updates the Docker info
func (hm *HeaderManager) UpdateDockerInfo() { hm.UpdateAll() }

// UpdateNavigation updates the navigation
func (hm *HeaderManager) UpdateNavigation() { hm.UpdateAll() }

// UpdateActions updates the actions
func (hm *HeaderManager) UpdateActions() { hm.UpdateAll() }

// createHeaderView creates a header view with appropriate layout
func (hm *HeaderManager) createHeaderView(title, content string, align int) tview.Primitive {
	lines := strings.Split(content, "\n")

	// Use table layout for long content (more than 7 lines)
	if len(lines) > 7 && title != "Docker Info" && title != "WhaleTui" {
		return hm.createTableLayout(lines, align, title)
	}

	return hm.createSimpleTextView(lines, align, title)
}

// createSimpleTextView creates a simple text view with styling
func (hm *HeaderManager) createSimpleTextView(
	lines []string,
	align int,
	title string,
) tview.Primitive {
	textView := tview.NewTextView()
	textView.SetText(strings.Join(lines, "\n"))
	textView.SetTitle(fmt.Sprintf(" %s ", title))
	textView.SetTextAlign(align)
	textView.SetBorder(true)
	textView.SetScrollable(false)
	textView.SetDynamicColors(false)
	textView.SetWordWrap(false)
	textView.SetWrap(false)

	// Apply theme colors
	theme := hm.ui.GetThemeManager()
	textView.SetTextColor(theme.GetTextColor())
	textView.SetBackgroundColor(theme.GetBackgroundColor())
	textView.SetBorderColor(theme.GetBorderColor())
	textView.SetTitleColor(theme.GetHeaderColor())

	return textView
}

// createTableLayout creates a table layout for long content
func (hm *HeaderManager) createTableLayout(
	lines []string,
	align int,
	title string,
) tview.Primitive {
	maxRows := 7
	numColumns := (len(lines) + maxRows - 1) / maxRows

	// Limit to max 3 columns for readability
	if numColumns > 3 {
		numColumns = 3
	}

	flex := tview.NewFlex()
	flex.SetDirection(tview.FlexColumn)
	flex.SetBorder(true)
	flex.SetTitle(fmt.Sprintf(" %s ", title))

	theme := hm.ui.GetThemeManager()
	flex.SetBorderColor(theme.GetBorderColor())
	flex.SetBackgroundColor(theme.GetBackgroundColor())

	// Add columns
	hm.addTableColumns(flex, lines, align, maxRows, numColumns)

	return flex
}

// addTableColumns adds columns to the table flex container
func (hm *HeaderManager) addTableColumns(
	flex *tview.Flex,
	lines []string,
	align, maxRows, numColumns int,
) {
	for col := 0; col < numColumns; col++ {
		columnText := hm.createTableColumn(lines, align, maxRows, col)
		flex.AddItem(columnText, 0, 1, false)
	}
}

// createTableColumn creates a single table column
func (hm *HeaderManager) createTableColumn(
	lines []string,
	align, maxRows, col int,
) *tview.TextView {
	columnText := tview.NewTextView()
	theme := hm.ui.GetThemeManager()

	columnText.SetTextColor(theme.GetTextColor())
	columnText.SetBackgroundColor(theme.GetBackgroundColor())
	columnText.SetTextAlign(align)
	columnText.SetBorder(false)
	columnText.SetScrollable(false)
	columnText.SetDynamicColors(false)
	columnText.SetWordWrap(false)
	columnText.SetWrap(false)

	// Set column content
	startIndex := col * maxRows
	endIndex := startIndex + maxRows
	if endIndex > len(lines) {
		endIndex = len(lines)
	}

	var columnLines []string
	for i := startIndex; i < endIndex; i++ {
		columnLines = append(columnLines, lines[i])
	}
	columnText.SetText(strings.Join(columnLines, "\n"))

	return columnText
}

// updateView updates a specific view with new content
func (hm *HeaderManager) updateView(index int, content string) {
	item := hm.headerFlex.GetItem(index)
	if textView, ok := item.(*tview.TextView); ok {
		textView.SetText(content)
	}
}

// getDockerInfoText returns the formatted Docker info text
func (hm *HeaderManager) getDockerInfoText() string {
	services := hm.ui.GetServices()
	if services == nil || !services.IsContainerServiceAvailable() {
		return fmt.Sprintf(constants.DockerInfoTemplate,
			"❌ Disconnected", "--", "--", "--", constants.AppVersion)
	}

	dockerInfoService := services.GetDockerInfoService()
	if dockerInfoService == nil {
		return fmt.Sprintf(constants.DockerInfoTemplate,
			"⚠️ Partial", "Available", "--", "--", constants.AppVersion)
	}

	ctx := context.Background()
	dockerInfoPtr, err := dockerInfoService.GetDockerInfo(ctx)
	if err != nil || dockerInfoPtr == nil {
		return fmt.Sprintf(constants.DockerInfoTemplate,
			"⚠️ Partial", "Available", "--", "--", constants.AppVersion)
	}

	// Dereference the pointer to interface
	dockerInfo := *dockerInfoPtr

	connectionStatus := "✅ Connected"
	if dockerInfo.GetVersion() == "" {
		connectionStatus = "❌ Disconnected"
	}

	return fmt.Sprintf(constants.DockerInfoTemplate,
		connectionStatus,
		dockerInfo.GetVersion(),
		dockerInfo.GetOperatingSystem(),
		dockerInfo.GetLoggingDriver(),
		constants.AppVersion)
}

// getNavigationText returns the navigation text
func (hm *HeaderManager) getNavigationText() string {
	if navigation := hm.getDynamicViewNavigation(); navigation != "No navigation available" {
		return navigation
	}

	if viewNavigation := hm.ui.GetCurrentViewNavigation(); viewNavigation != "" {
		return viewNavigation
	}

	return "No navigation available"
}

// getActionsText returns the actions text
func (hm *HeaderManager) getActionsText() string {
	// Handle special modes
	if hm.ui.IsInLogsMode() {
		return hm.getLogsActionsText()
	}

	if hm.ui.IsInDetailsMode() {
		return hm.getDetailsActionsText()
	}

	// Get dynamic actions from current view service
	if actions := hm.getDynamicViewActions(); actions != "No actions available" {
		return actions
	}

	// Fallback to view registry
	if viewActions := hm.ui.GetCurrentViewActions(); viewActions != "" {
		return viewActions
	}

	return "No actions available"
}

// getLogsActionsText returns logs-specific actions
func (hm *HeaderManager) getLogsActionsText() string {
	services := hm.ui.GetServices()
	if services == nil || !services.IsServiceAvailable("logs") {
		return "ESC/Enter: Back to table"
	}

	logsActions := services.GetLogsService().GetActions()
	var actionsText string
	for key, action := range logsActions {
		actionsText += fmt.Sprintf("<%c> %s\n", key, action)
	}
	actionsText += "ESC/Enter: Back to table"
	return actionsText
}

// getDetailsActionsText returns details-specific actions
func (hm *HeaderManager) getDetailsActionsText() string {
	currentActions := hm.ui.GetCurrentActions()
	if currentActions == nil {
		return "ESC/Enter: Back\n↑/↓: Scroll JSON\n<:> Command mode"
	}

	var actionsText string
	for key, action := range currentActions {
		actionsText += fmt.Sprintf("<%c> %s\n", key, action)
	}
	actionsText += "ESC/Enter: Back\n↑/↓: Scroll JSON\n<:> Command mode"
	return actionsText
}

// getDynamicViewActions returns actions from current view service
func (hm *HeaderManager) getDynamicViewActions() string {
	services := hm.ui.GetServices()
	viewRegistry := hm.ui.GetViewRegistry()
	if services == nil || viewRegistry == nil {
		return "No actions available"
	}

	currentViewName := hm.getCurrentViewName(viewRegistry)
	if currentViewName == "" {
		return "No actions available"
	}

	service := hm.getServiceForView(services, currentViewName)
	if service == nil {
		return "No actions available"
	}

	if actionService, ok := service.(interfaces.ServiceWithActions); ok {
		return actionService.GetActionsString()
	}

	return "No actions available"
}

// getDynamicViewNavigation returns navigation from current view service
func (hm *HeaderManager) getDynamicViewNavigation() string {
	services := hm.ui.GetServices()
	viewRegistry := hm.ui.GetViewRegistry()
	if services == nil || viewRegistry == nil {
		return "No navigation available"
	}

	currentViewName := hm.getCurrentViewName(viewRegistry)
	if currentViewName == "" {
		return "No navigation available"
	}

	service := hm.getServiceForView(services, currentViewName)
	if service == nil {
		return "No navigation available"
	}

	if navigationService, ok := service.(interfaces.ServiceWithNavigation); ok {
		return navigationService.GetNavigationString()
	}

	return "No navigation available"
}

// getCurrentViewName extracts the current view name from the registry
func (hm *HeaderManager) getCurrentViewName(viewRegistry any) string {
	registry, ok := viewRegistry.(interface{ GetCurrentName() string })
	if !ok {
		return ""
	}
	return registry.GetCurrentName()
}

// getServiceForView gets the appropriate service for the given view name
func (hm *HeaderManager) getServiceForView(
	services interfaces.ServiceFactoryInterface,
	viewName string,
) any {
	switch viewName {
	case "containers":
		return services.GetContainerService()
	case "images":
		return services.GetImageService()
	case "volumes":
		return services.GetVolumeService()
	case "networks":
		return services.GetNetworkService()
	default:
		return nil
	}
}
