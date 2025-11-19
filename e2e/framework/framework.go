// Package framework provides core testing utilities for WhaleTUI e2e tests.
package framework

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/docker/docker/client"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/require"
	"github.com/wikczerski/whaletui/internal/app"
)

// TestFramework provides the core testing infrastructure for e2e tests.
type TestFramework struct {
	t              *testing.T
	dockerClient   *client.Client
	app            *app.App
	tviewApp       *tview.Application
	ctx            context.Context
	cancel         context.CancelFunc
	dockerHelper   *DockerHelper
	tuiHelper      *TUIHelper
	cleanupFuncs   []func()
	testContainers []string
	testImages     []string
	testVolumes    []string
	testNetworks   []string
	testServices   []string
}

// NewTestFramework creates a new test framework instance.
func NewTestFramework(t *testing.T) *TestFramework {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)

	// Initialize Docker client
	dockerClient, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	require.NoError(t, err, "Failed to create Docker client")

	fw := &TestFramework{
		t:              t,
		dockerClient:   dockerClient,
		ctx:            ctx,
		cancel:         cancel,
		cleanupFuncs:   make([]func(), 0),
		testContainers: make([]string, 0),
		testImages:     make([]string, 0),
		testVolumes:    make([]string, 0),
		testNetworks:   make([]string, 0),
		testServices:   make([]string, 0),
	}

	// Initialize helpers
	fw.dockerHelper = NewDockerHelper(fw)
	fw.tuiHelper = NewTUIHelper(fw)

	// Register cleanup
	t.Cleanup(func() {
		fw.TearDown()
	})

	return fw
}

// SetupApp initializes the WhaleTUI application for testing.
// Note: Full TUI initialization is skipped for Docker-only tests.
func (fw *TestFramework) SetupApp() {
	fw.t.Helper()

	// For now, we skip full app initialization as most tests
	// only need Docker client functionality.
	// Full TUI tests would require additional setup.
	fw.t.Log("App setup skipped - using Docker client only")
}

// StartApp starts the application in a goroutine.
// Note: Currently not implemented for Docker-only tests.
func (fw *TestFramework) StartApp() {
	fw.t.Helper()
	fw.t.Log("StartApp skipped - using Docker client only")
}

// SimulateKey simulates a keyboard key press.
func (fw *TestFramework) SimulateKey(key tcell.Key, ch rune, mod tcell.ModMask) {
	fw.t.Helper()

	event := tcell.NewEventKey(key, ch, mod)
	fw.tviewApp.QueueEvent(event)

	// Allow time for event processing
	time.Sleep(100 * time.Millisecond)
}

// SimulateKeyRune simulates a rune key press.
func (fw *TestFramework) SimulateKeyRune(ch rune) {
	fw.SimulateKey(tcell.KeyRune, ch, tcell.ModNone)
}

// SimulateKeyPress simulates a special key press.
func (fw *TestFramework) SimulateKeyPress(key tcell.Key) {
	fw.SimulateKey(key, 0, tcell.ModNone)
}

// SimulateKeySequence simulates a sequence of key presses.
func (fw *TestFramework) SimulateKeySequence(keys []tcell.Key, runes []rune) {
	fw.t.Helper()

	for _, key := range keys {
		fw.SimulateKeyPress(key)
	}

	for _, ch := range runes {
		fw.SimulateKeyRune(ch)
	}
}

// WaitForCondition waits for a condition to be true or timeout.
func (fw *TestFramework) WaitForCondition(
	condition func() bool,
	timeout time.Duration,
	message string,
) {
	fw.t.Helper()

	deadline := time.Now().Add(timeout)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		if condition() {
			return
		}

		if time.Now().After(deadline) {
			fw.t.Fatalf("Timeout waiting for condition: %s", message)
		}

		<-ticker.C
	}
}

// WaitForConditionWithError waits for a condition and returns error instead of failing.
func (fw *TestFramework) WaitForConditionWithError(
	condition func() bool,
	timeout time.Duration,
) error {
	deadline := time.Now().Add(timeout)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		if condition() {
			return nil
		}

		if time.Now().After(deadline) {
			return errors.New("timeout waiting for condition")
		}

		<-ticker.C
	}
}

// AddCleanup adds a cleanup function to be called during teardown.
func (fw *TestFramework) AddCleanup(cleanup func()) {
	fw.cleanupFuncs = append(fw.cleanupFuncs, cleanup)
}

// TearDown cleans up all resources created during the test.
func (fw *TestFramework) TearDown() {
	fw.t.Helper()

	// Run cleanup functions in reverse order
	for i := len(fw.cleanupFuncs) - 1; i >= 0; i-- {
		fw.cleanupFuncs[i]()
	}

	// Cleanup Docker resources
	fw.dockerHelper.CleanupAll()

	// Cancel context
	if fw.cancel != nil {
		fw.cancel()
	}

	// Close Docker client
	if fw.dockerClient != nil {
		_ = fw.dockerClient.Close()
	}
}

// GetDockerClient returns the Docker client.
func (fw *TestFramework) GetDockerClient() *client.Client {
	return fw.dockerClient
}

// GetContext returns the test context.
func (fw *TestFramework) GetContext() context.Context {
	return fw.ctx
}

// GetApp returns the WhaleTUI application.
func (fw *TestFramework) GetApp() *app.App {
	return fw.app
}

// GetTviewApp returns the tview application.
func (fw *TestFramework) GetTviewApp() *tview.Application {
	return fw.tviewApp
}

// GetDockerHelper returns the Docker helper.
func (fw *TestFramework) GetDockerHelper() *DockerHelper {
	return fw.dockerHelper
}

// GetTUIHelper returns the TUI helper.
func (fw *TestFramework) GetTUIHelper() *TUIHelper {
	return fw.tuiHelper
}

// RegisterTestContainer registers a container for cleanup.
func (fw *TestFramework) RegisterTestContainer(id string) {
	fw.testContainers = append(fw.testContainers, id)
}

// RegisterTestImage registers an image for cleanup.
func (fw *TestFramework) RegisterTestImage(id string) {
	fw.testImages = append(fw.testImages, id)
}

// RegisterTestVolume registers a volume for cleanup.
func (fw *TestFramework) RegisterTestVolume(name string) {
	fw.testVolumes = append(fw.testVolumes, name)
}

// RegisterTestNetwork registers a network for cleanup.
func (fw *TestFramework) RegisterTestNetwork(id string) {
	fw.testNetworks = append(fw.testNetworks, id)
}

// RegisterTestService registers a swarm service for cleanup.
func (fw *TestFramework) RegisterTestService(id string) {
	fw.testServices = append(fw.testServices, id)
}

// Sleep pauses execution for the specified duration.
func (fw *TestFramework) Sleep(duration time.Duration) {
	time.Sleep(duration)
}

// Logf logs a formatted message.
func (fw *TestFramework) Logf(
	format string,
	args ...any,
) {
	fw.t.Logf(format, args...)
}
