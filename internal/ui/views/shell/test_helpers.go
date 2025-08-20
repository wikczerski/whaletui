package shell

import (
	"github.com/rivo/tview"
)

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

	view.view = tview.NewFlex().SetDirection(tview.FlexRow)
	view.view.AddItem(view.outputView, 0, 1, false)
	view.view.AddItem(view.inputField, 3, 0, true)

	return view
}
