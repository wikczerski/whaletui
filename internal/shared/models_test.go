package shared

import (
	"encoding/json"
	"testing"
)

func TestContainerDetails_JSON(t *testing.T) {
	containerDetails := ContainerDetails{
		Container: Container{
			ID:     "test-container-id",
			Name:   "test-container",
			Image:  "test-image:latest",
			Status: "running",
		},
		Command:     "/bin/bash",
		Args:        []string{"-c", "echo hello"},
		WorkingDir:  "/app",
		Entrypoint:  []string{"/bin/bash"},
		Environment: []string{"ENV=test", "DEBUG=true"},
		Labels: map[string]string{
			"app":     "test-app",
			"version": "1.0.0",
		},
		Mounts: []Mount{
			{
				Type:        "bind",
				Source:      "/host/path",
				Destination: "/container/path",
				ReadOnly:    false,
			},
		},
		NetworkSettings: NetworkSettings{
			IPAddress: "172.17.0.2",
			Gateway:   "172.17.0.1",
			Ports: map[string][]Port{
				"80/tcp": {
					{
						HostIP:   "0.0.0.0",
						HostPort: "8080",
					},
				},
			},
			Networks: map[string]Network{
				"bridge": {
					ID:   "bridge-id",
					Name: "bridge",
				},
			},
		},
	}

	// Test JSON marshaling
	data, err := json.Marshal(containerDetails)
	if err != nil {
		t.Fatalf("Failed to marshal ContainerDetails: %v", err)
	}

	// Test JSON unmarshaling
	var unmarshaled ContainerDetails
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal ContainerDetails: %v", err)
	}

	// Verify the data
	if unmarshaled.ID != containerDetails.ID {
		t.Errorf("ID mismatch: got %s, want %s", unmarshaled.ID, containerDetails.ID)
	}
	if unmarshaled.Command != containerDetails.Command {
		t.Errorf("Command mismatch: got %s, want %s", unmarshaled.Command, containerDetails.Command)
	}
	if len(unmarshaled.Args) != len(containerDetails.Args) {
		t.Errorf(
			"Args length mismatch: got %d, want %d",
			len(unmarshaled.Args),
			len(containerDetails.Args),
		)
	}
	if unmarshaled.WorkingDir != containerDetails.WorkingDir {
		t.Errorf(
			"WorkingDir mismatch: got %s, want %s",
			unmarshaled.WorkingDir,
			containerDetails.WorkingDir,
		)
	}
	if len(unmarshaled.Environment) != len(containerDetails.Environment) {
		t.Errorf(
			"Environment length mismatch: got %d, want %d",
			len(unmarshaled.Environment),
			len(containerDetails.Environment),
		)
	}
	if len(unmarshaled.Labels) != len(containerDetails.Labels) {
		t.Errorf(
			"Labels length mismatch: got %d, want %d",
			len(unmarshaled.Labels),
			len(containerDetails.Labels),
		)
	}
	if len(unmarshaled.Mounts) != len(containerDetails.Mounts) {
		t.Errorf(
			"Mounts length mismatch: got %d, want %d",
			len(unmarshaled.Mounts),
			len(containerDetails.Mounts),
		)
	}
}

func TestMount_JSON(t *testing.T) {
	mount := Mount{
		Type:        "bind",
		Source:      "/host/path",
		Destination: "/container/path",
		ReadOnly:    true,
	}

	// Test JSON marshaling
	data, err := json.Marshal(mount)
	if err != nil {
		t.Fatalf("Failed to marshal Mount: %v", err)
	}

	// Test JSON unmarshaling
	var unmarshaled Mount
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal Mount: %v", err)
	}

	// Verify the data
	if unmarshaled.Type != mount.Type {
		t.Errorf("Type mismatch: got %s, want %s", unmarshaled.Type, mount.Type)
	}
	if unmarshaled.Source != mount.Source {
		t.Errorf("Source mismatch: got %s, want %s", unmarshaled.Source, mount.Source)
	}
	if unmarshaled.Destination != mount.Destination {
		t.Errorf(
			"Destination mismatch: got %s, want %s",
			unmarshaled.Destination,
			mount.Destination,
		)
	}
	if unmarshaled.ReadOnly != mount.ReadOnly {
		t.Errorf("ReadOnly mismatch: got %t, want %t", unmarshaled.ReadOnly, mount.ReadOnly)
	}
}

func TestNetworkSettings_JSON(t *testing.T) {
	networkSettings := NetworkSettings{
		IPAddress: "172.17.0.2",
		Gateway:   "172.17.0.1",
		Ports: map[string][]Port{
			"80/tcp": {
				{
					HostIP:   "0.0.0.0",
					HostPort: "8080",
				},
				{
					HostIP:   "127.0.0.1",
					HostPort: "8081",
				},
			},
			"443/tcp": {
				{
					HostIP:   "0.0.0.0",
					HostPort: "8443",
				},
			},
		},
		Networks: map[string]Network{
			"bridge": {
				ID:   "bridge-id",
				Name: "bridge",
			},
			"custom": {
				ID:   "custom-id",
				Name: "custom",
			},
		},
	}

	// Test JSON marshaling
	data, err := json.Marshal(networkSettings)
	if err != nil {
		t.Fatalf("Failed to marshal NetworkSettings: %v", err)
	}

	// Test JSON unmarshaling
	var unmarshaled NetworkSettings
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal NetworkSettings: %v", err)
	}

	// Verify the data
	if unmarshaled.IPAddress != networkSettings.IPAddress {
		t.Errorf(
			"IPAddress mismatch: got %s, want %s",
			unmarshaled.IPAddress,
			networkSettings.IPAddress,
		)
	}
	if unmarshaled.Gateway != networkSettings.Gateway {
		t.Errorf("Gateway mismatch: got %s, want %s", unmarshaled.Gateway, networkSettings.Gateway)
	}
	if len(unmarshaled.Ports) != len(networkSettings.Ports) {
		t.Errorf(
			"Ports length mismatch: got %d, want %d",
			len(unmarshaled.Ports),
			len(networkSettings.Ports),
		)
	}
	if len(unmarshaled.Networks) != len(networkSettings.Networks) {
		t.Errorf(
			"Networks length mismatch: got %d, want %d",
			len(unmarshaled.Networks),
			len(networkSettings.Networks),
		)
	}
}

func TestPort_JSON(t *testing.T) {
	port := Port{
		HostIP:   "0.0.0.0",
		HostPort: "8080",
	}

	// Test JSON marshaling
	data, err := json.Marshal(port)
	if err != nil {
		t.Fatalf("Failed to marshal Port: %v", err)
	}

	// Test JSON unmarshaling
	var unmarshaled Port
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal Port: %v", err)
	}

	// Verify the data
	if unmarshaled.HostIP != port.HostIP {
		t.Errorf("HostIP mismatch: got %s, want %s", unmarshaled.HostIP, port.HostIP)
	}
	if unmarshaled.HostPort != port.HostPort {
		t.Errorf("HostPort mismatch: got %s, want %s", unmarshaled.HostPort, port.HostPort)
	}
}

func TestDockerInfo_JSON(t *testing.T) {
	dockerInfo := DockerInfo{
		Containers:         5,
		Images:             10,
		Volumes:            3,
		Networks:           2,
		Version:            "20.10.0",
		APIVersion:         "1.41",
		OS:                 "linux",
		Architecture:       "x86_64",
		KernelVersion:      "5.4.0",
		GoVersion:          "1.15.0",
		Driver:             "overlay2",
		DriverStatus:       [][]string{{"Backing Filesystem", "extfs"}},
		MemoryLimit:        true,
		SwapLimit:          true,
		KernelMemory:       false,
		CPUCfsQuota:        true,
		CPUCfsPeriod:       true,
		CPUShares:          true,
		CPUSet:             true,
		IPv4Forwarding:     true,
		BridgeNfIptables:   true,
		BridgeNfIP6tables:  true,
		Debug:              false,
		NFd:                20,
		NGoroutines:        50,
		SystemTime:         "2023-01-01T00:00:00Z",
		LoggingDriver:      "json-file",
		CgroupDriver:       "cgroupfs",
		OperatingSystem:    "Ubuntu 20.04.3 LTS",
		OSType:             "linux",
		IndexServerAddress: "https://index.docker.io/v1/",
		RegistryConfigs: map[string]any{
			"docker.io": map[string]any{
				"Mirrors": []string{"mirror1.example.com"},
			},
		},
		InsecureRegistries: []string{"localhost:5000"},
		RegistryMirrors:    []string{"https://mirror.example.com"},
		Experimental:       false,
		ServerVersion:      "20.10.0",
		ClusterStore:       "",
		ClusterAdvertise:   "",
		DefaultRuntime:     "runc",
		LiveRestoreEnabled: true,
		Isolation:          "default",
		InitBinary:         "docker-init",
		ContainerdCommit: Commit{
			ID:       "containerd-commit-id",
			Expected: "expected-containerd-commit-id",
		},
		RuncCommit: Commit{
			ID:       "runc-commit-id",
			Expected: "expected-runc-commit-id",
		},
		InitCommit: Commit{
			ID:       "init-commit-id",
			Expected: "expected-init-commit-id",
		},
		SecurityOptions: []string{"seccomp", "selinux"},
		ProductLicense:  "Community Engine",
		Warnings:        []string{"Warning 1", "Warning 2"},
	}

	// Test JSON marshaling
	data, err := json.Marshal(dockerInfo)
	if err != nil {
		t.Fatalf("Failed to marshal DockerInfo: %v", err)
	}

	// Test JSON unmarshaling
	var unmarshaled DockerInfo
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal DockerInfo: %v", err)
	}

	// Verify key fields
	if unmarshaled.Containers != dockerInfo.Containers {
		t.Errorf(
			"Containers mismatch: got %d, want %d",
			unmarshaled.Containers,
			dockerInfo.Containers,
		)
	}
	if unmarshaled.Images != dockerInfo.Images {
		t.Errorf("Images mismatch: got %d, want %d", unmarshaled.Images, dockerInfo.Images)
	}
	if unmarshaled.Version != dockerInfo.Version {
		t.Errorf("Version mismatch: got %s, want %s", unmarshaled.Version, dockerInfo.Version)
	}
	if unmarshaled.OS != dockerInfo.OS {
		t.Errorf("OS mismatch: got %s, want %s", unmarshaled.OS, dockerInfo.OS)
	}
	if unmarshaled.Architecture != dockerInfo.Architecture {
		t.Errorf(
			"Architecture mismatch: got %s, want %s",
			unmarshaled.Architecture,
			dockerInfo.Architecture,
		)
	}
	if unmarshaled.Driver != dockerInfo.Driver {
		t.Errorf("Driver mismatch: got %s, want %s", unmarshaled.Driver, dockerInfo.Driver)
	}
	if unmarshaled.MemoryLimit != dockerInfo.MemoryLimit {
		t.Errorf(
			"MemoryLimit mismatch: got %t, want %t",
			unmarshaled.MemoryLimit,
			dockerInfo.MemoryLimit,
		)
	}
	if unmarshaled.Debug != dockerInfo.Debug {
		t.Errorf("Debug mismatch: got %t, want %t", unmarshaled.Debug, dockerInfo.Debug)
	}
	if len(unmarshaled.SecurityOptions) != len(dockerInfo.SecurityOptions) {
		t.Errorf(
			"SecurityOptions length mismatch: got %d, want %d",
			len(unmarshaled.SecurityOptions),
			len(dockerInfo.SecurityOptions),
		)
	}
	if len(unmarshaled.Warnings) != len(dockerInfo.Warnings) {
		t.Errorf(
			"Warnings length mismatch: got %d, want %d",
			len(unmarshaled.Warnings),
			len(dockerInfo.Warnings),
		)
	}
}

func TestPlugins_JSON(t *testing.T) {
	plugins := Plugins{
		Volume:        []string{"local", "nfs"},
		Network:       []string{"bridge", "host", "overlay"},
		Authorization: []string{"authz"},
		Log:           []string{"json-file", "syslog"},
	}

	// Test JSON marshaling
	data, err := json.Marshal(plugins)
	if err != nil {
		t.Fatalf("Failed to marshal Plugins: %v", err)
	}

	// Test JSON unmarshaling
	var unmarshaled Plugins
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal Plugins: %v", err)
	}

	// Verify the data
	if len(unmarshaled.Volume) != len(plugins.Volume) {
		t.Errorf(
			"Volume plugins length mismatch: got %d, want %d",
			len(unmarshaled.Volume),
			len(plugins.Volume),
		)
	}
	if len(unmarshaled.Network) != len(plugins.Network) {
		t.Errorf(
			"Network plugins length mismatch: got %d, want %d",
			len(unmarshaled.Network),
			len(plugins.Network),
		)
	}
	if len(unmarshaled.Authorization) != len(plugins.Authorization) {
		t.Errorf(
			"Authorization plugins length mismatch: got %d, want %d",
			len(unmarshaled.Authorization),
			len(plugins.Authorization),
		)
	}
	if len(unmarshaled.Log) != len(plugins.Log) {
		t.Errorf(
			"Log plugins length mismatch: got %d, want %d",
			len(unmarshaled.Log),
			len(plugins.Log),
		)
	}
}

func TestCommit_JSON(t *testing.T) {
	commit := Commit{
		ID:       "abc123def456",
		Expected: "expected-commit-id",
	}

	// Test JSON marshaling
	data, err := json.Marshal(commit)
	if err != nil {
		t.Fatalf("Failed to marshal Commit: %v", err)
	}

	// Test JSON unmarshaling
	var unmarshaled Commit
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal Commit: %v", err)
	}

	// Verify the data
	if unmarshaled.ID != commit.ID {
		t.Errorf("ID mismatch: got %s, want %s", unmarshaled.ID, commit.ID)
	}
	if unmarshaled.Expected != commit.Expected {
		t.Errorf("Expected mismatch: got %s, want %s", unmarshaled.Expected, commit.Expected)
	}
}

func TestEmptyStructs(t *testing.T) {
	// Test empty structs can be marshaled/unmarshaled
	emptyContainerDetails := ContainerDetails{}
	emptyMount := Mount{}
	emptyNetworkSettings := NetworkSettings{}
	emptyPort := Port{}
	emptyDockerInfo := DockerInfo{}
	emptyPlugins := Plugins{}
	emptyCommit := Commit{}

	testCases := []struct {
		name string
		data any
	}{
		{"ContainerDetails", &emptyContainerDetails},
		{"Mount", &emptyMount},
		{"NetworkSettings", &emptyNetworkSettings},
		{"Port", &emptyPort},
		{"DockerInfo", &emptyDockerInfo},
		{"Plugins", &emptyPlugins},
		{"Commit", &emptyCommit},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test marshaling
			data, err := json.Marshal(tc.data)
			if err != nil {
				t.Errorf("Failed to marshal empty %s: %v", tc.name, err)
			}

			// Test unmarshaling
			err = json.Unmarshal(data, tc.data)
			if err != nil {
				t.Errorf("Failed to unmarshal empty %s: %v", tc.name, err)
			}
		})
	}
}
