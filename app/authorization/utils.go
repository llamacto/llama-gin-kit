package authorization

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func getUserIDFromContext(c *gin.Context) (uint, error) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		return 0, errors.New("user not authenticated: userID not found in context")
	}

	userID, ok := userIDVal.(uint)
	if !ok {
		return 0, errors.New("invalid userID format in context")
	}

	if userID == 0 {
		return 0, errors.New("invalid userID: 0")
	}

	return userID, nil
}
