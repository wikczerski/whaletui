// Package network provides Docker network management functionality for WhaleTUI.
package network

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

// ConfigFrom represents network configuration source
type ConfigFrom struct {
	Network string `json:"network"`
}

// Container represents a container connected to a network
type Container struct {
	Name        string `json:"name"`
	EndpointID  string `json:"endpoint_id"`
	MacAddress  string `json:"mac_address"`
	IPv4Address string `json:"ipv4_address"`
	IPv6Address string `json:"ipv6_address"`
}
