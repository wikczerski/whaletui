package shared

import "time"

// SwarmService represents a Docker Swarm service
type SwarmService struct {
	ID        string
	Name      string
	Image     string
	Mode      string // replicated or global
	Replicas  string // e.g., "3/3" for replicated, "global" for global
	Ports     []string
	CreatedAt time.Time
	UpdatedAt time.Time
	Status    string // running, updating, failed, etc.
}

// SwarmNode represents a Docker Swarm node
type SwarmNode struct {
	ID            string
	Hostname      string
	Role          string // manager or worker
	Availability  string // active, pause, drain
	Status        string // ready, down, unknown
	ManagerStatus string // leader, reachable, unreachable (for manager nodes)
	EngineVersion string
	Address       string
	CPUs          int64
	Memory        int64
	Labels        map[string]string
}
