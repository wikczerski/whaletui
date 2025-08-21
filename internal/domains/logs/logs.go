package logs

import (
	"context"
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/wikczerski/whaletui/internal/config"
	"github.com/wikczerski/whaletui/internal/ui/builders"
	"github.com/wikczerski/whaletui/internal/ui/constants"
	"github.com/wikczerski/whaletui/internal/ui/interfaces"
)

// LogsView represents a view for displaying logs from any Docker resource
type LogsView struct {
	ResourceType string
	ResourceID   string
	ResourceName string
	ui           interfaces.UIInterface
	view         *tview.Flex
	logsText     *tview.TextView
	themeManager *config.ThemeManager
}

// NewLogsView creates a new logs view for any resource type
func NewLogsView(ui interfaces.UIInterface, resourceType, resourceID, resourceName string) *LogsView {
	lv := &LogsView{
		ResourceType: resourceType,
		ResourceID:   resourceID,
		ResourceName: resourceName,
		ui:           ui,
		themeManager: ui.GetThemeManager(),
	}
	lv.createView()
	return lv
}

// GetView returns the logs view primitive
func (lv *LogsView) GetView() tview.Primitive {
	return lv.view
}

// LoadLogs loads logs from the specified Docker resource
func (lv *LogsView) LoadLogs() {
	ctx := context.Background()
	logs, err := lv.getResourceLogs(ctx)
	if err != nil {
		lv.logsText.SetText(fmt.Sprintf("Error loading logs: %v", err))
		return
	}

	lv.logsText.SetText(logs)
}

// GetActions returns the available actions for the logs view
func (lv *LogsView) GetActions() map[rune]string {
	services := lv.ui.GetServices()
	if !services.IsServiceAvailable("logs") {
		return map[rune]string{}
	}

	logsService := services.GetLogsService()
	if logsService == nil {
		return map[rune]string{}
	}

	return logsService.GetActions()
}

// createView creates the logs view UI components
func (lv *LogsView) createView() {
	componentBuilder := builders.NewComponentBuilderWithTheme(lv.themeManager)
	viewBuilder := builders.NewViewBuilder()

	logsFlex := viewBuilder.CreateView()

	lv.createTitleView(componentBuilder, logsFlex)
	lv.createLogsTextView(componentBuilder, logsFlex)
	lv.createBackButton(componentBuilder, logsFlex)
	lv.setupKeyBindings(logsFlex)

	lv.view = logsFlex
}

// createTitleView creates the title view for the logs
func (lv *LogsView) createTitleView(componentBuilder *builders.ComponentBuilder, logsFlex *tview.Flex) {
	displayName := lv.ResourceName
	if displayName == "" {
		displayName = lv.ResourceID
	}

	shortID := lv.ResourceID
	if len(shortID) > 12 {
		shortID = shortID[:12]
	}

	bottomTitleView := componentBuilder.CreateBorderedTextView(
		fmt.Sprintf(" %s<%s> (%s) ", displayName, shortID, lv.ResourceType),
		"",
		lv.themeManager.GetHeaderColor(),
	)
	bottomTitleView.SetTextAlign(tview.AlignCenter)

	logsFlex.AddItem(bottomTitleView, constants.TitleViewHeight, 0, false)
}

// createLogsTextView creates the logs text view
func (lv *LogsView) createLogsTextView(componentBuilder *builders.ComponentBuilder, logsFlex *tview.Flex) {
	lv.logsText = componentBuilder.CreateTextView("Loading logs...", tview.AlignLeft, lv.themeManager.GetTextColor())
	lv.logsText.SetDynamicColors(true)
	lv.logsText.SetScrollable(true)
	lv.logsText.SetBorder(true)
	lv.logsText.SetBorderColor(lv.themeManager.GetBorderColor())

	logsFlex.AddItem(lv.logsText, 0, 1, true)
}

// createBackButton creates the back button
func (lv *LogsView) createBackButton(componentBuilder *builders.ComponentBuilder, logsFlex *tview.Flex) {
	backButton := componentBuilder.CreateButton("Back to Table", func() {
		lv.ui.ShowCurrentView()
	})

	logsFlex.AddItem(backButton, constants.BackButtonHeight, 0, false)
}

// handleNavigationKeys handles navigation key events
func (lv *LogsView) handleNavigationKeys(event *tcell.EventKey) bool {
	switch event.Key() {
	case tcell.KeyEscape, tcell.KeyEnter:
		lv.ui.ShowCurrentView()
		return true
	}
	return false
}

// handleScrollingKeys handles scrolling key events
func (lv *LogsView) handleScrollingKeys(event *tcell.EventKey) bool {
	return event.Key() == tcell.KeyUp || event.Key() == tcell.KeyDown ||
		event.Key() == tcell.KeyPgUp || event.Key() == tcell.KeyPgDn ||
		event.Key() == tcell.KeyHome || event.Key() == tcell.KeyEnd
}

// setupKeyBindings sets up the key bindings for the logs view
func (lv *LogsView) setupKeyBindings(logsFlex *tview.Flex) {
	logsFlex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if lv.handleNavigationKeys(event) {
			return nil
		}

		if lv.handleScrollingKeys(event) {
			return event
		}

		return event
	})
}

func (lv *LogsView) getResourceLogs(ctx context.Context) (string, error) {
	services := lv.ui.GetServices()
	logsService := services.GetLogsService()
	if logsService == nil {
		return "", fmt.Errorf("logs service not available")
	}

	return logsService.GetLogs(ctx, lv.ResourceType, lv.ResourceID)
}
