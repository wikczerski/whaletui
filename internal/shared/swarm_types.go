package shared

import "time"

// SwarmService represents a Docker Swarm service
type SwarmService struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	ID        string
	Name      string
	Image     string
	Mode      string
	Replicas  string
	Status    string
	Ports     []string
}

// SwarmNode represents a Docker Swarm node
type SwarmNode struct {
	Labels        map[string]string
	ID            string
	Hostname      string
	Role          string
	Availability  string
	Status        string
	ManagerStatus string
	EngineVersion string
	Address       string
	CPUs          int64
	Memory        int64
}
