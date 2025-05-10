// Package tests provides tests for verifying the GRA framework migration
// NOTE: This test will pass only after the migration is complete
package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/lamboktulussimamora/gra/context"
	"github.com/lamboktulussimamora/gra/router"
	"github.com/lamboktulussimamora/gra/validator"
)

const validatePath = "/validate"

// TestGraFrameworkIntegration verifies that the GRA framework works correctly
func TestGraFrameworkIntegration(t *testing.T) {
	// Create a new router
	r := router.New()

	// Define test routes
	r.GET("/hello", func(c *context.Context) {
		c.Success(http.StatusOK, "Hello, World!", nil)
	})

	r.POST(validatePath, func(c *context.Context) {
		type TestData struct {
			Name  string `json:"name" validate:"required"`
			Email string `json:"email" validate:"required,email"`
		}

		var data TestData
		if err := c.BindJSON(&data); err != nil {
			c.Error(http.StatusBadRequest, "Invalid request format")
			return
		}

		// Validate the data
		v := validator.New()
		errors := v.Validate(&data)
		if len(errors) > 0 {
			c.JSON(http.StatusBadRequest, map[string]interface{}{
				"status": "error",
				"error":  "Validation failed",
				"errors": errors,
			})
			return
		}

		c.Success(http.StatusOK, "Validation successful", data)
	})

	// Test the hello endpoint
	t.Run("Hello Endpoint", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/hello", nil)
		w := httptest.NewRecorder()

		// Serve the request
		r.ServeHTTP(w, req)

		// Check the response
		resp := w.Result()
		assertStatus(t, resp.StatusCode, http.StatusOK, "Expected status %d, got %d")

		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)

		// Verify response fields
		if response["status"] != "success" {
			t.Errorf("Expected status 'success', got '%v'", response["status"])
		}

		if response["message"] != "Hello, World!" {
			t.Errorf("Expected message 'Hello, World!', got '%v'", response["message"])
		}
	})

	// Test validation success
	t.Run("Validation Success", func(t *testing.T) {
		jsonStr := `{"name":"John Doe","email":"john@example.com"}`
		req := httptest.NewRequest("POST", validatePath, strings.NewReader(jsonStr))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Serve the request
		r.ServeHTTP(w, req)

		// Check the response
		resp := w.Result()
		assertStatus(t, resp.StatusCode, http.StatusOK, "Expected status %d, got %d")
	})

	// Test validation failure
	t.Run("Validation Failure", func(t *testing.T) {
		jsonStr := `{"name":"","email":"invalid-email"}`
		req := httptest.NewRequest("POST", validatePath, strings.NewReader(jsonStr))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Serve the request
		r.ServeHTTP(w, req)

		// Check the response
		resp := w.Result()
		assertStatus(t, resp.StatusCode, http.StatusBadRequest, "Expected status %d, got %d")

		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)

		// Verify response fields
		if response["status"] != "error" {
			t.Errorf("Expected status 'error', got '%v'", response["status"])
		}

		if response["error"] != "Validation failed" {
			t.Errorf("Expected error 'Validation failed', got '%v'", response["error"])
		}

		errors, ok := response["errors"].([]interface{})
		if !ok || len(errors) == 0 {
			t.Error("Expected validation errors, got none")
		}
	})
}

// Helper function to assert HTTP status codes
func assertStatus(t *testing.T, got, want int, msg string) {
	t.Helper()
	if got != want {
		t.Errorf(msg, want, got)
	}
}
