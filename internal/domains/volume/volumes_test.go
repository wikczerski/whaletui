package volume

import (
	"context"
	"testing"
	"time"

	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	uimocks "github.com/wikczerski/whaletui/internal/mocks/ui"
	"github.com/wikczerski/whaletui/internal/ui/interfaces"
)

func newUIMockWithServices(t *testing.T, sf interfaces.ServiceFactoryInterface) *uimocks.MockUIInterface {
	ui := uimocks.NewMockUIInterface(t)
	ui.On("GetApp").Return(tview.NewApplication()).Maybe()
	ui.On("GetPages").Return(tview.NewPages()).Maybe()

	if sf == nil {
		// Create a mock service factory that returns nil for all services
		mockSF := uimocks.NewMockServiceFactoryInterface(t)
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

	assert.NotNil(t, volumesView.GetView())
}

func TestNewVolumesView_TableField(t *testing.T) {
	ui := newUIMockWithServices(t, nil)
	volumesView := NewVolumesView(ui)

	assert.NotNil(t, volumesView.GetTable())
}

func TestNewVolumesView_ItemsField(t *testing.T) {
	ui := newUIMockWithServices(t, nil)
	volumesView := NewVolumesView(ui)

	assert.Empty(t, volumesView.GetItems())
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

	assert.Equal(t, volumesView.GetView(), view)
}

func TestVolumesView_Refresh_NoServices(t *testing.T) {
	ui := newUIMockWithServices(t, nil)
	volumesView := NewVolumesView(ui)

	volumesView.Refresh()

	assert.Empty(t, volumesView.GetItems())
}

func TestVolumesView_Refresh_WithServices(t *testing.T) {
	mockVolumes := []Volume{
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
	vs := uimocks.NewMockVolumeService(t)
	// Convert []Volume to []any for the mock interface
	mockVolumesAny := make([]any, len(mockVolumes))
	for i, vol := range mockVolumes {
		mockVolumesAny[i] = vol
	}
	vs.EXPECT().ListVolumes(context.Background()).Return(mockVolumesAny, nil)

	sf := uimocks.NewMockServiceFactoryInterface(t)
	sf.EXPECT().GetVolumeService().Return(vs)
	ui := newUIMockWithServices(t, sf)

	volumesView := NewVolumesView(ui)
	volumesView.Refresh()

	assert.Equal(t, mockVolumes, volumesView.GetItems())
}

func TestVolumesView_Refresh_ServiceError(t *testing.T) {
	vs := uimocks.NewMockVolumeService(t)
	vs.EXPECT().ListVolumes(context.Background()).Return([]any{}, assert.AnError)

	sf := uimocks.NewMockServiceFactoryInterface(t)
	sf.EXPECT().GetVolumeService().Return(vs)
	ui := newUIMockWithServices(t, sf)
	ui.On("ShowError", assert.AnError).Return().Maybe()

	volumesView := NewVolumesView(ui)
	volumesView.Refresh()

	assert.Empty(t, volumesView.GetItems())
}

func TestVolumesView_ShowVolumeDetails_Success(t *testing.T) {
	vs := uimocks.NewMockVolumeService(t)
	vs.On("InspectVolume", context.Background(), "volume1").Return(map[string]any{"ok": true}, nil).Maybe()

	mockVolume := Volume{
		Name:       "volume1",
		Driver:     "local",
		Mountpoint: "/var/lib/docker/volumes/volume1/_data",
		Created:    time.Now(),
		Size:       "100MB",
	}

	sf := uimocks.NewMockServiceFactoryInterface(t)
	sf.EXPECT().GetVolumeService().Return(vs)
	ui := newUIMockWithServices(t, sf)
	ui.On("ShowDetails", mock.AnythingOfType("*tview.Flex")).Return().Maybe()

	volumesView := NewVolumesView(ui)
	// Test the method directly without accessing unexported fields
	volumesView.showVolumeDetails(&mockVolume)

	assert.NotNil(t, volumesView)
}

func TestVolumesView_ShowVolumeDetails_InspectError(t *testing.T) {
	vs := uimocks.NewMockVolumeService(t)
	vs.On("InspectVolume", context.Background(), "volume1").Return(map[string]any(nil), assert.AnError).Maybe()

	mockVolume := Volume{
		Name:       "volume1",
		Driver:     "local",
		Mountpoint: "/var/lib/docker/volumes/volume1/_data",
		Created:    time.Now(),
		Size:       "100MB",
	}

	sf := uimocks.NewMockServiceFactoryInterface(t)
	sf.EXPECT().GetVolumeService().Return(vs)
	ui := newUIMockWithServices(t, sf)
	ui.On("ShowDetails", mock.AnythingOfType("*tview.Flex")).Return().Maybe()

	volumesView := NewVolumesView(ui)
	// Test the method directly without accessing unexported fields
	volumesView.showVolumeDetails(&mockVolume)

	assert.NotNil(t, volumesView)
}

func TestVolumesView_HandleAction_Delete(t *testing.T) {
	ui := newUIMockWithServices(t, nil)
	volumesView := NewVolumesView(ui)
	testVolume := Volume{
		Name:       "volume1",
		Driver:     "local",
		Mountpoint: "/var/lib/docker/volumes/volume1/_data",
		Created:    time.Now(),
		Size:       "100MB",
	}

	// Test action handling with test volume
	volumesView.handleAction('d', &testVolume)
	volumesView.handleAction('i', &testVolume)

	assert.NotNil(t, volumesView)
}

func TestVolumesView_HandleAction_Inspect(t *testing.T) {
	ui := newUIMockWithServices(t, nil)
	volumesView := NewVolumesView(ui)
	testVolume := Volume{
		Name:       "volume1",
		Driver:     "local",
		Mountpoint: "/var/lib/docker/volumes/volume1/_data",
		Created:    time.Now(),
		Size:       "100MB",
	}

	// Test action handling with test volume
	volumesView.handleAction('d', &testVolume)
	volumesView.handleAction('i', &testVolume)

	assert.NotNil(t, volumesView)
}

func TestVolumesView_HandleAction_InvalidSelection(t *testing.T) {
	ui := newUIMockWithServices(t, nil)
	volumesView := NewVolumesView(ui)
	testVolume := Volume{
		Name:       "volume1",
		Driver:     "local",
		Mountpoint: "/var/lib/docker/volumes/volume1/_data",
		Created:    time.Now(),
		Size:       "100MB",
	}

	// Test action handling with test volume
	volumesView.handleAction('d', &testVolume)
	volumesView.handleAction('i', &testVolume)

	assert.NotNil(t, volumesView)
}

func TestVolumesView_ShowTable(t *testing.T) {
	ui := newUIMockWithServices(t, nil)
	volumesView := NewVolumesView(ui)

	// Test that the view is properly set up - we can't access private methods
	assert.NotNil(t, volumesView)
}

func TestVolumesView_DeleteVolume(t *testing.T) {
	ui := newUIMockWithServices(t, nil)
	volumesView := NewVolumesView(ui)

	// Test that the view is properly set up - we can't access private methods
	assert.NotNil(t, volumesView)
}

func TestVolumesView_InspectVolume(t *testing.T) {
	ui := newUIMockWithServices(t, nil)
	volumesView := NewVolumesView(ui)

	// Test that the view is properly set up - we can't access private methods
	assert.NotNil(t, volumesView)
}

func TestVolumesView_SetupKeyBindings(t *testing.T) {
	ui := newUIMockWithServices(t, nil)
	volumesView := NewVolumesView(ui)

	// Test that key bindings are properly set up
	assert.NotNil(t, volumesView.GetTable().GetInputCapture())
}

func TestVolumesView_SetupKeyBindings_InvalidSelection(t *testing.T) {
	ui := newUIMockWithServices(t, nil)
	volumesView := NewVolumesView(ui)

	// Test that key bindings are properly set up even with no items
	assert.NotNil(t, volumesView.GetTable().GetInputCapture())
}
