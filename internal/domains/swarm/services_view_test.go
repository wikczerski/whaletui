package swarm

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/wikczerski/whaletui/internal/config"
	uimocks "github.com/wikczerski/whaletui/internal/mocks/ui"
	"github.com/wikczerski/whaletui/internal/shared"
	uishared "github.com/wikczerski/whaletui/internal/ui/shared"
)

// TestNewServicesView tests the creation of a new services view
func TestNewServicesView(t *testing.T) {
	mockUI := uimocks.NewMockUIInterface(t)
	mockServiceService := &ServiceService{}
	mockModalManager := uimocks.NewMockModalManagerInterface(t)
	mockHeaderManager := uimocks.NewMockHeaderManagerInterface(t)

	// Mock GetThemeManager for character limits setup
	mockThemeManager := config.NewThemeManager("")
	mockUI.On("GetThemeManager").Return(mockThemeManager).Maybe()

	view := NewServicesView(mockUI, mockServiceService, mockModalManager, mockHeaderManager)

	assert.NotNil(t, view)
	assert.Equal(t, "swarm services", view.GetTitle())
	assert.Equal(
		t,
		[]string{"ID", "Name", "Image", "Mode", "Replicas", "Status", "Created"},
		view.GetHeaders(),
	)
}

// TestServicesView_GetActions tests the getActions method
func TestServicesView_GetActions(t *testing.T) {
	mockUI := uimocks.NewMockUIInterface(t)
	mockServiceService := &ServiceService{}
	mockModalManager := uimocks.NewMockModalManagerInterface(t)
	mockHeaderManager := uimocks.NewMockHeaderManagerInterface(t)

	// Mock GetThemeManager for character limits setup
	mockThemeManager := config.NewThemeManager("")
	mockUI.On("GetThemeManager").Return(mockThemeManager).Maybe()

	view := NewServicesView(mockUI, mockServiceService, mockModalManager, mockHeaderManager)

	actions := view.getActions()

	expectedActions := map[rune]string{
		'i': "Inspect",
		's': "Scale",
		'r': "Remove",
		'l': "Logs",
		'f': "Refresh",
	}

	assert.Equal(t, expectedActions, actions)
}

// TestServicesView_FormatServiceRow tests the service row formatting
func TestServicesView_FormatServiceRow(t *testing.T) {
	mockUI := uimocks.NewMockUIInterface(t)
	mockServiceService := &ServiceService{}
	mockModalManager := uimocks.NewMockModalManagerInterface(t)
	mockHeaderManager := uimocks.NewMockHeaderManagerInterface(t)

	// Mock GetThemeManager for character limits setup
	mockThemeManager := config.NewThemeManager("")
	mockUI.On("GetThemeManager").Return(mockThemeManager).Maybe()

	view := NewServicesView(mockUI, mockServiceService, mockModalManager, mockHeaderManager)

	service := shared.SwarmService{
		ID:        "test-service-id-123456789",
		Name:      "test-service",
		Image:     "nginx:latest",
		Mode:      "replicated",
		Replicas:  "3",
		Status:    "running",
		CreatedAt: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
	}

	row := view.FormatRow(service)

	expectedRow := []string{
		"test-service", // truncated ID
		"test-service",
		"nginx:latest",
		"replicated",
		"3",
		"running",
		"2023-01-01 12:00:00",
	}

	assert.Equal(t, expectedRow, row)
}

// TestServicesView_GetServiceID tests the service ID retrieval
func TestServicesView_GetServiceID(t *testing.T) {
	mockUI := uimocks.NewMockUIInterface(t)
	mockServiceService := &ServiceService{}
	mockModalManager := uimocks.NewMockModalManagerInterface(t)
	mockHeaderManager := uimocks.NewMockHeaderManagerInterface(t)

	// Mock GetThemeManager for character limits setup
	mockThemeManager := config.NewThemeManager("")
	mockUI.On("GetThemeManager").Return(mockThemeManager).Maybe()

	view := NewServicesView(mockUI, mockServiceService, mockModalManager, mockHeaderManager)

	service := shared.SwarmService{
		ID:   "test-service-id",
		Name: "test-service",
	}

	id := view.GetItemID(service)
	assert.Equal(t, "test-service-id", id)
}

// TestServicesView_GetServiceName tests the service name retrieval
func TestServicesView_GetServiceName(t *testing.T) {
	mockUI := uimocks.NewMockUIInterface(t)
	mockServiceService := &ServiceService{}
	mockModalManager := uimocks.NewMockModalManagerInterface(t)
	mockHeaderManager := uimocks.NewMockHeaderManagerInterface(t)

	// Mock GetThemeManager for character limits setup
	mockThemeManager := config.NewThemeManager("")
	mockUI.On("GetThemeManager").Return(mockThemeManager).Maybe()

	view := NewServicesView(mockUI, mockServiceService, mockModalManager, mockHeaderManager)

	service := shared.SwarmService{
		ID:   "test-service-id",
		Name: "test-service",
	}

	name := view.GetItemName(service)
	assert.Equal(t, "test-service", name)
}

// TestServicesView_SearchState tests search state management
func TestServicesView_SearchState(t *testing.T) {
	// Create a test view instance with proper initialization
	view := &ServicesView{}
	view.BaseView = uishared.NewBaseView[shared.SwarmService](
		nil, "Test Services", []string{"ID", "Name"},
	)

	// Test initial state
	assert.False(t, view.IsSearchActive())
	assert.Empty(t, view.GetSearchTerm())

	// Test setting search term
	view.Search("test")
	assert.True(t, view.IsSearchActive())
	assert.Equal(t, "test", view.GetSearchTerm())

	// Test clearing search
	view.ClearSearch()
	assert.False(t, view.IsSearchActive())
	assert.Empty(t, view.GetSearchTerm())
}

// TestServicesView_SearchWithMockData tests search functionality with mock data
func TestServicesView_SearchWithMockData(t *testing.T) {
	// Create a test view instance with proper initialization
	view := &ServicesView{}
	view.BaseView = uishared.NewBaseView[shared.SwarmService](
		nil, "Test Services", []string{"ID", "Name", "Image", "Mode", "Status"},
	)

	// Set up the callbacks that are needed for search to work
	view.FormatRow = func(s shared.SwarmService) []string { return view.formatServiceRow(&s) }
	view.GetItemID = func(s shared.SwarmService) string { return view.getServiceID(&s) }
	view.GetItemName = func(s shared.SwarmService) string { return view.getServiceName(&s) }

	// Set up the view with mock data
	testServices := []shared.SwarmService{
		{
			ID:        "service1",
			Name:      "web-service",
			Image:     "nginx:latest",
			Mode:      "replicated",
			Replicas:  "3",
			Status:    "running",
			CreatedAt: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
		},
		{
			ID:        "service2",
			Name:      "api-service",
			Image:     "node:18",
			Mode:      "global",
			Replicas:  "5",
			Status:    "running",
			CreatedAt: time.Date(2023, 1, 2, 12, 0, 0, 0, time.UTC),
		},
		{
			ID:        "service3",
			Name:      "database-service",
			Image:     "postgres:13",
			Mode:      "replicated",
			Replicas:  "1",
			Status:    "stopped",
			CreatedAt: time.Date(2023, 1, 3, 12, 0, 0, 0, time.UTC),
		},
	}

	// Set up the view with mock data by directly setting the items
	view.SetItems(testServices)

	// Test search by service name
	view.Search("web")
	assert.True(t, view.IsSearchActive())
	assert.Equal(t, "web", view.GetSearchTerm())

	// Get the filtered items to verify search worked
	filteredItems := view.GetItems()
	assert.Len(t, filteredItems, 1)
	assert.Equal(t, "web-service", view.getServiceName(&filteredItems[0]))

	// Test search by image
	view.Search("nginx")
	assert.True(t, view.IsSearchActive())
	assert.Equal(t, "nginx", view.GetSearchTerm())

	filteredItems = view.GetItems()
	assert.Len(t, filteredItems, 1)
	assert.Equal(t, "web-service", view.getServiceName(&filteredItems[0]))

	// Test search by mode
	view.Search("global")
	assert.True(t, view.IsSearchActive())
	assert.Equal(t, "global", view.GetSearchTerm())

	filteredItems = view.GetItems()
	assert.Len(t, filteredItems, 1)
	assert.Equal(t, "api-service", view.getServiceName(&filteredItems[0]))

	// Test search by status
	view.Search("stopped")
	assert.True(t, view.IsSearchActive())
	assert.Equal(t, "stopped", view.GetSearchTerm())

	filteredItems = view.GetItems()
	assert.Len(t, filteredItems, 1)
	assert.Equal(t, "database-service", view.getServiceName(&filteredItems[0]))

	// Test case insensitive search
	view.Search("WEB")
	assert.True(t, view.IsSearchActive())
	assert.Equal(t, "WEB", view.GetSearchTerm())

	filteredItems = view.GetItems()
	assert.Len(t, filteredItems, 1)
	assert.Equal(t, "web-service", view.getServiceName(&filteredItems[0]))

	// Test partial match
	view.Search("service")
	assert.True(t, view.IsSearchActive())
	assert.Equal(t, "service", view.GetSearchTerm())

	filteredItems = view.GetItems()
	assert.Len(t, filteredItems, 3)

	// Test no match
	view.Search("nonexistent")
	assert.True(t, view.IsSearchActive())
	assert.Equal(t, "nonexistent", view.GetSearchTerm())

	filteredItems = view.GetItems()
	assert.Len(t, filteredItems, 0)

	// Test clearing search
	view.ClearSearch()
	assert.False(t, view.IsSearchActive())
	assert.Empty(t, view.GetSearchTerm())

	filteredItems = view.GetItems()
	assert.Len(t, filteredItems, 3) // Should show all items again
}
