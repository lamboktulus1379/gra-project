package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// helloHandler responds to requests with a JSON message
func helloHandler(w http.ResponseWriter, r *http.Request) {
	// Create a response data structure
	response := map[string]string{
		"message": "Hello, World!",
		"status":  "success",
	}

	// Set content type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Encode the response as JSON and write to ResponseWriter
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		// If encoding fails, set appropriate status code and log the error
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		log.Printf("JSON encoding error: %v", err)
		return
	}
}

func main() {
	// Register the handler for the /hello endpoint
	http.HandleFunc("/hello", helloHandler)

	// Print a message indicating that the server is starting
	fmt.Println("Starting server on :8080")

	// Start the HTTP server on port 8080
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
