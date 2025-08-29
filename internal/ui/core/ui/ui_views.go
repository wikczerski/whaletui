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
		ui.services.GetSwarmServiceService().(*swarmDomain.ServiceService),
		ui.modalManager.(*managers.ModalManager),
		ui.headerManager.(*managers.HeaderManager),
	)
	ui.swarmNodesView = swarm.NewNodesView(
		ui,
		ui.services.GetSwarmNodeService().(*swarmDomain.NodeService),
		ui.modalManager.(*managers.ModalManager),
		ui.headerManager.(*managers.HeaderManager),
	)
}

// registerViewsWithActions registers views with their metadata and actions
func (ui *UI) registerViewsWithActions() {
	if ui.services != nil {
		ui.registerContainerView()
		ui.registerResourceViews()
	} else {
		// Register views without actions when services are not available
		ui.registerViewsWithoutServices()
	}
}

// registerContainerView registers the containers view with its actions
func (ui *UI) registerContainerView() {
	containerActions := ""
	containerNavigation := ""
	if ui.services != nil && ui.services.GetContainerService() != nil {
		if actionService, ok := ui.services.GetContainerService().(interfaces.ServiceWithActions); ok {
			containerActions = actionService.GetActionsString()
		}
		if navigationService, ok := ui.services.GetContainerService().(interfaces.ServiceWithNavigation); ok {
			containerNavigation = navigationService.GetNavigationString()
		}
	}
	ui.viewRegistry.Register(
		"containers",
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
			actions["images"] = actionService.GetActionsString()
		}
	}
}

// collectVolumeActions collects actions from the volume service
func (ui *UI) collectVolumeActions(actions map[string]string) {
	if volumeService := ui.services.GetVolumeService(); volumeService != nil {
		if actionService, ok := volumeService.(interfaces.ServiceWithActions); ok {
			actions["volumes"] = actionService.GetActionsString()
		}
	}
}

// collectNetworkActions collects actions from the network service
func (ui *UI) collectNetworkActions(actions map[string]string) {
	if networkService := ui.services.GetNetworkService(); networkService != nil {
		if actionService, ok := networkService.(interfaces.ServiceWithActions); ok {
			actions["networks"] = actionService.GetActionsString()
		}
	}
}

// collectSwarmServiceActions collects actions from the swarm service service
func (ui *UI) collectSwarmServiceActions(actions map[string]string) {
	if swarmServiceService := ui.services.GetSwarmServiceService(); swarmServiceService != nil {
		if actionService, ok := swarmServiceService.(interfaces.ServiceWithActions); ok {
			actions["swarmServices"] = actionService.GetActionsString()
		}
	}
}

// collectSwarmNodeActions collects actions from the swarm node service
func (ui *UI) collectSwarmNodeActions(actions map[string]string) {
	if swarmNodeService := ui.services.GetSwarmNodeService(); swarmNodeService != nil {
		if actionService, ok := swarmNodeService.(interfaces.ServiceWithActions); ok {
			actions["swarmNodes"] = actionService.GetActionsString()
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
		"images",
		"Images",
		'i',
		ui.imagesView.GetView(),
		ui.imagesView.Refresh,
		actions["images"],
		"",
	)
}

// registerVolumesView registers the volumes view
func (ui *UI) registerVolumesView(actions map[string]string) {
	ui.viewRegistry.Register(
		"volumes",
		"Volumes",
		'v',
		ui.volumesView.GetView(),
		ui.volumesView.Refresh,
		actions["volumes"],
		"",
	)
}

// registerNetworksView registers the networks view
func (ui *UI) registerNetworksView(actions map[string]string) {
	ui.viewRegistry.Register(
		"networks",
		"Networks",
		'n',
		ui.networksView.GetView(),
		ui.networksView.Refresh,
		actions["networks"],
		"",
	)
}

// registerSwarmServicesView registers the swarm services view
func (ui *UI) registerSwarmServicesView(actions map[string]string) {
	ui.viewRegistry.Register(
		"swarmServices",
		"Swarm Services",
		's',
		ui.swarmServicesView.GetView(),
		ui.swarmServicesView.Refresh,
		actions["swarmServices"],
		"",
	)
}

// registerSwarmNodesView registers the swarm nodes view
func (ui *UI) registerSwarmNodesView(actions map[string]string) {
	ui.viewRegistry.Register(
		"swarmNodes",
		"Swarm Nodes",
		'w',
		ui.swarmNodesView.GetView(),
		ui.swarmNodesView.Refresh,
		actions["swarmNodes"],
		"",
	)
}

// registerViewsWithoutServices registers views without service actions
func (ui *UI) registerViewsWithoutServices() {
	ui.registerContainersViewWithoutActions()
	ui.registerImagesViewWithoutActions()
	ui.registerVolumesViewWithoutActions()
	ui.registerNetworksViewWithoutActions()
	ui.registerSwarmServicesViewWithoutActions()
	ui.registerSwarmNodesViewWithoutActions()
}

// registerContainersViewWithoutActions registers the containers view without actions
func (ui *UI) registerContainersViewWithoutActions() {
	ui.viewRegistry.Register(
		"containers",
		"Containers",
		'c',
		ui.containersView.GetView(),
		ui.containersView.Refresh,
		"",
		"",
	)
}

// registerImagesViewWithoutActions registers the images view without actions
func (ui *UI) registerImagesViewWithoutActions() {
	ui.viewRegistry.Register(
		"images",
		"Images",
		'i',
		ui.imagesView.GetView(),
		ui.imagesView.Refresh,
		"",
		"",
	)
}

// registerVolumesViewWithoutActions registers the volumes view without actions
func (ui *UI) registerVolumesViewWithoutActions() {
	ui.viewRegistry.Register(
		"volumes",
		"Volumes",
		'v',
		ui.volumesView.GetView(),
		ui.volumesView.Refresh,
		"",
		"",
	)
}

// registerNetworksViewWithoutActions registers the networks view without actions
func (ui *UI) registerNetworksViewWithoutActions() {
	ui.viewRegistry.Register(
		"networks",
		"Networks",
		'n',
		ui.networksView.GetView(),
		ui.networksView.Refresh,
		"",
		"",
	)
}

// registerSwarmServicesViewWithoutActions registers the swarm services view without actions
func (ui *UI) registerSwarmServicesViewWithoutActions() {
	ui.viewRegistry.Register(
		"swarmServices",
		"Swarm Services",
		's',
		ui.swarmServicesView.GetView(),
		ui.swarmServicesView.Refresh,
		"",
		"",
	)
}

// registerSwarmNodesViewWithoutActions registers the swarm nodes view without actions
func (ui *UI) registerSwarmNodesViewWithoutActions() {
	ui.viewRegistry.Register(
		"swarmNodes",
		"Swarm Nodes",
		'w',
		ui.swarmNodesView.GetView(),
		ui.swarmNodesView.Refresh,
		"",
		"",
	)
}

// setDefaultView sets the default view for the application
func (ui *UI) setDefaultView() {
	ui.viewRegistry.SetCurrent(constants.DefaultView)

	// Set the default service to container for initial navigation
	if ui.services != nil {
		ui.services.SetCurrentService("container")
	}
}
