package network

import (
	"context"
	"fmt"

	"github.com/wikczerski/whaletui/internal/shared"
	"github.com/wikczerski/whaletui/internal/ui/builders"
	"github.com/wikczerski/whaletui/internal/ui/handlers"
	"github.com/wikczerski/whaletui/internal/ui/interfaces"
)

// NetworksView displays and manages Docker networks
type NetworksView struct {
	*shared.BaseView[Network]
	handlers *handlers.ActionHandlers
}

// NewNetworksView creates a new networks view
func NewNetworksView(ui interfaces.UIInterface) *NetworksView {
	headers := []string{"ID", "Name", "Driver", "Scope", "Created"}
	baseView := shared.NewBaseView[Network](ui, "networks", headers)

	nv := &NetworksView{
		BaseView: baseView,
		handlers: handlers.NewActionHandlers(ui),
	}

	// Set up callbacks
	nv.ListItems = nv.listNetworks
	nv.FormatRow = func(n Network) []string { return nv.formatNetworkRow(&n) }
	nv.GetItemID = func(n Network) string { return n.ID }
	nv.GetItemName = func(n Network) string { return n.Name }
	nv.HandleKeyPress = func(key rune, n Network) { nv.handleAction(key, &n) }
	nv.ShowDetails = func(n Network) { nv.showNetworkDetails(&n) }
	nv.GetActions = nv.getNetworkActions

	return nv
}

func (nv *NetworksView) listNetworks(ctx context.Context) ([]Network, error) {
	services := nv.GetUI().GetServices()
	if services == nil {
		return []Network{}, nil
	}

	if networkService := services.GetNetworkService(); networkService != nil {
		networks, err := networkService.ListNetworks(ctx)
		if err != nil {
			return nil, err
		}
		// Convert interface{} to Network
		result := make([]Network, 0, len(networks))
		for _, net := range networks {
			if network, ok := net.(Network); ok {
				result = append(result, network)
			}
		}
		return result, nil
	}

	return []Network{}, nil
}

func (nv *NetworksView) formatNetworkRow(network *Network) []string {
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

func (nv *NetworksView) handleAction(key rune, network *Network) {
	services := nv.GetUI().GetServices()
	if services == nil {
		return
	}

	if services.GetNetworkService() == nil {
		return
	}

	switch key {
	case 'd':
		nv.deleteNetwork(network.ID)
	case 'i':
		nv.inspectNetwork(network.ID)
	}
}

func (nv *NetworksView) showNetworkDetails(network *Network) {
	ctx := context.Background()
	services := nv.GetUI().GetServices()
	if services == nil {
		nv.ShowItemDetails(*network, nil, fmt.Errorf("network service not available"))
		return
	}

	if services.GetNetworkService() == nil {
		nv.ShowItemDetails(*network, nil, fmt.Errorf("network service not available"))
		return
	}

	inspectData, err := services.GetNetworkService().InspectNetwork(ctx, network.ID)
	nv.ShowItemDetails(*network, inspectData, err)
}

func (nv *NetworksView) deleteNetwork(id string) {
	services := nv.GetUI().GetServices()
	if services == nil || services.GetNetworkService() == nil {
		return
	}
	nv.handlers.HandleResourceAction('d', "network", id, "",
		services.GetNetworkService().InspectNetwork, nil, func() { nv.Refresh() })
}

func (nv *NetworksView) inspectNetwork(id string) {
	services := nv.GetUI().GetServices()
	if services == nil || services.GetNetworkService() == nil {
		return
	}
	nv.handlers.HandleResourceAction('i', "network", id, "",
		services.GetNetworkService().InspectNetwork, nil, func() { nv.Refresh() })
}
