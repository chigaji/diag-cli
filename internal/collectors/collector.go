package collectors

import "context"

type SystemData struct {
	Host         HostInfo      `json:"host" yaml:"host"`
	CPU          []CPUInfo     `json:"cpu" yaml:"cpu"`
	Memory       MemoryInfo    `json:"memory" yaml:"memory"`
	Disks        []DiskInfo    `json:"disks" yaml:"disks"`
	Temperatures []Temperature `json:"temperatures,omitempty" yaml:"temperatures,omitempty"`
}

type HostInfo struct {
	Hostname string
	OS       string
	Platform string
	Kernel   string
	Uptime   uint64
	BootTime uint64
}

type CPUInfo struct {
	Model  string
	Cores  int
	Mhz    float64
	User   float64
	System float64
	Idle   float64
}

type MemoryInfo struct {
	Total       uint64
	Used        uint64
	Free        uint64
	UsedPercent uint64
}

type DiskInfo struct {
	Device      string
	Mountpoint  string
	Fstype      string
	Total       uint64
	Used        uint64
	Free        uint64
	UsedPercent float64
}

type Temperature struct {
	Sensor  string
	Celsius float64
}

type NetworkData struct {
	Interfaces []NetIF `json:"interfaces" yaml:"interfaces"`
	IO         []NetIO `json:"io" yaml:"io"`
}

type NetIF struct {
	Name         string
	HardwareAddr string
	MTU          int
	Addrs        []string
}

type NetIO struct {
	Name        string
	BytesSent   uint64
	BytesRecv   uint64
	PacketsSent uint64
	PacketsRecv uint64
}

type ProcessRow struct {
	PID  int32
	Name string
	CPU  float64
	MEM  float32
}

type ProcData struct {
	Rows []ProcessRow
}

// collector abstract platform metrics provider

type Collector interface {
	System(ctx context.Context, all bool) (SystemData, error)
	Network(ctx context.Context, iface string) (NetworkData, error)
	Processes(ctx context.Context, top int, sortBy string) (ProcData, error)
}
