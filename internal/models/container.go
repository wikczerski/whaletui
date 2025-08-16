package models

import "github.com/wikczerski/D5r/internal/docker"

type Container = docker.Container

// ContainerDetails represents detailed container information
type ContainerDetails struct {
	Container
	Command         string            `json:"command"`
	Args            []string          `json:"args"`
	WorkingDir      string            `json:"working_dir"`
	Entrypoint      []string          `json:"entrypoint"`
	Environment     []string          `json:"environment"`
	Labels          map[string]string `json:"labels"`
	Mounts          []Mount           `json:"mounts"`
	NetworkSettings NetworkSettings   `json:"network_settings"`
}

// Mount represents a container mount
type Mount struct {
	Type        string `json:"type"`
	Source      string `json:"source"`
	Destination string `json:"destination"`
	ReadOnly    bool   `json:"read_only"`
}

// NetworkSettings represents container network configuration
type NetworkSettings struct {
	IPAddress string             `json:"ip_address"`
	Gateway   string             `json:"gateway"`
	Ports     map[string][]Port  `json:"ports"`
	Networks  map[string]Network `json:"networks"`
}

// Port represents a port binding
type Port struct {
	HostIP   string `json:"host_ip"`
	HostPort string `json:"host_port"`
}
