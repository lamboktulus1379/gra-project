package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/lamboktulussimamora/gra-project/internal/domain/auth"
	"github.com/lamboktulussimamora/gra-project/internal/interface/common"
)

// AuthMiddleware is a middleware that authenticates requests
type AuthMiddleware struct {
	jwtService auth.JWTService
}

// NewAuthMiddleware creates a new authentication middleware
func NewAuthMiddleware(jwtService auth.JWTService) *AuthMiddleware {
	return &AuthMiddleware{
		jwtService: jwtService,
	}
}

// Authenticate middleware checks for valid JWT token in Authorization header
func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			common.SendJSONResponse(w, http.StatusUnauthorized, common.APIResponse{
				Status: "error",
				Error:  "Authorization header is required",
			})
			return
		}

		// Check if the header has the correct format (Bearer <token>)
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			common.SendJSONResponse(w, http.StatusUnauthorized, common.APIResponse{
				Status: "error",
				Error:  "Authorization header format must be Bearer <token>",
			})
			return
		}

		// Extract the token
		tokenString := parts[1]

		// Validate the token
		claims, err := m.jwtService.ValidateToken(tokenString)
		if err != nil {
			var errorMsg string
			if err == auth.ErrExpiredToken {
				errorMsg = "Token has expired"
			} else {
				errorMsg = "Invalid token"
			}
			common.SendJSONResponse(w, http.StatusUnauthorized, common.APIResponse{
				Status: "error",
				Error:  errorMsg,
			})
			return
		}

		// Add claims to context using common UserClaimsKey
		ctx := context.WithValue(r.Context(), common.UserClaimsKey, claims)

		// Call the next handler with the updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
