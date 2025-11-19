package builders

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/wikczerski/whaletui/internal/ui/constants"
)

// DetailsViewBuilder creates details views with consistent styling
type DetailsViewBuilder struct {
	builder *ComponentBuilder
}

// CreateDetailsView is an exported function for backward compatibility
func CreateDetailsView(
	title, details string,
	actions map[rune]string,
	onAction func(rune),
	onBack func(),
) *tview.Flex {
	return NewDetailsViewBuilder().CreateDetailsView(title, details, actions, onAction, onBack)
}

// NewDetailsViewBuilder creates a new details view builder
func NewDetailsViewBuilder() *DetailsViewBuilder {
	return &DetailsViewBuilder{
		builder: NewComponentBuilder(),
	}
}

// CreateDetailsView creates a details view that can be displayed inline
func (dvb *DetailsViewBuilder) CreateDetailsView(
	title, details string,
	actions map[rune]string,
	onAction func(rune),
	onBack func(),
) *tview.Flex {
	detailsFlex := dvb.builder.CreateFlex(tview.FlexRow)

	dvb.addDetailsViewComponents(detailsFlex, title, details, actions, onBack)
	dvb.setupDetailsKeyBindings(detailsFlex, onAction, onBack)

	return detailsFlex
}

// addDetailsViewComponents adds the main components to the details view
func (dvb *DetailsViewBuilder) addDetailsViewComponents(
	detailsFlex *tview.Flex,
	title, details string,
	actions map[rune]string,
	onBack func(),
) {
	titleView := dvb.builder.CreateBorderedTextView("", title, constants.HeaderColor)
	titleView.SetTextAlign(tview.AlignCenter)

	detailsText := dvb.builder.CreateBorderedTextView(
		details+"\nActions:\n"+dvb.formatActions(actions),
		"",
		constants.BorderColor)
	detailsText.SetDynamicColors(true)
	detailsText.SetScrollable(true)

	backButton := dvb.builder.CreateButton("Back to Table", onBack)

	detailsFlex.AddItem(titleView, constants.TitleViewHeight, 0, false)
	detailsFlex.AddItem(detailsText, 0, 1, false)
	detailsFlex.AddItem(backButton, constants.BackButtonHeight, 0, true)
}

// setupDetailsKeyBindings sets up keyboard navigation for details view
func (dvb *DetailsViewBuilder) setupDetailsKeyBindings(
	detailsFlex *tview.Flex,
	onAction func(rune),
	onBack func(),
) {
	detailsFlex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape, tcell.KeyEnter, tcell.KeyBackspace:
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
