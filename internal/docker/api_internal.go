package docker

import (
	"context"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
)

// getContainerList retrieves the raw container list from Docker
func (c *Client) getContainerList(ctx context.Context, all bool) ([]container.Summary, error) {
	opts := container.ListOptions{All: all}
	return c.cli.ContainerList(ctx, opts)
}

// getImageList retrieves the raw image list from Docker
func (c *Client) getImageList(ctx context.Context) ([]image.Summary, error) {
	opts := image.ListOptions{}
	return c.cli.ImageList(ctx, opts)
}

// createExecInstance creates an exec instance in the container
func (c *Client) createExecInstance(
	ctx context.Context,
	id string,
	command []string,
) (*container.ExecCreateResponse, error) {
	execConfig := container.ExecOptions{
		Cmd:          command,
		Tty:          false, // Set to false to capture output
		AttachStdin:  false, // We don't need stdin for command execution
		AttachStdout: true,
		AttachStderr: true,
		Detach:       false,
	}

	execResp, err := c.cli.ContainerExecCreate(ctx, id, execConfig)
	if err != nil {
		return nil, err
	}

	return &execResp, nil
}

// executeAndCollectOutput executes the exec instance and collects the output
func (c *Client) executeAndCollectOutput(ctx context.Context, execID string) (string, error) {
	attachConfig := container.ExecStartOptions{
		Tty: false,
	}

	hijackedResp, err := c.cli.ContainerExecAttach(ctx, execID, attachConfig)
	if err != nil {
		return "", err
	}
	defer hijackedResp.Close()

	if err := c.cli.ContainerExecStart(ctx, execID, attachConfig); err != nil {
		return "", err
	}

	return c.readExecOutput(hijackedResp), nil
}
