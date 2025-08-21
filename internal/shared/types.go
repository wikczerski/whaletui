package shared

import (
	"time"

	"github.com/wikczerski/whaletui/internal/docker"
)

// Container represents a Docker container
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

// Image represents a Docker image
type Image = docker.Image

// ImageDetails represents detailed image information
type ImageDetails struct {
	Image
	Architecture string            `json:"architecture"`
	OS           string            `json:"os"`
	Author       string            `json:"author"`
	Comment      string            `json:"comment"`
	Config       ImageConfig       `json:"config"`
	History      []ImageHistory    `json:"history"`
	Labels       map[string]string `json:"labels"`
}

// ImageConfig represents image configuration
type ImageConfig struct {
	User         string              `json:"user"`
	WorkingDir   string              `json:"working_dir"`
	Entrypoint   []string            `json:"entrypoint"`
	Cmd          []string            `json:"cmd"`
	Environment  []string            `json:"environment"`
	ExposedPorts map[string]struct{} `json:"exposed_ports"`
	Volumes      map[string]struct{} `json:"volumes"`
	Labels       map[string]string   `json:"labels"`
}

// ImageHistory represents image layer history
type ImageHistory struct {
	Created    time.Time `json:"created"`
	CreatedBy  string    `json:"created_by"`
	Comment    string    `json:"comment"`
	EmptyLayer bool      `json:"empty_layer"`
}

// Volume represents a Docker volume
type Volume = docker.Volume

// VolumeDetails represents detailed volume information
type VolumeDetails struct {
	Volume
	Status map[string]interface{} `json:"status"`
}

// Network represents a Docker network
type Network = docker.Network

// NetworkDetails represents detailed network information
type NetworkDetails struct {
	Network
	Containers map[string]NetworkContainer `json:"containers"`
}

// NetworkContainer represents a container in a network
type NetworkContainer struct {
	Name        string            `json:"name"`
	EndpointID  string            `json:"endpoint_id"`
	MacAddress  string            `json:"mac_address"`
	IPv4Address string            `json:"ipv4_address"`
	IPv6Address string            `json:"ipv6_address"`
	Labels      map[string]string `json:"labels"`
	NetworkID   string            `json:"network_id"`
	DriverOpts  map[string]string `json:"driver_opts"`
	IPAMConfig  map[string]string `json:"ipam_config"`
}
