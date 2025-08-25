// Package network provides Docker network management functionality for WhaleTUI.
package network

import "github.com/wikczerski/whaletui/internal/shared"

// Network represents a Docker network
type Network = shared.Network

// Details represents detailed network information
type Details struct {
	Network
	IPAM       IPAMConfig           `json:"ipam"`
	ConfigFrom ConfigFrom           `json:"config_from"`
	ConfigOnly bool                 `json:"config_only"`
	Containers map[string]Container `json:"containers"`
	Options    map[string]string    `json:"options"`
}
