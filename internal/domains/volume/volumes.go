package volume

import (
	"context"
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

	// Set up callbacks
	vv.ListItems = vv.listVolumes
	vv.FormatRow = func(v shared.Volume) []string { return vv.formatVolumeRow(&v) }
	vv.GetItemID = func(v shared.Volume) string { return v.Name }
	vv.GetItemName = func(v shared.Volume) string { return v.Name }
	vv.HandleKeyPress = func(key rune, v shared.Volume) { vv.handleAction(key, &v) }
	vv.ShowDetails = func(v shared.Volume) { vv.showVolumeDetails(&v) }
	vv.GetActions = vv.getVolumeActions

	return vv
}

// createDeleteVolumeFunction creates a function to delete a volume
func (vv *VolumesView) createDeleteVolumeFunction(name string) func() error {
	return func() error {
		services := vv.GetUI().GetServicesAny()
		if services == nil {
			return fmt.Errorf("volume service not available")
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

		return fmt.Errorf("volume service not available")
	}
}

func (vv *VolumesView) listVolumes(ctx context.Context) ([]shared.Volume, error) {
	services := vv.GetUI().GetServicesAny()
	if services == nil {
		return []shared.Volume{}, nil
	}

	// Type assertion to get the service factory
	if serviceFactory, ok := services.(interfaces.ServiceFactoryInterface); ok {
		if volumeService := serviceFactory.GetVolumeService(); volumeService != nil {
			// Type assertion to get the ListVolumes method
			if volumeService != nil {
				volumes, err := volumeService.ListVolumes(ctx)
				if err != nil {
					return nil, err
				}
				return volumes, nil
			}
		}
	}

	return []shared.Volume{}, nil
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
		vv.ShowItemDetails(*volume, nil, fmt.Errorf("volume service not available"))
		return
	}

	// Type assertion to get the service factory
	if serviceFactory, ok := services.(interface{ GetVolumeService() any }); ok {
		if volumeService := serviceFactory.GetVolumeService(); volumeService != nil {
			// Type assertion to get the InspectVolume method
			if inspectService, ok := volumeService.(interface {
				InspectVolume(context.Context, string) (map[string]any, error)
			}); ok {
				inspectData, err := inspectService.InspectVolume(ctx, volume.Name)
				vv.ShowItemDetails(*volume, inspectData, err)
				return
			}
		}
	}

	vv.ShowItemDetails(*volume, nil, fmt.Errorf("volume service not available"))
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

	// Type assertion to get the service factory
	if serviceFactory, ok := services.(interface{ GetVolumeService() any }); ok {
		if volumeService := serviceFactory.GetVolumeService(); volumeService != nil {
			// Type assertion to get the InspectVolume method
			if inspectService, ok := volumeService.(interface {
				InspectVolume(context.Context, string) (map[string]any, error)
			}); ok {
				ctx := context.Background()
				inspectData, err := inspectService.InspectVolume(ctx, name)
				if err != nil {
					vv.GetUI().ShowError(err)
					return
				}

				vv.ShowItemDetails(shared.Volume{Name: name}, inspectData, err)
				return
			}
		}
	}
}
