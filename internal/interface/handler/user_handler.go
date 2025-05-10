package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/lamboktulussimamora/gra-project/internal/usecase"
)

// UserHandler handles HTTP requests related to users
type UserHandler struct {
	userUseCase *usecase.UserUseCase
}

// NewUserHandler creates a new user handler
func NewUserHandler(userUseCase *usecase.UserUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
	}
}

// RegisterUserRequest represents the user registration request data
type RegisterUserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

// UserResponseDTO represents the user data that is returned in API responses
type UserResponseDTO struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// APIResponse is a standardized response structure for all API handlers
type APIResponse struct {
	Status  string      `json:"status"`          // "success" or "error"
	Message string      `json:"message"`         // Human-readable message
	Data    interface{} `json:"data,omitempty"`  // Optional data payload
	Error   string      `json:"error,omitempty"` // Error message if status is "error"
}

// sendJSONResponse is a helper function to send an APIResponse
func sendJSONResponse(w http.ResponseWriter, status int, response APIResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// Register handles user registration requests
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	// Check if the request method is POST
	if r.Method != http.MethodPost {
		sendJSONResponse(w, http.StatusMethodNotAllowed, APIResponse{
			Status: "error",
			Error:  "Method not allowed",
		})
		return
	}

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		sendJSONResponse(w, http.StatusBadRequest, APIResponse{
			Status: "error",
			Error:  "Failed to read request body",
		})
		log.Printf("Error reading request body: %v", err)
		return
	}
	defer r.Body.Close()

	// Parse the JSON request into RegisterUserRequest struct
	var req RegisterUserRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		sendJSONResponse(w, http.StatusBadRequest, APIResponse{
			Status: "error",
			Error:  "Invalid request format",
		})
		log.Printf("Error unmarshaling JSON: %v", err)
		return
	}

	// Call the use case
	userResp, err := h.userUseCase.Register(req.FirstName, req.LastName, req.Email, req.Password)
	if err != nil {
		sendJSONResponse(w, http.StatusBadRequest, APIResponse{
			Status: "error",
			Error:  err.Error(),
		})
		return
	}

	// Convert domain response to DTO
	responseData := UserResponseDTO{
		FirstName: userResp.FirstName,
		LastName:  userResp.LastName,
		Email:     userResp.Email,
		CreatedAt: userResp.CreatedAt.Format(time.RFC3339),
		UpdatedAt: userResp.UpdatedAt.Format(time.RFC3339),
	}

	// Return success response
	sendJSONResponse(w, http.StatusCreated, APIResponse{
		Status:  "success",
		Message: "User registered successfully",
		Data:    responseData,
	})
}
