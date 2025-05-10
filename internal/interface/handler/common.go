package handler

import (
	"net/http"

	"github.com/lamboktulussimamora/gra-project/internal/interface/common"
)

// APIResponse is an alias for common.APIResponse for backward compatibility
type APIResponse = common.APIResponse

// SendJSONResponse is a wrapper for common.SendJSONResponse for backward compatibility
func SendJSONResponse(w http.ResponseWriter, status int, response common.APIResponse) {
	common.SendJSONResponse(w, status, response)
}
