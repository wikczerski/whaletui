package ui

import (
	"github.com/wikczerski/whaletui/internal/domains/container"
	"github.com/wikczerski/whaletui/internal/domains/image"
	"github.com/wikczerski/whaletui/internal/domains/network"
	swarmDomain "github.com/wikczerski/whaletui/internal/domains/swarm"
	"github.com/wikczerski/whaletui/internal/domains/volume"
	"github.com/wikczerski/whaletui/internal/ui/constants"
	"github.com/wikczerski/whaletui/internal/ui/interfaces"
	"github.com/wikczerski/whaletui/internal/ui/managers"
	"github.com/wikczerski/whaletui/internal/ui/views/swarm"
)

// createAndRegisterViews creates all views and registers them with the view registry
func (ui *UI) createAndRegisterViews() {
	ui.createResourceViews()
	ui.registerViewsWithActions()
	// Ensure views are fully registered before setting default view
	ui.setDefaultView()
}

// createResourceViews creates all the resource views
func (ui *UI) createResourceViews() {
	ui.containersView = container.NewContainersView(ui)
	ui.imagesView = image.NewImagesView(ui)
	ui.volumesView = volume.NewVolumesView(ui)
	ui.networksView = network.NewNetworksView(ui)

	// Create swarm views
	ui.swarmServicesView = swarm.NewServicesView(
		ui,
		ui.GetSwarmServiceService().(*swarmDomain.ServiceService),
		ui.modalManager.(*managers.ModalManager),
		ui.headerManager.(*managers.HeaderManager),
	)
	ui.swarmNodesView = swarm.NewNodesView(
		ui,
		ui.GetSwarmNodeService().(*swarmDomain.NodeService),
		ui.modalManager.(*managers.ModalManager),
		ui.headerManager.(*managers.HeaderManager),
	)
}

// registerViewsWithActions registers views with their metadata and actions
func (ui *UI) registerViewsWithActions() {
	ui.registerContainerView()
	ui.registerResourceViews()
}

// registerContainerView registers the containers view with its actions
func (ui *UI) registerContainerView() {
	containerActions := ""
	containerNavigation := ""
	if ui.services != nil && ui.GetContainerService() != nil {
		if actionService, ok := ui.GetContainerService().(interfaces.ServiceWithActions); ok {
			containerActions = actionService.GetActionsString()
		}
		if navigationService, ok := ui.GetContainerService().(interfaces.ServiceWithNavigation); ok {
			containerNavigation = navigationService.GetNavigationString()
		}
	}
	ui.viewRegistry.Register(
		constants.ViewContainers,
		"Containers",
		'c',
		ui.containersView.GetView(),
		ui.containersView.Refresh,
		containerActions,
		containerNavigation,
	)
}

// registerResourceViews registers the resource views with their actions
func (ui *UI) registerResourceViews() {
	actions := ui.collectServiceActions()
	ui.registerResourceViewsWithActions(actions)
}

// collectServiceActions collects actions from all available services
func (ui *UI) collectServiceActions() map[string]string {
	actions := make(map[string]string)

	if ui.services == nil {
		return actions
	}

	ui.collectImageActions(actions)
	ui.collectVolumeActions(actions)
	ui.collectNetworkActions(actions)
	ui.collectSwarmServiceActions(actions)
	ui.collectSwarmNodeActions(actions)

	return actions
}

// collectImageActions collects actions from the image service
func (ui *UI) collectImageActions(actions map[string]string) {
	if imageService := ui.services.GetImageService(); imageService != nil {
		if actionService, ok := imageService.(interfaces.ServiceWithActions); ok {
			actions[constants.ViewImages] = actionService.GetActionsString()
		}
	}
}

// collectVolumeActions collects actions from the volume service
func (ui *UI) collectVolumeActions(actions map[string]string) {
	if volumeService := ui.services.GetVolumeService(); volumeService != nil {
		if actionService, ok := volumeService.(interfaces.ServiceWithActions); ok {
			actions[constants.ViewVolumes] = actionService.GetActionsString()
		}
	}
}

// collectNetworkActions collects actions from the network service
func (ui *UI) collectNetworkActions(actions map[string]string) {
	if networkService := ui.services.GetNetworkService(); networkService != nil {
		if actionService, ok := networkService.(interfaces.ServiceWithActions); ok {
			actions[constants.ViewNetworks] = actionService.GetActionsString()
		}
	}
}

// collectSwarmServiceActions collects actions from the swarm service service
func (ui *UI) collectSwarmServiceActions(actions map[string]string) {
	if swarmServiceService := ui.GetSwarmServiceService(); swarmServiceService != nil {
		if actionService, ok := swarmServiceService.(interfaces.ServiceWithActions); ok {
			actions[constants.ViewSwarmServices] = actionService.GetActionsString()
		}
	}
}

// collectSwarmNodeActions collects actions from the swarm node service
func (ui *UI) collectSwarmNodeActions(actions map[string]string) {
	if swarmNodeService := ui.GetSwarmNodeService(); swarmNodeService != nil {
		if actionService, ok := swarmNodeService.(interfaces.ServiceWithActions); ok {
			actions[constants.ViewSwarmNodes] = actionService.GetActionsString()
		}
	}
}

// registerResourceViewsWithActions registers resource views with their collected actions
func (ui *UI) registerResourceViewsWithActions(actions map[string]string) {
	ui.registerImagesView(actions)
	ui.registerVolumesView(actions)
	ui.registerNetworksView(actions)
	ui.registerSwarmServicesView(actions)
	ui.registerSwarmNodesView(actions)
}

// registerImagesView registers the images view
func (ui *UI) registerImagesView(actions map[string]string) {
	ui.viewRegistry.Register(
		constants.ViewImages,
		"Images",
		'i',
		ui.imagesView.GetView(),
		ui.imagesView.Refresh,
		actions[constants.ViewImages],
		"",
	)
}

// registerVolumesView registers the volumes view
func (ui *UI) registerVolumesView(actions map[string]string) {
	ui.viewRegistry.Register(
		constants.ViewVolumes,
		"Volumes",
		'v',
		ui.volumesView.GetView(),
		ui.volumesView.Refresh,
		actions[constants.ViewVolumes],
		"",
	)
}

// registerNetworksView registers the networks view
func (ui *UI) registerNetworksView(actions map[string]string) {
	ui.viewRegistry.Register(
		constants.ViewNetworks,
		"Networks",
		'n',
		ui.networksView.GetView(),
		ui.networksView.Refresh,
		actions[constants.ViewNetworks],
		"",
	)
}

// registerSwarmServicesView registers the swarm services view
func (ui *UI) registerSwarmServicesView(actions map[string]string) {
	ui.viewRegistry.Register(
		constants.ViewSwarmServices,
		"Swarm Services",
		's',
		ui.swarmServicesView.GetView(),
		ui.swarmServicesView.Refresh,
		actions[constants.ViewSwarmServices],
		"",
	)
}

// registerSwarmNodesView registers the swarm nodes view
func (ui *UI) registerSwarmNodesView(actions map[string]string) {
	ui.viewRegistry.Register(
		constants.ViewSwarmNodes,
		"Swarm Nodes",
		'w',
		ui.swarmNodesView.GetView(),
		ui.swarmNodesView.Refresh,
		actions[constants.ViewSwarmNodes],
		"",
	)
}

// setDefaultView sets the default view for the application
func (ui *UI) setDefaultView() {
	ui.viewRegistry.SetCurrent(constants.DefaultView)

	// Set the default service to container for initial navigation
	if ui.services != nil {
		ui.services.SetCurrentService(constants.ViewContainers)
	}
}
