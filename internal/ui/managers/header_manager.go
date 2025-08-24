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
	// Create a simple flex layout that spans the full width
	headerFlex := tview.NewFlex()
	headerFlex.SetDirection(tview.FlexColumn)

	// Get Docker info
	dockerInfo := hm.getDockerInfoText()
	dockerLines := strings.Split(dockerInfo, "\n")
	// Add padding to make header taller
	dockerLines = append(dockerLines, "", "", "")

	// Get navigation info
	navigationInfo := hm.getNavigationText()
	navigationLines := strings.Split(navigationInfo, "\n")
	// Add padding to make header taller
	navigationLines = append(navigationLines, "", "", "")

	// Get actions info
	actionsInfo := hm.getActionsText()
	actionsLines := strings.Split(actionsInfo, "\n")
	// Add padding to make header taller
	actionsLines = append(actionsLines, "", "", "")

	// Get logo
	logoInfo := hm.getLogoText()
	logoLines := strings.Split(logoInfo, "\n")
	// Add padding to make header taller
	logoLines = append(logoLines, "", "", "")

	// Create enhanced text views for each section with more height
	dockerView := hm.createEnhancedTextView(dockerLines, tview.AlignLeft, "Docker Info")
	navigationView := hm.createEnhancedTextView(navigationLines, tview.AlignLeft, "Navigation")
	actionsView := hm.createEnhancedTextView(actionsLines, tview.AlignLeft, "Actions")
	logoView := hm.createEnhancedTextView(logoLines, tview.AlignCenter, "WhaleTui")

	// Add views to flex with specified proportions
	headerFlex.AddItem(dockerView, 0, 1, false)     // Docker Info: 1 part
	headerFlex.AddItem(navigationView, 0, 1, false) // Navigation: 1 part
	headerFlex.AddItem(actionsView, 0, 3, false)    // Actions: 3 parts
	headerFlex.AddItem(logoView, 0, 2, false)       // Logo: 2 parts

	// Store reference for updates
	hm.headerFlex = headerFlex

	return headerFlex
}

// UpdateAll updates all header content
func (hm *HeaderManager) UpdateAll() {
	if hm.headerFlex == nil {
		return
	}

	// Get updated content
	dockerInfo := hm.getDockerInfoText()
	dockerLines := strings.Split(dockerInfo, "\n")

	navigationInfo := hm.getNavigationText()
	navigationLines := strings.Split(navigationInfo, "\n")

	actionsInfo := hm.getActionsText()
	actionsLines := strings.Split(actionsInfo, "\n")

	logoInfo := hm.getLogoText()
	logoLines := strings.Split(logoInfo, "\n")

	// Update each view - handle both TextView and Grid types
	if hm.headerFlex.GetItemCount() >= 4 {
		// Update Docker view
		switch item := hm.headerFlex.GetItem(0).(type) {
		case *tview.TextView:
			item.SetText(strings.Join(dockerLines, "\n"))
		case *tview.Grid:
			// For now, just update text if it's a grid (simplified approach)
			// TODO: Implement proper grid recreation
		}

		// Update Navigation view
		switch item := hm.headerFlex.GetItem(1).(type) {
		case *tview.TextView:
			item.SetText(strings.Join(navigationLines, "\n"))
		case *tview.Grid:
			// For now, just update text if it's a grid (simplified approach)
			// TODO: Implement proper grid recreation
		}

		// Update Actions view
		switch item := hm.headerFlex.GetItem(2).(type) {
		case *tview.TextView:
			item.SetText(strings.Join(actionsLines, "\n"))
		case *tview.Grid:
			// For now, just update text if it's a grid (simplified approach)
			// TODO: Implement proper grid recreation
		}

		// Update Logo view
		switch item := hm.headerFlex.GetItem(3).(type) {
		case *tview.TextView:
			item.SetText(strings.Join(logoLines, "\n"))
		case *tview.Grid:
			// For now, just update text if it's a grid (simplified approach)
			// TODO: Implement proper grid recreation
		}
	}
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

// createEnhancedTextView creates an enhanced text view with border and title
func (hm *HeaderManager) createEnhancedTextView(lines []string, align int, title string) tview.Primitive {
	// Docker Info and Logo should never use table layout
	if title == "Docker Info" || title == "WhaleTui" {
		textView := tview.NewTextView()
		textView.SetText(strings.Join(lines, "\n"))
		textView.SetTextColor(hm.ui.GetThemeManager().GetTextColor())
		textView.SetBackgroundColor(hm.ui.GetThemeManager().GetBackgroundColor())
		textView.SetTextAlign(align)
		textView.SetBorder(true)
		textView.SetTitle(fmt.Sprintf(" %s ", title))
		textView.SetBorderColor(hm.ui.GetThemeManager().GetBorderColor())
		textView.SetTitleColor(hm.ui.GetThemeManager().GetHeaderColor())
		textView.SetScrollable(false)
		textView.SetDynamicColors(false)
		textView.SetWordWrap(false)
		textView.SetWrap(false)
		return textView
	}

	// For Navigation and Actions, use table layout if content is long
	if len(lines) > 7 {
		return hm.createTableLayout(lines, align, title)
	}

	// Simple text view for short content
	textView := tview.NewTextView()
	textView.SetText(strings.Join(lines, "\n"))
	textView.SetTextColor(hm.ui.GetThemeManager().GetTextColor())
	textView.SetBackgroundColor(hm.ui.GetThemeManager().GetBackgroundColor())
	textView.SetTextAlign(align)
	textView.SetBorder(true)
	textView.SetTitle(fmt.Sprintf(" %s ", title))
	textView.SetBorderColor(hm.ui.GetThemeManager().GetBorderColor())
	textView.SetTitleColor(hm.ui.GetThemeManager().GetHeaderColor())
	textView.SetScrollable(false)
	textView.SetDynamicColors(false)
	textView.SetWordWrap(false)
	textView.SetWrap(false)
	return textView
}

// createTableLayout creates a table-like layout for content longer than 7 lines
func (hm *HeaderManager) createTableLayout(lines []string, align int, title string) tview.Primitive {
	// Calculate how many columns we need
	maxRows := 7
	numColumns := (len(lines) + maxRows - 1) / maxRows // Ceiling division

	// Create a flex container for the table layout
	flex := tview.NewFlex()
	flex.SetDirection(tview.FlexColumn)
	flex.SetBorder(true)
	flex.SetTitle(fmt.Sprintf(" %s ", title))
	flex.SetBorderColor(hm.ui.GetThemeManager().GetBorderColor())
	flex.SetBackgroundColor(hm.ui.GetThemeManager().GetBackgroundColor())

	// Create columns
	for col := 0; col < numColumns; col++ {
		// Create a text view for this column
		columnText := tview.NewTextView()
		columnText.SetTextColor(hm.ui.GetThemeManager().GetTextColor())
		columnText.SetBackgroundColor(hm.ui.GetThemeManager().GetBackgroundColor())
		columnText.SetTextAlign(align)
		columnText.SetBorder(false)
		columnText.SetScrollable(false)
		columnText.SetDynamicColors(false)
		columnText.SetWordWrap(false)
		columnText.SetWrap(false)

		// Collect lines for this column
		var columnLines []string
		for i := col * maxRows; i < (col+1)*maxRows && i < len(lines); i++ {
			columnLines = append(columnLines, lines[i])
		}

		// Set the text for this column
		columnText.SetText(strings.Join(columnLines, "\n"))

		// Add column to flex container
		flex.AddItem(columnText, 0, 1, false)
	}

	return flex
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
func (hm *HeaderManager) getDockerInfoFromService(services interfaces.ServiceFactoryInterface) string {
	if !services.IsContainerServiceAvailable() {
		return hm.getDefaultDockerInfo()
	}

	// Try to get real Docker info
	dockerInfoService := services.GetDockerInfoService()
	if dockerInfoService != nil {
		ctx := context.Background()
		if dockerInfo, err := dockerInfoService.GetDockerInfo(ctx); err == nil {
			return hm.formatDockerInfo(*dockerInfo)
		}
	}

	// Fallback - if service is available but info fetch failed, show partial connection
	connectionStatus := "⚠️ Partial"
	if services.IsContainerServiceAvailable() {
		connectionStatus = "✅ Connected"
	}

	return fmt.Sprintf(constants.DockerInfoTemplate,
		connectionStatus,
		"Available",
		"--",
		"--",
		constants.AppVersion)
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
