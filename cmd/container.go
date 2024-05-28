package main

type Container struct {
	Id      string   `json:"Id"`
	Names   []string `json:"Names"`
	Image   string   `json:"Image"`
	ImageID string   `json:"ImageID"`
	State   string   `json:"State"`
	Status  string   `json:"Status"`
	Ports   []Port   `json:"Ports"`
	// Label []Label `json:"Label"`
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

// CreateDefaultContainer Test to create a default container
func CreateDefaultContainer() *Container {
	return &Container{
		Id:    "dummy_id",
		Names: []string{"nostalgic_beaver"},
		Image: "test-image:latest",
		Ports: []Port{},
		NetworkSettings: NetworkSettings{
			Networks: map[string]Network{},
		},
	}
}
