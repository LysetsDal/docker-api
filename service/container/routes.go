package container

import (
	"fmt"
	. "github.com/LysetsDal/docker-api/config"
	. "github.com/LysetsDal/docker-api/types"
	. "github.com/LysetsDal/docker-api/utils"
	"github.com/gorilla/mux"
	"io"
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

// RegisterRoutes Main controller (all handle functions added here)
func (h *Handler) RegisterRoutes(router *mux.Router) {
	// Multi container functions
	router.HandleFunc("/containers/list", MakeHttpHandleFunc(h.handleListContainers))
	router.HandleFunc("/containers/create", MakeHttpHandleFunc(h.handleCreateContainer))
	router.HandleFunc("/containers/stopall", MakeHttpHandleFunc(h.handleStopAllContainers))
	router.HandleFunc("/containers/prune", MakeHttpHandleFunc(h.handlePruneContainers))

	// Single container functions
	router.HandleFunc("/containers/{id}/json", MakeHttpHandleFunc(h.handleGetContainerById))
	router.HandleFunc("/containers/{id}/start", MakeHttpHandleFunc(h.handleStartContainer))
	router.HandleFunc("/containers/{id}/stop", MakeHttpHandleFunc(h.handleStopContainer))
	router.HandleFunc("/containers/{id}/top", MakeHttpHandleFunc(h.handleGetContainersProcesses))
}

// sendDockerGetRequest Send a get request to the Docker Socket
func (h *Handler) sendDockerGetRequest(w http.ResponseWriter, url string) (*http.Response, error) {
	response, err := h.DockerSock.Get(url)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err)
		return nil, err
	}

	return response, nil
}

// sendDockerGetRequest Send a get request to the Docker Socket
func (h *Handler) sendDockerGetRequestWithPayload(w http.ResponseWriter, url string, payload io.Reader) (*http.Response, error) {
	//response, err := h.DockerSock.Get(url)
	request, err := http.NewRequest("GET", url, payload)
	if err != nil {
		fmt.Printf("New request error: %s", err)
	}

	response, err := h.DockerSock.Do(request)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err)
		return nil, err
	}

	return response, nil
}

// sendDockerPostRequest Send a POST request to the Docker Socket
func (h *Handler) sendDockerPostRequest(w http.ResponseWriter, url string, reader io.Reader) (*http.Response, error) {
	request, err := h.DockerSock.Post(url, "application/json", reader)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err)
		return nil, err
	}

	return request, nil
}

// handleCreateContainer
// Send POST request to docker. Uses data from Request.Body as container specifications.
func (h *Handler) handleCreateContainer(w http.ResponseWriter, r *http.Request) error {
	url := fmt.Sprintf(UnixPrefix + "containers/create")

	createContainerResponse := CreateContainerResponse{}
	response, _ := h.sendDockerGetRequestWithPayload(w, url, r.Body)
	if response.Body == nil {
		return fmt.Errorf("missing request body")
	}
	defer response.Body.Close()

	switch response.StatusCode {
	case http.StatusCreated:
		if err := ReadJson(response.Body, &createContainerResponse); err != nil {
			fmt.Printf("Json read-error: %s", err)
		}
		return WriteJson(w, http.StatusCreated, createContainerResponse)

	case http.StatusNotModified:
		return WriteJson(w, http.StatusOK, ApiMessage{Message: "Container already started"})

	case http.StatusNotFound:
		return WriteJson(w, http.StatusNotFound, ApiError{Error: fmt.Sprintf("No such image: %s", response.Body)})

	default:
		return WriteJson(w, http.StatusInternalServerError, ApiError{Error: "Something went wrong"})
	}
}

// GET List of all containers
func (h *Handler) handleGetContainerById(w http.ResponseWriter, r *http.Request) error {
	pathVars := mux.Vars(r)
	url := fmt.Sprintf(UnixPrefix+"containers/%s/json", pathVars["id"])

	response, _ := h.sendDockerGetRequest(w, url)
	defer response.Body.Close()

	// Read and decode the JSON response
	inspectObject := InspectObject{}
	if err := ReadJson(response.Body, &inspectObject); err != nil {
		WriteError(w, http.StatusInternalServerError, err)
		return err
	}

	return WriteJson(w, response.StatusCode, inspectObject)
}

func (h *Handler) handleGetContainersProcesses(w http.ResponseWriter, r *http.Request) error {
	pathVars := mux.Vars(r)
	url := fmt.Sprintf(UnixPrefix+"containers/%s/top", pathVars["id"])

	response, _ := h.sendDockerGetRequest(w, url)
	defer response.Body.Close()

	containerProcs := Processes{}
	if err := ReadJson(response.Body, &containerProcs); err != nil {
		WriteError(w, http.StatusInternalServerError, err)
	}

	return WriteJson(w, http.StatusOK, containerProcs)
}

// GET List of all running containers
func (h *Handler) handleListContainers(w http.ResponseWriter, _ *http.Request) error {
	url := fmt.Sprintf(UnixPrefix + "containers/json")

	response, _ := h.sendDockerGetRequest(w, url)
	defer response.Body.Close()

	// Read and decode the JSON response
	containers := make([]Container, 0)
	if err := ReadJson(response.Body, &containers); err != nil {
		WriteError(w, http.StatusInternalServerError, err)
		return err
	}

	return WriteJson(w, response.StatusCode, containers)
}

// POST Start container
func (h *Handler) handleStartContainer(w http.ResponseWriter, r *http.Request) error {
	pathVars := mux.Vars(r)
	url := fmt.Sprintf(UnixPrefix+"containers/%s/start", pathVars["id"])

	request, _ := h.sendDockerGetRequest(w, url)
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
func (h *Handler) handleStopContainer(w http.ResponseWriter, r *http.Request) error {
	pathVars := mux.Vars(r)
	url := fmt.Sprintf(UnixPrefix+"containers/%s/stop", pathVars["id"])

	request, _ := h.sendDockerPostRequest(w, url, r.Body)
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
func (h *Handler) handleStopAllContainers(w http.ResponseWriter, r *http.Request) error {
	url := fmt.Sprintf(UnixPrefix + "containers/json")

	response, _ := h.sendDockerGetRequest(w, url)
	defer response.Body.Close()

	// Read and decode the JSON response
	var containers []Container
	if err := ReadJson(response.Body, &containers); err != nil {
		WriteError(w, http.StatusInternalServerError, err)
		return err
	}

	// Iterate over all containers, and send POST request to stop.
	for _, c := range containers {
		apiUrl := fmt.Sprintf(UnixPrefix+"containers/%s/stop", c.Id)
		request, _ := h.sendDockerPostRequest(w, apiUrl, r.Body)
		request.Body.Close()
	}

	return WriteJson(w, http.StatusOK, ApiMessage{Message: "all containers stopped"})
}

func (h *Handler) handlePruneContainers(w http.ResponseWriter, r *http.Request) error {
	url := fmt.Sprintf(UnixPrefix + "containers/prune")

	response, _ := h.sendDockerPostRequest(w, url, r.Body)
	defer response.Body.Close()

	deletedContainers := PruneResponse{}
	if err := ReadJson(response.Body, &deletedContainers); err != nil {
		WriteError(w, http.StatusInternalServerError, err)
		return err
	}

	return WriteJson(w, http.StatusOK, deletedContainers)
}
