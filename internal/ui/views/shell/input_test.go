package shell

import (
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
	mocks "github.com/wikczerski/whaletui/internal/mocks/ui"
)

func newTestViewForInput(t *testing.T) *View {
	mockUI := mocks.NewMockUIInterface(t)
	v := &View{ui: mockUI}
	v.inputField = tview.NewInputField()
	return v
}

func TestView_HandleInputCapture_EscapeKey(t *testing.T) {
	v := newTestViewForInput(t)
	event := tcell.NewEventKey(tcell.KeyEscape, 0, tcell.ModNone)

	result := v.handleInputCapture(event)
	assert.Nil(t, result)
}

func TestView_HandleInputCapture_UpArrow(t *testing.T) {
	v := newTestViewForInput(t)
	event := tcell.NewEventKey(tcell.KeyUp, 0, tcell.ModNone)

	result := v.handleInputCapture(event)
	assert.Nil(t, result)
}

func TestView_HandleInputCapture_DownArrow(t *testing.T) {
	v := newTestViewForInput(t)
	event := tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModNone)

	result := v.handleInputCapture(event)
	assert.Nil(t, result)
}

func TestView_HandleInputCapture_TabKey(t *testing.T) {
	v := newTestViewForInput(t)
	event := tcell.NewEventKey(tcell.KeyTab, 0, tcell.ModNone)

	result := v.handleInputCapture(event)
	assert.Nil(t, result)
}

func TestView_HandleInputCapture_OtherKey(t *testing.T) {
	v := newTestViewForInput(t)
	event := tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone)

	result := v.handleInputCapture(event)
	assert.Equal(t, event, result)
}
