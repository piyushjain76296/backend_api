package handlers

import (
	"encoding/json"
	"net/http"
)

// HealthResponse represents the response format for the health check.
type HealthResponse struct {
	Status string `json:"status"`
}

// HealthHandler responds with a basic 200 OK and a status JSON.
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp := HealthResponse{Status: "ok"}
	json.NewEncoder(w).Encode(resp)
}
