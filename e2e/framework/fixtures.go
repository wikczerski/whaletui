package framework

import (
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
)

// Fixtures provides predefined test data and configurations.

// ContainerFixtures contains predefined container configurations.
var ContainerFixtures = struct {
	Nginx    ContainerFixture
	Redis    ContainerFixture
	Postgres ContainerFixture
	Alpine   ContainerFixture
	Busybox  ContainerFixture
}{
	Nginx: ContainerFixture{
		Name:  "e2e-test-nginx",
		Image: "nginx:alpine",
		Config: &container.Config{
			Image: "nginx:alpine",
			ExposedPorts: nat.PortSet{
				"80/tcp": struct{}{},
			},
		},
		HostConfig: &container.HostConfig{
			PublishAllPorts: true,
		},
	},
	Redis: ContainerFixture{
		Name:  "e2e-test-redis",
		Image: "redis:alpine",
		Config: &container.Config{
			Image: "redis:alpine",
			ExposedPorts: nat.PortSet{
				"6379/tcp": struct{}{},
			},
		},
		HostConfig: &container.HostConfig{
			PublishAllPorts: true,
		},
	},
	Postgres: ContainerFixture{
		Name:  "e2e-test-postgres",
		Image: "postgres:13-alpine",
		Config: &container.Config{
			Image: "postgres:13-alpine",
			Env: []string{
				"POSTGRES_PASSWORD=test123",
				"POSTGRES_USER=testuser",
				"POSTGRES_DB=testdb",
			},
			ExposedPorts: nat.PortSet{
				"5432/tcp": struct{}{},
			},
		},
		HostConfig: &container.HostConfig{
			PublishAllPorts: true,
		},
	},
	Alpine: ContainerFixture{
		Name:  "e2e-test-alpine",
		Image: "alpine:latest",
		Config: &container.Config{
			Image: "alpine:latest",
			Cmd:   []string{"sleep", "3600"},
		},
		HostConfig: nil,
	},
	Busybox: ContainerFixture{
		Name:  "e2e-test-busybox",
		Image: "busybox:latest",
		Config: &container.Config{
			Image: "busybox:latest",
			Cmd:   []string{"sleep", "3600"},
		},
		HostConfig: nil,
	},
}

// ContainerFixture represents a container test fixture.
type ContainerFixture struct {
	Name       string
	Image      string
	Config     *container.Config
	HostConfig *container.HostConfig
}

// ImageFixtures contains predefined image names for testing.
var ImageFixtures = struct {
	Nginx    string
	Redis    string
	Postgres string
	Alpine   string
	Busybox  string
}{
	Nginx:    "nginx:alpine",
	Redis:    "redis:alpine",
	Postgres: "postgres:13-alpine",
	Alpine:   "alpine:latest",
	Busybox:  "busybox:latest",
}

// VolumeFixtures contains predefined volume configurations.
var VolumeFixtures = struct {
	Data1 VolumeFixture
	Data2 VolumeFixture
	Data3 VolumeFixture
}{
	Data1: VolumeFixture{
		Name:   "e2e-test-volume-1",
		Driver: "local",
	},
	Data2: VolumeFixture{
		Name:   "e2e-test-volume-2",
		Driver: "local",
	},
	Data3: VolumeFixture{
		Name:   "e2e-test-volume-3",
		Driver: "local",
	},
}

// VolumeFixture represents a volume test fixture.
type VolumeFixture struct {
	Name   string
	Driver string
}

// NetworkFixtures contains predefined network configurations.
var NetworkFixtures = struct {
	Bridge1 NetworkFixture
	Bridge2 NetworkFixture
}{
	Bridge1: NetworkFixture{
		Name:   "e2e-test-network-1",
		Driver: "bridge",
	},
	Bridge2: NetworkFixture{
		Name:   "e2e-test-network-2",
		Driver: "bridge",
	},
}

// NetworkFixture represents a network test fixture.
type NetworkFixture struct {
	Name   string
	Driver string
	Config *network.NetworkingConfig
}

// ServiceFixtures contains predefined swarm service configurations.
var ServiceFixtures = struct {
	WebService   ServiceFixture
	CacheService ServiceFixture
	DBService    ServiceFixture
}{
	WebService: ServiceFixture{
		Name:     "e2e-test-service-web",
		Image:    "nginx:alpine",
		Replicas: 2,
	},
	CacheService: ServiceFixture{
		Name:     "e2e-test-service-cache",
		Image:    "redis:alpine",
		Replicas: 1,
	},
	DBService: ServiceFixture{
		Name:     "e2e-test-service-db",
		Image:    "postgres:13-alpine",
		Replicas: 1,
	},
}

// ServiceFixture represents a swarm service test fixture.
type ServiceFixture struct {
	Name     string
	Image    string
	Replicas uint64
}

// TestData contains common test data.
var TestData = struct {
	SearchQueries   []string
	Commands        []string
	ShellCommands   []string
	InvalidCommands []string
}{
	SearchQueries: []string{
		"nginx",
		"test",
		"alpine",
		"redis",
	},
	Commands: []string{
		"containers",
		"images",
		"volumes",
		"networks",
		"services",
		"nodes",
		"quit",
		"help",
		"reload",
	},
	ShellCommands: []string{
		"ls",
		"pwd",
		"whoami",
		"echo 'test'",
		"cat /etc/os-release",
	},
	InvalidCommands: []string{
		"invalid",
		"xyz",
		"notacommand",
	},
}

// ErrorScenarios contains common error scenarios for testing.
var ErrorScenarios = struct {
	StartRunningContainer    string
	StopStoppedContainer     string
	DeleteRunningContainer   string
	DeleteImageInUse         string
	DeleteVolumeInUse        string
	DeleteNetworkInUse       string
	DeleteDefaultNetwork     string
	AttachToStoppedContainer string
	ExecInStoppedContainer   string
	ScaleToNegativeReplicas  string
}{
	StartRunningContainer:    "start_running",
	StopStoppedContainer:     "stop_stopped",
	DeleteRunningContainer:   "delete_running",
	DeleteImageInUse:         "delete_image_in_use",
	DeleteVolumeInUse:        "delete_volume_in_use",
	DeleteNetworkInUse:       "delete_network_in_use",
	DeleteDefaultNetwork:     "delete_default_network",
	AttachToStoppedContainer: "attach_stopped",
	ExecInStoppedContainer:   "exec_stopped",
	ScaleToNegativeReplicas:  "scale_negative",
}

// WaitTimeouts contains common timeout durations for tests.
var WaitTimeouts = struct {
	Short  int // milliseconds
	Medium int // milliseconds
	Long   int // milliseconds
}{
	Short:  1000,  // 1 second
	Medium: 5000,  // 5 seconds
	Long:   10000, // 10 seconds
}
