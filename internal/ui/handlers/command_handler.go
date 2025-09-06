package handlers

import (
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/wikczerski/whaletui/internal/config"
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
	mainFlex, ok := ch.ui.GetMainFlex().(*tview.Flex)
	if !ok {
		return
	}
	mainFlex.AddItem(ch.commandInput, 3, 1, true)
	app, ok := ch.ui.GetApp().(*tview.Application)
	if !ok {
		return
	}
	app.SetFocus(ch.commandInput)
}

// Exit deactivates command mode
func (ch *CommandHandler) Exit() {
	ch.isActive = false
	ch.clearError()
	ch.cleanupCommandInput()
	ch.restoreFocus()
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
	themeManager := ch.ui.GetThemeManager()
	ch.setupBasicStyling(themeManager)
	ch.setupAdvancedStyling(themeManager)
	ch.setupBehavior()
}

// setupBasicStyling sets up the basic styling for the command input
func (ch *CommandHandler) setupBasicStyling(themeManager *config.ThemeManager) {
	ch.commandInput.SetLabel(": ")
	ch.commandInput.SetLabelColor(themeManager.GetCommandModeLabelColor())
	ch.commandInput.SetFieldTextColor(themeManager.GetCommandModeTextColor())
	ch.commandInput.SetBorder(true)
	ch.commandInput.SetBorderColor(themeManager.GetCommandModeBorderColor())
}

// setupAdvancedStyling sets up the advanced styling for the command input
func (ch *CommandHandler) setupAdvancedStyling(themeManager *config.ThemeManager) {
	ch.commandInput.SetTitle(" Command Mode ")
	ch.commandInput.SetTitleColor(themeManager.GetCommandModeTitleColor())
	ch.commandInput.SetBackgroundColor(themeManager.GetCommandModeBackgroundColor())
	placeholder := "Type view name (containers, images, volumes, networks, swarm services, swarm nodes)"
	ch.commandInput.SetPlaceholder(placeholder)
	ch.commandInput.SetPlaceholderTextColor(themeManager.GetCommandModePlaceholderColor())
}

// setupBehavior sets up the behavior for the command input
func (ch *CommandHandler) setupBehavior() {
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
	if ch.isEmptyCommand(command) {
		return false
	}

	return ch.processCommandTypes(command)
}

// isEmptyCommand checks if the command is empty
func (ch *CommandHandler) isEmptyCommand(command string) bool {
	if strings.TrimSpace(command) == "" {
		ch.commandInput.SetText("")
		return true
	}
	return false
}

// processCommandTypes processes different types of commands
func (ch *CommandHandler) processCommandTypes(command string) bool {
	if ch.handleViewSwitchCommand(command) {
		return true
	}

	if ch.handleSystemCommand(command) {
		return true
	}

	if ch.handleHelpCommand(command) {
		return true
	}

	ch.handleUnknownCommand(command)
	return false
}

// handleViewSwitchCommand handles view switching commands
func (ch *CommandHandler) handleViewSwitchCommand(command string) bool {
	viewMappings := ch.getViewMappings()

	if viewName, exists := viewMappings[command]; exists {
		ch.ui.SwitchView(viewName)
		ch.Exit()
		return true
	}
	return false
}

// getViewMappings returns the mapping of command aliases to view names
func (ch *CommandHandler) getViewMappings() map[string]string {
	return map[string]string{
		"containers": constants.ViewContainers,
		"images":     constants.ViewImages,
		"volumes":    constants.ViewVolumes,
		"networks":   constants.ViewNetworks,
		"services":   constants.ViewSwarmServices,
		"nodes":      constants.ViewSwarmNodes,
	}
}

// handleSystemCommand handles system-level commands
func (ch *CommandHandler) handleSystemCommand(command string) bool {
	switch command {
	case "quit", "q", "exit":
		ch.handleQuitCommand()
		return true
	case "reload", "r":
		ch.handleReloadThemeCommand()
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

// handleReloadThemeCommand handles the reload theme command
func (ch *CommandHandler) handleReloadThemeCommand() {
	err := ch.ui.ReloadTheme()
	if err != nil {
		ch.showCommandError("Failed to reload theme: " + err.Error())
		return
	}

	ch.showCommandSuccess("Theme reloaded successfully")
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
			app, ok := ch.ui.GetApp().(*tview.Application)
			if !ok {
				return
			}
			app.QueueUpdateDraw(func() {
				if ch.isActive && ch.commandInput != nil {
					ch.commandInput.SetText("")
					ch.commandInput.SetFieldTextColor(originalTextColor)
				}
			})
		}
	})
}

// showCommandSuccess displays a success message in the command input with green text
func (ch *CommandHandler) showCommandSuccess(message string) {
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

	// Set success message with green text
	ch.commandInput.SetText(message)
	ch.commandInput.SetFieldTextColor(tcell.ColorGreen)

	// Set up a timer to clear the success message after 2 seconds
	ch.errorTimer = time.AfterFunc(2*time.Second, func() {
		// Use the UI's app to schedule the update on the main thread
		if ch.isActive && ch.commandInput != nil {
			app, ok := ch.ui.GetApp().(*tview.Application)
			if !ok {
				return
			}
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
		"quit", "q", "exit", "help", "?", "reload", "r",
	}

	var matches []string
	for _, suggestion := range suggestions {
		if strings.HasPrefix(strings.ToLower(suggestion), strings.ToLower(currentText)) {
			matches = append(matches, suggestion)
		}
	}
	return matches
}

// cleanupCommandInput cleans up the command input widget
func (ch *CommandHandler) cleanupCommandInput() {
	if ch.commandInput == nil {
		return
	}

	ch.hideCommandInput()
	mainFlex, ok := ch.ui.GetMainFlex().(*tview.Flex)
	if !ok {
		return
	}
	mainFlex.RemoveItem(ch.commandInput)
	ch.commandInput.SetText("")
}

// restoreFocus restores focus to the main view
func (ch *CommandHandler) restoreFocus() {
	if viewContainer := ch.ui.GetViewContainer(); viewContainer != nil {
		if vc, ok := viewContainer.(*tview.Flex); ok {
			app, ok := ch.ui.GetApp().(*tview.Application)
			if !ok {
				return
			}
			app.SetFocus(vc)
		}
	}
}
