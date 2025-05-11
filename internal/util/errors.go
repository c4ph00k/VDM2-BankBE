package util

import (
	"net/http"

	"github.com/gin-gonic/gin"
	// "github.com/pkg/errors"
)

// APIError represents a structured error response for the API
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Error implements the error interface
func (e APIError) Error() string {
	return e.Message
}

// NewAPIError creates a new APIError with the given code and message
func NewAPIError(code int, message string) *APIError {
	return &APIError{
		Code:    code,
		Message: message,
	}
}

// NewBadRequestError creates a new 400 Bad Request error
func NewBadRequestError(message string) *APIError {
	return NewAPIError(http.StatusBadRequest, message)
}

// NewUnauthorizedError creates a new 401 Unauthorized error
func NewUnauthorizedError(message string) *APIError {
	return NewAPIError(http.StatusUnauthorized, message)
}

// NewForbiddenError creates a new 403 Forbidden error
func NewForbiddenError(message string) *APIError {
	return NewAPIError(http.StatusForbidden, message)
}

// NewNotFoundError creates a new 404 Not Found error
func NewNotFoundError(message string) *APIError {
	return NewAPIError(http.StatusNotFound, message)
}

// NewInternalServerError creates a new 500 Internal Server Error
func NewInternalServerError(message string) *APIError {
	return NewAPIError(http.StatusInternalServerError, message)
}

// ErrorResponse wraps an APIError for HTTP responses
type ErrorResponse struct {
	Error *APIError `json:"error"`
}

// HandleError handles errors in Gin context
func HandleError(c *gin.Context, err error) {
	var apiErr *APIError
	switch e := err.(type) {
	case *APIError:
		apiErr = e
	default:
		apiErr = NewInternalServerError("internal server error")
	}

	resp := ErrorResponse{Error: apiErr}
	c.JSON(apiErr.Code, resp)
}
