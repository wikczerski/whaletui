package volume

import (
	"context"
	"errors"
	"fmt"

	"github.com/wikczerski/whaletui/internal/shared"
	"github.com/wikczerski/whaletui/internal/ui/builders"
	"github.com/wikczerski/whaletui/internal/ui/handlers"
	"github.com/wikczerski/whaletui/internal/ui/interfaces"
)

// VolumesView displays and manages Docker volumes
type VolumesView struct {
	*shared.BaseView[shared.Volume]
	executor *handlers.OperationExecutor
}

// NewVolumesView creates a new volumes view
func NewVolumesView(ui interfaces.UIInterface) *VolumesView {
	headers := []string{"Name", "Driver", "Mountpoint", "Created", "Size"}
	baseView := shared.NewBaseView[shared.Volume](ui, "volumes", headers)

	vv := &VolumesView{
		BaseView: baseView,
		executor: handlers.NewOperationExecutor(ui),
	}

	setupVolumeViewCallbacks(vv)
	return vv
}

// setupVolumeViewCallbacks sets up the callbacks for the volumes view
func setupVolumeViewCallbacks(vv *VolumesView) {
	vv.setupBasicCallbacks()
	vv.setupActionCallbacks()
}

// setupBasicCallbacks sets up the basic view callbacks
func (vv *VolumesView) setupBasicCallbacks() {
	vv.ListItems = vv.listVolumes
	vv.FormatRow = func(v shared.Volume) []string { return vv.formatVolumeRow(&v) }
	vv.GetItemID = func(v shared.Volume) string { return v.Name }
	vv.GetItemName = func(v shared.Volume) string { return v.Name }
}

// setupActionCallbacks sets up the action-related callbacks
func (vv *VolumesView) setupActionCallbacks() {
	vv.HandleKeyPress = func(key rune, v shared.Volume) { vv.handleAction(key, &v) }
	vv.ShowDetails = func(v shared.Volume) { vv.showVolumeDetails(&v) }
	vv.GetActions = vv.getVolumeActions
}

// createDeleteVolumeFunction creates a function to delete a volume
func (vv *VolumesView) createDeleteVolumeFunction(name string) func() error {
	return func() error {
		services := vv.GetUI().GetServicesAny()
		if services == nil {
			return errors.New("volume service not available")
		}

		// Type assertion to get the service factory
		if serviceFactory, ok := services.(interface{ GetVolumeService() any }); ok {
			if volumeService := serviceFactory.GetVolumeService(); volumeService != nil {
				// Type assertion to get the RemoveVolume method
				if removeService, ok := volumeService.(interface {
					RemoveVolume(context.Context, string, bool) error
				}); ok {
					ctx := context.Background()
					// Force removal to handle cases where volume might be in use
					return removeService.RemoveVolume(ctx, name, true)
				}
			}
		}

		return errors.New("volume service not available")
	}
}

func (vv *VolumesView) listVolumes(ctx context.Context) ([]shared.Volume, error) {
	services := vv.GetUI().GetServicesAny()
	if services == nil {
		return []shared.Volume{}, nil
	}

	return vv.getVolumesFromService(ctx, services)
}

// getVolumesFromService retrieves volumes from the service factory
func (vv *VolumesView) getVolumesFromService(
	ctx context.Context,
	services any,
) ([]shared.Volume, error) {
	serviceFactory, ok := services.(interfaces.ServiceFactoryInterface)
	if !ok {
		return []shared.Volume{}, nil
	}

	volumeService := serviceFactory.GetVolumeService()
	if volumeService == nil {
		return []shared.Volume{}, nil
	}

	return vv.executeVolumeList(ctx, volumeService)
}

// executeVolumeList executes the volume list operation
func (vv *VolumesView) executeVolumeList(
	ctx context.Context,
	volumeService any,
) ([]shared.Volume, error) {
	if volumeService == nil {
		return []shared.Volume{}, nil
	}

	volumes, err := volumeService.(interface {
		ListVolumes(context.Context) ([]shared.Volume, error)
	}).ListVolumes(ctx)
	if err != nil {
		return nil, err
	}

	return volumes, nil
}

func (vv *VolumesView) formatVolumeRow(volume *shared.Volume) []string {
	return []string{
		volume.Name,
		volume.Driver,
		volume.Mountpoint,
		builders.FormatTime(volume.CreatedAt),
		volume.Size,
	}
}

func (vv *VolumesView) getVolumeActions() map[rune]string {
	return map[rune]string{
		'd': "Delete",
		'i': "Inspect",
	}
}

func (vv *VolumesView) handleAction(key rune, volume *shared.Volume) {
	services := vv.GetUI().GetServicesAny()
	if services == nil {
		return
	}

	// Type assertion to get the service factory
	if serviceFactory, ok := services.(interface{ GetVolumeService() any }); ok {
		if serviceFactory.GetVolumeService() != nil {
			switch key {
			case 'd':
				vv.deleteVolume(volume.Name)
			case 'i':
				vv.inspectVolume(volume.Name)
			}
		}
	}
}

func (vv *VolumesView) showVolumeDetails(volume *shared.Volume) {
	ctx := context.Background()
	services := vv.GetUI().GetServicesAny()
	if services == nil {
		vv.ShowItemDetails(*volume, nil, errors.New("volume service not available"))
		return
	}

	inspectService := vv.getVolumeInspectService(services)
	if inspectService == nil {
		vv.ShowItemDetails(*volume, nil, errors.New("volume service not available"))
		return
	}

	vv.executeVolumeInspection(ctx, volume, inspectService)
}

// executeVolumeInspection performs the actual volume inspection
func (vv *VolumesView) executeVolumeInspection(
	ctx context.Context,
	volume *shared.Volume,
	inspectService interface {
		InspectVolume(context.Context, string) (map[string]any, error)
	},
) {
	inspectData, err := inspectService.InspectVolume(ctx, volume.Name)
	vv.ShowItemDetails(*volume, inspectData, err)
}

// getVolumeInspectService extracts the volume inspect service from the services interface
func (vv *VolumesView) getVolumeInspectService(services any) interface {
	InspectVolume(context.Context, string) (map[string]any, error)
} {
	serviceFactory, ok := services.(interface{ GetVolumeService() any })
	if !ok {
		return nil
	}

	volumeService := serviceFactory.GetVolumeService()
	if volumeService == nil {
		return nil
	}

	inspectService, ok := volumeService.(interface {
		InspectVolume(context.Context, string) (map[string]any, error)
	})
	if !ok {
		return nil
	}

	return inspectService
}

func (vv *VolumesView) deleteVolume(name string) {
	vv.executor.ExecuteWithConfirmation(
		fmt.Sprintf("Delete volume %s?", name),
		vv.createDeleteVolumeFunction(name),
		func() { vv.Refresh() },
	)
}

func (vv *VolumesView) inspectVolume(name string) {
	services := vv.GetUI().GetServicesAny()
	if services == nil {
		return
	}

	inspectService := vv.getVolumeInspectService(services)
	if inspectService == nil {
		return
	}

	vv.performVolumeInspection(name, inspectService)
}

// performVolumeInspection performs the actual volume inspection
func (vv *VolumesView) performVolumeInspection(name string, inspectService interface {
	InspectVolume(context.Context, string) (map[string]any, error)
},
) {
	ctx := context.Background()
	inspectData, err := inspectService.InspectVolume(ctx, name)
	if err != nil {
		vv.GetUI().ShowError(err)
		return
	}

	vv.ShowItemDetails(shared.Volume{Name: name}, inspectData, err)
}
