package managers

import (
	"os"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/wikczerski/D5r/internal/ui/interfaces"
)

// CommandHandler manages command mode functionality
type CommandHandler struct {
	ui           interfaces.UIInterface
	commandInput *tview.InputField
	isActive     bool
}

// NewCommandHandler creates a new command handler
func NewCommandHandler(ui interfaces.UIInterface) *CommandHandler {
	return &CommandHandler{ui: ui}
}

// CreateCommandInput creates and configures the command input field
func (ch *CommandHandler) CreateCommandInput() *tview.InputField {
	ch.commandInput = tview.NewInputField()
	ch.configureCommandInput()
	ch.hideCommandInput()
	return ch.commandInput
}

// configureCommandInput sets up the command input styling and behavior
func (ch *CommandHandler) configureCommandInput() {
	ch.commandInput.SetLabel(": ")
	ch.commandInput.SetLabelColor(tcell.ColorYellow)
	ch.commandInput.SetFieldTextColor(tcell.ColorWhite)
	ch.commandInput.SetBorder(true)
	ch.commandInput.SetBorderColor(tcell.ColorYellow)
	ch.commandInput.SetTitle(" Command Mode ")
	ch.commandInput.SetTitleColor(tcell.ColorYellow)
	ch.commandInput.SetBackgroundColor(tcell.ColorDarkGray)
	ch.commandInput.SetPlaceholder("Type view name (containers, images, volumes, networks)")
	ch.commandInput.SetPlaceholderTextColor(tcell.ColorGray)
	ch.commandInput.SetDoneFunc(ch.HandleInput)
	ch.commandInput.SetAutocompleteFunc(ch.getAutocomplete)
}

// hideCommandInput makes the command input completely invisible
func (ch *CommandHandler) hideCommandInput() {
	ch.commandInput.SetBorder(false)
	ch.commandInput.SetBackgroundColor(tcell.ColorDefault)
	ch.commandInput.SetFieldTextColor(tcell.ColorDefault)
	ch.commandInput.SetLabelColor(tcell.ColorDefault)
	ch.commandInput.SetPlaceholderTextColor(tcell.ColorDefault)
}

// showCommandInput makes the command input visible with proper styling
func (ch *CommandHandler) showCommandInput() {
	ch.commandInput.SetBorder(true)
	ch.commandInput.SetLabelColor(tcell.ColorYellow)
	ch.commandInput.SetFieldTextColor(tcell.ColorWhite)
	ch.commandInput.SetPlaceholderTextColor(tcell.ColorGray)
}

// Enter activates command mode
func (ch *CommandHandler) Enter() {
	ch.isActive = true
	ch.showCommandInput()
	mainFlex := ch.ui.GetMainFlex().(*tview.Flex)
	mainFlex.AddItem(ch.commandInput, 3, 1, true)
	app := ch.ui.GetApp().(*tview.Application)
	app.SetFocus(ch.commandInput)
}

// Exit deactivates command mode
func (ch *CommandHandler) Exit() {
	ch.isActive = false
	ch.hideCommandInput()
	mainFlex := ch.ui.GetMainFlex().(*tview.Flex)
	mainFlex.RemoveItem(ch.commandInput)
	ch.commandInput.SetText("")

	// Return focus to current view
	// For now, skip focus restoration to avoid complex type assertions
	// TODO: Implement proper focus restoration

	// ch.ui.log.Debug("Command mode deactivated")
}

// IsActive returns whether command mode is currently active
func (ch *CommandHandler) IsActive() bool {
	return ch.isActive
}

// GetInput returns the command input widget
func (ch *CommandHandler) GetInput() *tview.InputField {
	return ch.commandInput
}

// HandleInput processes command input
func (ch *CommandHandler) HandleInput(key tcell.Key) {
	switch key {
	case tcell.KeyEnter:
		command := ch.commandInput.GetText()
		ch.processCommand(command)
		ch.Exit()
	case tcell.KeyEscape:
		ch.Exit()
	}
}

// processCommand executes the given command
func (ch *CommandHandler) processCommand(command string) {
	switch command {
	case "containers", "c":
		ch.ui.SwitchView("containers")
	case "images", "i":
		ch.ui.SwitchView("images")
	case "volumes", "v":
		ch.ui.SwitchView("volumes")
	case "networks", "n":
		ch.ui.SwitchView("networks")
	case "quit", "q", "exit":
		os.Exit(0)
	case "help", "?":
		ch.ui.ShowHelp()
	default:
		// TODO: add feedback to user
		// ch.ui.log.Warn("Unknown command: %s", command)
	}
}

// getAutocomplete provides command autocomplete suggestions
func (ch *CommandHandler) getAutocomplete(currentText string) []string {
	suggestions := []string{
		"containers", "images", "volumes", "networks",
		"quit", "q", "exit", "help", "?",
	}

	var matches []string
	for _, suggestion := range suggestions {
		if strings.HasPrefix(strings.ToLower(suggestion), strings.ToLower(currentText)) {
			matches = append(matches, suggestion)
		}
	}
	return matches
}
