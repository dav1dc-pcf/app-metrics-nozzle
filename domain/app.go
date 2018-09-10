package domain

type InstanceCount struct {
	Configured int `json:"configured"`
	Running    int `json:"running"`
}

type Instances struct {
	CellIP        string  `json:"cell_ip"`
	CPUUsage      float64 `json:"cpu_usage"`  //CPUPercentage
	DiskUsage     uint64  `json:"disk_usage"` //DiskBytes
	GcStats       string  `json:"gc_stats"`
	InstanceIndex int64   `json:"index"`
	MemoryUsage   uint64  `json:"memory_usage"` //MemBytes
	Uptime        int32   `json:"uptime"`
	Since         int32   `json:"since"`
	State         string  `json:"state"`
}

type EnvironmentSummary struct {
	TotalCPU               float64 `json:"total_cpu"`
	TotalDiskConfigured    int32   `json:"total_disk_configured"`
	TotalDiskProvisioned   int32   `json:"total_disk_provisioned"`
	TotalDiskUsage         uint64  `json:"total_disk_usage"`
	TotalMemoryConfigured  int32   `json:"total_memory_configured"`
	TotalMemoryProvisioned int32   `json:"total_memory_provisioned"`
	TotalMemoryUsage       uint64  `json:"total_memory_usage"`
}
type App struct {
	Buildpack string `json:"buildpack"`
	Diego     bool   `json:"diego"`

	Environment        map[string]interface{} `json:"environment"`
	EnvironmentSummary EnvironmentSummary     `json:"environment_summary"`
	GUID               string                 `json:"guid"`
	InstanceCount      `json:"instance_count"`
	Instances          []Instances `json:"instances"`
	Name               string      `json:"name"`
	Organization       struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"organization"`
	EventCount            int64    `json:"event_count"`
	HTTPGoodCount         int64    `json:"http_good_count"`
	HTTPErrorCount        int64    `json:"http_error_count"`
	LastEventTime         int64    `json:"last_event_time"`
	RequestsPerSecond     float64  `json:"requests_per_second"`
	ElapsedSinceLastEvent int64    `json:"elapsed_since_last_event"`
	Routes                []string `json:"routes"`
	Space                 struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"space"`
	State string `json:"state"`
}
