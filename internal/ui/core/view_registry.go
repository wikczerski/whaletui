package core

import "github.com/rivo/tview"

// ViewInfo holds information about a UI view
type ViewInfo struct {
	Name       string
	Title      string
	Shortcut   rune
	View       tview.Primitive
	Refresh    func()
	Actions    string
	Navigation string
}

// ViewRegistry manages all available views and their metadata
type ViewRegistry struct {
	views       map[string]*ViewInfo
	currentView string
}

// NewViewRegistry creates a new view registry
func NewViewRegistry() *ViewRegistry {
	return &ViewRegistry{
		views: make(map[string]*ViewInfo),
	}
}

// Register adds a view to the registry
func (vr *ViewRegistry) Register(
	name, title string,
	shortcut rune,
	view tview.Primitive,
	refresh func(),
	actions, navigation string,
) {
	vr.views[name] = &ViewInfo{
		Name:       name,
		Title:      title,
		Shortcut:   shortcut,
		View:       view,
		Refresh:    refresh,
		Actions:    actions,
		Navigation: navigation,
	}
}

// GetCurrent returns the current view info
func (vr *ViewRegistry) GetCurrent() *ViewInfo {
	return vr.views[vr.currentView]
}

// SetCurrent sets the current view
func (vr *ViewRegistry) SetCurrent(name string) {
	if _, exists := vr.views[name]; exists {
		vr.currentView = name
	}
}

// GetCurrentActionsString returns the actions string from the current view
func (vr *ViewRegistry) GetCurrentActionsString() string {
	if currentView := vr.GetCurrent(); currentView != nil {
		return currentView.Actions
	}
	return ""
}

// GetCurrentNavigationString returns the navigation string from the current view
func (vr *ViewRegistry) GetCurrentNavigationString() string {
	if currentView := vr.GetCurrent(); currentView != nil {
		return currentView.Navigation
	}
	return ""
}

// Exists checks if a view exists
func (vr *ViewRegistry) Exists(name string) bool {
	_, exists := vr.views[name]
	return exists
}
