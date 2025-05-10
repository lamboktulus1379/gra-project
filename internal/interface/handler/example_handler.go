package handler

import (
	"net/http"

	"github.com/lamboktulussimamora/gra-project/internal/compatibility"
	"github.com/lamboktulussimamora/gra-project/internal/domain/auth"
	"github.com/lamboktulussimamora/gra/context"
	"github.com/lamboktulussimamora/gra/validator"
)

// ExampleHandler demonstrates using the core package
type ExampleHandler struct{}

// NewExampleHandler creates a new example handler
func NewExampleHandler() *ExampleHandler {
	return &ExampleHandler{}
}

// Hello demonstrates a simple GET handler
func (h *ExampleHandler) Hello(c *context.Context) {
	c.Success(http.StatusOK, "Hello, World!", nil)
}

// UserRegistrationRequest demonstrates request validation
type UserRegistrationRequest struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
}

// Register demonstrates request validation and processing
func (h *ExampleHandler) Register(c *context.Context) {
	var req UserRegistrationRequest
	if err := c.BindJSON(&req); err != nil {
		c.Error(http.StatusBadRequest, "Invalid request format")
		return
	}

	// Validate the request
	validator := validator.New()
	errors := validator.Validate(&req)
	if len(errors) > 0 {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status": "error",
			"error":  "Validation failed",
			"data":   errors,
		})
		return
	}

	// Process the request (example)
	c.Success(http.StatusCreated, "User registered successfully", map[string]string{
		"first_name": req.FirstName,
		"last_name":  req.LastName,
		"email":      req.Email,
	})
}

// Profile demonstrates accessing user claims from context
func (h *ExampleHandler) Profile(c *context.Context) {
	// Get user claims from context using the compatibility helper
	claimsVal, exists := compatibility.GetUserClaims(c)
	if !exists {
		c.Error(http.StatusUnauthorized, "Unauthorized")
		return
	}

	claims, ok := claimsVal.(*auth.Claims)
	if !ok {
		c.Error(http.StatusUnauthorized, "Invalid user claims")
		return
	}

	// Return user profile data
	c.Success(http.StatusOK, "Profile retrieved successfully", map[string]string{
		"email":      claims.Email,
		"first_name": claims.FirstName,
		"last_name":  claims.LastName,
	})
}
