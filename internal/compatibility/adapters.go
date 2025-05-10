// Package compatibility provides adapter functions for using the external gra framework
package compatibility

import (
	"net/http"
	"strings"

	"github.com/lamboktulussimamora/gra-project/internal/domain/auth"
	"github.com/lamboktulussimamora/gra/context"
	"github.com/lamboktulussimamora/gra/router"
)

// JWTAuthAdapter adapts our JWTService to implement middleware.JWTAuthenticator
type JWTAuthAdapter struct {
	jwtService auth.JWTService
}

// NewJWTAuthAdapter creates a new adapter for JWTService
func NewJWTAuthAdapter(jwtService auth.JWTService) *JWTAuthAdapter {
	return &JWTAuthAdapter{
		jwtService: jwtService,
	}
}

// ValidateToken implements the middleware.JWTAuthenticator interface
func (a *JWTAuthAdapter) ValidateToken(tokenString string) (interface{}, error) {
	return a.jwtService.ValidateToken(tokenString)
}

// AuthMiddleware creates an authentication middleware
func AuthMiddleware(jwtService auth.JWTService, claimsKey interface{}) router.Middleware {
	jwtAdapter := NewJWTAuthAdapter(jwtService)

	return func(next router.HandlerFunc) router.HandlerFunc {
		return func(c *context.Context) {
			// Get the Authorization header
			authHeader := c.Request.Header.Get("Authorization")
			if authHeader == "" {
				c.Error(http.StatusUnauthorized, "Authorization header is required")
				return
			}

			// Check if the header has the correct format (Bearer <token>)
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				c.Error(http.StatusUnauthorized, "Authorization header format must be Bearer <token>")
				return
			}

			// Extract the token
			tokenString := parts[1]

			// Validate the token
			claims, err := jwtAdapter.ValidateToken(tokenString)
			if err != nil {
				c.Error(http.StatusUnauthorized, "Invalid token")
				return
			}

			// Add claims to context
			c.WithValue(claimsKey, claims)

			// Call the next handler
			next(c)
		}
	}
}

// CORSMiddleware creates a CORS middleware
func CORSMiddleware(origin string) router.Middleware {
	return func(next router.HandlerFunc) router.HandlerFunc {
		return func(c *context.Context) {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			// Handle preflight requests
			if c.Request.Method == "OPTIONS" {
				c.Writer.WriteHeader(http.StatusOK)
				return
			}

			next(c)
		}
	}
}
