package main

import (
	"fmt"
	"log"
	"net/http"
)

// helloHandler responds to requests with the hello world message
func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
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
