package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/lamboktulussimamora/gra-project/internal/compatibility"
	"github.com/lamboktulussimamora/gra-project/internal/domain/auth"
	"github.com/lamboktulussimamora/gra-project/internal/interface/common"
	"github.com/lamboktulussimamora/gra-project/internal/interface/handler"
	"github.com/lamboktulussimamora/gra-project/internal/interface/repository"
	"github.com/lamboktulussimamora/gra-project/internal/usecase"
	"github.com/lamboktulussimamora/gra/context"
	"github.com/lamboktulussimamora/gra/middleware"
	"github.com/lamboktulussimamora/gra/router"
	"github.com/lamboktulussimamora/gra/validator"
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
	exampleHandler := handler.NewExampleHandler()

	// Create router
	r := router.New()

	// Apply global middleware
	r.Use(
		middleware.Logger(),
		middleware.Recovery(),
		compatibility.CORSMiddleware("*"),
	)

	// Register public routes
	r.GET("/hello", exampleHandler.Hello)
	r.POST("/register", exampleHandler.Register)
	r.POST("/login", func(c *context.Context) {
		var req handler.LoginRequest
		if err := c.BindJSON(&req); err != nil {
			c.Error(http.StatusBadRequest, "Invalid request format")
			return
		}

		// Validate request
		v := validator.New()
		errors := v.Validate(&req)
		if len(errors) > 0 {
			c.JSON(http.StatusBadRequest, map[string]interface{}{
				"status": "error",
				"error":  "Validation failed",
				"errors": errors,
			})
			return
		}

		// Call the use case
		authResp, err := userUseCase.Login(req.Email, req.Password)
		if err != nil {
			c.Error(http.StatusUnauthorized, err.Error())
			return
		}

		// Return success response
		c.Success(http.StatusOK, "Login successful", map[string]interface{}{
			"token": authResp.Token,
			"user": map[string]interface{}{
				"first_name": authResp.User.FirstName,
				"last_name":  authResp.User.LastName,
				"email":      authResp.User.Email,
			},
		})
	})

	// Create a group of protected routes
	protectedRouter := router.New()
	// Apply auth middleware to all routes in this router
	protectedRouter.Use(compatibility.AuthMiddleware(jwtService, common.UserClaimsKey))

	// Add protected routes
	protectedRouter.GET("/profile", exampleHandler.Profile)

	// Mount the protected router under the /api path
	r.Handle(http.MethodGet, "/api/profile", func(c *context.Context) {
		exampleHandler.Profile(c)
	})

	// Start server
	fmt.Println("Server started on :8082")
	log.Fatal(http.ListenAndServe(":8082", r))
}
