package handler

import (
	"net/http"
)

// HelloHandler handles HTTP requests for the hello endpoint
type HelloHandler struct{}

// NewHelloHandler creates a new hello handler
func NewHelloHandler() *HelloHandler {
	return &HelloHandler{}
}

// Hello handles requests to the hello endpoint
func (h *HelloHandler) Hello(w http.ResponseWriter, r *http.Request) {
	// Return a simple hello response
	sendJSONResponse(w, http.StatusOK, APIResponse{
		Status:  "success",
		Message: "Hello, World!",
	})
}
