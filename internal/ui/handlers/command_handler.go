package handlers

import (
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/wikczerski/whaletui/internal/ui/constants"
	"github.com/wikczerski/whaletui/internal/ui/interfaces"
)

// CommandHandler manages command mode functionality
type CommandHandler struct {
	ui           interfaces.UIInterface
	commandInput *tview.InputField
	isActive     bool
	errorTimer   *time.Timer // Timer for clearing error messages
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

	// Clear any error timer
	ch.clearError()

	// Only hide command input if it's been initialized
	if ch.commandInput != nil {
		ch.hideCommandInput()
		mainFlex := ch.ui.GetMainFlex().(*tview.Flex)
		mainFlex.RemoveItem(ch.commandInput)
		ch.commandInput.SetText("")
	}

	if viewContainer := ch.ui.GetViewContainer(); viewContainer != nil {
		if vc, ok := viewContainer.(*tview.Flex); ok {
			app := ch.ui.GetApp().(*tview.Application)
			app.SetFocus(vc)
		}
	}
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
		if ch.processCommand(command) {
			ch.Exit()
		}
		// If processCommand returns false, don't exit - let the error message show
	case tcell.KeyEscape:
		ch.Exit()
	case tcell.KeyRune:
		// User is typing - clear any error message
		ch.clearError()
	}
}

// configureCommandInput sets up the command input styling and behavior
func (ch *CommandHandler) configureCommandInput() {
	// Get theme manager for styling
	themeManager := ch.ui.GetThemeManager()

	ch.commandInput.SetLabel(": ")
	ch.commandInput.SetLabelColor(themeManager.GetCommandModeLabelColor())
	ch.commandInput.SetFieldTextColor(themeManager.GetCommandModeTextColor())
	ch.commandInput.SetBorder(true)
	ch.commandInput.SetBorderColor(themeManager.GetCommandModeBorderColor())
	ch.commandInput.SetTitle(" Command Mode ")
	ch.commandInput.SetTitleColor(themeManager.GetCommandModeTitleColor())
	ch.commandInput.SetBackgroundColor(themeManager.GetCommandModeBackgroundColor())
	ch.commandInput.SetPlaceholder("Type view name (containers, images, volumes, networks, swarm services, swarm nodes)")
	ch.commandInput.SetPlaceholderTextColor(themeManager.GetCommandModePlaceholderColor())
	ch.commandInput.SetDoneFunc(ch.HandleInput)
	ch.commandInput.SetAutocompleteFunc(ch.getAutocomplete)
}

// hideCommandInput makes the command input completely invisible
func (ch *CommandHandler) hideCommandInput() {
	if ch.commandInput == nil {
		return
	}

	ch.commandInput.SetBorder(false)
	ch.commandInput.SetBackgroundColor(constants.UIInvisibleColor)
	ch.commandInput.SetFieldTextColor(constants.UIInvisibleColor)
	ch.commandInput.SetLabelColor(constants.UIInvisibleColor)
	ch.commandInput.SetPlaceholderTextColor(constants.UIInvisibleColor)
}

// showCommandInput makes the command input visible with proper styling
func (ch *CommandHandler) showCommandInput() {
	if ch.commandInput == nil {
		return
	}

	// Get theme manager for styling
	themeManager := ch.ui.GetThemeManager()

	ch.commandInput.SetBorder(true)
	ch.commandInput.SetLabelColor(themeManager.GetCommandModeLabelColor())
	ch.commandInput.SetFieldTextColor(themeManager.GetCommandModeTextColor())
	ch.commandInput.SetPlaceholderTextColor(themeManager.GetCommandModePlaceholderColor())
}

// processCommand executes the given command
// Returns true if the command was successfully processed and the input should close
func (ch *CommandHandler) processCommand(command string) bool {
	// Handle empty command - just clear and stay open
	if strings.TrimSpace(command) == "" {
		ch.commandInput.SetText("")
		return false
	}

	if ch.handleViewSwitchCommand(command) {
		return true
	}

	if ch.handleSystemCommand(command) {
		return true
	}

	if ch.handleHelpCommand(command) {
		return true
	}

	// Handle unknown command
	ch.handleUnknownCommand(command)
	return false // Don't close input for unknown commands
}

// handleViewSwitchCommand handles view switching commands
func (ch *CommandHandler) handleViewSwitchCommand(command string) bool {
	switch command {
	case "containers", "c":
		ch.ui.SwitchView("containers")
		ch.Exit()
		return true
	case "images", "i":
		ch.ui.SwitchView("images")
		ch.Exit()
		return true
	case "volumes", "v":
		ch.ui.SwitchView("volumes")
		ch.Exit()
		return true
	case "networks", "n":
		ch.ui.SwitchView("networks")
		ch.Exit()
		return true
	case "swarm services", "swarm", "services", "s":
		ch.ui.SwitchView("swarmServices")
		ch.Exit()
		return true
	case "swarm nodes", "nodes", "w":
		ch.ui.SwitchView("swarmNodes")
		ch.Exit()
		return true
	}
	return false
}

// handleSystemCommand handles system-level commands
func (ch *CommandHandler) handleSystemCommand(command string) bool {
	switch command {
	case "quit", "q", "exit":
		ch.handleQuitCommand()
		return true
	}
	return false
}

// handleHelpCommand handles help-related commands
func (ch *CommandHandler) handleHelpCommand(command string) bool {
	switch command {
	case "help", "?":
		ch.ui.ShowHelp()
		ch.Exit()
		return true
	}
	return false
}

// handleQuitCommand handles the quit command by sending shutdown signal
func (ch *CommandHandler) handleQuitCommand() {
	shutdownChan := ch.ui.GetShutdownChan()
	select {
	case shutdownChan <- struct{}{}:
	default:
	}
	ch.Exit()
}

// handleUnknownCommand handles unknown commands by showing an error message in the command input
func (ch *CommandHandler) handleUnknownCommand(_ string) {
	ch.showCommandError("Wrong command")
	// Don't exit automatically - let the user see the error and continue typing
	// The error message will disappear after 3 seconds
}

// showCommandError displays an error message in the command input with red text
func (ch *CommandHandler) showCommandError(message string) {
	if ch.commandInput == nil {
		return
	}

	// Cancel any existing timer
	if ch.errorTimer != nil {
		ch.errorTimer.Stop()
	}

	// Get theme manager for styling
	themeManager := ch.ui.GetThemeManager()

	// Store original theme color
	originalTextColor := themeManager.GetCommandModeTextColor()

	// Set error message with red text
	ch.commandInput.SetText(message)
	ch.commandInput.SetFieldTextColor(tcell.ColorRed)

	// Set up a timer to clear the error message after 3 seconds
	ch.errorTimer = time.AfterFunc(3*time.Second, func() {
		// Use the UI's app to schedule the update on the main thread
		if ch.isActive && ch.commandInput != nil {
			app := ch.ui.GetApp().(*tview.Application)
			app.QueueUpdateDraw(func() {
				if ch.isActive && ch.commandInput != nil {
					ch.commandInput.SetText("")
					ch.commandInput.SetFieldTextColor(originalTextColor)
				}
			})
		}
	})
}

// clearError clears any displayed error message immediately
func (ch *CommandHandler) clearError() {
	if ch.errorTimer != nil {
		ch.errorTimer.Stop()
		ch.errorTimer = nil
	}

	if ch.commandInput != nil && ch.isActive {
		themeManager := ch.ui.GetThemeManager()
		ch.commandInput.SetText("")
		ch.commandInput.SetFieldTextColor(themeManager.GetCommandModeTextColor())
	}
}

// getAutocomplete provides command autocomplete suggestions
func (ch *CommandHandler) getAutocomplete(currentText string) []string {
	suggestions := []string{
		"containers", "images", "volumes", "networks",
		"swarm services", "swarm nodes", "services", "nodes",
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
