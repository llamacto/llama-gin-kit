package response

import (
	"errors"
	"net/http"
	
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ErrorResponse represents an error response structure
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// BadRequest sends a 400 Bad Request response
func BadRequest(c *gin.Context, message string, err error) {
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}
	
	c.JSON(http.StatusBadRequest, ErrorResponse{
		Code:    http.StatusBadRequest,
		Message: message,
		Error:   errMsg,
	})
}

// Unauthorized sends a 401 Unauthorized response
func Unauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, ErrorResponse{
		Code:    http.StatusUnauthorized,
		Message: message,
	})
}

// Forbidden sends a 403 Forbidden response
func Forbidden(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, ErrorResponse{
		Code:    http.StatusForbidden,
		Message: message,
	})
}

// NotFound sends a 404 Not Found response
func NotFound(c *gin.Context, message string, err error) {
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}
	
	c.JSON(http.StatusNotFound, ErrorResponse{
		Code:    http.StatusNotFound,
		Message: message,
		Error:   errMsg,
	})
}

// InternalServerError sends a 500 Internal Server Error response
func InternalServerError(c *gin.Context, message string, err error) {
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}
	
	c.JSON(http.StatusInternalServerError, ErrorResponse{
		Code:    http.StatusInternalServerError,
		Message: message,
		Error:   errMsg,
	})
}

// HandleError handles common errors and returns appropriate responses
func HandleError(c *gin.Context, message string, err error) {
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		NotFound(c, message, err)
	case err.Error() == "unauthorized":
		Unauthorized(c, message)
	default:
		InternalServerError(c, message, err)
	}
}
