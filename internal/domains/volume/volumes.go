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
	*shared.BaseView[Volume]
	executor *handlers.OperationExecutor
}

// NewVolumesView creates a new volumes view
func NewVolumesView(ui interfaces.UIInterface) *VolumesView {
	headers := []string{"Name", "Driver", "Mountpoint", "Created", "Size"}
	baseView := shared.NewBaseView[Volume](ui, "volumes", headers)

	vv := &VolumesView{
		BaseView: baseView,
		executor: handlers.NewOperationExecutor(ui),
	}

	// Set up callbacks
	vv.ListItems = vv.listVolumes
	vv.FormatRow = func(v Volume) []string { return vv.formatVolumeRow(&v) }
	vv.GetItemID = func(v Volume) string { return v.Name }
	vv.GetItemName = func(v Volume) string { return v.Name }
	vv.HandleKeyPress = func(key rune, v Volume) { vv.handleAction(key, &v) }
	vv.ShowDetails = func(v Volume) { vv.showVolumeDetails(&v) }
	vv.GetActions = vv.getVolumeActions

	return vv
}

// createDeleteVolumeFunction creates a function to delete a volume
func (vv *VolumesView) createDeleteVolumeFunction(name string) func() error {
	return func() error {
		services := vv.GetUI().GetServices()
		if services == nil || services.GetVolumeService() == nil {
			return fmt.Errorf("volume service not available")
		}
		ctx := context.Background()
		// Force removal to handle cases where volume might be in use
		return services.GetVolumeService().RemoveVolume(ctx, name, true)
	}
}

func (vv *VolumesView) listVolumes(ctx context.Context) ([]Volume, error) {
	services := vv.GetUI().GetServices()
	if services == nil {
		return []Volume{}, nil
	}

	if volumeService := services.GetVolumeService(); volumeService != nil {
		volumes, err := volumeService.ListVolumes(ctx)
		if err != nil {
			return nil, err
		}
		// Convert interface{} to Volume
		result := make([]Volume, 0, len(volumes))
		for _, vol := range volumes {
			if volume, ok := vol.(Volume); ok {
				result = append(result, volume)
			}
		}
		return result, nil
	}

	return []Volume{}, nil
}

func (vv *VolumesView) formatVolumeRow(volume *Volume) []string {
	return []string{
		volume.Name,
		volume.Driver,
		volume.Mountpoint,
		builders.FormatTime(volume.Created),
		volume.Size,
	}
}

func (vv *VolumesView) getVolumeActions() map[rune]string {
	return map[rune]string{
		'd': "Delete",
		'i': "Inspect",
	}
}

func (vv *VolumesView) handleAction(key rune, volume *Volume) {
	services := vv.GetUI().GetServices()
	if services == nil {
		return
	}

	if services.GetVolumeService() == nil {
		return
	}

	switch key {
	case 'd':
		vv.deleteVolume(volume.Name)
	case 'i':
		vv.inspectVolume(volume.Name)
	}
}

func (vv *VolumesView) showVolumeDetails(volume *Volume) {
	ctx := context.Background()
	services := vv.GetUI().GetServices()
	if services == nil {
		vv.ShowItemDetails(*volume, nil, fmt.Errorf("volume service not available"))
		return
	}

	if services.GetVolumeService() == nil {
		vv.ShowItemDetails(*volume, nil, fmt.Errorf("volume service not available"))
		return
	}

	inspectData, err := services.GetVolumeService().InspectVolume(ctx, volume.Name)
	vv.ShowItemDetails(*volume, inspectData, err)
}

func (vv *VolumesView) deleteVolume(name string) {
	vv.executor.ExecuteWithConfirmation(
		fmt.Sprintf("Delete volume %s?", name),
		vv.createDeleteVolumeFunction(name),
		func() { vv.Refresh() },
	)
}

func (vv *VolumesView) inspectVolume(name string) {
	services := vv.GetUI().GetServices()
	if services == nil || services.GetVolumeService() == nil {
		return
	}

	ctx := context.Background()
	inspectData, err := services.GetVolumeService().InspectVolume(ctx, name)
	if err != nil {
		vv.GetUI().ShowError(err)
		return
	}

	vv.ShowItemDetails(Volume{Name: name}, inspectData, err)
}
