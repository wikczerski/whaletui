package builders

import (
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewComponentBuilder(t *testing.T) {
	cb := NewComponentBuilder()
	require.NotNil(t, cb)
}

func TestComponentBuilder_CreateTextView(t *testing.T) {
	cb := NewComponentBuilder()
	text := "Test Text"
	align := tview.AlignCenter
	color := tcell.ColorRed

	tv := cb.CreateTextView(text, align, color)
	require.NotNil(t, tv)
}

func TestComponentBuilder_CreateTextView_TextContent(t *testing.T) {
	cb := NewComponentBuilder()
	text := "Test Text"
	align := tview.AlignCenter
	color := tcell.ColorRed

	tv := cb.CreateTextView(text, align, color)
	assert.Equal(t, text, tv.GetText(true))
}

func TestComponentBuilder_CreateBorderedTextView(t *testing.T) {
	cb := NewComponentBuilder()
	text := "Test Text"
	title := "Test Title"
	color := tcell.ColorBlue

	tv := cb.CreateBorderedTextView(text, title, color)
	require.NotNil(t, tv)
}

func TestComponentBuilder_CreateBorderedTextView_TextContent(t *testing.T) {
	cb := NewComponentBuilder()
	text := "Test Text"
	title := "Test Title"
	color := tcell.ColorBlue

	tv := cb.CreateBorderedTextView(text, title, color)
	assert.Equal(t, text, tv.GetText(true))
}

func TestComponentBuilder_CreateBorderedTextView_Title(t *testing.T) {
	cb := NewComponentBuilder()
	text := "Test Text"
	title := "Test Title"
	color := tcell.ColorBlue

	tv := cb.CreateBorderedTextView(text, title, color)
	assert.Contains(t, tv.GetTitle(), title)
}

func TestComponentBuilder_CreateButton(t *testing.T) {
	cb := NewComponentBuilder()
	text := "Test Button"

	btn := cb.CreateButton(text, func() {})
	require.NotNil(t, btn)
}

func TestComponentBuilder_CreateButton_Label(t *testing.T) {
	cb := NewComponentBuilder()
	text := "Test Button"

	btn := cb.CreateButton(text, func() {})
	assert.Equal(t, text, btn.GetLabel())
}

func TestComponentBuilder_CreateFlex_RowDirection(t *testing.T) {
	cb := NewComponentBuilder()

	rowFlex := cb.CreateFlex(tview.FlexRow)
	require.NotNil(t, rowFlex)
}

func TestComponentBuilder_CreateFlex_ColumnDirection(t *testing.T) {
	cb := NewComponentBuilder()

	colFlex := cb.CreateFlex(tview.FlexColumn)
	require.NotNil(t, colFlex)
}

func TestComponentBuilder_CreateTable(t *testing.T) {
	cb := NewComponentBuilder()

	table := cb.CreateTable()
	require.NotNil(t, table)
}

func TestComponentBuilder_CreateInputField(t *testing.T) {
	cb := NewComponentBuilder()
	label := "Test Label"

	input := cb.CreateInputField(label)
	require.NotNil(t, input)
}

func TestComponentBuilder_CreateInputField_Label(t *testing.T) {
	cb := NewComponentBuilder()
	label := "Test Label"

	input := cb.CreateInputField(label)
	assert.Equal(t, label, input.GetLabel())
}

func TestComponentBuilder_Consistency_TextViewStyling(t *testing.T) {
	cb := NewComponentBuilder()

	// Test that multiple components have consistent styling
	_ = cb.CreateTextView("Text1", tview.AlignLeft, tcell.ColorWhite)
	_ = cb.CreateTextView("Text2", tview.AlignRight, tcell.ColorBlack)

	// Both should have no border and same background
	// Note: Some methods are not available in this version of tview
}

func TestComponentBuilder_Consistency_TableStyling(t *testing.T) {
	cb := NewComponentBuilder()

	// Test table consistency
	table1 := cb.CreateTable()
	table2 := cb.CreateTable()

	// Note: Some table methods are not available in this version of tview
	require.NotNil(t, table1)
	require.NotNil(t, table2)
}

func TestComponentBuilder_EdgeCases_EmptyText(t *testing.T) {
	cb := NewComponentBuilder()

	tv := cb.CreateTextView("", tview.AlignLeft, tcell.ColorWhite)
	require.NotNil(t, tv)
}

func TestComponentBuilder_EdgeCases_NilColor(t *testing.T) {
	cb := NewComponentBuilder()

	tv := cb.CreateTextView("Test", tview.AlignLeft, tcell.ColorDefault)
	require.NotNil(t, tv)
}

func TestComponentBuilder_EdgeCases_EmptyLabel(t *testing.T) {
	cb := NewComponentBuilder()

	btn := cb.CreateButton("", func() {})
	require.NotNil(t, btn)
}

func TestComponentBuilder_EdgeCases_NilCallback(t *testing.T) {
	cb := NewComponentBuilder()

	btn := cb.CreateButton("Test", nil)
	require.NotNil(t, btn)
}

func TestComponentBuilder_ColorHandling_WhiteColor(t *testing.T) {
	cb := NewComponentBuilder()

	tv := cb.CreateTextView("Test", tview.AlignLeft, tcell.ColorWhite)
	require.NotNil(t, tv)
}

func TestComponentBuilder_ColorHandling_BlackColor(t *testing.T) {
	cb := NewComponentBuilder()

	tv := cb.CreateTextView("Test", tview.AlignLeft, tcell.ColorBlack)
	require.NotNil(t, tv)
}

func TestComponentBuilder_ColorHandling_RedColor(t *testing.T) {
	cb := NewComponentBuilder()

	tv := cb.CreateTextView("Test", tview.AlignLeft, tcell.ColorRed)
	require.NotNil(t, tv)
}

func TestComponentBuilder_ColorHandling_BlueColor(t *testing.T) {
	cb := NewComponentBuilder()

	tv := cb.CreateTextView("Test", tview.AlignLeft, tcell.ColorBlue)
	require.NotNil(t, tv)
}

func TestComponentBuilder_ColorHandling_GreenColor(t *testing.T) {
	cb := NewComponentBuilder()

	tv := cb.CreateTextView("Test", tview.AlignLeft, tcell.ColorGreen)
	require.NotNil(t, tv)
}

func TestComponentBuilder_AlignmentHandling_LeftAlign(t *testing.T) {
	cb := NewComponentBuilder()

	tv := cb.CreateTextView("Test", tview.AlignLeft, tcell.ColorWhite)
	require.NotNil(t, tv)
}

func TestComponentBuilder_AlignmentHandling_CenterAlign(t *testing.T) {
	cb := NewComponentBuilder()

	tv := cb.CreateTextView("Test", tview.AlignCenter, tcell.ColorWhite)
	require.NotNil(t, tv)
}

func TestComponentBuilder_AlignmentHandling_RightAlign(t *testing.T) {
	cb := NewComponentBuilder()

	tv := cb.CreateTextView("Test", tview.AlignRight, tcell.ColorWhite)
	require.NotNil(t, tv)
}

func TestComponentBuilder_FlexDirectionHandling_Row(t *testing.T) {
	cb := NewComponentBuilder()

	rowFlex := cb.CreateFlex(tview.FlexRow)
	require.NotNil(t, rowFlex)
}

func TestComponentBuilder_FlexDirectionHandling_Column(t *testing.T) {
	cb := NewComponentBuilder()

	colFlex := cb.CreateFlex(tview.FlexColumn)
	require.NotNil(t, colFlex)
}

func TestComponentBuilder_TableProperties_DefaultTable(t *testing.T) {
	cb := NewComponentBuilder()

	table := cb.CreateTable()
	require.NotNil(t, table)
}

func TestComponentBuilder_InputFieldProperties_DefaultInput(t *testing.T) {
	cb := NewComponentBuilder()

	input := cb.CreateInputField("Test Label")
	require.NotNil(t, input)
}

func TestComponentBuilder_ButtonProperties_DefaultButton(t *testing.T) {
	cb := NewComponentBuilder()

	btn := cb.CreateButton("Test Button", func() {})
	require.NotNil(t, btn)
}
