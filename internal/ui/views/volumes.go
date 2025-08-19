package views

import (
	"context"
	"fmt"

	"github.com/wikczerski/D5r/internal/models"
	"github.com/wikczerski/D5r/internal/ui/builders"
	"github.com/wikczerski/D5r/internal/ui/handlers"
	"github.com/wikczerski/D5r/internal/ui/interfaces"
)

// VolumesView displays and manages Docker volumes
type VolumesView struct {
	*BaseView[models.Volume]
	executor *handlers.OperationExecutor
}

// NewVolumesView creates a new volumes view
func NewVolumesView(ui interfaces.UIInterface) *VolumesView {
	headers := []string{"Name", "Driver", "Mountpoint", "Created", "Size"}
	baseView := NewBaseView[models.Volume](ui, "volumes", headers)

	vv := &VolumesView{
		BaseView: baseView,
		executor: handlers.NewOperationExecutor(ui),
	}

	// Set up callbacks
	vv.ListItems = vv.listVolumes
	vv.FormatRow = func(v models.Volume) []string { return vv.formatVolumeRow(&v) }
	vv.GetItemID = func(v models.Volume) string { return v.Name }
	vv.GetItemName = func(v models.Volume) string { return v.Name }
	vv.HandleKeyPress = func(key rune, v models.Volume) { vv.handleAction(key, &v) }
	vv.ShowDetails = func(v models.Volume) { vv.showVolumeDetails(&v) }
	vv.GetActions = vv.getVolumeActions

	return vv
}

func (vv *VolumesView) listVolumes(ctx context.Context) ([]models.Volume, error) {
	services := vv.ui.GetServices()
	if services == nil {
		return []models.Volume{}, nil
	}

	if volumeService := services.GetVolumeService(); volumeService != nil {
		return volumeService.ListVolumes(ctx)
	}

	return []models.Volume{}, nil
}

func (vv *VolumesView) formatVolumeRow(volume *models.Volume) []string {
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

func (vv *VolumesView) handleAction(key rune, volume *models.Volume) {
	services := vv.ui.GetServices()
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

func (vv *VolumesView) showVolumeDetails(volume *models.Volume) {
	ctx := context.Background()
	services := vv.ui.GetServices()
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

// createDeleteVolumeFunction creates a function to delete a volume
func (vv *VolumesView) createDeleteVolumeFunction(name string) func() error {
	return func() error {
		services := vv.ui.GetServices()
		if services == nil || services.GetVolumeService() == nil {
			return fmt.Errorf("volume service not available")
		}
		ctx := context.Background()
		// Force removal to handle cases where volume might be in use
		return services.GetVolumeService().RemoveVolume(ctx, name, true)
	}
}

func (vv *VolumesView) inspectVolume(name string) {
	services := vv.ui.GetServices()
	if services == nil || services.GetVolumeService() == nil {
		return
	}

	ctx := context.Background()
	inspectData, err := services.GetVolumeService().InspectVolume(ctx, name)
	if err != nil {
		vv.ui.ShowError(err)
		return
	}

	vv.ShowItemDetails(models.Volume{Name: name}, inspectData, err)
}
