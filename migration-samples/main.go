package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/lamboktulussimamora/gra/core" // Updated import
	"github.com/lamboktulussimamora/gra-project/internal/domain/auth"
	"github.com/lamboktulussimamora/gra-project/internal/interface/handler"
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
		SecretKey:  "your-secret-key",
		Issuer:     "gra-project",
		Audience:   "gra-users",
		ExpiryTime: 24 * time.Hour,
	}
	jwtService := auth.NewJWTService(jwtConfig)

	// Create use case
	userUseCase := usecase.NewUserUseCase(userRepo, passwordService, jwtService)

	// Create handlers
	userHandler := handler.NewUserHandler(userUseCase)
	exampleHandler := handler.NewExampleHandler()
	protectedHandler := handler.NewProtectedHandler()

	// Create router
	router := core.New()

	// Add global middleware
	router.Use(
		core.LoggingMiddleware(), // Using core package from external gra framework
		core.RecoveryMiddleware(),
	)

	// Public routes
	router.GET("/", exampleHandler.Hello)
	router.POST("/users", exampleHandler.Register)
	router.POST("/login", userHandler.Login)

	// Protected routes with auth middleware
	protectedRoutes := core.New()
	protectedRoutes.Use(core.AuthMiddleware(jwtService))
	protectedRoutes.GET("/protected", protectedHandler.Protected)

	// Mount protected routes
	router.GET("/api/*path", protectedRoutes.ServeHTTP)

	// Start the server
	fmt.Println("Server starting at http://localhost:8080")
	if err := core.Run(":8080", router); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
