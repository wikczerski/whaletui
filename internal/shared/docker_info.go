package shared

// DockerInfo represents Docker system information
type DockerInfo struct {
	RegistryConfigs    map[string]any `json:"registry_configs"`
	RuncCommit         Commit         `json:"runc_commit"`
	ContainerdCommit   Commit         `json:"containerd_commit"`
	InitCommit         Commit         `json:"init_commit"`
	ClusterStore       string         `json:"cluster_store"`
	Isolation          string         `json:"isolation"`
	OS                 string         `json:"os"`
	Architecture       string         `json:"architecture"`
	KernelVersion      string         `json:"kernel_version"`
	GoVersion          string         `json:"go_version"`
	Driver             string         `json:"driver"`
	Version            string         `json:"version"`
	LoggingDriver      string         `json:"logging_driver"`
	ProductLicense     string         `json:"product_license"`
	InitBinary         string         `json:"init_binary"`
	APIVersion         string         `json:"api_version"`
	DefaultRuntime     string         `json:"default_runtime"`
	ClusterAdvertise   string         `json:"cluster_advertise"`
	SystemTime         string         `json:"system_time"`
	ServerVersion      string         `json:"server_version"`
	ConnectionMethod   string         `json:"connection_method"`
	IndexServerAddress string         `json:"index_server_address"`
	OSType             string         `json:"os_type"`
	OperatingSystem    string         `json:"operating_system"`
	CgroupDriver       string         `json:"cgroup_driver"`
	Plugins            Plugins        `json:"plugins"`
	DriverStatus       [][]string     `json:"driver_status"`
	Warnings           []string       `json:"warnings"`
	SecurityOptions    []string       `json:"security_options"`
	InsecureRegistries []string       `json:"insecure_registries"`
	RegistryMirrors    []string       `json:"registry_mirrors"`
	NGoroutines        int            `json:"ngoroutines"`
	NFd                int            `json:"nfd"`
	Images             int            `json:"images"`
	Volumes            int            `json:"volumes"`
	Networks           int            `json:"networks"`
	Containers         int            `json:"containers"`
	CPUShares          bool           `json:"cpu_shares"`
	CPUCfsPeriod       bool           `json:"cpu_cfs_period"`
	CPUCfsQuota        bool           `json:"cpu_cfs_quota"`
	LiveRestoreEnabled bool           `json:"live_restore_enabled"`
	KernelMemory       bool           `json:"kernel_memory"`
	SwapLimit          bool           `json:"swap_limit"`
	MemoryLimit        bool           `json:"memory_limit"`
	CPUSet             bool           `json:"cpu_set"`
	Experimental       bool           `json:"experimental"`
	IPv4Forwarding     bool           `json:"ipv4_forwarding"`
	BridgeNfIptables   bool           `json:"bridge_nf_iptables"`
	BridgeNfIP6tables  bool           `json:"bridge_nf_ip6tables"`
	Debug              bool           `json:"debug"`
}

// Plugins represents Docker plugins information
type Plugins struct {
	Volume        []string `json:"volume"`
	Network       []string `json:"network"`
	Authorization []string `json:"authorization"`
	Log           []string `json:"log"`
}

// Commit represents a commit information
type Commit struct {
	ID       string `json:"id"`
	Expected string `json:"expected"`
}

// GetVersion returns the Docker version
func (d *DockerInfo) GetVersion() string {
	return d.Version
}

// GetOperatingSystem returns the operating system
func (d *DockerInfo) GetOperatingSystem() string {
	return d.OperatingSystem
}

// GetLoggingDriver returns the logging driver
func (d *DockerInfo) GetLoggingDriver() string {
	return d.LoggingDriver
}

// GetConnectionMethod returns the connection method used
func (d *DockerInfo) GetConnectionMethod() string {
	if d.ConnectionMethod == "" {
		return "Local Docker"
	}
	return d.ConnectionMethod
}
