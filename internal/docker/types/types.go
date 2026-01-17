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
	Created    time.Time
	ID         string
	Repository string
	Tag        string
	Size       string
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
	Created    time.Time
	Labels     map[string]string
	ID         string
	Name       string
	Driver     string
	Scope      string
	Containers int
	Internal   bool
	Attachable bool
	Ingress    bool
	IPv6       bool
	EnableIPv6 bool
}
