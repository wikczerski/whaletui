package swarm

import (
	"github.com/wikczerski/whaletui/internal/domains/swarm"
	uiinterfaces "github.com/wikczerski/whaletui/internal/ui/interfaces"
	"github.com/wikczerski/whaletui/internal/ui/managers"
)

// Registry manages swarm-related views
type Registry struct {
	servicesView *ServicesView
	nodesView    *NodesView
}

// NewRegistry creates a new swarm view registry
func NewRegistry(
	ui uiinterfaces.UIInterface,
	serviceService *swarm.ServiceService,
	nodeService *swarm.NodeService,
	modalManager *managers.ModalManager,
	headerManager *managers.HeaderManager,
) *Registry {
	return &Registry{
		servicesView: NewServicesView(ui, serviceService, modalManager, headerManager),
		nodesView:    NewNodesView(ui, nodeService, modalManager, headerManager),
	}
}

// GetServicesView returns the swarm services view
func (r *Registry) GetServicesView() *ServicesView {
	return r.servicesView
}

// GetNodesView returns the swarm nodes view
func (r *Registry) GetNodesView() *NodesView {
	return r.nodesView
}

// GetViewByName returns a view by name
func (r *Registry) GetViewByName(name string) interface{} {
	switch name {
	case "services":
		return r.servicesView
	case "nodes":
		return r.nodesView
	default:
		return nil
	}
}

// GetAllViews returns all swarm views
func (r *Registry) GetAllViews() map[string]interface{} {
	return map[string]interface{}{
		"services": r.servicesView,
		"nodes":    r.nodesView,
	}
}
