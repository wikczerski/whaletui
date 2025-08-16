package models

import "github.com/user/d5r/internal/docker"

type Network = docker.Network

// NetworkDetails represents detailed network information
type NetworkDetails struct {
	Network
	IPAM       IPAMConfig                  `json:"ipam"`
	ConfigFrom NetworkConfigFrom           `json:"config_from"`
	ConfigOnly bool                        `json:"config_only"`
	Containers map[string]NetworkContainer `json:"containers"`
	Options    map[string]string           `json:"options"`
}

// IPAMConfig represents IP address management configuration
type IPAMConfig struct {
	Driver  string            `json:"driver"`
	Options map[string]string `json:"options"`
	Config  []IPAMConfigEntry `json:"config"`
}

// IPAMConfigEntry represents an IPAM configuration entry
type IPAMConfigEntry struct {
	Subnet  string `json:"subnet"`
	Gateway string `json:"gateway"`
	IPRange string `json:"ip_range"`
}

// NetworkConfigFrom represents network configuration source
type NetworkConfigFrom struct {
	Network string `json:"network"`
}

// NetworkContainer represents a container connected to a network
type NetworkContainer struct {
	Name        string `json:"name"`
	EndpointID  string `json:"endpoint_id"`
	MacAddress  string `json:"mac_address"`
	IPv4Address string `json:"ipv4_address"`
	IPv6Address string `json:"ipv6_address"`
}
