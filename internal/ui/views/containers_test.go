package views

import (
	"testing"

	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	servicemocks "github.com/wikczerski/whaletui/internal/services/mocks"
	uimocks "github.com/wikczerski/whaletui/internal/ui/interfaces/mocks"
)

func newContainersUIMock(t *testing.T) *uimocks.MockUIInterface {
	ui := uimocks.NewMockUIInterface(t)
	ui.On("GetApp").Return(tview.NewApplication()).Maybe()
	ui.On("GetPages").Return(tview.NewPages()).Maybe()

	// Create a mock service factory that returns nil for all services
	mockSF := servicemocks.NewMockServiceFactoryInterface(t)
	mockSF.On("GetContainerService").Return(nil).Maybe()
	mockSF.On("GetImageService").Return(nil).Maybe()
	mockSF.On("GetVolumeService").Return(nil).Maybe()
	mockSF.On("GetNetworkService").Return(nil).Maybe()
	mockSF.On("GetDockerInfoService").Return(nil).Maybe()
	mockSF.On("GetLogsService").Return(nil).Maybe()
	mockSF.On("IsServiceAvailable", mock.AnythingOfType("string")).Return(false).Maybe()
	mockSF.On("IsContainerServiceAvailable").Return(false).Maybe()

	ui.On("GetServices").Return(mockSF).Maybe()
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
