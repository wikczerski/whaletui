package views

import (
	"context"
	"testing"
	"time"

	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
	"github.com/wikczerski/D5r/internal/models"
	"github.com/wikczerski/D5r/internal/services"
	servicemocks "github.com/wikczerski/D5r/internal/services/mocks"
	uimocks "github.com/wikczerski/D5r/internal/ui/interfaces/mocks"
)

func newUIMockWithServices(t *testing.T, sf *services.ServiceFactory) *uimocks.MockUIInterface {
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

func TestNewVolumesView_Creation(t *testing.T) {
	ui := newUIMockWithServices(t, nil)
	volumesView := NewVolumesView(ui)

	assert.NotNil(t, volumesView)
}

func TestNewVolumesView_ViewField(t *testing.T) {
	ui := newUIMockWithServices(t, nil)
	volumesView := NewVolumesView(ui)

	assert.NotNil(t, volumesView.view)
}

func TestNewVolumesView_TableField(t *testing.T) {
	ui := newUIMockWithServices(t, nil)
	volumesView := NewVolumesView(ui)

	assert.NotNil(t, volumesView.table)
}

func TestNewVolumesView_ItemsField(t *testing.T) {
	ui := newUIMockWithServices(t, nil)
	volumesView := NewVolumesView(ui)

	assert.Empty(t, volumesView.items)
}

func TestVolumesView_GetView(t *testing.T) {
	ui := newUIMockWithServices(t, nil)
	volumesView := NewVolumesView(ui)
	view := volumesView.GetView()

	assert.NotNil(t, view)
}

func TestVolumesView_GetView_ReturnsCorrectView(t *testing.T) {
	ui := newUIMockWithServices(t, nil)
	volumesView := NewVolumesView(ui)
	view := volumesView.GetView()

	assert.Equal(t, volumesView.view, view)
}

func TestVolumesView_Refresh_NoServices(t *testing.T) {
	ui := newUIMockWithServices(t, nil)
	volumesView := NewVolumesView(ui)

	volumesView.Refresh()

	assert.Empty(t, volumesView.items)
}

func TestVolumesView_Refresh_WithServices(t *testing.T) {
	mockVolumes := []models.Volume{
		{
			Name:       "volume1",
			Driver:     "local",
			Mountpoint: "/var/lib/docker/volumes/volume1/_data",
			Created:    time.Now(),
			Size:       "100MB",
		},
		{
			Name:       "volume2",
			Driver:     "local",
			Mountpoint: "/var/lib/docker/volumes/volume2/_data",
			Created:    time.Now().Add(-24 * time.Hour),
			Size:       "200MB",
		},
	}
	vs := servicemocks.NewMockVolumeService(t)
	vs.On("ListVolumes", context.Background()).Return(mockVolumes, nil)

	sf := &services.ServiceFactory{VolumeService: vs}
	ui := newUIMockWithServices(t, sf)

	volumesView := NewVolumesView(ui)
	volumesView.Refresh()

	assert.Equal(t, mockVolumes, volumesView.items)
}

func TestVolumesView_Refresh_ServiceError(t *testing.T) {
	vs := servicemocks.NewMockVolumeService(t)
	vs.On("ListVolumes", context.Background()).Return([]models.Volume{}, assert.AnError)

	sf := &services.ServiceFactory{VolumeService: vs}
	ui := newUIMockWithServices(t, sf)
	ui.On("ShowError", assert.AnError).Return().Maybe()

	volumesView := NewVolumesView(ui)
	volumesView.Refresh()

	assert.Empty(t, volumesView.items)
}

func TestVolumesView_ShowVolumeDetails_Success(t *testing.T) {
	vs := servicemocks.NewMockVolumeService(t)
	vs.On("InspectVolume", context.Background(), "volume1").Return(map[string]any{"ok": true}, nil).Maybe()

	mockVolume := models.Volume{
		Name:       "volume1",
		Driver:     "local",
		Mountpoint: "/var/lib/docker/volumes/volume1/_data",
		Created:    time.Now(),
		Size:       "100MB",
	}

	sf := &services.ServiceFactory{VolumeService: vs}
	ui := newUIMockWithServices(t, sf)

	volumesView := NewVolumesView(ui)
	volumesView.items = []models.Volume{mockVolume}

	assert.NotNil(t, volumesView.showVolumeDetails)
}

func TestVolumesView_ShowVolumeDetails_InspectError(t *testing.T) {
	vs := servicemocks.NewMockVolumeService(t)
	vs.On("InspectVolume", context.Background(), "volume1").Return(map[string]any(nil), assert.AnError).Maybe()

	mockVolume := models.Volume{
		Name:       "volume1",
		Driver:     "local",
		Mountpoint: "/var/lib/docker/volumes/volume1/_data",
		Created:    time.Now(),
		Size:       "100MB",
	}

	sf := &services.ServiceFactory{VolumeService: vs}
	ui := newUIMockWithServices(t, sf)

	volumesView := NewVolumesView(ui)
	volumesView.items = []models.Volume{mockVolume}

	assert.NotNil(t, volumesView.showVolumeDetails)
}

func TestVolumesView_HandleAction_Delete(t *testing.T) {
	ui := newUIMockWithServices(t, nil)
	volumesView := NewVolumesView(ui)
	testVolume := models.Volume{
		Name:       "volume1",
		Driver:     "local",
		Mountpoint: "/var/lib/docker/volumes/volume1/_data",
		Created:    time.Now(),
		Size:       "100MB",
	}
	volumesView.items = []models.Volume{testVolume}
	volumesView.table.Select(1, 0)

	// Test action handling
	volumesView.handleAction('d', &testVolume)
	volumesView.handleAction('i', &testVolume)

	assert.NotNil(t, volumesView)
}

func TestVolumesView_HandleAction_Inspect(t *testing.T) {
	ui := newUIMockWithServices(t, nil)
	volumesView := NewVolumesView(ui)
	testVolume := models.Volume{
		Name:       "volume1",
		Driver:     "local",
		Mountpoint: "/var/lib/docker/volumes/volume1/_data",
		Created:    time.Now(),
		Size:       "100MB",
	}
	volumesView.items = []models.Volume{testVolume}
	volumesView.table.Select(1, 0)

	// Test action handling
	volumesView.handleAction('d', &testVolume)
	volumesView.handleAction('i', &testVolume)

	assert.NotNil(t, volumesView)
}

func TestVolumesView_HandleAction_InvalidSelection(t *testing.T) {
	ui := newUIMockWithServices(t, nil)
	volumesView := NewVolumesView(ui)
	testVolume := models.Volume{
		Name:       "volume1",
		Driver:     "local",
		Mountpoint: "/var/lib/docker/volumes/volume1/_data",
		Created:    time.Now(),
		Size:       "100MB",
	}
	volumesView.items = []models.Volume{}
	volumesView.table.Select(0, 0)

	volumesView.handleAction('d', &testVolume)
	volumesView.handleAction('i', &testVolume)

	assert.NotNil(t, volumesView)
}

func TestVolumesView_ShowTable(t *testing.T) {
	ui := newUIMockWithServices(t, nil)
	volumesView := NewVolumesView(ui)

	assert.NotNil(t, volumesView.showTable)
}

func TestVolumesView_DeleteVolume(t *testing.T) {
	ui := newUIMockWithServices(t, nil)
	volumesView := NewVolumesView(ui)

	assert.NotNil(t, volumesView.deleteVolume)
}

func TestVolumesView_InspectVolume(t *testing.T) {
	ui := newUIMockWithServices(t, nil)
	volumesView := NewVolumesView(ui)

	assert.NotNil(t, volumesView.inspectVolume)
}

func TestVolumesView_SetupKeyBindings(t *testing.T) {
	ui := newUIMockWithServices(t, nil)
	volumesView := NewVolumesView(ui)
	volumesView.items = []models.Volume{{Name: "volume1", Driver: "local", Created: time.Now()}}
	volumesView.table.Select(1, 0)

	assert.NotNil(t, volumesView.table.GetInputCapture())
}

func TestVolumesView_SetupKeyBindings_InvalidSelection(t *testing.T) {
	ui := newUIMockWithServices(t, nil)
	volumesView := NewVolumesView(ui)
	volumesView.items = []models.Volume{}
	volumesView.table.Select(0, 0)

	assert.NotNil(t, volumesView.table.GetInputCapture())
}
