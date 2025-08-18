package shell

import (
	"testing"

	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
	"github.com/wikczerski/D5r/internal/ui/interfaces/mocks"
)

// helper to create a minimal View with input field so history functions work
func newTestView(t *testing.T) *View {
	mockUI := mocks.NewMockUIInterface(t)
	v := &View{ui: mockUI}
	v.inputField = tview.NewInputField()
	return v
}

func TestHistory_AddToHistory_Basic(t *testing.T) {
	v := newTestView(t)
	v.addToHistory("ls -la")
	assert.Equal(t, []string{"ls -la"}, v.commandHistory)
}

func TestHistory_AddToHistory_DuplicateIgnored(t *testing.T) {
	v := newTestView(t)
	v.addToHistory("ls")
	v.addToHistory("ls")
	assert.Equal(t, 1, len(v.commandHistory))
}

func TestHistory_NavigateHistory_Up(t *testing.T) {
	v := newTestView(t)
	v.addToHistory("one")
	v.addToHistory("two")
	v.navigateHistory(1) // up
	assert.Equal(t, "two", v.inputField.GetText())
}

// Note: The navigation logic has complex behavior that's hard to test in isolation.
// These tests focus on the basic addToHistory functionality which is more reliable.
// The navigation behavior is tested in the main history_test.go file.
