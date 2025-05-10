package common

import (
	"encoding/json"
	"log"
	"net/http"
)

// APIResponse is a standardized response structure for all API handlers
type APIResponse struct {
	Status  string      `json:"status"`          // "success" or "error"
	Message string      `json:"message"`         // Human-readable message
	Data    interface{} `json:"data,omitempty"`  // Optional data payload
	Error   string      `json:"error,omitempty"` // Error message if status is "error"
}

// SendJSONResponse is a helper function to send an APIResponse
func SendJSONResponse(w http.ResponseWriter, status int, response APIResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// UserContextKey is a key for user claims in the request context
type UserContextKey string

// UserClaimsKey is the context key for storing user claims
const UserClaimsKey UserContextKey = "userClaims"
