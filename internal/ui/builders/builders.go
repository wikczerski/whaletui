package builders

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/user/d5r/internal/ui/constants"
)

// ComponentBuilder provides methods to create consistent UI components
type ComponentBuilder struct{}

// NewComponentBuilder creates a new component builder
func NewComponentBuilder() *ComponentBuilder {
	return &ComponentBuilder{}
}

// CreateTextView creates a text view with consistent styling
func (cb *ComponentBuilder) CreateTextView(text string, align int, color tcell.Color) *tview.TextView {
	tv := tview.NewTextView()
	tv.SetText(text)
	tv.SetTextAlign(align)
	tv.SetTextColor(color)
	tv.SetBorder(false)
	tv.SetBackgroundColor(constants.BackgroundColor)
	return tv
}

// CreateBorderedTextView creates a bordered text view with consistent styling
func (cb *ComponentBuilder) CreateBorderedTextView(text, title string, color tcell.Color) *tview.TextView {
	tv := tview.NewTextView()
	tv.SetText(text)
	tv.SetTitle(fmt.Sprintf(" %s ", title))
	tv.SetBorder(true)
	tv.SetBorderColor(color)
	tv.SetBackgroundColor(constants.BackgroundColor)
	return tv
}

// CreateButton creates a button with consistent styling
func (cb *ComponentBuilder) CreateButton(text string, onSelected func()) *tview.Button {
	btn := tview.NewButton(text)
	btn.SetSelectedFunc(onSelected)
	return btn
}

// CreateFlex creates a flex container with consistent styling
func (cb *ComponentBuilder) CreateFlex(direction int) *tview.Flex {
	return tview.NewFlex().SetDirection(direction)
}

// CreateTable creates a table with consistent styling
func (cb *ComponentBuilder) CreateTable() *tview.Table {
	table := tview.NewTable()
	table.SetBorders(false)
	table.SetBorderPadding(0, 0, 0, 0)
	table.SetSelectable(true, false)
	table.SetFixed(1, 0)
	table.SetBorderColor(constants.BorderColor)
	return table
}

// CreateInputField creates an input field with consistent styling
func (cb *ComponentBuilder) CreateInputField(label string) *tview.InputField {
	input := tview.NewInputField()
	input.SetLabel(label)
	input.SetBorder(true)
	input.SetBorderColor(constants.BorderColor)
	input.SetBackgroundColor(constants.BackgroundColor)
	return input
}
