// Package network provides Docker network management functionality for WhaleTUI.
package network

import (
	"github.com/wikczerski/whaletui/internal/shared"
)

// Network represents a Docker network
type Network = shared.Network

// Details represents detailed network information
type Details struct {
	Labels      map[string]string `json:"labels"`
	Options     map[string]string `json:"options"`
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Driver      string            `json:"driver"`
	Scope       string            `json:"scope"`
	IPv6Enabled bool              `json:"ipv6_enabled"`
	Internal    bool              `json:"internal"`
}
