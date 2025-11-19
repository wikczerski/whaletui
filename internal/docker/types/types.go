package types

import "time"

// Container represents a Docker container
type Container struct {
	ID      string
	Name    string
	Image   string
	Status  string
	State   string
	Created time.Time
	Ports   string
	Size    string
}

// Image represents a Docker image
type Image struct {
	ID         string
	Repository string
	Tag        string
	Size       string
	Created    time.Time
	Containers int
}

// Volume represents a Docker volume
type Volume struct {
	Name       string
	Driver     string
	Mountpoint string
	Created    time.Time
	Labels     map[string]string
	Scope      string
	Size       string
}

// Network represents a Docker network
type Network struct {
	ID         string
	Name       string
	Driver     string
	Scope      string
	Created    time.Time
	Internal   bool
	Attachable bool
	Ingress    bool
	IPv6       bool
	EnableIPv6 bool
	Labels     map[string]string
	Containers int
}
