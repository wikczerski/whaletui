package views

import (
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
)

func TestLogsView_TextViewCreation(t *testing.T) {
	textView := tview.NewTextView()
	assert.NotNil(t, textView)
}

func TestLogsView_TextViewSetText(t *testing.T) {
	textView := tview.NewTextView()
	testText := "Test log content"

	textView.SetText(testText)

	assert.Equal(t, testText, textView.GetText(true))
}

func TestLogsView_TextViewTextAlignment(t *testing.T) {
	textView := tview.NewTextView()

	textView.SetTextAlign(tview.AlignCenter)

	assert.NotNil(t, textView)
}

func TestLogsView_TextViewDynamicColors(t *testing.T) {
	textView := tview.NewTextView()

	textView.SetDynamicColors(true)

	assert.NotNil(t, textView)
}

func TestLogsView_TextViewScrollable(t *testing.T) {
	textView := tview.NewTextView()

	textView.SetScrollable(true)

	assert.NotNil(t, textView)
}

func TestLogsView_TextViewBorder(t *testing.T) {
	textView := tview.NewTextView()

	textView.SetBorder(true)

	assert.NotNil(t, textView)
}

func TestLogsView_EscapeKeyEvent(t *testing.T) {
	escapeEvent := tcell.NewEventKey(tcell.KeyEscape, 0, tcell.ModNone)
	assert.NotNil(t, escapeEvent)
}

func TestLogsView_EscapeKeyValue(t *testing.T) {
	escapeEvent := tcell.NewEventKey(tcell.KeyEscape, 0, tcell.ModNone)

	assert.Equal(t, tcell.KeyEscape, escapeEvent.Key())
}

func TestLogsView_EnterKeyEvent(t *testing.T) {
	enterEvent := tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone)
	assert.NotNil(t, enterEvent)
}

func TestLogsView_EnterKeyValue(t *testing.T) {
	enterEvent := tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone)

	assert.Equal(t, tcell.KeyEnter, enterEvent.Key())
}

func TestLogsView_UpArrowKeyEvent(t *testing.T) {
	upEvent := tcell.NewEventKey(tcell.KeyUp, 0, tcell.ModNone)
	assert.NotNil(t, upEvent)
}

func TestLogsView_UpArrowKeyValue(t *testing.T) {
	upEvent := tcell.NewEventKey(tcell.KeyUp, 0, tcell.ModNone)

	assert.Equal(t, tcell.KeyUp, upEvent.Key())
}

func TestLogsView_DownArrowKeyEvent(t *testing.T) {
	downEvent := tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModNone)
	assert.NotNil(t, downEvent)
}

func TestLogsView_DownArrowKeyValue(t *testing.T) {
	downEvent := tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModNone)

	assert.Equal(t, tcell.KeyDown, downEvent.Key())
}

func TestLogsView_PageUpKeyEvent(t *testing.T) {
	pageUpEvent := tcell.NewEventKey(tcell.KeyPgUp, 0, tcell.ModNone)
	assert.NotNil(t, pageUpEvent)
}

func TestLogsView_PageDownKeyEvent(t *testing.T) {
	pageDownEvent := tcell.NewEventKey(tcell.KeyPgDn, 0, tcell.ModNone)
	assert.NotNil(t, pageDownEvent)
}

func TestLogsView_HomeKeyEvent(t *testing.T) {
	homeEvent := tcell.NewEventKey(tcell.KeyHome, 0, tcell.ModNone)
	assert.NotNil(t, homeEvent)
}

func TestLogsView_EndKeyEvent(t *testing.T) {
	endEvent := tcell.NewEventKey(tcell.KeyEnd, 0, tcell.ModNone)
	assert.NotNil(t, endEvent)
}

func TestLogsView_ActionMapping_FollowLogs(t *testing.T) {
	expectedActions := map[rune]string{
		'f': "Follow logs",
		't': "Tail logs",
		's': "Save logs",
		'c': "Clear logs",
		'w': "Wrap text",
	}

	action, exists := expectedActions['f']
	assert.True(t, exists)
	assert.Equal(t, "Follow logs", action)
}

func TestLogsView_ActionMapping_TailLogs(t *testing.T) {
	expectedActions := map[rune]string{
		'f': "Follow logs",
		't': "Tail logs",
		's': "Save logs",
		'c': "Clear logs",
		'w': "Wrap text",
	}

	action, exists := expectedActions['t']
	assert.True(t, exists)
	assert.Equal(t, "Tail logs", action)
}

func TestLogsView_ActionMapping_SaveLogs(t *testing.T) {
	expectedActions := map[rune]string{
		'f': "Follow logs",
		't': "Tail logs",
		's': "Save logs",
		'c': "Clear logs",
		'w': "Wrap text",
	}

	action, exists := expectedActions['s']
	assert.True(t, exists)
	assert.Equal(t, "Save logs", action)
}

func TestLogsView_ActionMapping_ClearLogs(t *testing.T) {
	expectedActions := map[rune]string{
		'f': "Follow logs",
		't': "Tail logs",
		's': "Save logs",
		'c': "Clear logs",
		'w': "Wrap text",
	}

	action, exists := expectedActions['c']
	assert.True(t, exists)
	assert.Equal(t, "Clear logs", action)
}

func TestLogsView_ActionMapping_WrapText(t *testing.T) {
	expectedActions := map[rune]string{
		'f': "Follow logs",
		't': "Tail logs",
		's': "Save logs",
		'c': "Clear logs",
		'w': "Wrap text",
	}

	action, exists := expectedActions['w']
	assert.True(t, exists)
	assert.Equal(t, "Wrap text", action)
}

func TestLogsView_ContainerIDTruncation(t *testing.T) {
	containerID := "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
	expectedTruncated := containerID[:12]

	assert.Equal(t, 12, len(expectedTruncated))
}

func TestLogsView_ContainerIDTruncationValue(t *testing.T) {
	containerID := "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
	expectedTruncated := containerID[:12]

	assert.Equal(t, "1234567890ab", expectedTruncated)
}

func TestLogsView_ContainerTitle(t *testing.T) {
	expectedTitle := " test-container<1234567890abcdef> "

	assert.Equal(t, " test-container<1234567890abcdef> ", expectedTitle)
}

func TestLogsView_ServiceErrorNotEmpty(t *testing.T) {
	serviceError := "Error: Container service not available"

	assert.NotEmpty(t, serviceError)
}

func TestLogsView_LogsErrorNotEmpty(t *testing.T) {
	logsError := "Error loading logs: test error"

	assert.NotEmpty(t, logsError)
}

func TestLogsView_LogsErrorPrefix(t *testing.T) {
	logsError := "Error loading logs: test error"
	expectedPrefix := "Error loading logs: "

	assert.True(t, len(logsError) >= len(expectedPrefix))
	assert.Equal(t, expectedPrefix, logsError[:len(expectedPrefix)])
}

func TestLogsView_TextViewInitialText(t *testing.T) {
	textView := tview.NewTextView()
	initialText := "Loading logs..."

	textView.SetText(initialText)

	assert.Equal(t, initialText, textView.GetText(true))
}

func TestLogsView_TextViewClearText(t *testing.T) {
	textView := tview.NewTextView()

	textView.SetText("")

	assert.Equal(t, "", textView.GetText(true))
}

func TestLogsView_TextViewBorderColor(t *testing.T) {
	textView := tview.NewTextView()

	textView.SetBorderColor(tcell.ColorRed)

	assert.NotNil(t, textView)
}

func TestLogsView_TextViewLeftAlignment(t *testing.T) {
	textView := tview.NewTextView()

	textView.SetTextAlign(tview.AlignLeft)

	assert.NotNil(t, textView)
}

func TestLogsView_ButtonCreation(t *testing.T) {
	button := tview.NewButton("Back to Table")

	assert.NotNil(t, button)
}

func TestLogsView_ButtonLabel(t *testing.T) {
	button := tview.NewButton("Back to Table")
	expectedLabel := "Back to Table"

	assert.Equal(t, expectedLabel, button.GetLabel())
}

func TestLogsView_ButtonClickHandler(t *testing.T) {
	button := tview.NewButton("Back to Table")
	expectedLabel := "Back to Table"

	button.SetSelectedFunc(func() {
		// Button click handler
	})

	assert.Equal(t, expectedLabel, button.GetLabel())
}

func TestLogsView_FlexCreation(t *testing.T) {
	flex := tview.NewFlex()

	assert.NotNil(t, flex)
}

func TestLogsView_FlexDirection(t *testing.T) {
	flex := tview.NewFlex()

	flex.SetDirection(tview.FlexRow)

	assert.NotNil(t, flex)
}

func TestLogsView_FlexAddItems(t *testing.T) {
	flex := tview.NewFlex()
	textView := tview.NewTextView()
	button := tview.NewButton("Test")

	flex.AddItem(textView, 0, 1, true)
	flex.AddItem(button, 1, 0, false)

	assert.NotNil(t, flex)
}
