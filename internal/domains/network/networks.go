package network

import (
	"context"
	"errors"

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

	setupNetworkViewCallbacks(nv)
	return nv
}

// setupNetworkViewCallbacks sets up the callbacks for the networks view
func setupNetworkViewCallbacks(nv *NetworksView) {
	nv.setupBasicCallbacks()
	nv.setupActionCallbacks()
}

// setupBasicCallbacks sets up the basic view callbacks
func (nv *NetworksView) setupBasicCallbacks() {
	nv.ListItems = nv.listNetworks
	nv.FormatRow = func(n shared.Network) []string { return nv.formatNetworkRow(&n) }
	nv.GetItemID = func(n shared.Network) string { return n.ID }
	nv.GetItemName = func(n shared.Network) string { return n.Name }
}

// setupActionCallbacks sets up the action-related callbacks
func (nv *NetworksView) setupActionCallbacks() {
	nv.HandleKeyPress = func(key rune, n shared.Network) { nv.handleAction(key, &n) }
	nv.ShowDetailsCallback = func(n shared.Network) { nv.showNetworkDetails(&n) }
	nv.GetActions = nv.getNetworkActions
}

func (nv *NetworksView) listNetworks(ctx context.Context) ([]shared.Network, error) {
	services := nv.GetUI().GetServicesAny()
	if services == nil {
		return []shared.Network{}, nil
	}

	networkService := nv.getNetworkService(services)
	if networkService == nil {
		return []shared.Network{}, nil
	}

	networks, err := networkService.ListNetworks(ctx)
	if err != nil {
		return nil, err
	}

	return networks, nil
}

// getNetworkService extracts the network service from the services interface
func (nv *NetworksView) getNetworkService(services any) interfaces.NetworkService {
	serviceFactory, ok := services.(interfaces.ServiceFactoryInterface)
	if !ok {
		return nil
	}

	networkService := serviceFactory.GetNetworkService()
	if networkService == nil {
		return nil
	}

	return networkService
}

// getNetworkInspectService extracts the network inspect service from the services interface
func (nv *NetworksView) getNetworkInspectService(services any) interface {
	InspectNetwork(context.Context, string) (map[string]any, error)
} {
	serviceFactory, ok := services.(interface{ GetNetworkService() any })
	if !ok {
		return nil
	}

	networkService := serviceFactory.GetNetworkService()
	if networkService == nil {
		return nil
	}

	inspectService, ok := networkService.(interface {
		InspectNetwork(context.Context, string) (map[string]any, error)
	})
	if !ok {
		return nil
	}

	return inspectService
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
		nv.ShowItemDetails(*network, nil, errors.New("network service not available"))
		return
	}

	inspectService := nv.getNetworkInspectService(services)
	if inspectService == nil {
		nv.ShowItemDetails(*network, nil, errors.New("network service not available"))
		return
	}

	nv.executeNetworkInspection(ctx, network, inspectService)
}

// executeNetworkInspection performs the actual network inspection
func (nv *NetworksView) executeNetworkInspection(
	ctx context.Context,
	network *shared.Network,
	inspectService interface {
		InspectNetwork(context.Context, string) (map[string]any, error)
	},
) {
	inspectData, err := inspectService.InspectNetwork(ctx, network.ID)
	nv.ShowItemDetails(*network, inspectData, err)
}

func (nv *NetworksView) deleteNetwork(id string) {
	services := nv.GetUI().GetServicesAny()
	if services == nil {
		return
	}

	inspectService := nv.getNetworkInspectService(services)
	if inspectService == nil {
		return
	}

	nv.handlers.HandleResourceAction('d', "network", id, "",
		inspectService.InspectNetwork, nil, func() { nv.Refresh() })
}

func (nv *NetworksView) inspectNetwork(id string) {
	services := nv.GetUI().GetServicesAny()
	if services == nil {
		return
	}

	inspectService := nv.getNetworkInspectService(services)
	if inspectService == nil {
		return
	}

	nv.handlers.HandleResourceAction('i', "network", id, "",
		inspectService.InspectNetwork, nil, func() { nv.Refresh() })
}
