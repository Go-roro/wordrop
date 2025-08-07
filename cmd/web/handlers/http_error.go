package handlers

import (
	"encoding/json"
	"net/http"
)

func NewHTTPError(w http.ResponseWriter, message string, code int) {
	response := map[string]string{
		"message": message,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode error response", http.StatusInternalServerError)
		return
	}
}
