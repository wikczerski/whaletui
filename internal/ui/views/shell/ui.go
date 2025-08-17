package shell

import (
	"fmt"

	"github.com/rivo/tview"
)

// createView creates the shell view UI components
func (sv *View) createView() {
	themeManager := sv.ui.GetThemeManager()

	sv.outputView = tview.NewTextView()
	sv.outputView.SetDynamicColors(true)
	sv.outputView.SetScrollable(true)
	sv.outputView.SetBorder(true)
	sv.outputView.SetTitle(fmt.Sprintf(" Shell - %s (%s) ", sv.containerName, sv.containerID[:12]))
	sv.outputView.SetTitleColor(themeManager.GetShellTitleColor())
	sv.outputView.SetBorderColor(themeManager.GetShellBorderColor())
	sv.outputView.SetTextColor(themeManager.GetShellTextColor())
	sv.outputView.SetBackgroundColor(themeManager.GetShellBackgroundColor())

	sv.inputField = tview.NewInputField()
	sv.inputField.SetLabel("$ ")
	sv.inputField.SetLabelColor(themeManager.GetShellCmdLabelColor())
	sv.inputField.SetFieldTextColor(themeManager.GetShellCmdTextColor())
	sv.inputField.SetBorder(true)
	sv.inputField.SetBorderColor(themeManager.GetShellCmdBorderColor())
	sv.inputField.SetBackgroundColor(themeManager.GetShellCmdBackgroundColor())
	sv.inputField.SetPlaceholder("Type command and press Enter (ESC to exit shell mode)")

	sv.inputField.SetDoneFunc(sv.handleCommand)
	sv.inputField.SetInputCapture(sv.handleInputCapture)

	sv.view = tview.NewFlex().SetDirection(tview.FlexRow)
	sv.view.AddItem(sv.outputView, 0, 1, false)
	sv.view.AddItem(sv.inputField, 3, 0, true)

	sv.addOutput(fmt.Sprintf("Welcome to shell for container: %s (%s)\n", sv.containerName, sv.containerID[:12]))
	sv.addOutput("Type 'exit' or press ESC to return to container view\n\n")
}
