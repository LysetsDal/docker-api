package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	// "io"
	"log"
	"net"
	"net/http"

	"github.com/gorilla/mux"
)

const docker_socket string = "/var/run/docker.sock"

// Decorator pattern to 'wrap' http.HandlerFunc
type apiFunc func(http.ResponseWriter, *http.Request) error

type apiError struct {
	Error string
}

// API
type APIServer struct {
	name string
	listenAddr string
	docker_sock http.Client
}

// Read the Docker daemons responses (format json)
func ReadJson(body io.Reader, v any) error {
	decoder := json.NewDecoder(body)
	return decoder.Decode(v)
}

// Write Json with standard header.
func WriteJson(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

// Make custom http.HandlerFunc
func makeHttpHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJson(w, http.StatusBadRequest, apiError{Error: err.Error()})
		} 
	}
}

// Create a UNIX Listener for the Docker socket
func connectDockerSocket() *http.Client {
	return &http.Client {
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", docker_socket)
			},
		},
	}
}

// Create new API-Server 
func NewAPIServer(listenAddr string) *APIServer {
	return &APIServer{
		name: "Docker-API Server",
		listenAddr: listenAddr,
		docker_sock: *connectDockerSocket(),
	}
}

// Start Listening 
func (s *APIServer) Run() {
	router := mux.NewRouter()
	router.HandleFunc("/containers/list", makeHttpHandleFunc(s.handleRequest))


	log.Printf("%s listening on %s\n", s.name, s.listenAddr)

	http.ListenAndServe(s.listenAddr, router)
}

// =======================================================
// ==================== Controllers ======================
// =======================================================
// Control logic 
func (s *APIServer) handleRequest(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" && r.RequestURI == "/containers/list"{
		return s.handleGetContainerList(w, r)
	}
	if r.Method == "POST" {
		return s.handlePostRequest(w, r)
	}

	return fmt.Errorf("method not allowed %s", r.Method)
}

// GET List of all containers 
func (s *APIServer) handleGetContainerList(w http.ResponseWriter, _ *http.Request) error {
	response, err := s.docker_sock.Get("http://unix/containers/json")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	defer response.Body.Close()

	// Read and decode the JSON response
	var containers []Container
	if err := ReadJson(response.Body, &containers); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	//@TODO: Handle logging in separate module
	fmt.Println("GET container/list")
	fmt.Println("Status: " + response.Status) 

	return WriteJson(w, response.StatusCode, containers)
}

// POST Create container 
func (s *APIServer) handlePostRequest(w http.ResponseWriter, r *http.Request) error {
	return WriteJson(w, http.StatusOK, r.GetBody)
}
