package types

type InspectObject struct {
	AppArmorProfile string   `json:"AppArmorProfile"`
	Args            []string `json:"Args"`
	Config          Config   `json:"Config"`
	Created         string   `json:"Created"`
	Driver          string   `json:"Driver"`
	Mounts          []Mount  `json:"Mounts"`
}

type Config struct {
	AttachStderr    bool                   `json:"AttachStderr"`
	AttachStdin     bool                   `json:"AttachStdin"`
	AttachStdout    bool                   `json:"AttachStdout"`
	Cmd             []string               `json:"Cmd"`
	DomainName      string                 `json:"Domainname"`
	Env             []string               `json:"Env"`
	Healthcheck     map[string]interface{} `json:"Healthcheck"`
	Hostname        string                 `json:"Hostname"`
	Image           string                 `json:"Image"`
	Labels          map[string]string      `json:"Labels"`
	MacAddress      string                 `json:"MacAddress"`
	NetworkDisabled bool                   `json:"NetworkDisabled"`
	OpenStdin       bool                   `json:"OpenStdin"`
	StdinOnce       bool                   `json:"StdinOnce"`
	Tty             bool                   `json:"Tty"`
	User            string                 `json:"User"`
	Volumes         map[string]interface{} `json:"Volumes"`
	WorkingDir      string                 `json:"WorkingDir"`
	StopSignal      string                 `json:"StopSignal"`
	StopTimeout     int                    `json:"StopTimeout"`
}

type Mount struct {
	Name        string `json:"Name"`
	Source      string `json:"Source"`
	Destination string `json:"Destination"`
	Driver      string `json:"Driver"`
	Mode        string `json:"Mode"`
	RW          bool   `json:"RW"`
	Propagation string `json:"Propagation"`
}
