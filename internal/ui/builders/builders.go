// Package builders provides UI component builders for WhaleTUI.
package builders

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/wikczerski/whaletui/internal/config"
)

// ComponentBuilder provides methods to create consistent UI components
type ComponentBuilder struct {
	themeManager *config.ThemeManager
}

// NewComponentBuilder creates a new component builder
func NewComponentBuilder() *ComponentBuilder {
	return &ComponentBuilder{}
}

// NewComponentBuilderWithTheme creates a new component builder with theme support
func NewComponentBuilderWithTheme(themeManager *config.ThemeManager) *ComponentBuilder {
	return &ComponentBuilder{
		themeManager: themeManager,
	}
}

// CreateTextView creates a text view with consistent styling
func (cb *ComponentBuilder) CreateTextView(
	text string,
	align int,
	color tcell.Color,
) *tview.TextView {
	tv := tview.NewTextView()
	tv.SetText(text)
	tv.SetTextAlign(align)
	tv.SetTextColor(color)
	tv.SetBorder(false)
	if cb.themeManager != nil {
		tv.SetBackgroundColor(cb.themeManager.GetBackgroundColor())
	}
	return tv
}

// CreateBorderedTextView creates a bordered text view with consistent styling
func (cb *ComponentBuilder) CreateBorderedTextView(
	text, title string,
	color tcell.Color,
) *tview.TextView {
	tv := tview.NewTextView()
	tv.SetText(text)
	tv.SetTitle(fmt.Sprintf(" %s ", title))
	tv.SetBorder(true)
	tv.SetBorderColor(color)
	if cb.themeManager != nil {
		tv.SetBackgroundColor(cb.themeManager.GetBackgroundColor())
	}
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
	if cb.themeManager != nil {
		table.SetBorderColor(cb.themeManager.GetBorderColor())
	}
	return table
}

// CreateInputField creates an input field with consistent styling
func (cb *ComponentBuilder) CreateInputField(label string) *tview.InputField {
	input := tview.NewInputField()
	input.SetLabel(label)
	input.SetBorder(true)
	if cb.themeManager != nil {
		input.SetBorderColor(cb.themeManager.GetBorderColor())
		input.SetBackgroundColor(cb.themeManager.GetBackgroundColor())
	}
	return input
}
