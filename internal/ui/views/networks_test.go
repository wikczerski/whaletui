package views

import (
	"context"
	"testing"
	"time"

	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/wikczerski/whaletui/internal/models"
	"github.com/wikczerski/whaletui/internal/services"
	servicemocks "github.com/wikczerski/whaletui/internal/services/mocks"
	uimocks "github.com/wikczerski/whaletui/internal/ui/interfaces/mocks"
)

func newNetworksUIMockWithServices(t *testing.T, sf *services.ServiceFactory) *uimocks.MockUIInterface {
	ui := uimocks.NewMockUIInterface(t)
	ui.On("GetApp").Return(tview.NewApplication()).Maybe()
	ui.On("GetPages").Return(tview.NewPages()).Maybe()

	if sf == nil {
		// Create a mock service factory that returns nil for all services
		mockSF := servicemocks.NewMockServiceFactoryInterface(t)
		mockSF.On("GetImageService").Return(nil).Maybe()
		mockSF.On("GetContainerService").Return(nil).Maybe()
		mockSF.On("GetVolumeService").Return(nil).Maybe()
		mockSF.On("GetNetworkService").Return(nil).Maybe()
		mockSF.On("GetDockerInfoService").Return(nil).Maybe()
		mockSF.On("GetLogsService").Return(nil).Maybe()
		ui.On("GetServices").Return(mockSF).Maybe()
	} else {
		ui.On("GetServices").Return(sf).Maybe()
	}

	return ui
}

func TestNewNetworksView_Creation(t *testing.T) {
	ui := newNetworksUIMockWithServices(t, nil)
	networksView := NewNetworksView(ui)

	assert.NotNil(t, networksView)
}

func TestNewNetworksView_ViewField(t *testing.T) {
	ui := newNetworksUIMockWithServices(t, nil)
	networksView := NewNetworksView(ui)

	assert.NotNil(t, networksView.view)
}

func TestNewNetworksView_TableField(t *testing.T) {
	ui := newNetworksUIMockWithServices(t, nil)
	networksView := NewNetworksView(ui)

	assert.NotNil(t, networksView.table)
}

func TestNewNetworksView_ItemsField(t *testing.T) {
	ui := newNetworksUIMockWithServices(t, nil)
	networksView := NewNetworksView(ui)

	assert.Empty(t, networksView.items)
}

func TestNetworksView_GetView(t *testing.T) {
	ui := newNetworksUIMockWithServices(t, nil)
	networksView := NewNetworksView(ui)
	view := networksView.GetView()

	assert.NotNil(t, view)
}

func TestNetworksView_GetView_ReturnsCorrectView(t *testing.T) {
	ui := newNetworksUIMockWithServices(t, nil)
	networksView := NewNetworksView(ui)
	view := networksView.GetView()

	assert.Equal(t, networksView.view, view)
}

func TestNetworksView_Refresh_NoServices(t *testing.T) {
	ui := newNetworksUIMockWithServices(t, nil)
	networksView := NewNetworksView(ui)

	networksView.Refresh()

	assert.Empty(t, networksView.items)
}

func TestNetworksView_Refresh_WithServices(t *testing.T) {
	mockNetworks := []models.Network{
		{
			ID:      "network1",
			Name:    "bridge",
			Driver:  "bridge",
			Scope:   "local",
			Created: time.Now(),
		},
		{
			ID:      "network2",
			Name:    "host",
			Driver:  "host",
			Scope:   "local",
			Created: time.Now().Add(-24 * time.Hour),
		},
	}

	ns := servicemocks.NewMockNetworkService(t)
	ns.On("ListNetworks", context.Background()).Return(mockNetworks, nil)

	sf := &services.ServiceFactory{NetworkService: ns}
	ui := newNetworksUIMockWithServices(t, sf)

	networksView := NewNetworksView(ui)
	networksView.Refresh()

	assert.Equal(t, mockNetworks, networksView.items)
}

func TestNetworksView_Refresh_ServiceError(t *testing.T) {
	ns := servicemocks.NewMockNetworkService(t)
	ns.On("ListNetworks", context.Background()).Return([]models.Network{}, assert.AnError)

	sf := &services.ServiceFactory{NetworkService: ns}
	ui := newNetworksUIMockWithServices(t, sf)
	ui.On("ShowError", assert.AnError).Return().Maybe()

	networksView := NewNetworksView(ui)
	networksView.Refresh()

	assert.Empty(t, networksView.items)
}

func TestNetworksView_ShowNetworkDetails_Success(t *testing.T) {
	ns := servicemocks.NewMockNetworkService(t)
	ns.On("InspectNetwork", context.Background(), "network1").Return(map[string]any{"ok": true}, nil).Maybe()

	mockNetwork := models.Network{
		ID:      "network1",
		Name:    "bridge",
		Driver:  "bridge",
		Scope:   "local",
		Created: time.Now(),
	}

	sf := &services.ServiceFactory{NetworkService: ns}
	ui := newNetworksUIMockWithServices(t, sf)
	ui.On("ShowDetails", mock.AnythingOfType("*tview.Flex")).Return().Maybe()

	networksView := NewNetworksView(ui)
	networksView.items = []models.Network{mockNetwork}

	assert.NotNil(t, networksView.showNetworkDetails)
}

func TestNetworksView_ShowNetworkDetails_InspectError(t *testing.T) {
	ns := servicemocks.NewMockNetworkService(t)
	ns.On("InspectNetwork", context.Background(), "network1").Return(map[string]any(nil), assert.AnError).Maybe()

	mockNetwork := models.Network{
		ID:      "network1",
		Name:    "bridge",
		Driver:  "bridge",
		Scope:   "local",
		Created: time.Now(),
	}

	sf := &services.ServiceFactory{NetworkService: ns}
	ui := newNetworksUIMockWithServices(t, sf)
	ui.On("ShowDetails", mock.AnythingOfType("*tview.Flex")).Return().Maybe()

	networksView := NewNetworksView(ui)
	networksView.items = []models.Network{mockNetwork}

	assert.NotNil(t, networksView.showNetworkDetails)
}

func TestNetworksView_HandleAction_Delete(t *testing.T) {
	ui := newNetworksUIMockWithServices(t, nil)
	networksView := NewNetworksView(ui)
	testNetwork := models.Network{
		ID:      "network1",
		Name:    "test-network",
		Driver:  "bridge",
		Scope:   "local",
		Created: time.Now(),
	}
	networksView.items = []models.Network{testNetwork}
	networksView.table.Select(1, 0)

	// Test action handling
	networksView.handleAction('d', &testNetwork)
	networksView.handleAction('i', &testNetwork)

	assert.NotNil(t, networksView)
}

func TestNetworksView_HandleAction_Inspect(t *testing.T) {
	ui := newNetworksUIMockWithServices(t, nil)
	networksView := NewNetworksView(ui)
	testNetwork := models.Network{
		ID:      "network1",
		Name:    "bridge",
		Driver:  "bridge",
		Scope:   "local",
		Created: time.Now(),
	}
	networksView.items = []models.Network{testNetwork}
	networksView.table.Select(1, 0)

	// Test action handling
	networksView.handleAction('d', &testNetwork)
	networksView.handleAction('i', &testNetwork)

	assert.NotNil(t, networksView)
}

func TestNetworksView_HandleAction_InvalidSelection(t *testing.T) {
	ui := newNetworksUIMockWithServices(t, nil)
	networksView := NewNetworksView(ui)
	testNetwork := models.Network{
		ID:      "network1",
		Name:    "test-network",
		Driver:  "bridge",
		Scope:   "local",
		Created: time.Now(),
	}
	networksView.items = []models.Network{}
	networksView.table.Select(0, 0)

	networksView.handleAction('d', &testNetwork)
	networksView.handleAction('i', &testNetwork)

	assert.NotNil(t, networksView)
}

func TestNetworksView_ShowTable(t *testing.T) {
	ui := newNetworksUIMockWithServices(t, nil)
	networksView := NewNetworksView(ui)

	assert.NotNil(t, networksView.showTable)
}

func TestNetworksView_DeleteNetwork(t *testing.T) {
	ui := newNetworksUIMockWithServices(t, nil)
	networksView := NewNetworksView(ui)

	assert.NotNil(t, networksView.deleteNetwork)
}

func TestNetworksView_InspectNetwork(t *testing.T) {
	ui := newNetworksUIMockWithServices(t, nil)
	networksView := NewNetworksView(ui)

	assert.NotNil(t, networksView.inspectNetwork)
}

func TestNetworksView_SetupKeyBindings(t *testing.T) {
	ui := newNetworksUIMockWithServices(t, nil)
	networksView := NewNetworksView(ui)
	networksView.items = []models.Network{
		{
			ID:      "network1",
			Name:    "bridge",
			Driver:  "bridge",
			Scope:   "local",
			Created: time.Now(),
		},
	}
	networksView.table.Select(1, 0)

	assert.NotNil(t, networksView.table.GetInputCapture())
}

func TestNetworksView_SetupKeyBindings_InvalidSelection(t *testing.T) {
	ui := newNetworksUIMockWithServices(t, nil)
	networksView := NewNetworksView(ui)
	networksView.items = []models.Network{}
	networksView.table.Select(0, 0)

	assert.NotNil(t, networksView.table.GetInputCapture())
}
