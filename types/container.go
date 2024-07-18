package types

type Container struct {
	Id              string          `json:"Id"`
	Names           []string        `json:"Names"`
	Image           string          `json:"Image"`
	ImageID         string          `json:"ImageID"`
	State           string          `json:"State"`
	Status          string          `json:"Status"`
	Ports           []Port          `json:"Ports"`
	NetworkSettings NetworkSettings `json:"NetworkSettings"`
}

type Port struct {
	PrivatePort int    `json:"PrivatePort"`
	PublicPort  int    `json:"PublicPort"`
	Type        string `json:"Type"`
}

type Label struct {
	ComExampleVendor  string `json:"com.example.vendor"`
	ComExampleLicense string `json:"com.example.license"`
	ComExampleVersion string `json:"com.example.version"`
}

type NetworkSettings struct {
	Networks map[string]Network `json:"Networks"`
}

type Network struct {
	Links       *[]string `json:"Links"`
	Aliases     *[]string `json:"Aliases"`
	MacAddress  string    `json:"MacAddress"`
	NetworkID   string    `json:"NetworkID"`
	EndpointID  string    `json:"EndpointID"`
	Gateway     string    `json:"Gateway"`
	IPAddress   string    `json:"IPAddress"`
	IPPrefixLen int       `json:"IPPrefixlen"`
	DNSNames    *[]string `json:"DNSNames"`
}

type StopParams struct {
	Signal string `json:"Signal"`
	T      int    `json:"T"`
}

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

// Processes  TOP struct for processes
type Processes struct {
	Processes [][]string `json:"Processes"`
}
