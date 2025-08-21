package shared

// DockerInfo represents Docker system information
type DockerInfo struct {
	Containers         int            `json:"containers"`
	Images             int            `json:"images"`
	Volumes            int            `json:"volumes"`
	Networks           int            `json:"networks"`
	Version            string         `json:"version"`
	APIVersion         string         `json:"api_version"`
	OS                 string         `json:"os"`
	Architecture       string         `json:"architecture"`
	KernelVersion      string         `json:"kernel_version"`
	GoVersion          string         `json:"go_version"`
	Driver             string         `json:"driver"`
	DriverStatus       [][]string     `json:"driver_status"`
	Plugins            Plugins        `json:"plugins"`
	MemoryLimit        bool           `json:"memory_limit"`
	SwapLimit          bool           `json:"swap_limit"`
	KernelMemory       bool           `json:"kernel_memory"`
	CPUCfsQuota        bool           `json:"cpu_cfs_quota"`
	CPUCfsPeriod       bool           `json:"cpu_cfs_period"`
	CPUShares          bool           `json:"cpu_shares"`
	CPUSet             bool           `json:"cpu_set"`
	IPv4Forwarding     bool           `json:"ipv4_forwarding"`
	BridgeNfIptables   bool           `json:"bridge_nf_iptables"`
	BridgeNfIP6tables  bool           `json:"bridge_nf_ip6tables"`
	Debug              bool           `json:"debug"`
	NFd                int            `json:"nfd"`
	NGoroutines        int            `json:"ngoroutines"`
	SystemTime         string         `json:"system_time"`
	LoggingDriver      string         `json:"logging_driver"`
	CgroupDriver       string         `json:"cgroup_driver"`
	OperatingSystem    string         `json:"operating_system"`
	OSType             string         `json:"os_type"`
	IndexServerAddress string         `json:"index_server_address"`
	RegistryConfigs    map[string]any `json:"registry_configs"`
	InsecureRegistries []string       `json:"insecure_registries"`
	RegistryMirrors    []string       `json:"registry_mirrors"`
	Experimental       bool           `json:"experimental"`
	ServerVersion      string         `json:"server_version"`
	ClusterStore       string         `json:"cluster_store"`
	ClusterAdvertise   string         `json:"cluster_advertise"`
	DefaultRuntime     string         `json:"default_runtime"`
	LiveRestoreEnabled bool           `json:"live_restore_enabled"`
	Isolation          string         `json:"isolation"`
	InitBinary         string         `json:"init_binary"`
	ContainerdCommit   Commit         `json:"containerd_commit"`
	RuncCommit         Commit         `json:"runc_commit"`
	InitCommit         Commit         `json:"init_commit"`
	SecurityOptions    []string       `json:"security_options"`
	ProductLicense     string         `json:"product_license"`
	Warnings           []string       `json:"warnings"`
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
