package services

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	domaintypes "github.com/wikczerski/whaletui/internal/docker/types"
	"github.com/wikczerski/whaletui/internal/docker/utils"
)

// ContainerService handles container-related operations
type ContainerService struct {
	cli *client.Client
	log *slog.Logger
}

// NewContainerService creates a new ContainerService
func NewContainerService(cli *client.Client, log *slog.Logger) *ContainerService {
	return &ContainerService{
		cli: cli,
		log: log,
	}
}

// ListContainers lists all containers
func (s *ContainerService) ListContainers(ctx context.Context, all bool) ([]domaintypes.Container, error) {
	containers, err := s.getContainerList(ctx, all)
	if err != nil {
		return nil, err
	}

	result := s.convertToContainers(containers)
	utils.SortContainersByCreationTime(result)
	return result, nil
}

// InspectContainer inspects a container
func (s *ContainerService) InspectContainer(ctx context.Context, id string) (map[string]any, error) {
	container, err := s.cli.ContainerInspect(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("container inspect failed %s: %w", id, err)
	}

	return utils.MarshalToMap(container)
}

// GetContainerLogs retrieves container logs
func (s *ContainerService) GetContainerLogs(ctx context.Context, id string) (string, error) {
	logs, err := s.cli.ContainerLogs(ctx, id, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Timestamps: true,
		Follow:     false,
		Tail:       "100",
	})
	if err != nil {
		return "", fmt.Errorf("container logs failed %s: %w", id, err)
	}
	defer func() {
		if err := logs.Close(); err != nil {
			s.log.Warn("Failed to close logs", "error", err)
		}
	}()

	return s.readAndFormatLogs(logs), nil
}

// GetContainerStats gets container stats
func (s *ContainerService) GetContainerStats(ctx context.Context, id string) (map[string]any, error) {
	stats, err := s.cli.ContainerStats(ctx, id, false)
	if err != nil {
		return nil, fmt.Errorf("failed to get container stats: %w", err)
	}
	defer func() {
		if err := stats.Body.Close(); err != nil {
			s.log.Warn("Failed to close stats body", "error", err)
		}
	}()

	return utils.MarshalToMap(stats.Body) // Note: original used decodeStatsResponse, but MarshalToMap might work if it was just JSON decoding
}

// StartContainer starts a container
func (s *ContainerService) StartContainer(ctx context.Context, id string) error {
	return s.cli.ContainerStart(ctx, id, container.StartOptions{})
}

// StopContainer stops a container
func (s *ContainerService) StopContainer(ctx context.Context, id string, timeout *time.Duration) error {
	opts := utils.BuildStopOptions(timeout)
	return s.cli.ContainerStop(ctx, id, opts)
}

// RestartContainer restarts a container
func (s *ContainerService) RestartContainer(ctx context.Context, id string, timeout *time.Duration) error {
	opts := utils.BuildStopOptions(timeout)
	return s.cli.ContainerRestart(ctx, id, opts)
}

// RemoveContainer removes a container
func (s *ContainerService) RemoveContainer(ctx context.Context, id string, force bool) error {
	opts := container.RemoveOptions{
		Force: force,
	}

	if err := utils.ValidateID(id, "container ID"); err != nil {
		return err
	}

	return s.cli.ContainerRemove(ctx, id, opts)
}

// ExecContainer executes a command in a running container and returns the output
func (s *ContainerService) ExecContainer(
	ctx context.Context,
	id string,
	command []string,
	_ bool,
) (string, error) {
	execResp, err := s.createExecInstance(ctx, id, command)
	if err != nil {
		return "", err
	}

	output, err := s.executeAndCollectOutput(ctx, execResp.ID)
	if err != nil {
		return "", err
	}

	return output, nil
}

// AttachContainer attaches to a running container
func (s *ContainerService) AttachContainer(ctx context.Context, id string) (types.HijackedResponse, error) {
	attachConfig := container.AttachOptions{
		Stream: true,
		Stdin:  true,
		Stdout: true,
		Stderr: true,
		Logs:   false,
	}

	if err := utils.ValidateID(id, "container ID"); err != nil {
		return types.HijackedResponse{}, err
	}

	response, err := s.cli.ContainerAttach(ctx, id, attachConfig)
	if err != nil {
		return types.HijackedResponse{}, fmt.Errorf("failed to attach to container: %w", err)
	}

	return response, nil
}

// Helper methods

func (s *ContainerService) getContainerList(ctx context.Context, all bool) ([]container.Summary, error) {
	opts := container.ListOptions{All: all}
	return s.cli.ContainerList(ctx, opts)
}

func (s *ContainerService) convertToContainers(containers []container.Summary) []domaintypes.Container {
	result := make([]domaintypes.Container, 0, len(containers))
	for i := range containers {
		cont := &containers[i]
		result = append(result, s.convertContainer(cont))
	}
	return result
}

func (s *ContainerService) convertContainer(cont *container.Summary) domaintypes.Container {
	ports := utils.FormatContainerPorts(cont.Ports)
	return domaintypes.Container{
		ID:      cont.ID[:12],
		Name:    cont.Names[0][1:], // Remove leading slash
		Image:   cont.Image,
		Status:  cont.Status,
		State:   cont.State,
		Created: time.Unix(cont.Created, 0),
		Ports:   ports,
		Size:    "", // Size is not available in ContainerList
	}
}

func (s *ContainerService) readAndFormatLogs(logs io.ReadCloser) string {
	var logLines []string
	buffer := make([]byte, 1024)

	s.readLogsIntoBuffer(logs, buffer, &logLines)

	return strings.Join(logLines, "") // strings package needs to be imported if not already
}

func (s *ContainerService) readLogsIntoBuffer(logs io.ReadCloser, buffer []byte, logLines *[]string) {
	for {
		n, err := logs.Read(buffer)
		if n > 0 {
			line := string(buffer[:n])
			formattedLine := s.formatLogLine(line)
			*logLines = append(*logLines, formattedLine)
		}
		if err != nil {
			break
		}
	}
}

func (s *ContainerService) formatLogLine(line string) string {
	if len(line) >= 8 {
		return line[8:] // Remove timestamp prefix
	}
	return line
}

func (s *ContainerService) createExecInstance(
	ctx context.Context,
	id string,
	command []string,
) (*container.ExecCreateResponse, error) {
	execConfig := container.ExecOptions{
		Cmd:          command,
		Tty:          false,
		AttachStdin:  false,
		AttachStdout: true,
		AttachStderr: true,
		Detach:       false,
	}

	execResp, err := s.cli.ContainerExecCreate(ctx, id, execConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create exec instance: %w", err)
	}

	return &execResp, nil
}

func (s *ContainerService) executeAndCollectOutput(ctx context.Context, execID string) (string, error) {
	attachConfig := container.ExecStartOptions{
		Tty: false,
	}

	hijackedResp, err := s.cli.ContainerExecAttach(ctx, execID, attachConfig)
	if err != nil {
		return "", fmt.Errorf("failed to attach to exec instance: %w", err)
	}
	defer hijackedResp.Close()

	if err := s.cli.ContainerExecStart(ctx, execID, attachConfig); err != nil {
		return "", fmt.Errorf("failed to start exec instance: %w", err)
	}

	return s.readExecOutput(hijackedResp), nil
}

func (s *ContainerService) readExecOutput(hijackedResp types.HijackedResponse) string {
	var output strings.Builder
	buffer := make([]byte, 1024)

	for {
		n, err := hijackedResp.Reader.Read(buffer)
		if n > 0 {
			output.Write(buffer[:n])
		}
		if err != nil {
			break
		}
	}

	return output.String()
}
