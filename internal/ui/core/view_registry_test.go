package core

import (
	"testing"

	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
)

func TestViewRegistry_NewViewRegistry_Creation(t *testing.T) {
	// Test creating a new view registry
	vr := NewViewRegistry()
	assert.NotNil(t, vr)
}

func TestViewRegistry_NewViewRegistry_ViewsMapInitialized(t *testing.T) {
	// Test creating a new view registry
	vr := NewViewRegistry()

	// Test initial state
	assert.NotNil(t, vr.views)
}

func TestViewRegistry_NewViewRegistry_ViewsMapEmpty(t *testing.T) {
	// Test creating a new view registry
	vr := NewViewRegistry()

	if len(vr.views) != 0 {
		t.Errorf("Expected 0 views initially, got %d", len(vr.views))
	}
}

func TestViewRegistry_NewViewRegistry_CurrentViewEmpty(t *testing.T) {
	// Test creating a new view registry
	vr := NewViewRegistry()

	if vr.currentView != "" {
		t.Errorf("Expected empty current view initially, got '%s'", vr.currentView)
	}
}

func TestViewRegistry_Register_Count(t *testing.T) {
	// Test registering views
	vr := NewViewRegistry()

	// Create a mock view
	mockView := tview.NewTextView()

	// Test registering a view
	vr.Register("test", "Test View", 't', mockView, func() {}, "Test actions", "Test navigation")

	// Verify view was registered
	assert.Equal(t, 1, len(vr.views))
}

func TestViewRegistry_Register_ViewInfoExists(t *testing.T) {
	// Test registering views
	vr := NewViewRegistry()

	// Create a mock view
	mockView := tview.NewTextView()

	// Test registering a view
	vr.Register("test", "Test View", 't', mockView, func() {}, "Test actions", "Test navigation")

	// Verify view info
	viewInfo := vr.views["test"]
	assert.NotNil(t, viewInfo)
}

func TestViewRegistry_Register_ViewInfoName(t *testing.T) {
	// Test registering views
	vr := NewViewRegistry()

	// Create a mock view
	mockView := tview.NewTextView()

	// Test registering a view
	vr.Register("test", "Test View", 't', mockView, func() {}, "Test actions", "Test navigation")

	// Verify view info
	viewInfo := vr.views["test"]
	assert.Equal(t, "test", viewInfo.Name)
}

func TestViewRegistry_Register_ViewInfoTitle(t *testing.T) {
	// Test registering views
	vr := NewViewRegistry()

	// Create a mock view
	mockView := tview.NewTextView()

	// Test registering a view
	vr.Register("test", "Test View", 't', mockView, func() {}, "Test actions", "Test navigation")

	// Verify view info
	viewInfo := vr.views["test"]
	assert.Equal(t, "Test View", viewInfo.Title)
}

func TestViewRegistry_Register_ViewInfoShortcut(t *testing.T) {
	// Test registering views
	vr := NewViewRegistry()

	// Create a mock view
	mockView := tview.NewTextView()

	// Test registering a view
	vr.Register("test", "Test View", 't', mockView, func() {}, "Test actions", "Test navigation")

	// Verify view info
	viewInfo := vr.views["test"]
	assert.Equal(t, 't', viewInfo.Shortcut)
}

func TestViewRegistry_Register_ViewInfoView(t *testing.T) {
	// Test registering views
	vr := NewViewRegistry()

	// Create a mock view
	mockView := tview.NewTextView()

	// Test registering a view
	vr.Register("test", "Test View", 't', mockView, func() {}, "Test actions", "Test navigation")

	// Verify view info
	viewInfo := vr.views["test"]
	assert.Equal(t, mockView, viewInfo.View)
}

func TestViewRegistry_Register_ViewInfoActions(t *testing.T) {
	// Test registering views
	vr := NewViewRegistry()

	// Create a mock view
	mockView := tview.NewTextView()

	// Test registering a view
	vr.Register("test", "Test View", 't', mockView, func() {}, "Test actions", "Test navigation")

	// Verify view info
	viewInfo := vr.views["test"]
	assert.Equal(t, "Test actions", viewInfo.Actions)
}

func TestViewRegistry_SetCurrent_ExistingView1(t *testing.T) {
	// Test setting current view
	vr := NewViewRegistry()

	// Register views
	mockView1 := tview.NewTextView()
	mockView2 := tview.NewTextView()
	vr.Register("view1", "View 1", '1', mockView1, func() {}, "Actions 1", "Navigation 1")
	vr.Register("view2", "View 2", '2', mockView2, func() {}, "Actions 2", "Navigation 2")

	// Test setting current view to existing view
	vr.SetCurrent("view1")
	assert.Equal(t, "view1", vr.currentView)
}

func TestViewRegistry_SetCurrent_ExistingView2(t *testing.T) {
	// Test setting current view
	vr := NewViewRegistry()

	// Register views
	mockView1 := tview.NewTextView()
	mockView2 := tview.NewTextView()
	vr.Register("view1", "View 1", '1', mockView1, func() {}, "Actions 1", "Navigation 1")
	vr.Register("view2", "View 2", '2', mockView2, func() {}, "Actions 2", "Navigation 2")

	// Test setting current view to another existing view
	vr.SetCurrent("view2")
	assert.Equal(t, "view2", vr.currentView)
}

func TestViewRegistry_SetCurrent_NonExistentView(t *testing.T) {
	// Test setting current view
	vr := NewViewRegistry()

	// Register views
	mockView1 := tview.NewTextView()
	mockView2 := tview.NewTextView()
	vr.Register("view1", "View 1", '1', mockView1, func() {}, "Actions 1", "Navigation 1")
	vr.Register("view2", "View 2", '2', mockView2, func() {}, "Actions 2", "Navigation 2")

	// Set current view to view2 first
	vr.SetCurrent("view2")
	assert.Equal(t, "view2", vr.currentView)

	// Test setting current view to non-existent view (should not change)
	vr.SetCurrent("nonexistent")
	assert.Equal(t, "view2", vr.currentView)
}

func TestViewRegistry_GetCurrent_InitiallyNil(t *testing.T) {
	// Test getting current view info
	vr := NewViewRegistry()

	// Initially no current view
	currentView := vr.GetCurrent()
	assert.Nil(t, currentView)
}

func TestViewRegistry_GetCurrent_AfterSetCurrent(t *testing.T) {
	// Test getting current view info
	vr := NewViewRegistry()

	// Register and set current view
	mockView := tview.NewTextView()
	vr.Register("test", "Test View", 't', mockView, func() {}, "Test actions", "Test navigation")
	vr.SetCurrent("test")

	// Test getting current view
	currentView := vr.GetCurrent()
	assert.NotNil(t, currentView)
}

func TestViewRegistry_GetCurrent_AfterSetCurrent_Name(t *testing.T) {
	// Test getting current view info
	vr := NewViewRegistry()

	// Register and set current view
	mockView := tview.NewTextView()
	vr.Register("test", "Test View", 't', mockView, func() {}, "Test actions", "Test navigation")
	vr.SetCurrent("test")

	// Test getting current view
	currentView := vr.GetCurrent()
	assert.Equal(t, "test", currentView.Name)
}

func TestViewRegistry_GetCurrentActionsString_InitiallyEmpty(t *testing.T) {
	// Test getting current view actions string
	vr := NewViewRegistry()

	// Initially no current view
	actions := vr.GetCurrentActionsString()
	assert.Equal(t, "", actions)
}

func TestViewRegistry_GetCurrentActionsString_AfterSetCurrent(t *testing.T) {
	// Test getting current view actions string
	vr := NewViewRegistry()

	// Register and set current view
	mockView := tview.NewTextView()
	vr.Register("test", "Test View", 't', mockView, func() {}, "Test actions", "Test navigation")
	vr.SetCurrent("test")

	// Test getting current view actions
	actions := vr.GetCurrentActionsString()
	assert.Equal(t, "Test actions", actions)
}

func TestViewRegistry_Exists_InitiallyFalse(t *testing.T) {
	// Test checking if views exist
	vr := NewViewRegistry()

	// Initially no views exist
	assert.False(t, vr.Exists("test"))
}

func TestViewRegistry_Exists_AfterRegistration(t *testing.T) {
	// Test checking if views exist
	vr := NewViewRegistry()

	// Register a view
	mockView := tview.NewTextView()
	vr.Register("test", "Test View", 't', mockView, func() {}, "Test actions", "Test navigation")

	// Test existing view
	assert.True(t, vr.Exists("test"))
}

func TestViewRegistry_Exists_NonExistentView(t *testing.T) {
	// Test checking if views exist
	vr := NewViewRegistry()

	// Register a view
	mockView := tview.NewTextView()
	vr.Register("test", "Test View", 't', mockView, func() {}, "Test actions", "Test navigation")

	// Test non-existent view
	assert.False(t, vr.Exists("nonexistent"))
}

func TestViewRegistry_MultipleRegistrations_Count(t *testing.T) {
	// Test multiple view registrations
	vr := NewViewRegistry()

	// Register multiple views
	mockViews := []tview.Primitive{
		tview.NewTextView(),
		tview.NewTextView(),
		tview.NewTextView(),
		tview.NewTextView(),
	}

	viewNames := []string{"view1", "view2", "view3", "view4"}
	shortcuts := []rune{'1', '2', '3', '4'}
	actions := []string{"Actions 1", "Actions 2", "Actions 3", "Actions 4"}
	navigations := []string{"Navigation 1", "Navigation 2", "Navigation 3", "Navigation 4"}

	for i, name := range viewNames {
		vr.Register(
			name,
			"View "+string(rune('1'+i)),
			shortcuts[i],
			mockViews[i],
			func() {},
			actions[i],
			navigations[i],
		)
	}

	// Verify all views were registered
	assert.Equal(t, 4, len(vr.views))
}

func TestViewRegistry_EdgeCases_EmptyName(t *testing.T) {
	// Test edge cases and boundary conditions
	vr := NewViewRegistry()

	// Test registering view with empty name
	vr.Register(
		"",
		"Empty Name",
		' ',
		tview.NewTextView(),
		func() {},
		"Empty Actions",
		"Empty Navigation",
	)

	assert.True(t, vr.Exists(""))
}
