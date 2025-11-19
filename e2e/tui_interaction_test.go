package e2e

import (
	"testing"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
	"github.com/wikczerski/whaletui/e2e/framework"
)

// TestTUIBasicRendering tests that the TUI renders basic components.
func TestTUIBasicRendering(t *testing.T) {
	fw := framework.NewTUITestFramework(t)

	// Create a simple table
	table := tview.NewTable()
	table.SetBorder(true)
	table.SetTitle(" Test Table ")

	// Add headers
	headers := []string{"ID", "Name", "Status"}
	for col, header := range headers {
		table.SetCell(0, col, tview.NewTableCell(header).SetTextColor(tcell.ColorYellow))
	}

	// Add a row
	table.SetCell(1, 0, tview.NewTableCell("123"))
	table.SetCell(1, 1, tview.NewTableCell("test-container"))
	table.SetCell(1, 2, tview.NewTableCell("running"))

	// Start app with table
	fw.StartApp(table)

	// Verify title is rendered
	assert.True(t, fw.WaitForText("Test Table", 2*time.Second), "Table title should be visible")

	// Verify headers are rendered
	assert.True(t, fw.VerifyTextContains("ID"), "ID header should be visible")
	assert.True(t, fw.VerifyTextContains("Name"), "Name header should be visible")
	assert.True(t, fw.VerifyTextContains("Status"), "Status header should be visible")

	// Verify data is rendered
	assert.True(t, fw.VerifyTextContains("test-container"), "Container name should be visible")
	assert.True(t, fw.VerifyTextContains("running"), "Status should be visible")
}

// TestTUIKeyboardNavigation tests keyboard navigation in a table.
func TestTUIKeyboardNavigation(t *testing.T) {
	fw := framework.NewTUITestFramework(t)

	// Create a table with multiple rows
	table := tview.NewTable()
	table.SetBorder(true)
	table.SetSelectable(true, false)

	// Add rows
	rows := []string{"Row 1", "Row 2", "Row 3"}
	for i, row := range rows {
		table.SetCell(i, 0, tview.NewTableCell(row))
	}

	// Track selected row
	selectedRow := 0
	table.SetSelectionChangedFunc(func(row, col int) {
		selectedRow = row
	})

	fw.StartApp(table)

	// Wait for initial render
	time.Sleep(200 * time.Millisecond)

	// Initial selection should be row 0
	assert.Equal(t, 0, selectedRow, "Initial selection should be row 0")

	// Press down arrow
	fw.InjectKeyPress(tcell.KeyDown)
	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, 1, selectedRow, "Selection should move to row 1")

	// Press down arrow again
	fw.InjectKeyPress(tcell.KeyDown)
	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, 2, selectedRow, "Selection should move to row 2")

	// Press up arrow
	fw.InjectKeyPress(tcell.KeyUp)
	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, 1, selectedRow, "Selection should move back to row 1")
}

// TestTUIInputField tests input field interaction.
func TestTUIInputField(t *testing.T) {
	fw := framework.NewTUITestFramework(t)

	// Create input field
	inputField := tview.NewInputField()
	inputField.SetLabel("Search: ")
	inputField.SetFieldWidth(20)

	fw.StartApp(inputField)

	// Wait for render
	assert.True(t, fw.WaitForText("Search:", 2*time.Second), "Label should be visible")

	// Type some text
	fw.InjectString("test")
	time.Sleep(200 * time.Millisecond)

	// Verify text was entered
	assert.Equal(t, "test", inputField.GetText(), "Input field should contain typed text")
}

// TestTUIModal tests modal dialog interaction.
func TestTUIModal(t *testing.T) {
	fw := framework.NewTUITestFramework(t)

	// Create modal
	modal := tview.NewModal()
	modal.SetText("Are you sure?")
	modal.AddButtons([]string{"Yes", "No"})

	buttonPressed := ""
	modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		buttonPressed = buttonLabel
	})

	fw.StartApp(modal)

	// Wait for modal to render
	assert.True(t, fw.WaitForText("Are you sure?", 2*time.Second), "Modal text should be visible")
	assert.True(t, fw.VerifyTextContains("Yes"), "Yes button should be visible")
	assert.True(t, fw.VerifyTextContains("No"), "No button should be visible")

	// Press Enter to select first button (Yes)
	fw.InjectKeyPress(tcell.KeyEnter)
	time.Sleep(100 * time.Millisecond)

	assert.Equal(t, "Yes", buttonPressed, "Yes button should be pressed")
}

// TestTUIViewSwitching tests switching between different views.
func TestTUIViewSwitching(t *testing.T) {
	fw := framework.NewTUITestFramework(t)

	// Create pages
	pages := tview.NewPages()

	// Create two views
	view1 := tview.NewTextView()
	view1.SetText("View 1 Content")
	view1.SetBorder(true)
	view1.SetTitle(" View 1 ")

	view2 := tview.NewTextView()
	view2.SetText("View 2 Content")
	view2.SetBorder(true)
	view2.SetTitle(" View 2 ")

	pages.AddPage("view1", view1, true, true)
	pages.AddPage("view2", view2, true, false)

	fw.StartApp(pages)

	// Verify view 1 is shown
	assert.True(t, fw.WaitForText("View 1", 2*time.Second), "View 1 should be visible")
	assert.True(t, fw.VerifyTextContains("View 1 Content"), "View 1 content should be visible")

	// Switch to view 2
	pages.SwitchToPage("view2")
	fw.Sync()
	time.Sleep(300 * time.Millisecond) // Give more time for view switch

	// Verify view 2 is shown
	assert.True(t, fw.VerifyTextContains("View 2"), "View 2 should be visible")
	assert.True(t, fw.VerifyTextContains("View 2 Content"), "View 2 content should be visible")
}

// TestTUISearchFunctionality tests search/filter functionality.
func TestTUISearchFunctionality(t *testing.T) {
	fw := framework.NewTUITestFramework(t)

	// Create a flex layout with table and search input
	flex := tview.NewFlex()
	flex.SetDirection(tview.FlexRow)

	// Create table with data
	table := tview.NewTable()
	table.SetBorder(true)
	items := []string{"nginx", "redis", "postgres", "alpine"}
	for i, item := range items {
		table.SetCell(i, 0, tview.NewTableCell(item))
	}

	// Create search input
	searchInput := tview.NewInputField()
	searchInput.SetLabel("Search: ")

	// Implement search filtering
	searchInput.SetChangedFunc(func(text string) {
		// Simple filter: hide rows that don't match
		for i := 0; i < len(items); i++ {
			cell := table.GetCell(i, 0)
			if text == "" || contains(items[i], text) {
				cell.SetText(items[i])
			} else {
				cell.SetText("")
			}
		}
	})

	flex.AddItem(table, 0, 1, false)
	flex.AddItem(searchInput, 1, 0, true)

	fw.StartApp(flex)

	// Verify all items are visible initially
	assert.True(t, fw.WaitForText("nginx", 2*time.Second), "nginx should be visible")
	assert.True(t, fw.VerifyTextContains("redis"), "redis should be visible")
	assert.True(t, fw.VerifyTextContains("postgres"), "postgres should be visible")

	// Type search query
	fw.InjectString("ng")
	time.Sleep(200 * time.Millisecond)

	// Verify filtering works (nginx should still be visible)
	assert.True(t, fw.VerifyTextContains("nginx"), "nginx should still be visible after search")
}

// TestTUIColorRendering tests that colors are rendered correctly.
func TestTUIColorRendering(t *testing.T) {
	fw := framework.NewTUITestFramework(t)

	// Create table with colored cells
	table := tview.NewTable()
	table.SetBorder(true)

	// Add cells with different colors
	table.SetCell(0, 0, tview.NewTableCell("Running").SetTextColor(tcell.ColorGreen))
	table.SetCell(1, 0, tview.NewTableCell("Stopped").SetTextColor(tcell.ColorRed))
	table.SetCell(2, 0, tview.NewTableCell("Created").SetTextColor(tcell.ColorYellow))

	fw.StartApp(table)

	// Wait for render
	time.Sleep(200 * time.Millisecond)

	// Verify text is present (color verification would require checking style)
	assert.True(t, fw.VerifyTextContains("Running"), "Running should be visible")
	assert.True(t, fw.VerifyTextContains("Stopped"), "Stopped should be visible")
	assert.True(t, fw.VerifyTextContains("Created"), "Created should be visible")

	// Note: Full color verification would require checking cell styles
	// which is possible with fw.GetCellStyle(x, y)
}

// Helper function for contains check
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
