package views

import (
	"context"
	"fmt"

	"github.com/rivo/tview"
	"github.com/wikczerski/D5r/internal/models"
	"github.com/wikczerski/D5r/internal/ui/builders"
	"github.com/wikczerski/D5r/internal/ui/handlers"
	"github.com/wikczerski/D5r/internal/ui/interfaces"
)

type VolumesView struct {
	*BaseView[models.Volume]
	executor *handlers.OperationExecutor
}

func NewVolumesView(ui interfaces.UIInterface) *VolumesView {
	headers := []string{"Name", "Driver", "Mountpoint", "Created", "Size"}
	baseView := NewBaseView[models.Volume](ui, "volumes", headers)

	vv := &VolumesView{
		BaseView: baseView,
		executor: handlers.NewOperationExecutor(ui),
	}

	// Set up callbacks
	vv.ListItems = vv.listVolumes
	vv.FormatRow = vv.formatVolumeRow
	vv.GetItemID = func(v models.Volume) string { return v.Name }
	vv.GetItemName = func(v models.Volume) string { return v.Name }
	vv.HandleKeyPress = vv.handleVolumeKey
	vv.ShowDetails = vv.showVolumeDetails
	vv.GetActions = vv.getVolumeActions

	return vv
}

func (vv *VolumesView) listVolumes(ctx context.Context) ([]models.Volume, error) {
	services := vv.ui.GetServices()
	if services == nil || services.VolumeService == nil {
		return []models.Volume{}, nil
	}
	return services.VolumeService.ListVolumes(ctx)
}

func (vv *VolumesView) formatVolumeRow(volume models.Volume) []string {
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

func (vv *VolumesView) handleVolumeKey(key rune, volume models.Volume) {
	switch key {
	case 'd':
		vv.deleteVolume(volume.Name)
	case 'i':
		vv.inspectVolume(volume.Name)
	}
}

func (vv *VolumesView) showVolumeDetails(volume models.Volume) {
	ctx := context.Background()
	services := vv.ui.GetServices()
	inspectData, err := services.VolumeService.InspectVolume(ctx, volume.Name)
	vv.ShowItemDetails(volume, inspectData, err)
}

func (vv *VolumesView) deleteVolume(name string) {
	vv.executor.ExecuteWithConfirmation(
		fmt.Sprintf("Delete volume %s?", name),
		func() error {
			services := vv.ui.GetServices()
			if services == nil || services.VolumeService == nil {
				return fmt.Errorf("volume service not available")
			}

			ctx := context.Background()
			// Force removal to handle cases where volume might be in use
			return services.VolumeService.RemoveVolume(ctx, name, true)
		},
		func() { vv.Refresh() },
	)
}

func (vv *VolumesView) inspectVolume(name string) {
	inspectView, inspectFlex := builders.CreateInspectView(fmt.Sprintf("Inspect: %s", name))

	inspectFlex.GetItem(1).(*tview.Button).SetSelectedFunc(func() {
		pages := vv.ui.GetPages().(*tview.Pages)
		pages.RemovePage("inspect")
	})

	pages := vv.ui.GetPages().(*tview.Pages)
	pages.AddPage("inspect", inspectFlex, true, true)
	inspectView.SetText("Volume inspection not implemented yet.")
}
