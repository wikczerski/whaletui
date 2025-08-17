package shell

// addToHistory adds a command to the history
func (sv *ShellView) addToHistory(command string) {
	if command == "" || (len(sv.commandHistory) > 0 && sv.commandHistory[len(sv.commandHistory)-1] == command) {
		return
	}

	sv.commandHistory = append(sv.commandHistory, command)
	sv.historyIndex = len(sv.commandHistory)
}

// navigateHistory navigates through command history
func (sv *ShellView) navigateHistory(direction int) {
	if len(sv.commandHistory) == 0 {
		return
	}

	if direction > 0 { // Up arrow - go back in history
		if sv.historyIndex > 0 {
			sv.historyIndex--
		}
	} else { // Down arrow - go forward in history
		if sv.historyIndex < len(sv.commandHistory) {
			sv.historyIndex++
		}
	}

	// Set the input field text based on history position
	if sv.historyIndex == len(sv.commandHistory) {
		sv.inputField.SetText(sv.currentInput)
	} else if sv.historyIndex >= 0 {
		sv.inputField.SetText(sv.commandHistory[sv.historyIndex])
	}
}
