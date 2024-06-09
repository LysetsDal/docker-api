package types

import "net/http"

// ApiFunc Decorator pattern to 'wrap' http.HandlerFunc
type ApiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string `json:"error"`
}

type ApiMessage struct {
	Message string `json:"message"`
}

type PruneResponse struct {
	ContainersDeleted []string `json:"ContainersDeleted"`
	SpaceReclaimed    int      `json:"SpaceReclaimed"`
}

type CreateContainerResponse struct {
	Id       string   `json:"Id"`
	Warnings []string `json:"Warnings"`
}

type CreateContainerRequest struct {
	Image        string     `json:"Image"`
	AttachStdin  bool       `json:"AttachStdin"`
	AttachStdout bool       `json:"AttachStdout"`
	AttachStderr bool       `json:"AttachStderr"`
	Tty          bool       `json:"Tty"`
	OpenStdin    bool       `json:"OpenStdin"`
	StdinOnce    bool       `json:"StdinOnce"`
	Cmd          []string   `json:"Cmd"`
	HostConfig   HostConfig `json:"HostConfig"`
}

type Filter struct {
	Status []string `json:"Status"`
}
