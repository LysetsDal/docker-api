package api

import (
	"context"
	"github.com/LysetsDal/docker-api/service/container"
	. "github.com/LysetsDal/docker-api/utils"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	. "github.com/klauspost/cpuid/v2"
)

const dockerSocket string = "/var/run/docker.sock"

// APIServer API struct
type APIServer struct {
	Name           string
	ServerCPU      string
	ServerCPUCores int
	ListenAddr     string
	StartTime      time.Time
	DockerSock     http.Client
}

type VersionData struct {
	Name           string `json:"Name"`
	ServerCPU      string `json:"ServerCPU"`
	ServerCPUCores int    `json:"ServerCPUCores"`
	ListenAddr     string `json:"ListenAddr"`
	CurrentTime    string `json:"Time"`
	DockerSock     string `json:"DockerSocket"`
}

// Create a UNIX Listener for the Docker socket
func connectDockerSocket() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", dockerSocket)
			},
		},
	}
}

// NewAPIServer Create new API-Server
func NewAPIServer(listenAddr string) *APIServer {
	return &APIServer{
		Name:           "Docker-API Server",
		ServerCPU:      CPU.BrandName,
		ServerCPUCores: CPU.PhysicalCores,
		ListenAddr:     listenAddr,
		StartTime:      time.Now(),
		DockerSock:     *connectDockerSocket(),
	}
}

// Run Start Listening
func (s *APIServer) Run() {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	containerHandler := container.NewHandler(s.DockerSock)
	containerHandler.RegisterRoutes(subrouter)

	subrouter.HandleFunc("/", MakeHttpHandleFunc(s.HomeHandler))

	log.Printf("%s listening on %s\n", s.Name, s.ListenAddr)

	subrouter.Use(logMW)

	// Run server in go func so it doesn't block
	go func() {}()
	if err := http.ListenAndServe(s.ListenAddr, router); err != nil {
		return
	}

	// Graceful shutdown on ctrl + c:
	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt)

	<-c

}

func logMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s - %s (%s)", r.Method, r.URL.Path, r.RemoteAddr)

		// compare the return-value to the authMW
		next.ServeHTTP(w, r)
	})
}

// =======================================================
// ==================== Controllers ======================
// =======================================================

// HomeHandler Displays version info
func (s *APIServer) HomeHandler(w http.ResponseWriter, _ *http.Request) error {

	data := VersionData{
		Name:           s.Name,
		ServerCPU:      s.ServerCPU,
		ServerCPUCores: s.ServerCPUCores,
		ListenAddr:     s.ListenAddr,
		CurrentTime:    time.Now().String(),
		DockerSock:     dockerSocket,
	}

	return WriteJson(w, http.StatusOK, data)
}

func (s *APIServer) handleStopAllContainerRequest(w http.ResponseWriter, r *http.Request) error {

	return nil
}
