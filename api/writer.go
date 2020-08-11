package api

import (
	"encoding/json"
	"net/http"
)

// WriteError writes specified error to ResponseWriter
func WriteError(w http.ResponseWriter, statusCode int, body interface{}) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(body)
}

// WriteSuccess writes a success response to ResponseWriter
func WriteSuccess(w http.ResponseWriter, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(body)
}
