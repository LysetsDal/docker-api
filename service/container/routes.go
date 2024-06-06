package container

import (
	"bytes"
	"encoding/json"
	"fmt"
	. "github.com/LysetsDal/docker-api/config"
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
	router.HandleFunc("/containers/stopall", MakeHttpHandleFunc(h.handleStopAllContainerRequest))
	router.HandleFunc("/containers/prune", MakeHttpHandleFunc(h.handlePruneContainersRequest))
}

// GET List of all containers
func (h *Handler) handleGetContainerById(w http.ResponseWriter, r *http.Request) error {
	pathVars := mux.Vars(r)
	fmt.Println(pathVars["id"])

	uri := fmt.Sprintf(UnixPrefix+"containers/%s/json", pathVars["id"])

	response, err := h.DockerSock.Get(uri)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err)
		return err
	}
	defer response.Body.Close()

	// Read and decode the JSON response
	var inspectObject InspectObject
	if err := ReadJson(response.Body, &inspectObject); err != nil {
		WriteError(w, http.StatusInternalServerError, err)
		return err
	}

	return WriteJson(w, response.StatusCode, inspectObject)
}

// GET List of all containers
func (h *Handler) handleGetContainerList(w http.ResponseWriter, _ *http.Request) error {
	response, err := h.DockerSock.Get(UnixPrefix + "containers/json")
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err)
		return err
	}
	defer response.Body.Close()

	// Read and decode the JSON response
	var containers []Container
	if err := ReadJson(response.Body, &containers); err != nil {
		WriteError(w, http.StatusInternalServerError, err)
		return err
	}

	return WriteJson(w, response.StatusCode, containers)
}

// POST Start container
func (h *Handler) handleStartContainerRequest(w http.ResponseWriter, r *http.Request) error {
	pathVars := mux.Vars(r)
	uri := fmt.Sprintf(UnixPrefix+"containers/%s/start", pathVars["id"])

	request, err := h.DockerSock.Post(uri, "application/json", r.Body)
	if err != nil {
		WriteError(w, http.StatusNotFound, err)
		return err
	}
	defer request.Body.Close()

	switch request.StatusCode {
	case http.StatusNoContent:
		return WriteJson(w, http.StatusOK, ApiMessage{Message: "Container started"})
	case http.StatusNotModified:
		return WriteJson(w, http.StatusOK, ApiMessage{Message: "Container already started"})
	case http.StatusNotFound:
		return WriteJson(w, http.StatusNotFound, ApiError{Error: "No such container"})
	default:
		return WriteJson(w, http.StatusInternalServerError, ApiError{Error: "Something went wrong"})
	}
}

// POST Stop container
func (h *Handler) handleStopContainerRequest(w http.ResponseWriter, r *http.Request) error {
	pathVars := mux.Vars(r)
	uri := fmt.Sprintf(UnixPrefix+"containers/%s/stop", pathVars["id"])

	request, err := h.DockerSock.Post(uri, "application/json", r.Body)
	if err != nil {
		WriteError(w, http.StatusNotFound, err)
		return err
	}
	defer request.Body.Close()

	fmt.Printf("r.Body: %s\n", r.Body)

	switch request.StatusCode {
	case http.StatusNoContent:
		return WriteJson(w, http.StatusOK, ApiMessage{Message: "Container stopped"})
	case http.StatusNotModified:
		return WriteJson(w, http.StatusOK, ApiMessage{Message: "Container already stopped"})
	case http.StatusNotFound:
		return WriteJson(w, http.StatusNotFound, ApiError{Error: "No such container"})
	default:
		return WriteJson(w, http.StatusInternalServerError, ApiError{Error: "Something went wrong"})
	}
}

// POST Stop all containers
func (h *Handler) handleStopAllContainerRequest(w http.ResponseWriter, r *http.Request) error {
	// get list of all container ids
	response, err := h.DockerSock.Get(UnixPrefix + "containers/json")
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err)
		return err
	}
	defer response.Body.Close()

	// Read and decode the JSON response
	var containers []Container
	if err := ReadJson(response.Body, &containers); err != nil {
		WriteError(w, http.StatusInternalServerError, err)
		return err
	}

	params := StopParams{Signal: "SIGINT", T: 5}
	paramsJson, err := json.Marshal(params)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err)
	}

	// Iterate over all containers, and send POST request to stop.
	for _, c := range containers {
		apiUrl := fmt.Sprintf(UnixPrefix+"containers/%s/stop", c.Id)

		request, err := h.DockerSock.Post(apiUrl, "application/json", bytes.NewBuffer(paramsJson))
		if err != nil {
			WriteError(w, http.StatusInternalServerError, err)
			return err
		}
		request.Body.Close()
	}

	return WriteJson(w, http.StatusOK, ApiMessage{Message: "all containers stopped"})
}

func (h *Handler) handlePruneContainersRequest(w http.ResponseWriter, r *http.Request) error {
	type pruneResponse struct {
		ContainersDeleted []string `json:"ContainersDeleted"`
		SpaceReclaimed    int      `json:"SpaceReclaimed"`
	}

	uri := fmt.Sprintf(UnixPrefix + "containers/prune")

	response, err := h.DockerSock.Post(uri, "application/json", r.Body)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err)
		return err
	}
	defer response.Body.Close()

	var deletedContainers pruneResponse
	if err := ReadJson(response.Body, &deletedContainers); err != nil {
		WriteError(w, http.StatusInternalServerError, err)
		return err
	}

	return WriteJson(w, http.StatusOK, deletedContainers)
}
