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
	if key == tcell.KeyEnter {
		command := sv.inputField.GetText()
		if command == "" {
			return
		}

		if sv.isMultiLineCommand(command) {
			sv.addToMultiLineBuffer(command)
			sv.addOutput(fmt.Sprintf("> %s\n", command))
			sv.inputField.SetText("")
			return
		}

		if len(sv.multiLineBuffer) > 0 {
			sv.addToMultiLineBuffer(command)
			command = sv.getMultiLineCommand()
			sv.addOutput(fmt.Sprintf("$ %s\n", command))
			sv.clearMultiLineBuffer()
		} else {
			sv.addOutput(fmt.Sprintf("$ %s\n", command))
		}

		sv.currentInput = ""
		sv.historyIndex = len(sv.commandHistory)

		sv.addToHistory(command)

		if sv.handleBuiltInCommand(command) {
			sv.inputField.SetText("")
			return
		}

		sv.executeDockerCommand(command)

		sv.inputField.SetText("")
	}
}

// handleBuiltInCommand handles built-in shell commands
func (sv *View) handleBuiltInCommand(command string) bool {
	switch command {
	case "exit", "quit":
		sv.addOutput("Exiting shell...\n")
		sv.exitShell()
		return true
	case "clear":
		sv.outputView.Clear()
		sv.addOutput(fmt.Sprintf("Welcome to shell for container: %s (%s)\n", sv.containerName, shared.TruncName(sv.containerID, 12)))
		sv.addOutput("Type 'exit' or press ESC to return to container view\n\n")
		return true
	case "help":
		sv.showHelp()
		return true
	}
	return false
}

// executeDockerCommand executes a command in the container
func (sv *View) executeDockerCommand(command string) {
	if sv.execFunc == nil {
		sv.addOutput("Error: command execution not available\n")
		return
	}

	if sv.isInteractiveCommand(command) {
		sv.showInteractiveCommandWarning()
		return
	}

	args := sv.parseCommandArgs(command)
	if len(args) == 0 {
		sv.addOutput("Error: invalid command\n")
		return
	}

	ctx := context.Background()
	output, err := sv.execFunc(ctx, sv.containerID, args, false)
	if err != nil {
		sv.addOutput(fmt.Sprintf("Error: %s\n", err.Error()))
	} else if output != "" {
		sv.addOutput(output)
		if !strings.HasSuffix(output, "\n") {
			sv.addOutput("\n")
		}
	}
}

// parseCommandArgs parses a command string into arguments
func (sv *View) parseCommandArgs(command string) []string {
	if strings.Contains(command, "|") || strings.Contains(command, ">") || strings.Contains(command, "<") || strings.Contains(command, "&&") || strings.Contains(command, "||") {
		return []string{"/bin/sh", "-c", command}
	}

	return strings.Fields(command)
}

// showInteractiveCommandWarning displays a warning about interactive commands
func (sv *View) showInteractiveCommandWarning() {
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
