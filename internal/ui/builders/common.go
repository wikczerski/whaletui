package builders

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/wikczerski/whaletui/internal/ui/constants"
)

// TimeFormatter provides time formatting utilities
type TimeFormatter struct{}

// NewTimeFormatter creates a new time formatter
func NewTimeFormatter() *TimeFormatter {
	return &TimeFormatter{}
}

// FormatTime formats a time.Time to a human-readable string
func (tf *TimeFormatter) FormatTime(t time.Time) string {
	if time.Since(t) < constants.TimeThreshold24h {
		return fmt.Sprintf("%s %s", tf.formatDuration(time.Since(t)), constants.TimeFormatRelative)
	}
	return t.Format(constants.TimeFormatAbsolute)
}

// formatDuration formats a duration to a human-readable string
func (tf *TimeFormatter) formatDuration(d time.Duration) string {
	seconds := int(d.Seconds())

	switch {
	case seconds < 60:
		return fmt.Sprintf("%ds", seconds)
	case seconds < 3600:
		return fmt.Sprintf("%dm %ds", seconds/60, seconds%60)
	case seconds < 86400:
		return fmt.Sprintf("%dh %dm", seconds/3600, (seconds%3600)/60)
	default:
		return fmt.Sprintf("%dd %dh", seconds/86400, (seconds%86400)/3600)
	}
}

// DetailsViewBuilder creates details views with consistent styling
type DetailsViewBuilder struct {
	builder *ComponentBuilder
}

// NewDetailsViewBuilder creates a new details view builder
func NewDetailsViewBuilder() *DetailsViewBuilder {
	return &DetailsViewBuilder{
		builder: NewComponentBuilder(),
	}
}

// CreateDetailsView creates a details view that can be displayed inline
func (dvb *DetailsViewBuilder) CreateDetailsView(title, details string, actions map[rune]string, onAction func(rune), onBack func()) *tview.Flex {
	detailsFlex := dvb.builder.CreateFlex(tview.FlexRow)

	titleView := dvb.builder.CreateBorderedTextView("", title, constants.HeaderColor)
	titleView.SetTextAlign(tview.AlignCenter)

	detailsText := dvb.builder.CreateBorderedTextView(details+"\nActions:\n"+dvb.formatActions(actions), "", constants.BorderColor)
	detailsText.SetDynamicColors(true)
	detailsText.SetScrollable(true)

	backButton := dvb.builder.CreateButton("Back to Table", onBack)

	detailsFlex.AddItem(titleView, constants.TitleViewHeight, 0, false)
	detailsFlex.AddItem(detailsText, 0, 1, false)
	detailsFlex.AddItem(backButton, constants.BackButtonHeight, 0, true)

	dvb.setupDetailsKeyBindings(detailsFlex, onAction, onBack)

	return detailsFlex
}

// setupDetailsKeyBindings sets up keyboard navigation for details view
func (dvb *DetailsViewBuilder) setupDetailsKeyBindings(detailsFlex *tview.Flex, onAction func(rune), onBack func()) {
	detailsFlex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape, tcell.KeyEnter:
			if onBack != nil {
				onBack()
			}
			return nil
		}

		if onAction != nil {
			onAction(event.Rune())
		}
		return nil
	})
}

// formatActions formats the actions map into a readable string
func (dvb *DetailsViewBuilder) formatActions(actions map[rune]string) string {
	var result string
	for key, action := range actions {
		result += fmt.Sprintf("%c: %s\n", key, action)
	}
	return result
}

// TableBuilder provides methods to create and configure tables
type TableBuilder struct {
	builder *ComponentBuilder
}

// NewTableBuilder creates a new table builder
func NewTableBuilder() *TableBuilder {
	return &TableBuilder{
		builder: NewComponentBuilder(),
	}
}

// CreateTable creates a new table with consistent styling
func (tb *TableBuilder) CreateTable() *tview.Table {
	return tb.builder.CreateTable()
}

// SetupHeaders sets up table headers with consistent styling
func (tb *TableBuilder) SetupHeaders(table *tview.Table, headers []string) {
	for i, header := range headers {
		cell := tview.NewTableCell(header).
			SetTextColor(constants.HeaderColor).
			SetSelectable(false).
			SetAlign(tview.AlignCenter).
			SetExpansion(1)
		table.SetCell(0, i, cell)
	}
}

// SetupRow sets up a table row with consistent styling
func (tb *TableBuilder) SetupRow(table *tview.Table, row int, cells []string, textColor tcell.Color) {
	for i, cell := range cells {
		tableCell := tview.NewTableCell(cell).
			SetTextColor(textColor).
			SetAlign(tview.AlignLeft).
			SetExpansion(1)
		table.SetCell(row, i, tableCell)
	}
}

// ViewBuilder provides methods to create views
type ViewBuilder struct {
	builder *ComponentBuilder
}

// NewViewBuilder creates a new view builder
func NewViewBuilder() *ViewBuilder {
	return &ViewBuilder{
		builder: NewComponentBuilder(),
	}
}

// CreateView creates a new view with consistent setup
func (vb *ViewBuilder) CreateView() *tview.Flex {
	return vb.builder.CreateFlex(tview.FlexRow)
}

// FormatTime is an exported function for backward compatibility
func FormatTime(t time.Time) string {
	return NewTimeFormatter().FormatTime(t)
}

// CreateDetailsView is an exported function for backward compatibility
func CreateDetailsView(title, details string, actions map[rune]string, onAction func(rune), onBack func()) *tview.Flex {
	return NewDetailsViewBuilder().CreateDetailsView(title, details, actions, onAction, onBack)
}

// CreateInspectDetailsView is an exported function for backward compatibility
func CreateInspectDetailsView(title string, inspectData map[string]any, actions map[rune]string, onAction func(rune), onBack func()) *tview.Flex {
	return createInspectDetailsView(title, inspectData, actions, onAction, onBack)
}

// createInspectView creates a reusable inspection view
func createInspectView(title string) (*tview.TextView, *tview.Flex) {
	inspectView := tview.NewTextView().SetDynamicColors(true).SetScrollable(true)
	inspectView.SetTitle(fmt.Sprintf(" %s ", title)).SetBorder(true)

	backButton := tview.NewButton("Back").SetSelectedFunc(func() {
		// This will be handled by the caller
	})

	// Create flex for inspect view
	inspectFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	inspectFlex.AddItem(inspectView, 0, 1, false)
	inspectFlex.AddItem(backButton, 1, 0, true)

	return inspectView, inspectFlex
}

// CreateInspectView is an exported function for backward compatibility
func CreateInspectView(title string) (*tview.TextView, *tview.Flex) {
	return createInspectView(title)
}

// createInspectDetailsView creates a details view that displays Docker inspect data in condensed JSON
func createInspectDetailsView(title string, inspectData map[string]any, actions map[rune]string, onAction func(rune), onBack func()) *tview.Flex {
	detailsFlex := tview.NewFlex().SetDirection(tview.FlexRow)

	titleView := createInspectTitleView(title)
	inspectText := createInspectTextView(inspectData, actions)
	backButton := createInspectBackButton(onBack)

	// Add components to flex
	detailsFlex.AddItem(titleView, constants.TitleViewHeight, 0, false)
	detailsFlex.AddItem(inspectText, 0, 1, true) // Set to true to make it focusable and scrollable
	detailsFlex.AddItem(backButton, constants.BackButtonHeight, 0, false)

	// Set up key bindings for the details view
	setupInspectDetailsKeyBindings(detailsFlex, inspectText, onAction, onBack)

	return detailsFlex
}

// createInspectTitleView creates the title view for the inspect details
func createInspectTitleView(title string) *tview.TextView {
	titleView := tview.NewTextView().SetText(fmt.Sprintf(" %s ", title)).SetTextAlign(tview.AlignCenter)
	titleView.SetBorder(true).SetBorderColor(constants.HeaderColor)
	return titleView
}

// createInspectTextView creates the main text view for displaying inspect data
func createInspectTextView(inspectData map[string]any, actions map[rune]string) *tview.TextView {
	inspectText := tview.NewTextView()
	inspectText.SetDynamicColors(true)
	inspectText.SetScrollable(true)
	inspectText.SetBorder(true)
	inspectText.SetBorderColor(constants.BorderColor)

	// Configure spacebar for half-page scrolling
	setupInspectTextScrolling(inspectText)

	// Format the inspect data as condensed JSON
	condensedJSON := formatInspectData(inspectData)

	// Add actions if provided
	if len(actions) > 0 {
		condensedJSON += "\n\nActions:\n" + formatActionsText(actions)
	}

	inspectText.SetText(condensedJSON)
	return inspectText
}

// createInspectBackButton creates the back button for the inspect details
func createInspectBackButton(onBack func()) *tview.Button {
	return tview.NewButton("Back to Table").SetSelectedFunc(onBack)
}

// setupInspectTextScrolling configures the scrolling behavior for the inspect text view
func setupInspectTextScrolling(inspectText *tview.TextView) {
	inspectText.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if setupInspectTextScrollingHandleSpacebar(event, inspectText) {
			return nil
		}
		return event
	})
}

// setupInspectTextScrollingHandleSpacebar handles spacebar for half-page scrolling
func setupInspectTextScrollingHandleSpacebar(event *tcell.EventKey, inspectText *tview.TextView) bool {
	if event.Key() == tcell.KeyRune && event.Rune() == ' ' {
		setupInspectTextScrollingPerformHalfPageScroll(inspectText)
		return true
	}
	return false
}

// setupInspectTextScrollingPerformHalfPageScroll performs the actual half-page scrolling
func setupInspectTextScrollingPerformHalfPageScroll(inspectText *tview.TextView) {
	// Get current scroll position
	_, currentLine := inspectText.GetScrollOffset()

	// Get the visible area height
	_, _, _, visibleHeight := inspectText.GetInnerRect()

	// Calculate half of the visible area
	halfView := setupInspectTextScrollingCalculateHalfView(visibleHeight)

	// Calculate new scroll position
	newLine := currentLine + halfView

	// Scroll to the new position
	inspectText.ScrollTo(newLine, 0)
}

// setupInspectTextScrollingCalculateHalfView calculates half of the visible area
func setupInspectTextScrollingCalculateHalfView(visibleHeight int) int {
	halfView := visibleHeight / 2
	if halfView < 1 {
		halfView = 1
	}
	return halfView
}

// setupInspectDetailsKeyBindings sets up the key bindings for the inspect details view
func setupInspectDetailsKeyBindings(detailsFlex *tview.Flex, _ *tview.TextView, onAction func(rune), onBack func()) {
	detailsFlex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if setupInspectDetailsKeyBindingsHandleNavigationKeys(event, onBack) {
			return nil
		}

		if setupInspectDetailsKeyBindingsHandleScrollingKeys(event) {
			return event
		}

		if setupInspectDetailsKeyBindingsHandleActionKeys(event, onAction) {
			return nil
		}

		return event
	})
}

// setupInspectDetailsKeyBindingsHandleNavigationKeys handles navigation keys (Escape, Enter)
func setupInspectDetailsKeyBindingsHandleNavigationKeys(event *tcell.EventKey, onBack func()) bool {
	switch event.Key() {
	case tcell.KeyEscape, tcell.KeyEnter:
		if onBack != nil {
			onBack()
		}
		return true
	}
	return false
}

// setupInspectDetailsKeyBindingsHandleScrollingKeys handles scrolling keys and spacebar
func setupInspectDetailsKeyBindingsHandleScrollingKeys(event *tcell.EventKey) bool {
	// Handle spacebar for half-page scrolling
	if event.Key() == tcell.KeyRune && event.Rune() == ' ' {
		return true
	}

	// Handle scrolling keys
	if event.Key() == tcell.KeyUp || event.Key() == tcell.KeyDown ||
		event.Key() == tcell.KeyPgUp || event.Key() == tcell.KeyPgDn ||
		event.Key() == tcell.KeyHome || event.Key() == tcell.KeyEnd {
		return true
	}

	return false
}

// setupInspectDetailsKeyBindingsHandleActionKeys handles action keys (rune characters)
func setupInspectDetailsKeyBindingsHandleActionKeys(event *tcell.EventKey, onAction func(rune)) bool {
	if event.Key() == tcell.KeyRune {
		if onAction != nil {
			onAction(event.Rune())
		}
		return true
	}
	return false
}

// formatInspectData formats Docker inspect data in a condensed, readable JSON format
func formatInspectData(data map[string]any) string {
	if data == nil {
		return "No inspect data available"
	}

	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Sprintf("Error formatting data: %v", err)
	}

	jsonStr := string(jsonBytes)

	// For very long JSON, we could truncate or show only key sections
	// For now, return the full formatted JSON
	return jsonStr
}

// formatActionsText formats the actions map into a readable string
func formatActionsText(actions map[rune]string) string {
	var result string
	for key, action := range actions {
		result += fmt.Sprintf("%c: %s\n", key, action)
	}
	return result
}
