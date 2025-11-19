package e2e

import (
	"testing"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wikczerski/whaletui/e2e/framework"
)

// TestSwarmServiceList tests listing swarm services.
func TestSwarmServiceList(t *testing.T) {
	fw := framework.NewTestFramework(t)
	dh := fw.GetDockerHelper()
	ctx := fw.GetContext()
	client := fw.GetDockerClient()

	// Create test services
	svc1 := dh.CreateTestService(
		framework.ServiceFixtures.WebService.Name,
		framework.ServiceFixtures.WebService.Image,
		framework.ServiceFixtures.WebService.Replicas,
	)
	svc2 := dh.CreateTestService(
		framework.ServiceFixtures.CacheService.Name,
		framework.ServiceFixtures.CacheService.Image,
		framework.ServiceFixtures.CacheService.Replicas,
	)

	// List services
	services, err := client.ServiceList(ctx, types.ServiceListOptions{})
	require.NoError(t, err, "Failed to list services")

	// Verify services exist
	serviceIDs := make([]string, 0)
	for _, s := range services {
		serviceIDs = append(serviceIDs, s.ID)
	}

	assert.Contains(t, serviceIDs, svc1, "Service 1 should be in list")
	assert.Contains(t, serviceIDs, svc2, "Service 2 should be in list")
}

// TestSwarmServiceScale tests scaling a swarm service.
func TestSwarmServiceScale(t *testing.T) {
	fw := framework.NewTestFramework(t)
	dh := fw.GetDockerHelper()

	// Create service with 1 replica
	serviceID := dh.CreateTestService(
		"e2e-test-service-scale",
		framework.ImageFixtures.Nginx,
		1,
	)

	// Verify initial replica count
	replicas := dh.GetServiceReplicas(serviceID)
	assert.Equal(t, uint64(1), replicas, "Service should start with 1 replica")

	// Scale to 3 replicas
	dh.ScaleService(serviceID, 3)
	time.Sleep(2 * time.Second) // Allow time for scaling

	// Verify new replica count
	replicas = dh.GetServiceReplicas(serviceID)
	assert.Equal(t, uint64(3), replicas, "Service should be scaled to 3 replicas")

	// Scale down to 1
	dh.ScaleService(serviceID, 1)
	time.Sleep(2 * time.Second)

	replicas = dh.GetServiceReplicas(serviceID)
	assert.Equal(t, uint64(1), replicas, "Service should be scaled down to 1 replica")
}

// TestSwarmServiceScaleToZero tests scaling a service to zero replicas.
func TestSwarmServiceScaleToZero(t *testing.T) {
	fw := framework.NewTestFramework(t)
	dh := fw.GetDockerHelper()

	// Create service
	serviceID := dh.CreateTestService(
		"e2e-test-service-scale-zero",
		framework.ImageFixtures.Redis,
		2,
	)

	// Scale to 0
	dh.ScaleService(serviceID, 0)
	time.Sleep(2 * time.Second)

	// Verify replica count is 0
	replicas := dh.GetServiceReplicas(serviceID)
	assert.Equal(t, uint64(0), replicas, "Service should be scaled to 0 replicas")
}

// TestSwarmServiceRemove tests removing a swarm service.
func TestSwarmServiceRemove(t *testing.T) {
	fw := framework.NewTestFramework(t)
	dh := fw.GetDockerHelper()
	ctx := fw.GetContext()
	client := fw.GetDockerClient()

	// Create service
	serviceID := dh.CreateTestService(
		"e2e-test-service-remove",
		framework.ImageFixtures.Alpine,
		1,
	)

	// Remove service
	dh.RemoveService(serviceID)

	// Verify service is gone
	_, _, err := client.ServiceInspectWithRaw(ctx, serviceID, types.ServiceInspectOptions{})
	assert.Error(t, err, "Service should be removed")
}

// TestSwarmServiceInspect tests inspecting a swarm service.
func TestSwarmServiceInspect(t *testing.T) {
	fw := framework.NewTestFramework(t)
	dh := fw.GetDockerHelper()
	ctx := fw.GetContext()
	client := fw.GetDockerClient()

	// Create service
	serviceName := "e2e-test-service-inspect"
	serviceID := dh.CreateTestService(serviceName, framework.ImageFixtures.Nginx, 2)

	// Inspect service
	service, _, err := client.ServiceInspectWithRaw(ctx, serviceID, types.ServiceInspectOptions{})
	require.NoError(t, err, "Failed to inspect service")

	// Verify inspect data
	assert.Equal(t, serviceID, service.ID, "Service ID should match")
	assert.Equal(t, serviceName, service.Spec.Name, "Service name should match")
	assert.Equal(t, uint64(2), *service.Spec.Mode.Replicated.Replicas, "Replicas should match")
}

// TestSwarmServiceUpdate tests updating a swarm service.
func TestSwarmServiceUpdate(t *testing.T) {
	fw := framework.NewTestFramework(t)
	dh := fw.GetDockerHelper()
	ctx := fw.GetContext()
	client := fw.GetDockerClient()

	// Create service
	serviceID := dh.CreateTestService(
		"e2e-test-service-update",
		framework.ImageFixtures.Nginx,
		1,
	)

	// Get current service spec
	service, _, err := client.ServiceInspectWithRaw(ctx, serviceID, types.ServiceInspectOptions{})
	require.NoError(t, err, "Failed to inspect service")

	// Update service (scale to 2)
	newReplicas := uint64(2)
	service.Spec.Mode.Replicated.Replicas = &newReplicas

	_, err = client.ServiceUpdate(ctx, serviceID, service.Version, service.Spec, types.ServiceUpdateOptions{})
	require.NoError(t, err, "Failed to update service")

	// Verify update
	time.Sleep(2 * time.Second)
	replicas := dh.GetServiceReplicas(serviceID)
	assert.Equal(t, uint64(2), replicas, "Service should be updated to 2 replicas")
}
