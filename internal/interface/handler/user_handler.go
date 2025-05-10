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

// LoginRequest represents the login request data
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthResponseDTO represents the authentication response data
type AuthResponseDTO struct {
	User  UserResponseDTO `json:"user"`
	Token string          `json:"token"`
}

// Register handles user registration requests
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	// Check if the request method is POST
	if r.Method != http.MethodPost {
		SendJSONResponse(w, http.StatusMethodNotAllowed, APIResponse{
			Status: "error",
			Error:  "Method not allowed",
		})
		return
	}

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		SendJSONResponse(w, http.StatusBadRequest, APIResponse{
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
		SendJSONResponse(w, http.StatusBadRequest, APIResponse{
			Status: "error",
			Error:  "Invalid request format",
		})
		log.Printf("Error unmarshaling JSON: %v", err)
		return
	}

	// Call the use case
	userResp, err := h.userUseCase.Register(req.FirstName, req.LastName, req.Email, req.Password)
	if err != nil {
		SendJSONResponse(w, http.StatusBadRequest, APIResponse{
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
	SendJSONResponse(w, http.StatusCreated, APIResponse{
		Status:  "success",
		Message: "User registered successfully",
		Data:    responseData,
	})
}

// Login handles user login requests
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	// Check if the request method is POST
	if r.Method != http.MethodPost {
		SendJSONResponse(w, http.StatusMethodNotAllowed, APIResponse{
			Status: "error",
			Error:  "Method not allowed",
		})
		return
	}

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		SendJSONResponse(w, http.StatusBadRequest, APIResponse{
			Status: "error",
			Error:  "Failed to read request body",
		})
		log.Printf("Error reading request body: %v", err)
		return
	}
	defer r.Body.Close()

	// Parse the JSON request into LoginRequest struct
	var req LoginRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		SendJSONResponse(w, http.StatusBadRequest, APIResponse{
			Status: "error",
			Error:  "Invalid request format",
		})
		log.Printf("Error unmarshaling JSON: %v", err)
		return
	}

	// Call the use case
	authResp, err := h.userUseCase.Login(req.Email, req.Password)
	if err != nil {
		SendJSONResponse(w, http.StatusUnauthorized, APIResponse{
			Status: "error",
			Error:  err.Error(),
		})
		return
	}

	// Convert domain response to DTO
	responseData := AuthResponseDTO{
		User: UserResponseDTO{
			FirstName: authResp.User.FirstName,
			LastName:  authResp.User.LastName,
			Email:     authResp.User.Email,
			CreatedAt: authResp.User.CreatedAt.Format(time.RFC3339),
			UpdatedAt: authResp.User.UpdatedAt.Format(time.RFC3339),
		},
		Token: authResp.Token,
	}

	// Return success response
	SendJSONResponse(w, http.StatusOK, APIResponse{
		Status:  "success",
		Message: "Login successful",
		Data:    responseData,
	})
}
