package e2e

import (
	"fmt"
	"testing"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
	"github.com/wikczerski/whaletui/e2e/framework"
)

// TestTUIContainerListWorkflow tests the complete container list and interaction workflow.
func TestTUIContainerListWorkflow(t *testing.T) {
	fw := framework.NewTUITestFramework(t)

	// Create a container list view
	table := tview.NewTable()
	table.SetBorder(true)
	table.SetTitle(" Containers ")
	table.SetSelectable(true, false)

	// Add headers
	headers := []string{"ID", "Name", "Image", "Status", "State"}
	for col, header := range headers {
		table.SetCell(0, col, tview.NewTableCell(header).
			SetTextColor(tcell.ColorYellow).
			SetSelectable(false))
	}

	// Add container data
	containers := []struct {
		id, name, image, status, state string
		color                          tcell.Color
	}{
		{"abc123", "web-server", "nginx:alpine", "Up 2 hours", "running", tcell.ColorGreen},
		{"def456", "cache", "redis:alpine", "Exited (0)", "exited", tcell.ColorRed},
		{"ghi789", "database", "postgres:13", "Up 1 day", "running", tcell.ColorGreen},
	}

	for i, c := range containers {
		row := i + 1
		table.SetCell(row, 0, tview.NewTableCell(c.id))
		table.SetCell(row, 1, tview.NewTableCell(c.name))
		table.SetCell(row, 2, tview.NewTableCell(c.image))
		table.SetCell(row, 3, tview.NewTableCell(c.status))
		table.SetCell(row, 4, tview.NewTableCell(c.state).SetTextColor(c.color))
	}

	fw.StartApp(table)

	// Verify container list is displayed
	assert.True(t, fw.WaitForText("Containers", 2*time.Second),
		"Container view title should be visible")
	assert.True(t, fw.VerifyTextContains("web-server"), "Container name should be visible")
	assert.True(t, fw.VerifyTextContains("nginx:alpine"), "Image name should be visible")
	assert.True(t, fw.VerifyTextContains("running"), "Container state should be visible")

	// Set initial selection to first data row (row 1, after header row 0)
	table.Select(1, 0)

	// Navigate through containers
	fw.InjectKeyPress(tcell.KeyDown) // Move to second container (row 2)
	time.Sleep(100 * time.Millisecond)

	// Verify navigation worked (should be on row 2 now)
	row, _ := table.GetSelection()
	assert.Equal(t, 2, row, "Should be on second container row")
}

// TestTUICommandModeWorkflow tests command mode interaction.
func TestTUICommandModeWorkflow(t *testing.T) {
	fw := framework.NewTUITestFramework(t)

	// Create main view with command input
	pages := tview.NewPages()

	// Main view
	mainView := tview.NewTextView()
	mainView.SetText("Main View - Press : for command mode")
	mainView.SetBorder(true)

	// Command input (initially hidden)
	commandInput := tview.NewInputField()
	commandInput.SetLabel(": ")
	commandInput.SetFieldWidth(30)

	// Layout
	flex := tview.NewFlex()
	flex.SetDirection(tview.FlexRow)
	flex.AddItem(mainView, 0, 1, true)
	flex.AddItem(commandInput, 1, 0, false)

	pages.AddPage("main", flex, true, true)

	// Set up input capture for command mode
	commandMode := false
	flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == ':' && !commandMode {
			commandMode = true
			fw.GetApp().SetFocus(commandInput)
			return nil
		}
		if event.Key() == tcell.KeyEscape && commandMode {
			commandMode = false
			commandInput.SetText("")
			fw.GetApp().SetFocus(mainView)
			return nil
		}
		return event
	})

	fw.StartApp(pages)

	// Verify main view is shown
	assert.True(t, fw.WaitForText("Main View", 2*time.Second), "Main view should be visible")

	// Enter command mode
	fw.InjectKeyRune(':')
	time.Sleep(200 * time.Millisecond)

	// Type a command
	fw.InjectString("containers")
	time.Sleep(200 * time.Millisecond)

	// Verify command was typed
	assert.Equal(t, "containers", commandInput.GetText(), "Command should be entered")

	// Exit command mode
	fw.InjectKeyPress(tcell.KeyEscape)
	time.Sleep(200 * time.Millisecond)

	// Verify command input is cleared
	assert.Equal(t, "", commandInput.GetText(), "Command input should be cleared")
}

// TestTUISearchWorkflow tests search/filter workflow.
func TestTUISearchWorkflow(t *testing.T) {
	fw := framework.NewTUITestFramework(t)

	// Create searchable list
	flex := tview.NewFlex()
	flex.SetDirection(tview.FlexRow)

	// List of items
	list := tview.NewList()
	list.SetBorder(true)
	list.SetTitle(" Containers ")

	items := []string{"nginx-web", "redis-cache", "postgres-db", "alpine-test"}
	for _, item := range items {
		list.AddItem(item, "", 0, nil)
	}

	// Search input
	searchInput := tview.NewInputField()
	searchInput.SetLabel("/ ")
	searchInput.SetFieldWidth(30)

	// Implement search
	searchInput.SetChangedFunc(func(text string) {
		list.Clear()
		for _, item := range items {
			if text == "" || containsSubstring(item, text) {
				list.AddItem(item, "", 0, nil)
			}
		}
	})

	flex.AddItem(list, 0, 1, true)
	flex.AddItem(searchInput, 1, 0, false)

	fw.StartApp(flex)

	// Verify all items are visible
	assert.True(t, fw.WaitForText("nginx-web", 2*time.Second), "nginx-web should be visible")
	assert.True(t, fw.VerifyTextContains("redis-cache"), "redis-cache should be visible")

	// Focus search input and type
	fw.GetApp().SetFocus(searchInput)
	time.Sleep(100 * time.Millisecond)

	fw.InjectString("redis")
	time.Sleep(200 * time.Millisecond)

	// Verify filtering (redis-cache should still be visible, others filtered out)
	assert.True(t, fw.VerifyTextContains("redis-cache"),
		"redis-cache should be visible after search")
	assert.Equal(t, 1, list.GetItemCount(), "Only 1 item should match")
}

// TestTUIModalConfirmationWorkflow tests modal confirmation workflow.
func TestTUIModalConfirmationWorkflow(t *testing.T) {
	fw := framework.NewTUITestFramework(t)

	pages := tview.NewPages()

	// Main view
	mainView := tview.NewTextView()
	mainView.SetText("Press 'd' to delete")
	mainView.SetBorder(true)

	pages.AddPage("main", mainView, true, true)

	// Track if delete was confirmed
	deleteConfirmed := false

	// Set up input capture
	mainView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'd' {
			// Show confirmation modal
			modal := tview.NewModal()
			modal.SetText("Are you sure you want to delete?")
			modal.AddButtons([]string{"Yes", "No"})
			modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
				if buttonLabel == "Yes" {
					deleteConfirmed = true
				}
				pages.RemovePage("modal")
			})

			pages.AddPage("modal", modal, true, true)
			return nil
		}
		return event
	})

	fw.StartApp(pages)

	// Verify main view
	assert.True(t, fw.WaitForText("Press 'd' to delete", 2*time.Second),
		"Main view should be visible")

	// Press 'd' to trigger delete
	fw.InjectKeyRune('d')
	time.Sleep(200 * time.Millisecond)

	// Verify modal is shown
	assert.True(t, fw.WaitForText("Are you sure", 2*time.Second),
		"Confirmation modal should appear")

	// Press Enter to confirm (Yes button)
	fw.InjectKeyPress(tcell.KeyEnter)
	time.Sleep(200 * time.Millisecond)

	// Verify delete was confirmed
	assert.True(t, deleteConfirmed, "Delete should be confirmed")
}

// TestTUIDetailsViewWorkflow tests viewing details of an item.
func TestTUIDetailsViewWorkflow(t *testing.T) {
	fw := framework.NewTUITestFramework(t)

	pages := tview.NewPages()

	// List view
	list := tview.NewList()
	list.SetBorder(true)
	list.SetTitle(" Containers ")
	list.AddItem("nginx-web", "Running", 0, nil)
	list.AddItem("redis-cache", "Stopped", 0, nil)

	pages.AddPage("list", list, true, true)

	// Set up Enter key to show details
	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEnter {
			index := list.GetCurrentItem()
			itemName, _ := list.GetItemText(index)

			// Create details view
			details := tview.NewTextView()
			details.SetBorder(true)
			details.SetTitle(fmt.Sprintf(" Details: %s ", itemName))
			details.SetText(fmt.Sprintf(
				"Container: %s\nStatus: Running\nImage: nginx:alpine\nPorts: 80:8080",
				itemName,
			))

			// Set up Backspace to return
			details.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
				if event.Key() == tcell.KeyBackspace || event.Key() == tcell.KeyEscape {
					pages.SwitchToPage("list")
					return nil
				}
				return event
			})

			pages.AddPage("details", details, true, true)
			return nil
		}
		return event
	})

	fw.StartApp(pages)

	// Verify list view
	assert.True(t, fw.WaitForText("Containers", 2*time.Second), "Container list should be visible")
	assert.True(t, fw.VerifyTextContains("nginx-web"), "Container should be visible")

	// Press Enter to view details
	fw.InjectKeyPress(tcell.KeyEnter)
	time.Sleep(200 * time.Millisecond)

	// Verify details view
	assert.True(t, fw.WaitForText("Details:", 2*time.Second), "Details view should be shown")
	assert.True(t, fw.VerifyTextContains("nginx:alpine"), "Image details should be visible")

	// Press Backspace to return
	fw.InjectKeyPress(tcell.KeyBackspace)
	time.Sleep(200 * time.Millisecond)

	// Verify back to list view
	assert.True(t, fw.VerifyTextContains("Containers"), "Should return to container list")
}

// TestTUIMultiViewNavigation tests navigating between multiple views.
func TestTUIMultiViewNavigation(t *testing.T) {
	fw := framework.NewTUITestFramework(t)

	pages := tview.NewPages()

	// Create multiple views
	views := map[string]string{
		"containers": "Containers View",
		"images":     "Images View",
		"volumes":    "Volumes View",
	}

	for name, content := range views {
		view := tview.NewTextView()
		view.SetText(content)
		view.SetBorder(true)
		view.SetTitle(fmt.Sprintf(" %s ", content))
		pages.AddPage(name, view, true, name == "containers")
	}

	fw.StartApp(pages)

	// Verify initial view (containers)
	assert.True(t, fw.WaitForText("Containers View", 2*time.Second),
		"Containers view should be shown")

	// Switch to images
	pages.SwitchToPage("images")
	fw.GetApp().Draw() // Force redraw
	fw.Sync()

	assert.True(t, fw.WaitForText("Images View", 2*time.Second), "Images view should be shown")

	// Switch to volumes
	pages.SwitchToPage("volumes")
	fw.GetApp().Draw() // Force redraw
	fw.Sync()

	assert.True(t, fw.WaitForText("Volumes View", 2*time.Second), "Volumes view should be shown")
}

// Helper function
func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
