package views

import (
	"context"
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/wikczerski/D5r/internal/config"
	"github.com/wikczerski/D5r/internal/models"
	"github.com/wikczerski/D5r/internal/ui/builders"
	"github.com/wikczerski/D5r/internal/ui/constants"
	"github.com/wikczerski/D5r/internal/ui/interfaces"
)

// LogsView represents a view for displaying container logs
type LogsView struct {
	ContainerID   string
	ContainerName string
	ui            interfaces.UIInterface
	view          *tview.Flex
	logsText      *tview.TextView
	themeManager  *config.ThemeManager
}

// NewLogsView creates a new logs view
func NewLogsView(ui interfaces.UIInterface, containerID, containerName string) *LogsView {
	lv := &LogsView{
		ContainerID:   containerID,
		ContainerName: containerName,
		ui:            ui,
		themeManager:  ui.GetThemeManager(),
	}
	lv.createView()
	return lv
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
	bottomTitleView := componentBuilder.CreateBorderedTextView(
		fmt.Sprintf(" %s<%s> ", lv.ContainerName, lv.ContainerID[:12]),
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

// GetView returns the logs view primitive
func (lv *LogsView) GetView() tview.Primitive {
	return lv.view
}

// LoadLogs loads container logs from Docker
func (lv *LogsView) LoadLogs() {
	ctx := context.Background()
	logs, err := lv.getContainerLogs(ctx)
	if err != nil {
		lv.logsText.SetText(fmt.Sprintf("Error loading logs: %v", err))
		return
	}

	lv.logsText.SetText(logs)
}

func (lv *LogsView) listContainers(ctx context.Context) ([]models.Container, error) {
	services := lv.ui.GetServices()
	if services == nil || services.GetContainerService() == nil {
		return []models.Container{}, nil
	}
	return services.GetContainerService().ListContainers(ctx)
}

func (lv *LogsView) getContainerLogs(ctx context.Context) (string, error) {
	services := lv.ui.GetServices()
	if services == nil || services.GetContainerService() == nil {
		return "", fmt.Errorf("container service not available")
	}
	return services.GetContainerService().GetContainerLogs(ctx, lv.ContainerID)
}

// GetActions returns the available actions for the logs view
func (lv *LogsView) GetActions() map[rune]string {
	return map[rune]string{
		'f': "Follow logs",
		't': "Tail logs",
		's': "Save logs",
		'c': "Clear logs",
		'w': "Wrap text",
	}
}
