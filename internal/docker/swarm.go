package docker

import (
	"context"
	"fmt"
	"io"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/swarm"
	"github.com/wikczerski/whaletui/internal/docker/utils"
)

// ListSwarmServices lists all swarm services
func (c *Client) ListSwarmServices(ctx context.Context) ([]swarm.Service, error) {
	if err := c.validateClient(); err != nil {
		return nil, err
	}

	services, err := c.cli.ServiceList(ctx, swarm.ServiceListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list swarm services: %w", err)
	}

	return services, nil
}

// InspectSwarmService inspects a swarm service
func (c *Client) InspectSwarmService(ctx context.Context, id string) (swarm.Service, error) {
	if err := c.validateClient(); err != nil {
		return swarm.Service{}, err
	}

	if err := utils.ValidateID(id, "service ID"); err != nil {
		return swarm.Service{}, err
	}

	service, _, err := c.cli.ServiceInspectWithRaw(ctx, id, swarm.ServiceInspectOptions{})
	if err != nil {
		return swarm.Service{}, fmt.Errorf("failed to inspect swarm service %s: %w", id, err)
	}

	return service, nil
}

// UpdateSwarmService updates a swarm service
// nolint:gocritic // Docker API requires value parameter for compatibility
func (c *Client) UpdateSwarmService(
	ctx context.Context,
	id string,
	version swarm.Version,
	spec swarm.ServiceSpec,
) error {
	if err := c.validateClient(); err != nil {
		return err
	}

	if err := utils.ValidateID(id, "service ID"); err != nil {
		return err
	}

	response, err := c.cli.ServiceUpdate(ctx, id, version, spec, swarm.ServiceUpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to update swarm service %s: %w", id, err)
	}

	if len(response.Warnings) > 0 {
		c.log.Warn("Service update warnings", "service_id", id, "warnings", response.Warnings)
	}

	return nil
}

// RemoveSwarmService removes a swarm service
func (c *Client) RemoveSwarmService(ctx context.Context, id string) error {
	if err := c.validateClient(); err != nil {
		return err
	}

	if err := utils.ValidateID(id, "service ID"); err != nil {
		return err
	}

	if err := c.cli.ServiceRemove(ctx, id); err != nil {
		return fmt.Errorf("failed to remove swarm service %s: %w", id, err)
	}

	return nil
}

// GetSwarmServiceLogs gets logs for a swarm service
func (c *Client) GetSwarmServiceLogs(ctx context.Context, id string) (string, error) {
	if err := c.validateServiceLogsRequest(id); err != nil {
		return "", err
	}

	response, err := c.getServiceLogsResponse(ctx, id)
	if err != nil {
		return "", err
	}
	defer c.closeServiceLogsResponse(response)

	return c.readServiceLogs(response), nil
}

// ListSwarmNodes lists all swarm nodes
func (c *Client) ListSwarmNodes(ctx context.Context) ([]swarm.Node, error) {
	if err := c.validateClient(); err != nil {
		return nil, err
	}

	nodes, err := c.cli.NodeList(ctx, swarm.NodeListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list swarm nodes: %w", err)
	}

	return nodes, nil
}

// InspectSwarmNode inspects a swarm node
func (c *Client) InspectSwarmNode(ctx context.Context, id string) (swarm.Node, error) {
	if err := c.validateClient(); err != nil {
		return swarm.Node{}, err
	}

	if err := utils.ValidateID(id, "node ID"); err != nil {
		return swarm.Node{}, err
	}

	node, _, err := c.cli.NodeInspectWithRaw(ctx, id)
	if err != nil {
		return swarm.Node{}, fmt.Errorf("failed to inspect swarm node %s: %w", id, err)
	}

	return node, nil
}

// UpdateSwarmNode updates a swarm node
func (c *Client) UpdateSwarmNode(
	ctx context.Context,
	id string,
	version swarm.Version,
	spec swarm.NodeSpec,
) error {
	if err := c.validateClient(); err != nil {
		return err
	}

	if err := utils.ValidateID(id, "node ID"); err != nil {
		return err
	}

	if err := c.cli.NodeUpdate(ctx, id, version, spec); err != nil {
		return fmt.Errorf("failed to update swarm node %s: %w", id, err)
	}

	return nil
}

// RemoveSwarmNode removes a swarm node
func (c *Client) RemoveSwarmNode(ctx context.Context, id string, force bool) error {
	if err := c.validateClient(); err != nil {
		return err
	}

	if err := utils.ValidateID(id, "node ID"); err != nil {
		return err
	}

	options := swarm.NodeRemoveOptions{
		Force: force,
	}

	if err := c.cli.NodeRemove(ctx, id, options); err != nil {
		return fmt.Errorf("failed to remove swarm node %s: %w", id, err)
	}

	return nil
}

// validateServiceLogsRequest validates the service logs request
func (c *Client) validateServiceLogsRequest(id string) error {
	if err := c.validateClient(); err != nil {
		return err
	}
	return utils.ValidateID(id, "service ID")
}

// getServiceLogsResponse gets the service logs response from Docker
func (c *Client) getServiceLogsResponse(ctx context.Context, id string) (io.ReadCloser, error) {
	options := container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Timestamps: true,
		Follow:     false,
	}

	response, err := c.cli.ServiceLogs(ctx, id, options)
	if err != nil {
		return nil, fmt.Errorf("failed to get logs for swarm service %s: %w", id, err)
	}
	return response, nil
}

// closeServiceLogsResponse safely closes the service logs response
func (c *Client) closeServiceLogsResponse(response io.ReadCloser) {
	if err := response.Close(); err != nil {
		c.log.Warn("Failed to close response", "error", err)
	}
}
