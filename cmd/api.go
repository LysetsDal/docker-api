package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type apiFunc func(http.ResponseWriter, *http.Request) error

type apiError struct {
	Error string
}

type APIServer struct {
	listenAddr string
}

func WriteJson(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

func makeHttpHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJson(w, http.StatusBadRequest, apiError{Error: err.Error()})
		}
	}
}

func NewAPIServer(listenAddr string) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()
	router.HandleFunc("/test", makeHttpHandleFunc(s.handleRequest))

	log.Printf("DockerAPI Server listening on %s\n", s.listenAddr)

	http.ListenAndServe(s.listenAddr, router)
}

func (s *APIServer) handleRequest(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetRequest(w, r)
	}
	if r.Method == "POST" {
		return s.handlePostRequest(w, r)
	}

	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *APIServer) handleGetRequest(w http.ResponseWriter, r *http.Request) error {
	fmt.Println("Get Request")
	return nil
}

func (s *APIServer) handlePostRequest(w http.ResponseWriter, r *http.Request) error {
	fmt.Println("Post Request")
	return nil
}
