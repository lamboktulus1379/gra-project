package handler

import (
	"net/http"

	"github.com/lamboktulussimamora/gra/core" // Updated import
)

// ExampleHandler demonstrates using the core package
type ExampleHandler struct{}

// NewExampleHandler creates a new example handler
func NewExampleHandler() *ExampleHandler {
	return &ExampleHandler{}
}

// Hello demonstrates a simple GET handler
func (h *ExampleHandler) Hello(c *core.Context) {
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
func (h *ExampleHandler) Register(c *core.Context) {
	// Parse request body
	var req UserRegistrationRequest
	if err := c.BindJSON(&req); err != nil {
		c.Error(http.StatusBadRequest, "Invalid request format")
		return
	}

	// Validate the request
	validator := core.NewValidator()
	errors := validator.Validate(req)
	if len(errors) > 0 {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status": "error",
			"errors": errors,
		})
		return
	}

	// Process the validated request (just sending back for demonstration)
	c.Success(http.StatusCreated, "User registered", map[string]string{
		"first_name": req.FirstName,
		"last_name":  req.LastName,
		"email":      req.Email,
	})
}
