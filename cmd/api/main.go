package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/lamboktulussimamora/gra-project/internal/domain/auth"
	"github.com/lamboktulussimamora/gra-project/internal/interface/handler"
	"github.com/lamboktulussimamora/gra-project/internal/interface/middleware"
	"github.com/lamboktulussimamora/gra-project/internal/interface/repository"
	"github.com/lamboktulussimamora/gra-project/internal/usecase"
)

func main() {
	// Create repository
	userRepo := repository.NewInMemoryUserRepository()

	// Initialize password service
	passwordParams := auth.ArgonParams{
		Memory:      64 * 1024,
		Iterations:  3,
		Parallelism: 2,
		SaltLength:  16,
		KeyLength:   32,
	}
	passwordService := auth.NewPasswordService(passwordParams)

	// Initialize JWT service
	jwtConfig := auth.JWTConfig{
		SecretKey:     "your-secret-key-here", // In production, use environment variables for secrets
		TokenDuration: time.Hour * 24,         // 24 hours token validity
	}
	jwtService := auth.NewJWTService(jwtConfig)

	// Create use cases
	userUseCase := usecase.NewUserUseCase(userRepo, passwordService, jwtService)

	// Create handlers
	userHandler := handler.NewUserHandler(userUseCase)
	helloHandler := handler.NewHelloHandler()
	protectedHandler := handler.NewProtectedHandler()

	// Create middleware
	authMiddleware := middleware.NewAuthMiddleware(jwtService)

	// Register public endpoints
	http.HandleFunc("/hello", helloHandler.Hello)
	http.HandleFunc("/register", userHandler.Register)
	http.HandleFunc("/login", userHandler.Login)

	// Register protected endpoints with auth middleware
	http.Handle("/profile", authMiddleware.Authenticate(http.HandlerFunc(protectedHandler.Profile)))

	// Print a message indicating that the server is starting
	fmt.Println("Starting server on :8080")

	// Start the HTTP server on port 8080
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
