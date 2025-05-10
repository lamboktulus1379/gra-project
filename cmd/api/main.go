package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/lamboktulussimamora/gra-project/internal/interface/handler"
	"github.com/lamboktulussimamora/gra-project/internal/interface/repository"
	"github.com/lamboktulussimamora/gra-project/internal/usecase"
)

func main() {
	// Create repository
	userRepo := repository.NewInMemoryUserRepository()

	// Create use cases
	userUseCase := usecase.NewUserUseCase(userRepo)

	// Create handlers
	userHandler := handler.NewUserHandler(userUseCase)
	helloHandler := handler.NewHelloHandler()

	// Register the handlers for the endpoints
	http.HandleFunc("/hello", helloHandler.Hello)
	http.HandleFunc("/register", userHandler.Register)

	// Print a message indicating that the server is starting
	fmt.Println("Starting server on :8080")

	// Start the HTTP server on port 8080
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
