package shell

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestView_InitialState(t *testing.T) {
	view := createTestView()
	assert.NotNil(t, view.commandHistory)
}

func TestView_InitialState_Empty(t *testing.T) {
	view := createTestView()
	assert.Empty(t, view.commandHistory)
}

func TestView_InitialState_HistoryIndex(t *testing.T) {
	view := createTestView()
	assert.Equal(t, 0, view.historyIndex)
}

func TestView_AddToHistory_Single(t *testing.T) {
	view := createTestView()
	command := "ls -la"

	view.addToHistory(command)
	assert.Len(t, view.commandHistory, 1)
}

func TestView_AddToHistory_FirstCommand(t *testing.T) {
	view := createTestView()
	command := "ls -la"

	view.addToHistory(command)
	assert.Equal(t, command, view.commandHistory[0])
}

func TestView_AddToHistory_Multiple(t *testing.T) {
	view := createTestView()
	commands := []string{"ls -la", "cd /tmp", "pwd"}

	for _, cmd := range commands {
		view.addToHistory(cmd)
	}
	assert.Len(t, view.commandHistory, len(commands))
}

func TestView_AddToHistory_Duplicate(t *testing.T) {
	view := createTestView()
	command := "ls -la"

	view.addToHistory(command)
	view.addToHistory(command)
	assert.Len(t, view.commandHistory, 1)
}

func TestView_AddToHistory_Duplicate_Order(t *testing.T) {
	view := createTestView()
	command := "ls -la"

	view.addToHistory(command)
	view.addToHistory("cd /tmp")
	view.addToHistory(command)
	assert.Equal(t, command, view.commandHistory[0])
}

func TestView_AddToHistory_Empty(t *testing.T) {
	view := createTestView()
	command := ""

	view.addToHistory(command)
	assert.Len(t, view.commandHistory, 0)
}

func TestView_AddToHistory_Whitespace(t *testing.T) {
	view := createTestView()
	command := "   "

	view.addToHistory(command)
	// The actual implementation only checks for empty strings, not whitespace-only
	// So whitespace-only strings are added to history
	assert.Len(t, view.commandHistory, 1)
	assert.Equal(t, "   ", view.commandHistory[0])
}

func TestView_AddToHistory_HistoryIndex(t *testing.T) {
	view := createTestView()
	view.addToHistory("ls -la")

	assert.Equal(t, 1, view.historyIndex)
}

func TestView_NavigateHistory_Up(t *testing.T) {
	view := createTestView()
	view.commandHistory = []string{"ls -la", "cd /tmp", "pwd"}
	view.historyIndex = 3

	view.navigateHistory(1)
	assert.Equal(t, 2, view.historyIndex)
}

func TestView_NavigateHistory_Up_AtBeginning(t *testing.T) {
	view := createTestView()
	view.commandHistory = []string{"ls -la", "cd /tmp", "pwd"}
	view.historyIndex = 0

	view.navigateHistory(1)
	assert.Equal(t, 0, view.historyIndex)
}

func TestView_NavigateHistory_Down(t *testing.T) {
	view := createTestView()
	view.commandHistory = []string{"ls -la", "cd /tmp", "pwd"}
	view.historyIndex = 1

	view.navigateHistory(0)
	assert.Equal(t, 2, view.historyIndex)
}

func TestView_NavigateHistory_Down_AtEnd(t *testing.T) {
	view := createTestView()
	view.commandHistory = []string{"ls -la", "cd /tmp", "pwd"}
	view.historyIndex = 3

	view.navigateHistory(0)
	assert.Equal(t, 3, view.historyIndex)
}

func TestView_NavigateHistory_Empty(t *testing.T) {
	view := createTestView()
	view.commandHistory = []string{}
	view.historyIndex = 0

	view.navigateHistory(1)
	assert.Equal(t, 0, view.historyIndex)
}

func TestView_NavigateHistory_Empty_Down(t *testing.T) {
	view := createTestView()
	view.commandHistory = []string{}
	view.historyIndex = 0

	view.navigateHistory(0)
	assert.Equal(t, 0, view.historyIndex)
}

func TestView_NavigateHistory_InputFieldText_Current(t *testing.T) {
	view := createTestView()
	view.commandHistory = []string{"ls -la", "cd /tmp", "pwd"}
	view.historyIndex = 3
	view.currentInput = "current command"

	view.navigateHistory(1)
	assert.Equal(t, "pwd", view.inputField.GetText())
}

func TestView_NavigateHistory_InputFieldText_History(t *testing.T) {
	view := createTestView()
	view.commandHistory = []string{"ls -la", "cd /tmp", "pwd"}
	view.historyIndex = 1

	view.navigateHistory(0)
	assert.Equal(t, "pwd", view.inputField.GetText())
}

func TestView_AddToHistory_LongCommand(t *testing.T) {
	view := createTestView()
	longCommand := "echo '" + string(make([]byte, 1000)) + "'"

	view.addToHistory(longCommand)
	assert.Len(t, view.commandHistory, 1)
}

func TestView_AddToHistory_SpecialChars(t *testing.T) {
	view := createTestView()
	specialCommand := "echo 'Hello World!'"

	view.addToHistory(specialCommand)
	assert.Equal(t, specialCommand, view.commandHistory[0])
}

func TestView_AddToHistory_Unicode(t *testing.T) {
	view := createTestView()
	unicodeCommand := "echo 世界"

	view.addToHistory(unicodeCommand)
	assert.Equal(t, unicodeCommand, view.commandHistory[0])
}
