//nolint:max-public-structs
package shared

import (
	"time"
)

// Container represents a Docker container
type Container struct {
	Created     time.Time         `json:"created"`
	Labels      map[string]string `json:"labels"`
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Image       string            `json:"image"`
	Status      string            `json:"status"`
	State       string            `json:"state"`
	NetworkMode string            `json:"network_mode"`
	Ports       []string          `json:"ports"`
	Mounts      []string          `json:"mounts"`
	SizeRw      int64             `json:"size_rw"`
	SizeRootFs  int64             `json:"size_root_fs"`
}

// ContainerDetails represents detailed container information
type ContainerDetails struct {
	Labels          map[string]string `json:"labels"`
	NetworkSettings NetworkSettings   `json:"network_settings"`
	Command         string            `json:"command"`
	WorkingDir      string            `json:"working_dir"`
	Args            []string          `json:"args"`
	Entrypoint      []string          `json:"entrypoint"`
	Environment     []string          `json:"environment"`
	Mounts          []Mount           `json:"mounts"`
	Container
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
	Ports     map[string][]Port  `json:"ports"`
	Networks  map[string]Network `json:"networks"`
	IPAddress string             `json:"ip_address"`
	Gateway   string             `json:"gateway"`
}

// Port represents a port binding
type Port struct {
	HostIP   string `json:"host_ip"`
	HostPort string `json:"host_port"`
}

// Image represents a Docker image
type Image struct {
	Created     time.Time         `json:"created"`
	Labels      map[string]string `json:"labels"`
	ID          string            `json:"id"`
	Size        string            `json:"size"`
	ParentID    string            `json:"parent_id"`
	RepoTags    []string          `json:"repo_tags"`
	SharedSize  int64             `json:"shared_size"`
	VirtualSize int64             `json:"virtual_size"`
}

// ImageDetails represents detailed image information
type ImageDetails struct {
	Labels       map[string]string `json:"labels"`
	Architecture string            `json:"architecture"`
	OS           string            `json:"os"`
	Author       string            `json:"author"`
	Comment      string            `json:"comment"`
	Config       ImageConfig       `json:"config"`
	History      []ImageHistory    `json:"history"`
	Image
}

// ImageConfig represents image configuration
type ImageConfig struct {
	ExposedPorts map[string]struct{} `json:"exposed_ports"`
	Volumes      map[string]struct{} `json:"volumes"`
	Labels       map[string]string   `json:"labels"`
	User         string              `json:"user"`
	WorkingDir   string              `json:"working_dir"`
	Entrypoint   []string            `json:"entrypoint"`
	Cmd          []string            `json:"cmd"`
	Environment  []string            `json:"environment"`
}

// ImageHistory represents image layer history
type ImageHistory struct {
	Created    time.Time `json:"created"`
	CreatedBy  string    `json:"created_by"`
	Comment    string    `json:"comment"`
	EmptyLayer bool      `json:"empty_layer"`
}

// Volume represents a Docker volume
type Volume struct {
	Name       string            `json:"name"`
	Driver     string            `json:"driver"`
	Mountpoint string            `json:"mountpoint"`
	CreatedAt  time.Time         `json:"created_at"`
	Status     map[string]any    `json:"status"`
	Labels     map[string]string `json:"labels"`
	Scope      string            `json:"scope"`
	Size       string            `json:"size"`
}

// VolumeDetails represents detailed volume information
type VolumeDetails struct {
	Status map[string]any `json:"status"`
	Volume
}

// Network represents a Docker network
type Network struct {
	Created    time.Time         `json:"created"`
	Options    map[string]string `json:"options"`
	Labels     map[string]string `json:"labels"`
	ID         string            `json:"id"`
	Name       string            `json:"name"`
	Driver     string            `json:"driver"`
	Scope      string            `json:"scope"`
	IPAM       IPAM              `json:"ipam"`
	Internal   bool              `json:"internal"`
	Attachable bool              `json:"attachable"`
	Ingress    bool              `json:"ingress"`
	EnableIPv6 bool              `json:"enable_ipv6"`
}

// IPAM represents IP Address Management configuration
type IPAM struct {
	Driver  string            `json:"driver"`
	Options map[string]string `json:"options"`
	Config  []IPAMConfig      `json:"config"`
}

// IPAMConfig represents IPAM configuration
type IPAMConfig struct {
	Subnet  string `json:"subnet"`
	IPRange string `json:"ip_range"`
	Gateway string `json:"gateway"`
}

// NetworkDetails represents detailed network information
type NetworkDetails struct {
	Containers map[string]NetworkContainer `json:"containers"`
	Network
}

// NetworkContainer represents a container in a network
type NetworkContainer struct {
	Labels      map[string]string `json:"labels"`
	DriverOpts  map[string]string `json:"driver_opts"`
	IPAMConfig  map[string]string `json:"ipam_config"`
	Name        string            `json:"name"`
	EndpointID  string            `json:"endpoint_id"`
	MacAddress  string            `json:"mac_address"`
	IPv4Address string            `json:"ipv4_address"`
	IPv6Address string            `json:"ipv6_address"`
	NetworkID   string            `json:"network_id"`
}
