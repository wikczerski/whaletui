// Package shell provides shell command functionality for WhaleTUI.
package shell

import (
	"context"
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/wikczerski/whaletui/internal/shared"
)

// handleCommand processes the entered command
func (sv *View) handleCommand(key tcell.Key) {
	if key != tcell.KeyEnter {
		return
	}

	command := sv.inputField.GetText()
	if command == "" {
		return
	}

	sv.processCommand(command)
	sv.inputField.SetText("")
}

// processCommand processes the entered command
func (sv *View) processCommand(command string) {
	if sv.isMultiLineCommand(command) {
		sv.handleMultiLineCommand(command)
		return
	}

	sv.handleSingleLineCommand(command)
}

// handleMultiLineCommand handles multi-line command input
func (sv *View) handleMultiLineCommand(command string) {
	sv.addToMultiLineBuffer(command)
	sv.addOutput(fmt.Sprintf("> %s\n", command))
}

// handleSingleLineCommand handles single-line command input
func (sv *View) handleSingleLineCommand(command string) {
	sv.processMultiLineCommand(command)
	sv.updateCommandState(command)

	if sv.handleBuiltInCommand(command) {
		return
	}

	sv.executeDockerCommand(command)
}

// processMultiLineCommand processes multi-line command input
func (sv *View) processMultiLineCommand(command string) {
	if len(sv.multiLineBuffer) > 0 {
		sv.addToMultiLineBuffer(command)
		command = sv.getMultiLineCommand()
		sv.addOutput(fmt.Sprintf("$ %s\n", command))
		sv.clearMultiLineBuffer()
	} else {
		sv.addOutput(fmt.Sprintf("$ %s\n", command))
	}
}

// updateCommandState updates the command state after processing
func (sv *View) updateCommandState(command string) {
	sv.currentInput = ""
	sv.historyIndex = len(sv.commandHistory)
	sv.addToHistory(command)
}

// handleBuiltInCommand handles built-in shell commands
func (sv *View) handleBuiltInCommand(command string) bool {
	switch command {
	case "exit", "quit":
		return sv.handleExitCommand()
	case "clear":
		return sv.handleClearCommand()
	case "help":
		return sv.handleHelpCommand()
	}
	return false
}

// handleExitCommand handles the exit command
func (sv *View) handleExitCommand() bool {
	sv.addOutput("Exiting shell...\n")
	sv.exitShell()
	return true
}

// handleClearCommand handles the clear command
func (sv *View) handleClearCommand() bool {
	sv.outputView.Clear()
	sv.addOutput(fmt.Sprintf("Welcome to shell for container: %s (%s)\n",
		sv.containerName, shared.TruncName(sv.containerID, 12)))
	sv.addOutput("Type 'exit' or press ESC to return to container view\n\n")
	return true
}

// handleHelpCommand handles the help command
func (sv *View) handleHelpCommand() bool {
	sv.showHelp()
	return true
}

// executeDockerCommand executes a command in the container
func (sv *View) executeDockerCommand(command string) {
	if !sv.validateCommandExecution(command) {
		return
	}

	args := sv.parseCommandArgs(command)
	if len(args) == 0 {
		sv.addOutput("Error: invalid command\n")
		return
	}

	sv.executeCommandWithArgs(args)
}

// validateCommandExecution validates if command execution is possible
func (sv *View) validateCommandExecution(command string) bool {
	if sv.execFunc == nil {
		sv.addOutput("Error: command execution not available\n")
		return false
	}

	if sv.isInteractiveCommand(command) {
		sv.showInteractiveCommandWarning()
		return false
	}

	return true
}

// executeCommandWithArgs executes the command with the parsed arguments
func (sv *View) executeCommandWithArgs(args []string) {
	ctx := context.Background()
	output, err := sv.execFunc(ctx, sv.containerID, args, false)
	if err != nil {
		sv.addOutput(fmt.Sprintf("Error: %s\n", err.Error()))
		return
	}

	sv.handleCommandOutput(output)
}

// handleCommandOutput handles the command output
func (sv *View) handleCommandOutput(output string) {
	if output == "" {
		return
	}

	sv.addOutput(output)
	if !strings.HasSuffix(output, "\n") {
		sv.addOutput("\n")
	}
}

// parseCommandArgs parses a command string into arguments
func (sv *View) parseCommandArgs(command string) []string {
	if strings.Contains(command, "|") || strings.Contains(command, ">") ||
		strings.Contains(command, "<") || strings.Contains(command, "&&") ||
		strings.Contains(command, "||") {
		return []string{"/bin/sh", "-c", command}
	}

	return strings.Fields(command)
}

// showInteractiveCommandWarning displays a warning about interactive commands
func (sv *View) showInteractiveCommandWarning() {
	sv.showWarningHeader()
	sv.showFreezingCommands()
	sv.showAlternativeCommands()
}

// showWarningHeader shows the warning header
func (sv *View) showWarningHeader() {
	sv.addOutput("⚠️  Warning: This is an interactive command that will freeze the TUI.\n")
	sv.addOutput("   Interactive commands require a real terminal with TTY support.\n")
}

// showFreezingCommands shows commands that will freeze the TUI
func (sv *View) showFreezingCommands() {
	sv.addOutput("   Commands that will freeze:\n")
	sv.addOutput("   - top, htop, vim, nano, less, more\n")
	sv.addOutput("   - Any command that requires TTY input\n")
	sv.addOutput("   - Commands that update the screen continuously\n\n")
}

// showAlternativeCommands shows alternative non-interactive commands
func (sv *View) showAlternativeCommands() {
	sv.addOutput("   For non-interactive alternatives, try:\n")
	sv.addOutput("   - 'ps aux' instead of 'top'\n")
	sv.addOutput("   - 'cat' instead of 'less' or 'more'\n")
	sv.addOutput("   - 'ls -la' instead of interactive file managers\n")
	sv.addOutput("   - 'free -h' instead of interactive system monitors\n\n")
}

// isInteractiveCommand checks if a command will freeze the TUI
func (sv *View) isInteractiveCommand(command string) bool {
	words := strings.Fields(command)
	if len(words) == 0 {
		return false
	}

	baseCommand := strings.ToLower(words[0])

	// Commands that require TTY and will freeze the TUI
	interactiveCommands := []string{
		"top", "htop", "vim", "nano", "emacs", "less", "more",
		"watch", "man", "info", "ncurses", "dialog", "whiptail",
		"screen", "tmux", "byobu", "mosh", "ssh", "telnet",
		"ftp", "sftp", "ncftp", "lynx", "links", "w3m",
		"mysql", "psql", "sqlite3", "redis-cli", "mongo",
		"irb", "python", "node", "gdb", "lldb", "perf",
		"strace", "ltrace", "valgrind", "gprof",
	}

	for _, cmd := range interactiveCommands {
		if baseCommand == cmd {
			return true
		}
	}

	return false
}

// showHelp displays shell help information
func (sv *View) showHelp() {
	sv.showHelpHeader()
	sv.showBuiltInCommands()
	sv.showNavigationHelp()
	sv.showMultiLineHelp()
	sv.showInteractiveWarning()
	sv.showCommandExecutionInfo()
}

// showHelpHeader shows the help header
func (sv *View) showHelpHeader() {
	sv.addOutput("Shell Commands:\n")
	sv.addOutput("==============\n")
}

// showBuiltInCommands shows built-in command help
func (sv *View) showBuiltInCommands() {
	sv.addOutput("Built-in commands:\n")
	sv.addOutput("  exit, quit    - Exit shell and return to container view\n")
	sv.addOutput("  clear         - Clear shell output\n")
	sv.addOutput("  help          - Show this help message\n")
	sv.addOutput("\n")
}

// showNavigationHelp shows navigation help
func (sv *View) showNavigationHelp() {
	sv.addOutput("Navigation:\n")
	sv.addOutput("  Up/Down arrows - Navigate command history\n")
	sv.addOutput("  Tab           - Smart tab completion\n")
	sv.addOutput("  ESC           - Exit shell\n")
	sv.addOutput("\n")
}

// showMultiLineHelp shows multi-line command help
func (sv *View) showMultiLineHelp() {
	sv.addOutput("Multi-line commands:\n")
	sv.addOutput("  End line with \\ to continue on next line\n")
	sv.addOutput("  Example: echo 'Hello' \\\n")
	sv.addOutput("           && echo 'World'\n")
	sv.addOutput("\n")
}

// showInteractiveWarning shows interactive command warning
func (sv *View) showInteractiveWarning() {
	sv.addOutput("⚠️  Interactive Commands:\n")
	sv.addOutput("  Commands like 'top', 'vim', 'less' will freeze the TUI\n")
	sv.addOutput("  These require a real terminal with TTY support\n")
	sv.addOutput("  Use alternatives: 'ps aux' instead of 'top'\n")
	sv.addOutput("\n")
}

// showCommandExecutionInfo shows command execution information
func (sv *View) showCommandExecutionInfo() {
	sv.addOutput("Any other command will be executed in the container.\n\n")
}
