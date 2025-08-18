package views

import (
	"testing"

	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
	uimocks "github.com/wikczerski/D5r/internal/ui/interfaces/mocks"
)

func newContainersUIMock(t *testing.T) *uimocks.MockUIInterface {
	ui := uimocks.NewMockUIInterface(t)
	ui.On("GetApp").Return(tview.NewApplication()).Maybe()
	ui.On("GetPages").Return(tview.NewPages()).Maybe()
	ui.On("GetServices").Return(nil).Maybe()
	return ui
}

func TestNewContainersView(t *testing.T) {
	ui := newContainersUIMock(t)
	containersView := NewContainersView(ui)

	assert.NotNil(t, containersView)
	assert.NotNil(t, containersView.view)
	assert.NotNil(t, containersView.table)
	assert.Empty(t, containersView.items)
}

func TestContainersView_GetView(t *testing.T) {
	ui := newContainersUIMock(t)
	containersView := NewContainersView(ui)
	view := containersView.GetView()

	assert.NotNil(t, view)
	assert.Equal(t, containersView.view, view)
}

func TestContainersView_Refresh_NoServices(t *testing.T) {
	ui := newContainersUIMock(t)
	containersView := NewContainersView(ui)
	containersView.Refresh()
	assert.Empty(t, containersView.items)
}

func TestContainersView_SetupKeyBindings(t *testing.T) {
	ui := newContainersUIMock(t)
	containersView := NewContainersView(ui)
	assert.NotNil(t, containersView.table.GetInputCapture())
}
