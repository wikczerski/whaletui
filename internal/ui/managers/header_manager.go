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
	return &HeaderManager{
		ui: ui,
	}
}

// CreateHeaderSection creates the header section with Docker info, navigation, and actions
func (hm *HeaderManager) CreateHeaderSection() tview.Primitive {
	headerFlex := hm.createHeaderFlex()
	hm.addHeaderViews(headerFlex)
	hm.headerFlex = headerFlex
	return headerFlex
}

// UpdateAll updates all header content
func (hm *HeaderManager) UpdateAll() {
	if hm.headerFlex == nil {
		return
	}

	contentLines := hm.getUpdatedContentLines()
	hm.updateHeaderViews(contentLines)
}

// UpdateDockerInfo updates the Docker info (kept for backward compatibility)
func (hm *HeaderManager) UpdateDockerInfo() {
	hm.UpdateAll()
}

// UpdateNavigation updates the navigation (kept for backward compatibility)
func (hm *HeaderManager) UpdateNavigation() {
	hm.UpdateAll()
}

// UpdateActions updates the actions (kept for backward compatibility)
func (hm *HeaderManager) UpdateActions() {
	hm.UpdateAll()
}

// GetDockerInfoCol returns nil (kept for backward compatibility)
func (hm *HeaderManager) GetDockerInfoCol() *tview.TextView {
	return nil
}

// GetNavCol returns nil (kept for backward compatibility)
func (hm *HeaderManager) GetNavCol() *tview.TextView {
	return nil
}

// GetActionsCol returns nil (kept for backward compatibility)
func (hm *HeaderManager) GetActionsCol() *tview.TextView {
	return nil
}

// createHeaderFlex creates the main header flex container
func (hm *HeaderManager) createHeaderFlex() *tview.Flex {
	headerFlex := tview.NewFlex()
	headerFlex.SetDirection(tview.FlexColumn)
	return headerFlex
}

// addHeaderViews adds all the header views to the flex container
func (hm *HeaderManager) addHeaderViews(headerFlex *tview.Flex) {
	views := hm.createHeaderViews()
	hm.addViewsToFlex(headerFlex, views)
}

// createHeaderViews creates all the header view components
func (hm *HeaderManager) createHeaderViews() map[string]tview.Primitive {
	dockerLines := hm.getPaddedLines(hm.getDockerInfoText())
	navigationLines := hm.getPaddedLines(hm.getNavigationText())
	actionsLines := hm.getPaddedLines(hm.getActionsText())
	logoLines := hm.getPaddedLines(hm.getLogoText())

	return map[string]tview.Primitive{
		"docker":     hm.createEnhancedTextView(dockerLines, tview.AlignLeft, "Docker Info"),
		"navigation": hm.createEnhancedTextView(navigationLines, tview.AlignLeft, "Navigation"),
		"actions":    hm.createEnhancedTextView(actionsLines, tview.AlignLeft, "Actions"),
		"logo":       hm.createEnhancedTextView(logoLines, tview.AlignCenter, "WhaleTui"),
	}
}

// getPaddedLines adds padding to make header sections taller
func (hm *HeaderManager) getPaddedLines(content string) []string {
	lines := strings.Split(content, "\n")
	return append(lines, "", "", "")
}

// addViewsToFlex adds views to the flex container with specified proportions
func (hm *HeaderManager) addViewsToFlex(headerFlex *tview.Flex, views map[string]tview.Primitive) {
	headerFlex.AddItem(views["docker"], 0, 1, false)     // Docker Info: 1 part
	headerFlex.AddItem(views["navigation"], 0, 1, false) // Navigation: 1 part
	headerFlex.AddItem(views["actions"], 0, 3, false)    // Actions: 3 parts
	headerFlex.AddItem(views["logo"], 0, 2, false)       // Logo: 2 parts
}

// getUpdatedContentLines gets all the updated content lines
func (hm *HeaderManager) getUpdatedContentLines() map[string][]string {
	return map[string][]string{
		"docker":     strings.Split(hm.getDockerInfoText(), "\n"),
		"navigation": strings.Split(hm.getNavigationText(), "\n"),
		"actions":    strings.Split(hm.getActionsText(), "\n"),
		"logo":       strings.Split(hm.getLogoText(), "\n"),
	}
}

// updateHeaderViews updates all the header views with new content
func (hm *HeaderManager) updateHeaderViews(contentLines map[string][]string) {
	if hm.headerFlex.GetItemCount() < 4 {
		return
	}

	hm.updateView(0, contentLines["docker"])
	hm.updateView(1, contentLines["navigation"])
	hm.updateView(2, contentLines["actions"])
	hm.updateView(3, contentLines["logo"])
}

// updateView updates a specific view with new content
func (hm *HeaderManager) updateView(index int, lines []string) {
	item := hm.headerFlex.GetItem(index)
	switch view := item.(type) {
	case *tview.TextView:
		view.SetText(strings.Join(lines, "\n"))
	case *tview.Grid:
		// For now, just update text if it's a grid (simplified approach)
	}
}

// createEnhancedTextView creates an enhanced text view with border and title
func (hm *HeaderManager) createEnhancedTextView(
	lines []string,
	align int,
	title string,
) tview.Primitive {
	if hm.shouldUseTableLayout(title, lines) {
		return hm.createTableLayout(lines, align, title)
	}

	return hm.createSimpleTextView(lines, align, title)
}

// shouldUseTableLayout determines if table layout should be used
func (hm *HeaderManager) shouldUseTableLayout(title string, lines []string) bool {
	// Docker Info and Logo should never use table layout
	if title == "Docker Info" || title == "WhaleTui" {
		return false
	}

	// For Navigation and Actions, use table layout if content is long
	return len(lines) > 7
}

// createSimpleTextView creates a simple text view with styling
func (hm *HeaderManager) createSimpleTextView(
	lines []string,
	align int,
	title string,
) tview.Primitive {
	textView := tview.NewTextView()
	hm.setupTextViewStyling(textView, lines, align, title)
	return textView
}

// setupTextViewStyling sets up the styling for a text view
func (hm *HeaderManager) setupTextViewStyling(
	textView *tview.TextView,
	lines []string,
	align int,
	title string,
) {
	hm.setupTextViewContent(textView, lines, title)
	hm.setupTextViewColors(textView)
	hm.setupTextViewBehavior(textView, align)
}

// setupTextViewContent sets up the content and title for a text view
func (hm *HeaderManager) setupTextViewContent(
	textView *tview.TextView,
	lines []string,
	title string,
) {
	textView.SetText(strings.Join(lines, "\n"))
	textView.SetTitle(fmt.Sprintf(" %s ", title))
}

// setupTextViewColors sets up the colors for a text view
func (hm *HeaderManager) setupTextViewColors(textView *tview.TextView) {
	themeManager := hm.ui.GetThemeManager()
	textView.SetTextColor(themeManager.GetTextColor())
	textView.SetBackgroundColor(themeManager.GetBackgroundColor())
	textView.SetBorderColor(themeManager.GetBorderColor())
	textView.SetTitleColor(themeManager.GetHeaderColor())
}

// setupTextViewBehavior sets up the behavior and layout for a text view
func (hm *HeaderManager) setupTextViewBehavior(textView *tview.TextView, align int) {
	textView.SetTextAlign(align)
	textView.SetBorder(true)
	textView.SetScrollable(false)
	textView.SetDynamicColors(false)
	textView.SetWordWrap(false)
	textView.SetWrap(false)
}

// createTableLayout creates a table-like layout for content longer than 7 lines
func (hm *HeaderManager) createTableLayout(
	lines []string,
	align int,
	title string,
) tview.Primitive {
	maxRows := 7
	numColumns := hm.calculateColumnCount(len(lines), maxRows)

	flex := hm.createTableFlex(title)
	hm.addTableColumns(flex, lines, align, maxRows, numColumns)

	return flex
}

// calculateColumnCount calculates the number of columns needed
func (hm *HeaderManager) calculateColumnCount(linesCount, maxRows int) int {
	return (linesCount + maxRows - 1) / maxRows // Ceiling division
}

// createTableFlex creates the flex container for table layout
func (hm *HeaderManager) createTableFlex(title string) *tview.Flex {
	flex := tview.NewFlex()
	flex.SetDirection(tview.FlexColumn)
	flex.SetBorder(true)
	flex.SetTitle(fmt.Sprintf(" %s ", title))
	flex.SetBorderColor(hm.ui.GetThemeManager().GetBorderColor())
	flex.SetBackgroundColor(hm.ui.GetThemeManager().GetBackgroundColor())
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
	hm.setupTableColumnStyling(columnText, align)
	hm.setTableColumnContent(columnText, lines, maxRows, col)
	return columnText
}

// setupTableColumnStyling sets up the styling for a table column
func (hm *HeaderManager) setupTableColumnStyling(columnText *tview.TextView, align int) {
	columnText.SetTextColor(hm.ui.GetThemeManager().GetTextColor())
	columnText.SetBackgroundColor(hm.ui.GetThemeManager().GetBackgroundColor())
	columnText.SetTextAlign(align)
	columnText.SetBorder(false)
	columnText.SetScrollable(false)
	columnText.SetDynamicColors(false)
	columnText.SetWordWrap(false)
	columnText.SetWrap(false)
}

// setTableColumnContent sets the content for a table column
func (hm *HeaderManager) setTableColumnContent(
	columnText *tview.TextView,
	lines []string,
	maxRows, col int,
) {
	var columnLines []string
	for i := col * maxRows; i < (col+1)*maxRows && i < len(lines); i++ {
		columnLines = append(columnLines, lines[i])
	}
	columnText.SetText(strings.Join(columnLines, "\n"))
}

// getDockerInfoText returns the formatted Docker info text
func (hm *HeaderManager) getDockerInfoText() string {
	// Try to get Docker info from service first
	services := hm.ui.GetServices()
	if services != nil && services.IsContainerServiceAvailable() {
		return hm.getDockerInfoFromService(services)
	}

	// Fallback to default info
	return hm.getDefaultDockerInfo()
}

// getDefaultDockerInfo returns default Docker info when service is not available
func (hm *HeaderManager) getDefaultDockerInfo() string {
	return fmt.Sprintf(constants.DockerInfoTemplate,
		"❌ Disconnected",
		"--",
		"--",
		"--",
		constants.AppVersion)
}

// getDockerInfoFromService gets Docker info from the service
func (hm *HeaderManager) getDockerInfoFromService(
	services interfaces.ServiceFactoryInterface,
) string {
	if !services.IsContainerServiceAvailable() {
		return hm.getDefaultDockerInfo()
	}

	if dockerInfo := hm.tryGetDockerInfo(services); dockerInfo != nil {
		return hm.formatDockerInfo(*dockerInfo)
	}

	return hm.getFallbackDockerInfo(services)
}

// tryGetDockerInfo attempts to get Docker info from the service
func (hm *HeaderManager) tryGetDockerInfo(
	services interfaces.ServiceFactoryInterface,
) *interfaces.DockerInfo {
	dockerInfoService := services.GetDockerInfoService()
	if dockerInfoService == nil {
		return nil
	}

	ctx := context.Background()
	dockerInfo, err := dockerInfoService.GetDockerInfo(ctx)
	if err != nil {
		return nil
	}

	return dockerInfo
}

// getFallbackDockerInfo gets fallback Docker info when service fetch fails
func (hm *HeaderManager) getFallbackDockerInfo(services interfaces.ServiceFactoryInterface) string {
	connectionStatus := hm.getConnectionStatus(services)
	return fmt.Sprintf(constants.DockerInfoTemplate,
		connectionStatus,
		"Available",
		"--",
		"--",
		constants.AppVersion)
}

// getConnectionStatus determines the connection status for fallback info
func (hm *HeaderManager) getConnectionStatus(services interfaces.ServiceFactoryInterface) string {
	if services.IsContainerServiceAvailable() {
		return "✅ Connected"
	}
	return "⚠️ Partial"
}

// formatDockerInfo formats Docker info data into displayable text
func (hm *HeaderManager) formatDockerInfo(dockerInfo interfaces.DockerInfo) string {
	if dockerInfo == nil {
		return hm.getDefaultDockerInfo()
	}

	// Dynamic connection status based on Docker info availability
	connectionStatus := "✅ Connected"
	if dockerInfo.GetVersion() == "" {
		connectionStatus = "❌ Disconnected"
	}

	// Use the template from constants
	return fmt.Sprintf(constants.DockerInfoTemplate,
		connectionStatus,
		dockerInfo.GetVersion(),
		dockerInfo.GetOperatingSystem(),
		dockerInfo.GetLoggingDriver(),
		constants.AppVersion)
}

// getNavigationText returns the navigation text based on current context
func (hm *HeaderManager) getNavigationText() string {
	// Get navigation for current view from registry
	viewNavigation := hm.ui.GetCurrentViewNavigation()
	if viewNavigation != "" {
		return viewNavigation
	}

	// Get navigation dynamically from the current view's service
	return hm.getDynamicViewNavigation()
}

// getActionsText returns the actions text based on current view and mode
func (hm *HeaderManager) getActionsText() string {
	// If in logs mode, show logs actions
	if hm.ui.IsInLogsMode() {
		return hm.getLogsActionsText()
	}

	// If in details mode, show current actions
	if hm.ui.IsInDetailsMode() {
		return hm.getDetailsActionsText()
	}

	// Get actions for current view from registry
	viewActions := hm.ui.GetCurrentViewActions()
	if viewActions != "" {
		return viewActions
	}

	// Get actions dynamically from the current view's service
	return hm.getDynamicViewActions()
}

// getLogsActionsText returns the logs-specific actions text
func (hm *HeaderManager) getLogsActionsText() string {
	services := hm.ui.GetServices()
	if !services.IsServiceAvailable("logs") {
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

// getDetailsActionsText returns the details-specific actions text
func (hm *HeaderManager) getDetailsActionsText() string {
	currentActions := hm.ui.GetCurrentActions()
	if currentActions != nil {
		var actionsText string
		for key, action := range currentActions {
			actionsText += fmt.Sprintf("<%c> %s\n", key, action)
		}
		actionsText += "ESC/Enter: Back\n↑/↓: Scroll JSON\n<:> Command mode"
		return actionsText
	}
	// Fallback to default actions
	return "ESC/Enter: Back\n↑/↓: Scroll JSON\n<:> Command mode"
}

// getDynamicViewActions returns actions from the current view's service
func (hm *HeaderManager) getDynamicViewActions() string {
	services := hm.ui.GetServices()
	if services == nil {
		return "No services available"
	}

	// Get the currently active service
	if currentService := services.GetCurrentService(); currentService != nil {
		if actionService, ok := currentService.(interfaces.ServiceWithActions); ok {
			return actionService.GetActionsString()
		}
	}

	return "No actions available"
}

// getDynamicViewNavigation returns navigation from the current view's service
func (hm *HeaderManager) getDynamicViewNavigation() string {
	// Try to get navigation from the current view's service
	services := hm.ui.GetServices()
	if services == nil {
		return "No services available"
	}

	// Get the currently active service
	if currentService := services.GetCurrentService(); currentService != nil {
		if navigationService, ok := currentService.(interfaces.ServiceWithNavigation); ok {
			return navigationService.GetNavigationString()
		}
	}

	return "No actions available"
}

// getLogoText returns the logo text
func (hm *HeaderManager) getLogoText() string {
	return constants.WhaleTuiLogo
}
