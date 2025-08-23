package shared

import (
	"time"
)

// Container represents a Docker container
type Container struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Image       string            `json:"image"`
	Status      string            `json:"status"`
	Created     time.Time         `json:"created"`
	Ports       []string          `json:"ports"`
	SizeRw      int64             `json:"size_rw"`
	SizeRootFs  int64             `json:"size_root_fs"`
	Labels      map[string]string `json:"labels"`
	State       string            `json:"state"`
	NetworkMode string            `json:"network_mode"`
	Mounts      []string          `json:"mounts"`
}

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
type Image struct {
	ID          string            `json:"id"`
	RepoTags    []string          `json:"repo_tags"`
	Created     time.Time         `json:"created"`
	Size        string            `json:"size"`
	SharedSize  int64             `json:"shared_size"`
	VirtualSize int64             `json:"virtual_size"`
	Labels      map[string]string `json:"labels"`
	ParentID    string            `json:"parent_id"`
}

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
	Volume
	Status map[string]interface{} `json:"status"`
}

// Network represents a Docker network
type Network struct {
	ID         string            `json:"id"`
	Name       string            `json:"name"`
	Driver     string            `json:"driver"`
	Scope      string            `json:"scope"`
	IPAM       IPAM              `json:"ipam"`
	Internal   bool              `json:"internal"`
	Attachable bool              `json:"attachable"`
	Ingress    bool              `json:"ingress"`
	EnableIPv6 bool              `json:"enable_ipv6"`
	Options    map[string]string `json:"options"`
	Labels     map[string]string `json:"labels"`
	Created    time.Time         `json:"created"`
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
