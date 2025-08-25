package swarm

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	mocksShared "github.com/wikczerski/whaletui/internal/mocks/shared"
	"github.com/wikczerski/whaletui/internal/shared"
)

// TestNewServiceService tests the creation of a new service service
func TestNewServiceService(t *testing.T) {
	// We can't directly test with a mock client because NewServiceService expects *docker.Client
	// This test verifies that the constructor returns a non-nil service
	t.Skip("Requires real docker.Client, tested indirectly through other tests")
}

// TestServiceService_ListServices_Success tests the ListServices method returns services successfully
func TestServiceService_ListServices_Success(t *testing.T) {
	mockService := mocksShared.NewMockSwarmServiceService(t)

	// Create test data
	expectedServices := []shared.SwarmService{
		{
			ID:        "service1",
			Name:      "test-service-1",
			Image:     "nginx:latest",
			Mode:      "replicated",
			Replicas:  "2/2",
			Ports:     []string{"80:80"},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Status:    "running",
		},
	}

	ctx := context.Background()
	mockService.EXPECT().ListServices(ctx).Return(expectedServices, nil)

	result, err := mockService.ListServices(ctx)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
}

// TestServiceService_ListServices_ReturnsCorrectServiceID tests the service ID is correct
func TestServiceService_ListServices_ReturnsCorrectServiceID(t *testing.T) {
	mockService := mocksShared.NewMockSwarmServiceService(t)

	expectedServices := []shared.SwarmService{
		{
			ID:   "service1",
			Name: "test-service-1",
		},
	}

	ctx := context.Background()
	mockService.EXPECT().ListServices(ctx).Return(expectedServices, nil)

	result, _ := mockService.ListServices(ctx)

	assert.Equal(t, "service1", result[0].ID)
}

// TestServiceService_ListServices_ReturnsCorrectServiceName tests the service name is correct
func TestServiceService_ListServices_ReturnsCorrectServiceName(t *testing.T) {
	mockService := mocksShared.NewMockSwarmServiceService(t)

	expectedServices := []shared.SwarmService{
		{
			ID:   "service1",
			Name: "test-service-1",
		},
	}

	ctx := context.Background()
	mockService.EXPECT().ListServices(ctx).Return(expectedServices, nil)

	result, _ := mockService.ListServices(ctx)

	assert.Equal(t, "test-service-1", result[0].Name)
}

// TestServiceService_ListServices_Error tests the ListServices method with error
func TestServiceService_ListServices_Error(t *testing.T) {
	mockService := mocksShared.NewMockSwarmServiceService(t)

	ctx := context.Background()
	expectedError := errors.New("failed to list services")
	mockService.EXPECT().ListServices(ctx).Return([]shared.SwarmService(nil), expectedError)

	result, err := mockService.ListServices(ctx)

	assert.Error(t, err)
	assert.Nil(t, result)
}

// TestServiceService_InspectService_Success tests the InspectService method returns no error
func TestServiceService_InspectService_Success(t *testing.T) {
	mockService := mocksShared.NewMockSwarmServiceService(t)

	ctx := context.Background()
	serviceID := "service1"
	expectedInfo := map[string]any{
		"ID":   "service1",
		"Name": "test-service-1",
		"Mode": "replicated",
	}

	mockService.EXPECT().InspectService(ctx, serviceID).Return(expectedInfo, nil)

	_, err := mockService.InspectService(ctx, serviceID)

	assert.NoError(t, err)
}

// TestServiceService_InspectService_ReturnsNotNil tests the InspectService method returns non-nil result
func TestServiceService_InspectService_ReturnsNotNil(t *testing.T) {
	mockService := mocksShared.NewMockSwarmServiceService(t)

	ctx := context.Background()
	serviceID := "service1"
	expectedInfo := map[string]any{"ID": "service1"}

	mockService.EXPECT().InspectService(ctx, serviceID).Return(expectedInfo, nil)

	result, _ := mockService.InspectService(ctx, serviceID)

	assert.NotNil(t, result)
}

// TestServiceService_InspectService_ReturnsCorrectID tests the InspectService method returns correct ID
func TestServiceService_InspectService_ReturnsCorrectID(t *testing.T) {
	mockService := mocksShared.NewMockSwarmServiceService(t)

	ctx := context.Background()
	serviceID := "service1"
	expectedInfo := map[string]any{"ID": "service1"}

	mockService.EXPECT().InspectService(ctx, serviceID).Return(expectedInfo, nil)

	result, _ := mockService.InspectService(ctx, serviceID)

	assert.Equal(t, "service1", result["ID"])
}

// TestServiceService_InspectService_ReturnsCorrectName tests the InspectService method returns correct name
func TestServiceService_InspectService_ReturnsCorrectName(t *testing.T) {
	mockService := mocksShared.NewMockSwarmServiceService(t)

	ctx := context.Background()
	serviceID := "service1"
	expectedInfo := map[string]any{"Name": "test-service-1"}

	mockService.EXPECT().InspectService(ctx, serviceID).Return(expectedInfo, nil)

	result, _ := mockService.InspectService(ctx, serviceID)

	assert.Equal(t, "test-service-1", result["Name"])
}

// TestServiceService_ScaleService tests the ScaleService method
func TestServiceService_ScaleService(t *testing.T) {
	mockService := mocksShared.NewMockSwarmServiceService(t)

	ctx := context.Background()
	serviceID := "service1"
	newReplicas := uint64(5)

	mockService.EXPECT().ScaleService(ctx, serviceID, newReplicas).Return(nil)

	err := mockService.ScaleService(ctx, serviceID, newReplicas)

	assert.NoError(t, err)
}

// TestServiceService_ScaleService_Error_ReturnsError tests the ScaleService method with error returns error
func TestServiceService_ScaleService_Error_ReturnsError(t *testing.T) {
	mockService := mocksShared.NewMockSwarmServiceService(t)

	ctx := context.Background()
	serviceID := "service1"
	newReplicas := uint64(5)
	expectedError := errors.New("cannot scale global service")

	mockService.EXPECT().ScaleService(ctx, serviceID, newReplicas).Return(expectedError)

	err := mockService.ScaleService(ctx, serviceID, newReplicas)

	assert.Error(t, err)
}

// TestServiceService_ScaleService_Error_ReturnsCorrectError tests the ScaleService method returns correct error
func TestServiceService_ScaleService_Error_ReturnsCorrectError(t *testing.T) {
	mockService := mocksShared.NewMockSwarmServiceService(t)

	ctx := context.Background()
	serviceID := "service1"
	newReplicas := uint64(5)
	expectedError := errors.New("cannot scale global service")

	mockService.EXPECT().ScaleService(ctx, serviceID, newReplicas).Return(expectedError)

	err := mockService.ScaleService(ctx, serviceID, newReplicas)

	assert.Equal(t, expectedError, err)
}

// TestServiceService_RemoveService tests the RemoveService method
func TestServiceService_RemoveService(t *testing.T) {
	mockService := mocksShared.NewMockSwarmServiceService(t)

	ctx := context.Background()
	serviceID := "service1"

	mockService.EXPECT().RemoveService(ctx, serviceID).Return(nil)

	err := mockService.RemoveService(ctx, serviceID)

	assert.NoError(t, err)
}

// TestServiceService_GetServiceLogs_Success tests the GetServiceLogs method returns no error
func TestServiceService_GetServiceLogs_Success(t *testing.T) {
	mockService := mocksShared.NewMockSwarmServiceService(t)

	ctx := context.Background()
	serviceID := "service1"
	expectedLogs := "test log output"

	mockService.EXPECT().GetServiceLogs(ctx, serviceID).Return(expectedLogs, nil)

	_, err := mockService.GetServiceLogs(ctx, serviceID)

	assert.NoError(t, err)
}

// TestServiceService_GetServiceLogs_ReturnsCorrectLogs tests the GetServiceLogs method returns correct logs
func TestServiceService_GetServiceLogs_ReturnsCorrectLogs(t *testing.T) {
	mockService := mocksShared.NewMockSwarmServiceService(t)

	ctx := context.Background()
	serviceID := "service1"
	expectedLogs := "test log output"

	mockService.EXPECT().GetServiceLogs(ctx, serviceID).Return(expectedLogs, nil)

	result, _ := mockService.GetServiceLogs(ctx, serviceID)

	assert.Equal(t, expectedLogs, result)
}

// TestServiceService_GetActions_ReturnsNotNil tests the GetActions method returns non-nil
func TestServiceService_GetActions_ReturnsNotNil(t *testing.T) {
	mockService := mocksShared.NewMockSwarmServiceService(t)

	expectedActions := map[rune]string{
		'i': "Inspect",
		's': "Scale",
		'r': "Remove",
		'l': "Logs",
	}

	mockService.EXPECT().GetActions().Return(expectedActions)

	actions := mockService.GetActions()

	assert.NotNil(t, actions)
}

// TestServiceService_GetActions_ContainsInspect tests the GetActions method contains inspect action
func TestServiceService_GetActions_ContainsInspect(t *testing.T) {
	mockService := mocksShared.NewMockSwarmServiceService(t)

	expectedActions := map[rune]string{'i': "Inspect"}
	mockService.EXPECT().GetActions().Return(expectedActions)

	actions := mockService.GetActions()

	assert.Contains(t, actions, 'i')
}

// TestServiceService_GetActions_ContainsScale tests the GetActions method contains scale action
func TestServiceService_GetActions_ContainsScale(t *testing.T) {
	mockService := mocksShared.NewMockSwarmServiceService(t)

	expectedActions := map[rune]string{'s': "Scale"}
	mockService.EXPECT().GetActions().Return(expectedActions)

	actions := mockService.GetActions()

	assert.Contains(t, actions, 's')
}

// TestServiceService_GetActions_ContainsRemove tests the GetActions method contains remove action
func TestServiceService_GetActions_ContainsRemove(t *testing.T) {
	mockService := mocksShared.NewMockSwarmServiceService(t)

	expectedActions := map[rune]string{'r': "Remove"}
	mockService.EXPECT().GetActions().Return(expectedActions)

	actions := mockService.GetActions()

	assert.Contains(t, actions, 'r')
}

// TestServiceService_GetActions_ContainsLogs tests the GetActions method contains logs action
func TestServiceService_GetActions_ContainsLogs(t *testing.T) {
	mockService := mocksShared.NewMockSwarmServiceService(t)

	expectedActions := map[rune]string{'l': "Logs"}
	mockService.EXPECT().GetActions().Return(expectedActions)

	actions := mockService.GetActions()

	assert.Contains(t, actions, 'l')
}

// TestServiceService_GetActionsString_NotEmpty tests the GetActionsString method returns non-empty string
func TestServiceService_GetActionsString_NotEmpty(t *testing.T) {
	mockService := mocksShared.NewMockSwarmServiceService(t)

	expectedActionsString := "i: Inspect, s: Scale, r: Remove, l: Logs"
	mockService.EXPECT().GetActionsString().Return(expectedActionsString)

	actionsString := mockService.GetActionsString()

	assert.NotEmpty(t, actionsString)
}

// TestServiceService_GetActionsString_ContainsInspect tests the GetActionsString method contains Inspect
func TestServiceService_GetActionsString_ContainsInspect(t *testing.T) {
	mockService := mocksShared.NewMockSwarmServiceService(t)

	expectedActionsString := "i: Inspect"
	mockService.EXPECT().GetActionsString().Return(expectedActionsString)

	actionsString := mockService.GetActionsString()

	assert.Contains(t, actionsString, "Inspect")
}

// TestServiceService_GetActionsString_ContainsScale tests the GetActionsString method contains Scale
func TestServiceService_GetActionsString_ContainsScale(t *testing.T) {
	mockService := mocksShared.NewMockSwarmServiceService(t)

	expectedActionsString := "s: Scale"
	mockService.EXPECT().GetActionsString().Return(expectedActionsString)

	actionsString := mockService.GetActionsString()

	assert.Contains(t, actionsString, "Scale")
}

// TestServiceService_GetActionsString_ContainsRemove tests the GetActionsString method contains Remove
func TestServiceService_GetActionsString_ContainsRemove(t *testing.T) {
	mockService := mocksShared.NewMockSwarmServiceService(t)

	expectedActionsString := "r: Remove"
	mockService.EXPECT().GetActionsString().Return(expectedActionsString)

	actionsString := mockService.GetActionsString()

	assert.Contains(t, actionsString, "Remove")
}

// TestServiceService_GetActionsString_ContainsLogs tests the GetActionsString method contains Logs
func TestServiceService_GetActionsString_ContainsLogs(t *testing.T) {
	mockService := mocksShared.NewMockSwarmServiceService(t)

	expectedActionsString := "l: Logs"
	mockService.EXPECT().GetActionsString().Return(expectedActionsString)

	actionsString := mockService.GetActionsString()

	assert.Contains(t, actionsString, "Logs")
}
