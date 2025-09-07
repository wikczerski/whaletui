package container

import (
	"context"
	"log/slog"
	"time"

	"github.com/wikczerski/whaletui/internal/docker"
	"github.com/wikczerski/whaletui/internal/logger"
	"github.com/wikczerski/whaletui/internal/shared"
	"github.com/wikczerski/whaletui/internal/shared/interfaces"
)

type containerService struct {
	*shared.BaseService[shared.Container]
	operations *shared.CommonOperations
	log        *slog.Logger
}

// NewContainerService creates a new container service
func NewContainerService(client *docker.Client) interfaces.ContainerService {
	base := shared.NewBaseService[shared.Container](client, "container")
	ops := shared.NewCommonOperations(client)

	base.ListFunc = createListFunction()
	base.RemoveFunc = createRemoveFunction()
	base.InspectFunc = createInspectFunction()

	return &containerService{
		BaseService: base,
		operations:  ops,
		log:         logger.GetLogger(),
	}
}

// createListFunction creates the list function for containers
func createListFunction() func(client *docker.Client, ctx context.Context) ([]shared.Container, error) {
	return func(client *docker.Client, ctx context.Context) ([]shared.Container, error) {
		dockerContainers, err := client.ListContainers(ctx, true)
		if err != nil {
			return nil, err
		}

		return convertDockerContainers(dockerContainers), nil
	}
}

// convertDockerContainers converts Docker API containers to shared containers
func convertDockerContainers(dockerContainers []docker.Container) []shared.Container {
	result := make([]shared.Container, len(dockerContainers))
	for i := range dockerContainers {
		result[i] = convertDockerContainer(dockerContainers[i])
	}
	return result
}

// convertDockerContainer converts a single Docker API container to shared container
func convertDockerContainer(dockerContainer docker.Container) shared.Container {
	return shared.Container{
		ID:          dockerContainer.ID,
		Name:        dockerContainer.Name,
		Image:       dockerContainer.Image,
		Status:      dockerContainer.Status,
		Created:     dockerContainer.Created,
		Ports:       []string{dockerContainer.Ports},
		SizeRw:      0, // docker.Container.Size is string, shared.Container.SizeRw is int64
		SizeRootFs:  0, // docker.Container doesn't have SizeRootFs
		Labels:      make(map[string]string),
		State:       dockerContainer.State,
		NetworkMode: "",
		Mounts:      []string{},
	}
}

// createRemoveFunction creates the remove function for containers
func createRemoveFunction() func(client *docker.Client, ctx context.Context, id string, force bool) error {
	return func(client *docker.Client, ctx context.Context, id string, force bool) error {
		return client.RemoveContainer(ctx, id, force)
	}
}

// createInspectFunction creates the inspect function for containers
func createInspectFunction() func(client *docker.Client, ctx context.Context, id string) (map[string]any, error) {
	return func(client *docker.Client, ctx context.Context, id string) (map[string]any, error) {
		return client.InspectContainer(ctx, id)
	}
}

func (s *containerService) ListContainers(ctx context.Context) ([]shared.Container, error) {
	return s.List(ctx)
}

func (s *containerService) StartContainer(ctx context.Context, id string) error {
	return s.operations.StartContainer(ctx, id)
}

func (s *containerService) StopContainer(
	ctx context.Context,
	id string,
	timeout *time.Duration,
) error {
	return s.operations.StopContainer(ctx, id, timeout)
}

func (s *containerService) RestartContainer(
	ctx context.Context,
	id string,
	timeout *time.Duration,
) error {
	return s.operations.RestartContainer(ctx, id, timeout)
}

func (s *containerService) RemoveContainer(ctx context.Context, id string, force bool) error {
	return s.Remove(ctx, id, force)
}

func (s *containerService) GetContainerLogs(ctx context.Context, id string) (string, error) {
	return s.operations.GetContainerLogs(ctx, id)
}

func (s *containerService) InspectContainer(
	ctx context.Context,
	id string,
) (map[string]any, error) {
	return s.Inspect(ctx, id)
}

func (s *containerService) ExecContainer(
	ctx context.Context,
	id string,
	command []string,
	tty bool,
) (string, error) {
	return s.operations.ExecContainer(ctx, id, command, tty)
}

func (s *containerService) AttachContainer(ctx context.Context, id string) (any, error) {
	return s.operations.AttachContainer(ctx, id)
}

// GetActions returns the available actions for containers as a map
func (s *containerService) GetActions() map[rune]string {
	return map[rune]string{
		's': "Start",
		'S': "Stop",
		'r': "Restart",
		'd': "Delete",
		'a': "Attach",
		'l': "View Logs",
		'i': "Inspect",
		'n': "New",
		'e': "Exec",
		'f': "Filter",
		't': "Sort",
		'h': "History",
	}
}

// GetActionsString returns the available actions for containers as a formatted string
func (s *containerService) GetActionsString() string {
	actions := "<s> Start\n" +
		"<S> Stop\n" +
		"<r> Restart\n" +
		"<d> Delete\n" +
		"<a> Attach\n" +
		"<l> Logs\n" +
		"<i> Inspect\n" +
		"<n> New\n" +
		"<e> Exec\n" +
		"<f> Filter\n" +
		"<t> Sort\n" +
		"<h> History\n" +
		"<enter> Details\n" +
		"<:> Command"
	return actions
}

// GetNavigation returns the available navigation options for containers as a map
// This is where navigation for the current view starts
func (s *containerService) GetNavigation() map[rune]string {
	return map[rune]string{
		'↑': "Up",
		'↓': "Down",
		':': "Command",
		'/': "Filter",
	}
}

// GetNavigationString returns the available navigation options for containers as a formatted string
// This is where navigation for the current view starts
func (s *containerService) GetNavigationString() string {
	return "↑/↓: Navigate\n<:> Command mode\n/: Filter"
}
