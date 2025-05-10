// Package compatibility provides helper functions for compatibility between internal core and external gra framework
package compatibility

import (
	"github.com/lamboktulussimamora/gra-project/internal/interface/common"
	"github.com/lamboktulussimamora/gra/context"
)

// GetUserClaims gets user claims from context
func GetUserClaims(c *context.Context) (interface{}, bool) {
	value := c.Value(common.UserClaimsKey)
	if value == nil {
		return nil, false
	}
	return value, true
}
