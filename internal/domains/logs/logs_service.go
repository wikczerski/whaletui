package logs

import (
	"context"
	"errors"
	"fmt"

	"github.com/wikczerski/whaletui/internal/shared/interfaces"
)

// LogsService implements logs operations
type logsService struct {
	containerService interfaces.ContainerService
	swarmService     interfaces.SwarmServiceService
}

// NewLogsService creates a new logs service
func NewLogsService(
	containerService interfaces.ContainerService,
	swarmService interfaces.SwarmServiceService,
) interfaces.LogsService {
	return &logsService{
		containerService: containerService,
		swarmService:     swarmService,
	}
}

// GetLogs retrieves logs for a specific resource type and ID
func (s *logsService) GetLogs(
	ctx context.Context,
	resourceType, resourceID string,
) (string, error) {
	switch resourceType {
	case "container":
		if s.containerService == nil {
			return "", errors.New("container service not available")
		}
		return s.containerService.GetContainerLogs(ctx, resourceID)
	case "service":
		if s.swarmService == nil {
			return "", errors.New("swarm service not available")
		}
		return s.swarmService.GetServiceLogs(ctx, resourceID)
	default:
		return "", fmt.Errorf("unsupported resource type: %s", resourceType)
	}
}

// GetActions returns the available actions for logs as a map
func (s *logsService) GetActions() map[rune]string {
	return map[rune]string{
		'f': "Follow logs",
		't': "Tail logs",
		's': "Save logs",
		'c': "Clear logs",
		'w': "Wrap text",
	}
}

// GetActionsString returns the available actions for logs as a formatted string
func (s *logsService) GetActionsString() string {
	return "<f> Follow logs\n<t> Tail logs\n<s> Save logs\n<c> Clear logs\n<w> Wrap text"
}
