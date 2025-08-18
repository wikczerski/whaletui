package views

import (
	"context"

	"github.com/wikczerski/D5r/internal/models"
	"github.com/wikczerski/D5r/internal/ui/builders"
	"github.com/wikczerski/D5r/internal/ui/handlers"
	"github.com/wikczerski/D5r/internal/ui/interfaces"
)

// NetworksView displays and manages Docker networks
type NetworksView struct {
	*BaseView[models.Network]
	handlers *handlers.ActionHandlers
}

// NewNetworksView creates a new networks view
func NewNetworksView(ui interfaces.UIInterface) *NetworksView {
	headers := []string{"ID", "Name", "Driver", "Scope", "Created"}
	baseView := NewBaseView[models.Network](ui, "networks", headers)

	nv := &NetworksView{
		BaseView: baseView,
		handlers: handlers.NewActionHandlers(ui),
	}

	// Set up callbacks
	nv.ListItems = nv.listNetworks
	nv.FormatRow = func(n models.Network) []string { return nv.formatNetworkRow(&n) }
	nv.GetItemID = func(n models.Network) string { return n.ID }
	nv.GetItemName = func(n models.Network) string { return n.Name }
	nv.HandleKeyPress = func(key rune, n models.Network) { nv.handleAction(key, &n) }
	nv.ShowDetails = func(n models.Network) { nv.showNetworkDetails(&n) }
	nv.GetActions = nv.getNetworkActions

	return nv
}

func (nv *NetworksView) listNetworks(ctx context.Context) ([]models.Network, error) {
	services := nv.ui.GetServices()
	if services == nil || services.GetNetworkService() == nil {
		return []models.Network{}, nil
	}
	return services.GetNetworkService().ListNetworks(ctx)
}

func (nv *NetworksView) formatNetworkRow(network *models.Network) []string {
	return []string{
		network.ID,
		network.Name,
		network.Driver,
		network.Scope,
		builders.FormatTime(network.Created),
	}
}

func (nv *NetworksView) getNetworkActions() map[rune]string {
	return map[rune]string{
		'd': "Delete",
		'i': "Inspect",
	}
}

func (nv *NetworksView) handleAction(key rune, network *models.Network) {
	services := nv.ui.GetServices()
	if services == nil || services.GetNetworkService() == nil {
		return
	}

	switch key {
	case 'd':
		nv.deleteNetwork(network.ID)
	case 'i':
		nv.inspectNetwork(network.ID)
	}
}

func (nv *NetworksView) showNetworkDetails(network *models.Network) {
	ctx := context.Background()
	services := nv.ui.GetServices()
	inspectData, err := services.GetNetworkService().InspectNetwork(ctx, network.ID)
	nv.ShowItemDetails(*network, inspectData, err)
}

func (nv *NetworksView) deleteNetwork(id string) {
	services := nv.ui.GetServices()
	nv.handlers.HandleResourceAction('d', "network", id, "",
		services.GetNetworkService().InspectNetwork, nil, func() { nv.Refresh() })
}

func (nv *NetworksView) inspectNetwork(id string) {
	services := nv.ui.GetServices()
	nv.handlers.HandleResourceAction('i', "network", id, "",
		services.GetNetworkService().InspectNetwork, nil, func() { nv.Refresh() })
}
