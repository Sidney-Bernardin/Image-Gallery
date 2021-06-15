package server

import (
	"encoding/json"
	"net/http"
)

func (s *server) respond(w http.ResponseWriter, status int, data interface{}) {

	// Serialize the data.
	b, err := json.Marshal(data)
	if err != nil {
		s.err(w, err, http.StatusInternalServerError)
		return
	}

	// Respond.
	w.WriteHeader(status)
	_, err = w.Write(b)
	if err != nil {
		s.err(w, err, http.StatusInternalServerError)
		return
	}
}

func (s *server) err(w http.ResponseWriter, e error, status int) {

	// If its a server error, log the error.
	if status >= 500 {
		s.logger.Warnf("Internal Server Error: %v", e)
	}

	// Serialize the error.
	b, err := json.Marshal(map[string]string{"error": e.Error()})
	if err != nil {
		s.logger.Warnf("Internal Server Error: %v", err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	// Response.
	w.WriteHeader(status)
	_, err = w.Write(b)
	if err != nil {
		s.logger.Warnf("Internal Server Error: %v", err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
}
