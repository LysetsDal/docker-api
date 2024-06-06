package utils

import (
	"encoding/json"
	. "github.com/LysetsDal/docker-api/types"
	"io"
	"net/http"
)

// ReadJson Read the Docker daemons responses (format json)
func ReadJson(body io.Reader, v any) error {
	decoder := json.NewDecoder(body)
	return decoder.Decode(v)
}

// WriteJson Write Json with standard header.
func WriteJson(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

// MakeHttpHandleFunc Make custom http.HandlerFunc
func MakeHttpHandleFunc(f ApiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			err := WriteJson(w, http.StatusBadRequest, ApiError{Error: err.Error()})
			if err != nil {
				return
			}
		}
	}
}
