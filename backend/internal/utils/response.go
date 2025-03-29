package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

func WriteErrorResponse(w http.ResponseWriter, code int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	response := map[string]string{"error": err.Error()}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Failed to encode error response: %v", err)
	}
}

func WriteJSONResponse(w http.ResponseWriter, response map[string]interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
