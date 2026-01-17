package docker

import (
	"context"
	"errors"

	"github.com/docker/docker/api/types/swarm"
)

// ListSwarmServices lists all swarm services
func (c *Client) ListSwarmServices(ctx context.Context) ([]swarm.Service, error) {
	if c.Swarm == nil {
		return nil, errors.New("swarm service not initialized")
	}
	return c.Swarm.ListSwarmServices(ctx)
}

// InspectSwarmService inspects a swarm service
func (c *Client) InspectSwarmService(ctx context.Context, id string) (swarm.Service, error) {
	if c.Swarm == nil {
		return swarm.Service{}, errors.New("swarm service not initialized")
	}
	return c.Swarm.InspectSwarmService(ctx, id)
}

// UpdateSwarmService updates a swarm service
func (c *Client) UpdateSwarmService(
	ctx context.Context,
	id string,
	version swarm.Version,
	spec swarm.ServiceSpec,
) error {
	if c.Swarm == nil {
		return errors.New("swarm service not initialized")
	}
	return c.Swarm.UpdateSwarmService(ctx, id, version, spec)
}

// RemoveSwarmService removes a swarm service
func (c *Client) RemoveSwarmService(ctx context.Context, id string) error {
	if c.Swarm == nil {
		return errors.New("swarm service not initialized")
	}
	return c.Swarm.RemoveSwarmService(ctx, id)
}

// GetSwarmServiceLogs gets logs for a swarm service
func (c *Client) GetSwarmServiceLogs(ctx context.Context, id string) (string, error) {
	if c.Swarm == nil {
		return "", errors.New("swarm service not initialized")
	}
	return c.Swarm.GetSwarmServiceLogs(ctx, id)
}

// ListSwarmNodes lists all swarm nodes
func (c *Client) ListSwarmNodes(ctx context.Context) ([]swarm.Node, error) {
	if c.Swarm == nil {
		return nil, errors.New("swarm service not initialized")
	}
	return c.Swarm.ListSwarmNodes(ctx)
}

// InspectSwarmNode inspects a swarm node
func (c *Client) InspectSwarmNode(ctx context.Context, id string) (swarm.Node, error) {
	if c.Swarm == nil {
		return swarm.Node{}, errors.New("swarm service not initialized")
	}
	return c.Swarm.InspectSwarmNode(ctx, id)
}

// UpdateSwarmNode updates a swarm node
func (c *Client) UpdateSwarmNode(
	ctx context.Context,
	id string,
	version swarm.Version,
	spec swarm.NodeSpec,
) error {
	if c.Swarm == nil {
		return errors.New("swarm service not initialized")
	}
	return c.Swarm.UpdateSwarmNode(ctx, id, version, spec)
}

// RemoveSwarmNode removes a swarm node
func (c *Client) RemoveSwarmNode(ctx context.Context, id string, force bool) error {
	if c.Swarm == nil {
		return errors.New("swarm service not initialized")
	}
	return c.Swarm.RemoveSwarmNode(ctx, id, force)
}
