package views

import (
	"context"
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/wikczerski/D5r/internal/config"
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

	bottomTitleView := componentBuilder.CreateBorderedTextView(
		fmt.Sprintf(" %s<%s> ", lv.ContainerName, lv.ContainerID[:12]),
		"",
		lv.themeManager.GetHeaderColor(),
	)
	bottomTitleView.SetTextAlign(tview.AlignCenter)

	lv.logsText = componentBuilder.CreateTextView("Loading logs...", tview.AlignLeft, lv.themeManager.GetTextColor())
	lv.logsText.SetDynamicColors(true)
	lv.logsText.SetScrollable(true)
	lv.logsText.SetBorder(true)
	lv.logsText.SetBorderColor(lv.themeManager.GetBorderColor())

	backButton := componentBuilder.CreateButton("Back to Table", func() {
		lv.ui.ShowCurrentView()
	})

	logsFlex.AddItem(bottomTitleView, constants.TitleViewHeight, 0, false)
	logsFlex.AddItem(lv.logsText, 0, 1, true)
	logsFlex.AddItem(backButton, constants.BackButtonHeight, 0, false)

	// Set up key bindings
	logsFlex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape, tcell.KeyEnter:
			lv.ui.ShowCurrentView()
			return nil
		}

		// Let the text view handle scrolling keys
		if event.Key() == tcell.KeyUp || event.Key() == tcell.KeyDown ||
			event.Key() == tcell.KeyPgUp || event.Key() == tcell.KeyPgDn ||
			event.Key() == tcell.KeyHome || event.Key() == tcell.KeyEnd {
			return event
		}

		return event
	})

	lv.view = logsFlex
}

// GetView returns the logs view primitive
func (lv *LogsView) GetView() tview.Primitive {
	return lv.view
}

// LoadLogs loads container logs from Docker
func (lv *LogsView) LoadLogs() {
	ctx := context.Background()
	services := lv.ui.GetServices()
	if services == nil || services.ContainerService == nil {
		lv.logsText.SetText("Error: Container service not available")
		return
	}

	logs, err := services.ContainerService.GetContainerLogs(ctx, lv.ContainerID)
	if err != nil {
		lv.logsText.SetText(fmt.Sprintf("Error loading logs: %v", err))
		return
	}

	lv.logsText.SetText(logs)
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
