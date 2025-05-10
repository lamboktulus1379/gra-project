package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// User represents the user registration data
type User struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

// UserResponse represents the user data that is safe to return in API responses
type UserResponse struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// APIResponse is a standardized response structure for all API handlers
type APIResponse struct {
	Status  string `json:"status"`          // "success" or "error"
	Message string `json:"message"`         // Human-readable message
	Data    any    `json:"data,omitempty"`  // Optional data payload
	Error   string `json:"error,omitempty"` // Error message if status is "error"
}

// sendJSONResponse is a helper function to send an APIResponse
func sendJSONResponse(w http.ResponseWriter, status int, response APIResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// helloHandler responds to requests with a JSON message
func helloHandler(w http.ResponseWriter, r *http.Request) {
	// Create a standardized response
	response := APIResponse{
		Status:  "success",
		Message: "Hello, World!",
	}

	// Send the response
	sendJSONResponse(w, http.StatusOK, response)
}

// registerHandler handles user registration requests
func registerHandler(w http.ResponseWriter, r *http.Request) {
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

	// Parse the JSON request into User struct
	var user User
	err = json.Unmarshal(body, &user)
	if err != nil {
		sendJSONResponse(w, http.StatusBadRequest, APIResponse{
			Status: "error",
			Error:  "Invalid request format",
		})
		log.Printf("Error unmarshaling JSON: %v", err)
		return
	}

	// Validate required fields
	if user.FirstName == "" || user.LastName == "" || user.Email == "" || user.Password == "" {
		sendJSONResponse(w, http.StatusBadRequest, APIResponse{
			Status: "error",
			Error:  "Missing required fields",
		})
		return
	}

	// In a real application, you would:
	// 1. Hash the password
	// 2. Save user to database
	// 3. Generate authentication token

	// Create user data for response using the UserResponse struct
	userData := UserResponse{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
	}

	// Create standardized response
	response := APIResponse{
		Status:  "success",
		Message: "User registered successfully",
		Data:    userData,
	}

	// Send the response
	sendJSONResponse(w, http.StatusCreated, response)
}

func main() {
	// Register the handler for the /hello endpoint
	http.HandleFunc("/hello", helloHandler)

	// Register the handler for the /register endpoint
	http.HandleFunc("/register", registerHandler)

	// Print a message indicating that the server is starting
	fmt.Println("Starting server on :8080")

	// Start the HTTP server on port 8080
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
