package container

import (
	"fmt"
	. "github.com/LysetsDal/docker-api/types"
	. "github.com/LysetsDal/docker-api/utils"
	"github.com/gorilla/mux"
	"net/http"
)

type Handler struct {
	DockerSock http.Client
}

func NewHandler(sock http.Client) *Handler {
	return &Handler{
		DockerSock: sock,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/containers/list", MakeHttpHandleFunc(h.handleGetContainerList))
	router.HandleFunc("/containers/{id}/json", MakeHttpHandleFunc(h.handleGetContainerList))
	router.HandleFunc("/containers/{id}/start", MakeHttpHandleFunc(h.handleStartContainerRequest))
	router.HandleFunc("/containers/{id}/stop", MakeHttpHandleFunc(h.handleStopContainerRequest))
	router.HandleFunc("/containers/prune", MakeHttpHandleFunc(h.handlePruneContainersRequest))
}

// GET List of all containers
func (h *Handler) handleGetContainerById(w http.ResponseWriter, r *http.Request) error {
	pathVars := mux.Vars(r)
	fmt.Println(pathVars["id"])

	uri := fmt.Sprintf("http://unix/containers/%s/json", pathVars["id"])

	response, err := h.DockerSock.Get(uri)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	defer response.Body.Close()

	// Read and decode the JSON response
	var inspectObject InspectObject
	if err := ReadJson(response.Body, &inspectObject); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	return WriteJson(w, response.StatusCode, inspectObject)
}

// GET List of all containers
func (h *Handler) handleGetContainerList(w http.ResponseWriter, _ *http.Request) error {
	response, err := h.DockerSock.Get("http://unix/containers/json")
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

	return WriteJson(w, response.StatusCode, containers)
}

// POST Start container
func (h *Handler) handleStartContainerRequest(w http.ResponseWriter, r *http.Request) error {
	pathVars := mux.Vars(r)
	uri := fmt.Sprintf("http://unix/containers/%s/start", pathVars["id"])

	request, err := h.DockerSock.Post(uri, "application/json", r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return err
	}
	defer request.Body.Close()

	if request.StatusCode == http.StatusNoContent {
		return WriteJson(w, http.StatusOK, ApiMessage{Message: "Container started"})

	} else if request.StatusCode == http.StatusNotModified {
		return WriteJson(w, http.StatusOK, ApiMessage{Message: "Container already started"})

	} else if request.StatusCode == http.StatusNotFound {
		return WriteJson(w, http.StatusNotFound, ApiError{Error: "No such container"})

	} else {
		return WriteJson(w, http.StatusInternalServerError, ApiError{Error: "Something went wrong"})
	}
}

// POST Stop container
func (h *Handler) handleStopContainerRequest(w http.ResponseWriter, r *http.Request) error {
	pathVars := mux.Vars(r)
	uri := fmt.Sprintf("http://unix/containers/%s/stop", pathVars["id"])

	request, err := h.DockerSock.Post(uri, "application/json", r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return err
	}
	defer request.Body.Close()

	fmt.Printf("r.Body: %s\n", r.Body)

	if request.StatusCode == http.StatusNoContent {
		return WriteJson(w, http.StatusOK, ApiMessage{Message: "Container stopped"})

	} else if request.StatusCode == http.StatusNotModified {
		return WriteJson(w, http.StatusOK, ApiMessage{Message: "Container already stopped"})

	} else if request.StatusCode == http.StatusNotFound {
		return WriteJson(w, http.StatusNotFound, ApiError{Error: "No such container"})

	} else {
		return WriteJson(w, http.StatusInternalServerError, ApiError{Error: "Something went wrong"})
	}
}

func (h *Handler) handlePruneContainersRequest(w http.ResponseWriter, r *http.Request) error {
	type pruneResponse struct {
		ContainersDeleted []string `json:"ContainersDeleted"`
		SpaceReclaimed    int      `json:"SpaceReclaimed"`
	}

	uri := fmt.Sprintf("http://unix/containers/prune")

	response, err := h.DockerSock.Post(uri, "application/json", r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	defer response.Body.Close()

	var deletedContainers pruneResponse
	if err := ReadJson(response.Body, &deletedContainers); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	return WriteJson(w, http.StatusOK, deletedContainers)
}
