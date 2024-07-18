package types

// Payload Define the Payload struct with all nested structs
type Payload struct {
	Hostname         string              `json:"Hostname"`
	Domainname       string              `json:"Domainname"`
	User             string              `json:"User"`
	AttachStdin      bool                `json:"AttachStdin"`
	AttachStdout     bool                `json:"AttachStdout"`
	AttachStderr     bool                `json:"AttachStderr"`
	Tty              bool                `json:"Tty"`
	OpenStdin        bool                `json:"OpenStdin"`
	StdinOnce        bool                `json:"StdinOnce"`
	Env              []string            `json:"Env"`
	Cmd              []string            `json:"Cmd"`
	Entrypoint       string              `json:"Entrypoint"`
	Image            string              `json:"Image"`
	Labels           map[string]string   `json:"Labels"`
	Volumes          map[string]struct{} `json:"Volumes"`
	WorkingDir       string              `json:"WorkingDir"`
	NetworkDisabled  bool                `json:"NetworkDisabled"`
	MacAddress       string              `json:"MacAddress"`
	ExposedPorts     map[string]struct{} `json:"ExposedPorts"`
	StopSignal       string              `json:"StopSignal"`
	StopTimeout      int                 `json:"StopTimeout"`
	HostConfig       HostConfig          `json:"HostConfig"`
	NetworkingConfig NetworkingConfig    `json:"NetworkingConfig"`
}

type HostConfig struct {
	Binds                []string                 `json:"Binds"`
	Links                []string                 `json:"Links"`
	Memory               int                      `json:"Memory"`
	MemorySwap           int                      `json:"MemorySwap"`
	MemoryReservation    int                      `json:"MemoryReservation"`
	NanoCpus             int64                    `json:"NanoCpus"`
	CpuPercent           int                      `json:"CpuPercent"`
	CpuShares            int                      `json:"CpuShares"`
	CpuPeriod            int                      `json:"CpuPeriod"`
	CpuRealtimePeriod    int                      `json:"CpuRealtimePeriod"`
	CpuRealtimeRuntime   int                      `json:"CpuRealtimeRuntime"`
	CpuQuota             int                      `json:"CpuQuota"`
	CpusetCpus           string                   `json:"CpusetCpus"`
	CpusetMems           string                   `json:"CpusetMems"`
	MaximumIOps          int                      `json:"MaximumIOps"`
	MaximumIOBps         int                      `json:"MaximumIOBps"`
	BlkioWeight          int                      `json:"BlkioWeight"`
	BlkioWeightDevice    []BlkioDevice            `json:"BlkioWeightDevice"`
	BlkioDeviceReadBps   []BlkioDevice            `json:"BlkioDeviceReadBps"`
	BlkioDeviceReadIOps  []BlkioDevice            `json:"BlkioDeviceReadIOps"`
	BlkioDeviceWriteBps  []BlkioDevice            `json:"BlkioDeviceWriteBps"`
	BlkioDeviceWriteIOps []BlkioDevice            `json:"BlkioDeviceWriteIOps"`
	DeviceRequests       []DeviceRequest          `json:"DeviceRequests"`
	MemorySwappiness     int                      `json:"MemorySwappiness"`
	OomKillDisable       bool                     `json:"OomKillDisable"`
	OomScoreAdj          int                      `json:"OomScoreAdj"`
	PidMode              string                   `json:"PidMode"`
	PidsLimit            int                      `json:"PidsLimit"`
	PortBindings         map[string][]PortBinding `json:"PortBindings"`
	PublishAllPorts      bool                     `json:"PublishAllPorts"`
	Privileged           bool                     `json:"Privileged"`
	ReadonlyRootfs       bool                     `json:"ReadonlyRootfs"`
	Dns                  []string                 `json:"Dns"`
	DnsOptions           []string                 `json:"DnsOptions"`
	DnsSearch            []string                 `json:"DnsSearch"`
	VolumesFrom          []string                 `json:"VolumesFrom"`
	CapAdd               []string                 `json:"CapAdd"`
	CapDrop              []string                 `json:"CapDrop"`
	GroupAdd             []string                 `json:"GroupAdd"`
	RestartPolicy        RestartPolicy            `json:"RestartPolicy"`
	AutoRemove           bool                     `json:"AutoRemove"`
	NetworkMode          string                   `json:"NetworkMode"`
	Devices              []Device                 `json:"Devices"`
	Ulimits              []Ulimit                 `json:"Ulimits"`
	LogConfig            LogConfig                `json:"LogConfig"`
	SecurityOpt          []string                 `json:"SecurityOpt"`
	StorageOpt           map[string]string        `json:"StorageOpt"`
	CgroupParent         string                   `json:"CgroupParent"`
	VolumeDriver         string                   `json:"VolumeDriver"`
	ShmSize              int                      `json:"ShmSize"`
}

type BlkioDevice struct {
	// Add appropriate fields here
}

type DeviceRequest struct {
	Driver       string            `json:"Driver"`
	Count        int               `json:"Count"`
	DeviceIDs    []string          `json:"DeviceIDs"`
	Capabilities [][]string        `json:"Capabilities"`
	Options      map[string]string `json:"Options"`
}

type PortBinding struct {
	HostPort string `json:"HostPort"`
}

type RestartPolicy struct {
	Name              string `json:"Name"`
	MaximumRetryCount int    `json:"MaximumRetryCount"`
}

type Device struct {
	// Add appropriate fields here
}

type Ulimit struct {
	// Add appropriate fields here
}

type LogConfig struct {
	Type   string            `json:"Type"`
	Config map[string]string `json:"Config"`
}

type NetworkingConfig struct {
	EndpointsConfig map[string]EndpointConfig `json:"EndpointsConfig"`
}

type EndpointConfig struct {
	IPAMConfig IPAMConfig `json:"IPAMConfig"`
	Links      []string   `json:"Links"`
	Aliases    []string   `json:"Aliases"`
}

type IPAMConfig struct {
	IPv4Address  string   `json:"IPv4Address"`
	IPv6Address  string   `json:"IPv6Address"`
	LinkLocalIPs []string `json:"LinkLocalIPs"`
}
