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
	*shared.BaseView[shared.Network]
	handlers *handlers.ActionHandlers
}

// NewNetworksView creates a new networks view
func NewNetworksView(ui interfaces.UIInterface) *NetworksView {
	headers := []string{"ID", "Name", "Driver", "Scope", "Created"}
	baseView := shared.NewBaseView[shared.Network](ui, "networks", headers)

	nv := &NetworksView{
		BaseView: baseView,
		handlers: handlers.NewActionHandlers(ui),
	}

	// Set up callbacks
	nv.ListItems = nv.listNetworks
	nv.FormatRow = func(n shared.Network) []string { return nv.formatNetworkRow(&n) }
	nv.GetItemID = func(n shared.Network) string { return n.ID }
	nv.GetItemName = func(n shared.Network) string { return n.Name }
	nv.HandleKeyPress = func(key rune, n shared.Network) { nv.handleAction(key, &n) }
	nv.ShowDetails = func(n shared.Network) { nv.showNetworkDetails(&n) }
	nv.GetActions = nv.getNetworkActions

	return nv
}

func (nv *NetworksView) listNetworks(ctx context.Context) ([]shared.Network, error) {
	services := nv.GetUI().GetServicesAny()
	if services == nil {
		return []shared.Network{}, nil
	}

	// Type assertion to get the service factory
	if serviceFactory, ok := services.(interfaces.ServiceFactoryInterface); ok {
		if networkService := serviceFactory.GetNetworkService(); networkService != nil {
			// Type assertion to get the ListNetworks method
			if networkService != nil {
				networks, err := networkService.ListNetworks(ctx)
				if err != nil {
					return nil, err
				}
				return networks, nil
			}
		}
	}

	return []shared.Network{}, nil
}

func (nv *NetworksView) formatNetworkRow(network *shared.Network) []string {
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

func (nv *NetworksView) handleAction(key rune, network *shared.Network) {
	services := nv.GetUI().GetServicesAny()
	if services == nil {
		return
	}

	// Type assertion to get the service factory
	if serviceFactory, ok := services.(interface{ GetNetworkService() any }); ok {
		if serviceFactory.GetNetworkService() != nil {
			switch key {
			case 'd':
				nv.deleteNetwork(network.ID)
			case 'i':
				nv.inspectNetwork(network.ID)
			}
		}
	}
}

func (nv *NetworksView) showNetworkDetails(network *shared.Network) {
	ctx := context.Background()
	services := nv.GetUI().GetServicesAny()
	if services == nil {
		nv.ShowItemDetails(*network, nil, fmt.Errorf("network service not available"))
		return
	}

	// Type assertion to get the service factory
	if serviceFactory, ok := services.(interface{ GetNetworkService() any }); ok {
		if networkService := serviceFactory.GetNetworkService(); networkService != nil {
			// Type assertion to get the InspectNetwork method
			if inspectService, ok := networkService.(interface {
				InspectNetwork(context.Context, string) (map[string]any, error)
			}); ok {
				inspectData, err := inspectService.InspectNetwork(ctx, network.ID)
				nv.ShowItemDetails(*network, inspectData, err)
				return
			}
		}
	}

	nv.ShowItemDetails(*network, nil, fmt.Errorf("network service not available"))
}

func (nv *NetworksView) deleteNetwork(id string) {
	services := nv.GetUI().GetServicesAny()
	if services == nil {
		return
	}

	// Type assertion to get the service factory
	if serviceFactory, ok := services.(interface{ GetNetworkService() any }); ok {
		if networkService := serviceFactory.GetNetworkService(); networkService != nil {
			// Type assertion to get the InspectNetwork method
			if inspectService, ok := networkService.(interface {
				InspectNetwork(context.Context, string) (map[string]any, error)
			}); ok {
				nv.handlers.HandleResourceAction('d', "network", id, "",
					inspectService.InspectNetwork, nil, func() { nv.Refresh() })
			}
		}
	}
}

func (nv *NetworksView) inspectNetwork(id string) {
	services := nv.GetUI().GetServicesAny()
	if services == nil {
		return
	}

	// Type assertion to get the service factory
	if serviceFactory, ok := services.(interface{ GetNetworkService() any }); ok {
		if networkService := serviceFactory.GetNetworkService(); networkService != nil {
			// Type assertion to get the InspectNetwork method
			if inspectService, ok := networkService.(interface {
				InspectNetwork(context.Context, string) (map[string]any, error)
			}); ok {
				nv.handlers.HandleResourceAction('i', "network", id, "",
					inspectService.InspectNetwork, nil, func() { nv.Refresh() })
			}
		}
	}
}
