package shell

import (
	"fmt"

	"github.com/rivo/tview"
	"github.com/wikczerski/whaletui/internal/config"
	"github.com/wikczerski/whaletui/internal/shared"
)

// createView creates the shell view UI components
func (sv *View) createView() {
	sv.createOutputView()
	sv.createInputField()
	sv.createMainLayout()
	sv.addWelcomeMessage()
}

// createOutputView creates and configures the output view
func (sv *View) createOutputView() {
	themeManager := sv.ui.GetThemeManager()

	sv.outputView = tview.NewTextView()
	sv.outputView.SetDynamicColors(true)
	sv.outputView.SetScrollable(true)
	sv.outputView.SetBorder(true)
	sv.outputView.SetTitle(
		fmt.Sprintf(" Shell - %s (%s) ", sv.containerName, shared.TruncName(sv.containerID, 12)),
	)
	sv.outputView.SetTitleColor(themeManager.GetShellTitleColor())
	sv.outputView.SetBorderColor(themeManager.GetShellBorderColor())
	sv.outputView.SetTextColor(themeManager.GetShellTextColor())
	sv.outputView.SetBackgroundColor(themeManager.GetShellBackgroundColor())
}

// createInputField creates and configures the input field
func (sv *View) createInputField() {
	themeManager := sv.ui.GetThemeManager()

	sv.inputField = tview.NewInputField()
	sv.setupInputFieldStyling(themeManager)
	sv.setupInputFieldBehavior()
}

// setupInputFieldStyling sets up the visual styling of the input field
func (sv *View) setupInputFieldStyling(themeManager *config.ThemeManager) {
	sv.inputField.SetLabel("$ ")
	sv.inputField.SetLabelColor(themeManager.GetShellCmdLabelColor())
	sv.inputField.SetFieldTextColor(themeManager.GetShellCmdTextColor())
	sv.inputField.SetBorder(true)
	sv.inputField.SetBorderColor(themeManager.GetShellCmdBorderColor())
	sv.inputField.SetBackgroundColor(themeManager.GetShellCmdBackgroundColor())
	sv.inputField.SetPlaceholder("Type command and press Enter (ESC to exit shell mode)")
}

// setupInputFieldBehavior sets up the behavior of the input field
func (sv *View) setupInputFieldBehavior() {
	sv.inputField.SetDoneFunc(sv.handleCommand)
	sv.inputField.SetInputCapture(sv.handleInputCapture)
}

// createMainLayout creates the main layout and adds components
func (sv *View) createMainLayout() {
	sv.view = tview.NewFlex().SetDirection(tview.FlexRow)
	sv.view.AddItem(sv.outputView, 0, 1, false)
	sv.view.AddItem(sv.inputField, 3, 0, true)
}

// addWelcomeMessage adds the initial welcome message to the output view
func (sv *View) addWelcomeMessage() {
	sv.addOutput(fmt.Sprintf("Welcome to shell for container: %s (%s)\n",
		sv.containerName,
		shared.TruncName(sv.containerID, 12)))
	sv.addOutput("Type 'exit' or press ESC to return to container view\n\n")
}
