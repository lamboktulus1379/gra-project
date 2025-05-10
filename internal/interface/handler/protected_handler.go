package handler

import (
	"net/http"

	"github.com/lamboktulussimamora/gra-project/internal/domain/auth"
	"github.com/lamboktulussimamora/gra-project/internal/interface/common"
)

// ProtectedHandler handles HTTP requests for protected endpoints
type ProtectedHandler struct{}

// NewProtectedHandler creates a new protected handler
func NewProtectedHandler() *ProtectedHandler {
	return &ProtectedHandler{}
}

// Profile handles requests to the protected profile endpoint
func (h *ProtectedHandler) Profile(w http.ResponseWriter, r *http.Request) {
	// Get user claims from context
	claims, ok := r.Context().Value(common.UserClaimsKey).(*auth.Claims)
	if !ok {
		common.SendJSONResponse(w, http.StatusUnauthorized, common.APIResponse{
			Status: "error",
			Error:  "Unauthorized",
		})
		return
	}

	// Return user profile data
	common.SendJSONResponse(w, http.StatusOK, common.APIResponse{
		Status:  "success",
		Message: "Profile retrieved successfully",
		Data: map[string]string{
			"email":      claims.Email,
			"first_name": claims.FirstName,
			"last_name":  claims.LastName,
		},
	})
}
