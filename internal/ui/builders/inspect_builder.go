package builders

import (
	"encoding/json"
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/wikczerski/whaletui/internal/ui/constants"
)

// CreateInspectDetailsView is an exported function for backward compatibility
func CreateInspectDetailsView(
	title string,
	inspectData map[string]any,
	actions map[rune]string,
	onAction func(rune),
	onBack func(),
) *tview.Flex {
	detailsFlex := tview.NewFlex().SetDirection(tview.FlexRow)

	titleView := createInspectTitleView(title)
	inspectText := createInspectTextView(inspectData, actions)
	backButton := tview.NewButton("Back to Table").SetSelectedFunc(onBack)

	// Add components to flex
	detailsFlex.AddItem(titleView, constants.TitleViewHeight, 0, false)
	detailsFlex.AddItem(inspectText, 0, 1, true) // Set to true to make it focusable and scrollable
	detailsFlex.AddItem(backButton, constants.BackButtonHeight, 0, false)

	// Set up key bindings for the details view
	setupInspectDetailsKeyBindings(detailsFlex, inspectText, onAction, onBack)

	return detailsFlex
}

// CreateInspectView is an exported function for backward compatibility
func CreateInspectView(title string) (*tview.TextView, *tview.Flex) {
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

// createInspectTitleView creates the title view for the inspect details
func createInspectTitleView(title string) *tview.TextView {
	titleView := tview.NewTextView().
		SetText(fmt.Sprintf(" %s ", title)).
		SetTextAlign(tview.AlignCenter)
	titleView.SetBorder(true).SetBorderColor(constants.HeaderColor)
	return titleView
}

// createInspectTextView creates the main text view for displaying inspect data
func createInspectTextView(inspectData map[string]any, actions map[rune]string) *tview.TextView {
	inspectText := createBaseInspectTextView()
	setupInspectTextScrolling(inspectText)

	content := buildInspectContent(inspectData, actions)
	inspectText.SetText(content)

	return inspectText
}

// createBaseInspectTextView creates the base inspect text view with common settings
func createBaseInspectTextView() *tview.TextView {
	inspectText := tview.NewTextView()
	inspectText.SetDynamicColors(true)
	inspectText.SetScrollable(true)
	inspectText.SetBorder(true)
	inspectText.SetBorderColor(constants.BorderColor)
	return inspectText
}

// buildInspectContent builds the content for the inspect text view
func buildInspectContent(inspectData map[string]any, actions map[rune]string) string {
	condensedJSON := formatInspectData(inspectData)

	if len(actions) > 0 {
		condensedJSON += "\n\nActions:\n" + formatActionsText(actions)
	}

	return condensedJSON
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
func setupInspectTextScrollingHandleSpacebar(
	event *tcell.EventKey,
	inspectText *tview.TextView,
) bool {
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
func setupInspectDetailsKeyBindings(
	detailsFlex *tview.Flex,
	_ *tview.TextView,
	onAction func(rune),
	onBack func(),
) {
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

// setupInspectDetailsKeyBindingsHandleNavigationKeys handles navigation keys (Escape, Enter, Backspace)
func setupInspectDetailsKeyBindingsHandleNavigationKeys(event *tcell.EventKey, onBack func()) bool {
	switch event.Key() {
	case tcell.KeyEscape, tcell.KeyEnter, tcell.KeyBackspace:
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
func setupInspectDetailsKeyBindingsHandleActionKeys(
	event *tcell.EventKey,
	onAction func(rune),
) bool {
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
