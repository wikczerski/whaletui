package shell

import (
	"context"
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
)

// handleCommand processes the entered command
func (sv *ShellView) handleCommand(key tcell.Key) {
	if key == tcell.KeyEnter {
		command := sv.inputField.GetText()
		if command == "" {
			return
		}

		// Check if this is a multi-line command continuation
		if sv.isMultiLineCommand(command) {
			sv.addToMultiLineBuffer(command)
			sv.addOutput(fmt.Sprintf("> %s\n", command))
			sv.inputField.SetText("")
			return
		}

		// If we have a multi-line buffer, add the final line and execute
		if len(sv.multiLineBuffer) > 0 {
			sv.addToMultiLineBuffer(command)
			command = sv.getMultiLineCommand()
			sv.addOutput(fmt.Sprintf("$ %s\n", command))
			sv.clearMultiLineBuffer()
		} else {
			// Single line command
			sv.addOutput(fmt.Sprintf("$ %s\n", command))
		}

		// Store current input before clearing
		sv.currentInput = ""
		sv.historyIndex = len(sv.commandHistory)

		// Add command to history
		sv.addToHistory(command)

		// Handle built-in commands
		if sv.handleBuiltInCommand(command) {
			sv.inputField.SetText("")
			return
		}

		// Execute command using Docker exec
		sv.executeDockerCommand(command)

		// Clear input field
		sv.inputField.SetText("")
	}
}

// handleBuiltInCommand handles built-in shell commands
func (sv *ShellView) handleBuiltInCommand(command string) bool {
	switch command {
	case "exit", "quit":
		sv.addOutput("Exiting shell...\n")
		sv.exitShell()
		return true

	case "clear":
		sv.outputView.Clear()
		sv.addOutput(fmt.Sprintf("Welcome to shell for container: %s (%s)\n", sv.containerName, sv.containerID[:12]))
		sv.addOutput("Type 'exit' or press ESC to return to container view\n\n")
		return true
	case "help":
		sv.showHelp()
		return true
	}
	return false
}

// executeDockerCommand executes a command in the container
func (sv *ShellView) executeDockerCommand(command string) {
	if sv.execFunc == nil {
		sv.addOutput("Error: command execution not available\n")
		return
	}

	// Check if this is an interactive command that will freeze the TUI
	if sv.isInteractiveCommand(command) {
		sv.showInteractiveCommandWarning()
		return
	}

	// Parse command into arguments
	args := sv.parseCommandArgs(command)
	if len(args) == 0 {
		sv.addOutput("Error: invalid command\n")
		return
	}

	// Execute the command
	ctx := context.Background()
	output, err := sv.execFunc(ctx, sv.containerID, args, false)
	if err != nil {
		sv.addOutput(fmt.Sprintf("Error: %s\n", err.Error()))
	} else {
		// Display command output
		if output != "" {
			sv.addOutput(output)
			if !strings.HasSuffix(output, "\n") {
				sv.addOutput("\n")
			}
		}
	}
}

// parseCommandArgs parses a command string into arguments
func (sv *ShellView) parseCommandArgs(command string) []string {
	if strings.Contains(command, "|") || strings.Contains(command, ">") || strings.Contains(command, "<") || strings.Contains(command, "&&") || strings.Contains(command, "||") {
		// Execute through shell to handle pipes, redirects, etc.
		return []string{"/bin/sh", "-c", command}
	}

	// Parse simple command into arguments
	return strings.Fields(command)
}

// showInteractiveCommandWarning displays a warning about interactive commands
func (sv *ShellView) showInteractiveCommandWarning() {
	sv.addOutput("⚠️  Warning: This is an interactive command that will freeze the TUI.\n")
	sv.addOutput("   Interactive commands require a real terminal with TTY support.\n")
	sv.addOutput("   Commands that will freeze:\n")
	sv.addOutput("   - top, htop, vim, nano, less, more\n")
	sv.addOutput("   - Any command that requires TTY input\n")
	sv.addOutput("   - Commands that update the screen continuously\n\n")
	sv.addOutput("   For non-interactive alternatives, try:\n")
	sv.addOutput("   - 'ps aux' instead of 'top'\n")
	sv.addOutput("   - 'cat' instead of 'less' or 'more'\n")
	sv.addOutput("   - 'ls -la' instead of interactive file managers\n")
	sv.addOutput("   - 'free -h' instead of interactive system monitors\n\n")
}

// isInteractiveCommand checks if a command will freeze the TUI
func (sv *ShellView) isInteractiveCommand(command string) bool {
	// Get the base command (first word)
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
func (sv *ShellView) showHelp() {
	sv.addOutput("Shell Commands:\n")
	sv.addOutput("==============\n")
	sv.addOutput("Built-in commands:\n")
	sv.addOutput("  exit, quit    - Exit shell and return to container view\n")
	sv.addOutput("  clear         - Clear shell output\n")
	sv.addOutput("  help          - Show this help message\n")
	sv.addOutput("\n")
	sv.addOutput("Navigation:\n")
	sv.addOutput("  Up/Down arrows - Navigate command history\n")
	sv.addOutput("  Tab           - Smart tab completion\n")
	sv.addOutput("  ESC           - Exit shell\n")
	sv.addOutput("\n")
	sv.addOutput("Multi-line commands:\n")
	sv.addOutput("  End line with \\ to continue on next line\n")
	sv.addOutput("  Example: echo 'Hello' \\\n")
	sv.addOutput("           && echo 'World'\n")
	sv.addOutput("\n")
	sv.addOutput("⚠️  Interactive Commands:\n")
	sv.addOutput("  Commands like 'top', 'vim', 'less' will freeze the TUI\n")
	sv.addOutput("  These require a real terminal with TTY support\n")
	sv.addOutput("  Use alternatives: 'ps aux' instead of 'top'\n")
	sv.addOutput("\n")
	sv.addOutput("Any other command will be executed in the container.\n\n")
}
