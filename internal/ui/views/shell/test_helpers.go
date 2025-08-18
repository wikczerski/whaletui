package shell

import (
	"github.com/rivo/tview"
)

// createTestView creates a properly initialized test view
// Note: This does not call createView() to avoid theme manager dependencies
func createTestView() *View {
	return &View{
		containerID:     "test-container",
		containerName:   "test-container",
		inputField:      tview.NewInputField(),
		outputView:      tview.NewTextView(),
		commandHistory:  make([]string, 0),
		multiLineBuffer: make([]string, 0),
		historyIndex:    0,
		currentInput:    "",
		isMultiLine:     false,
	}
}

// createFullTestView creates a test view with all fields initialized
// This is for tests that need a complete view structure
func createFullTestView() *View {
	view := &View{
		containerID:     "test-container",
		containerName:   "test-container",
		inputField:      tview.NewInputField(),
		outputView:      tview.NewTextView(),
		commandHistory:  make([]string, 0),
		multiLineBuffer: make([]string, 0),
		historyIndex:    0,
		currentInput:    "",
		isMultiLine:     false,
	}

	// Create a basic flex layout without theme manager
	view.view = tview.NewFlex().SetDirection(tview.FlexRow)
	view.view.AddItem(view.outputView, 0, 1, false)
	view.view.AddItem(view.inputField, 3, 0, true)

	return view
}
