package swarm

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/wikczerski/whaletui/internal/domains/swarm"
	uimocks "github.com/wikczerski/whaletui/internal/mocks/ui"
	"github.com/wikczerski/whaletui/internal/shared"
)

// TestNewServicesView tests the creation of a new services view
func TestNewServicesView(t *testing.T) {
	mockUI := uimocks.NewMockUIInterface(t)
	mockServiceService := &swarm.ServiceService{}
	mockModalManager := uimocks.NewMockModalManagerInterface(t)
	mockHeaderManager := uimocks.NewMockHeaderManagerInterface(t)

	view := NewServicesView(mockUI, mockServiceService, mockModalManager, mockHeaderManager)

	assert.NotNil(t, view)
	assert.Equal(t, "Swarm Services", view.GetTitle())
	assert.Equal(t, []string{"ID", "Name", "Image", "Mode", "Replicas", "Status", "Created"}, view.GetHeaders())
}

// TestServicesView_GetActions tests the getActions method
func TestServicesView_GetActions(t *testing.T) {
	mockUI := uimocks.NewMockUIInterface(t)
	mockServiceService := &swarm.ServiceService{}
	mockModalManager := uimocks.NewMockModalManagerInterface(t)
	mockHeaderManager := uimocks.NewMockHeaderManagerInterface(t)

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
	mockServiceService := &swarm.ServiceService{}
	mockModalManager := uimocks.NewMockModalManagerInterface(t)
	mockHeaderManager := uimocks.NewMockHeaderManagerInterface(t)

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
	mockServiceService := &swarm.ServiceService{}
	mockModalManager := uimocks.NewMockModalManagerInterface(t)
	mockHeaderManager := uimocks.NewMockHeaderManagerInterface(t)

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
	mockServiceService := &swarm.ServiceService{}
	mockModalManager := uimocks.NewMockModalManagerInterface(t)
	mockHeaderManager := uimocks.NewMockHeaderManagerInterface(t)

	view := NewServicesView(mockUI, mockServiceService, mockModalManager, mockHeaderManager)

	service := shared.SwarmService{
		ID:   "test-service-id",
		Name: "test-service",
	}

	name := view.GetItemName(service)
	assert.Equal(t, "test-service", name)
}

// TestServicesView_HandleInput_Refresh tests the refresh input handling
func TestServicesView_HandleInput_Refresh(t *testing.T) {
	mockUI := uimocks.NewMockUIInterface(t)
	mockServiceService := &swarm.ServiceService{}
	mockModalManager := uimocks.NewMockModalManagerInterface(t)
	mockHeaderManager := uimocks.NewMockHeaderManagerInterface(t)

	view := NewServicesView(mockUI, mockServiceService, mockModalManager, mockHeaderManager)

	ctx := context.Background()
	result, err := view.HandleInput(ctx, 'f')

	assert.NoError(t, err)
	assert.Equal(t, view, result)
}

// TestServicesView_HandleInput_Help tests the help input handling
func TestServicesView_HandleInput_Help(t *testing.T) {
	mockUI := uimocks.NewMockUIInterface(t)
	mockServiceService := &swarm.ServiceService{}
	mockModalManager := uimocks.NewMockModalManagerInterface(t)
	mockHeaderManager := uimocks.NewMockHeaderManagerInterface(t)

	// Set up mock expectations
	mockUI.EXPECT().ShowContextualHelp("swarm_services", "").Return()

	view := NewServicesView(mockUI, mockServiceService, mockModalManager, mockHeaderManager)

	ctx := context.Background()
	result, err := view.HandleInput(ctx, 'h')

	assert.NoError(t, err)
	assert.Equal(t, view, result)
}

// TestServicesView_HandleInput_Quit tests the quit input handling
func TestServicesView_HandleInput_Quit(t *testing.T) {
	mockUI := uimocks.NewMockUIInterface(t)
	mockServiceService := &swarm.ServiceService{}
	mockModalManager := uimocks.NewMockModalManagerInterface(t)
	mockHeaderManager := uimocks.NewMockHeaderManagerInterface(t)

	view := NewServicesView(mockUI, mockServiceService, mockModalManager, mockHeaderManager)

	ctx := context.Background()
	result, err := view.HandleInput(ctx, 'q')

	assert.Error(t, err)
	assert.Equal(t, "main menu view not implemented yet", err.Error())
	assert.Equal(t, view, result)
}

// TestServicesView_HandleInput_NavigateToNodes tests the navigation input handling
func TestServicesView_HandleInput_NavigateToNodes(t *testing.T) {
	mockUI := uimocks.NewMockUIInterface(t)
	mockServiceService := &swarm.ServiceService{}
	mockModalManager := uimocks.NewMockModalManagerInterface(t)
	mockHeaderManager := uimocks.NewMockHeaderManagerInterface(t)

	view := NewServicesView(mockUI, mockServiceService, mockModalManager, mockHeaderManager)

	ctx := context.Background()
	result, err := view.HandleInput(ctx, 'n')

	assert.Error(t, err)
	assert.Equal(t, "nodes view not implemented yet", err.Error())
	assert.Equal(t, view, result)
}

// TestServicesView_HandleInput_Default tests the default input handling
func TestServicesView_HandleInput_Default(t *testing.T) {
	mockUI := uimocks.NewMockUIInterface(t)
	mockServiceService := &swarm.ServiceService{}
	mockModalManager := uimocks.NewMockModalManagerInterface(t)
	mockHeaderManager := uimocks.NewMockHeaderManagerInterface(t)

	view := NewServicesView(mockUI, mockServiceService, mockModalManager, mockHeaderManager)

	ctx := context.Background()
	result, err := view.HandleInput(ctx, 'x') // unknown key

	assert.NoError(t, err)
	assert.Equal(t, view, result)
}
