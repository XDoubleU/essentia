// Package errors contains reusable error messages
// and other helpers for dealing with errors.
package errors

import (
	"errors"
	"net/http"
)

var (
	// ErrResourceNotFound is an error with value "resource not found".
	ErrResourceNotFound = errors.New("resource not found")
	// ErrResourceUniqueValue is an error with value "resource unique value already used".
	ErrResourceConflict = errors.New("resource conflicts with existing resource")
	ErrFailedValidation = errors.New("failed validation")
)

//nolint:lll // can't make these lines shorter
const (
	MessageInternalServerError = "the server encountered a problem and could not process your request"
	MessageTooManyRequests     = "rate limit exceeded"
	MessageForbidden           = "user has no access to this resource"
)

// ErrorDto is used to return the error back to the client.
type ErrorDto struct {
	Status  int    `json:"status"`
	Error   string `json:"error"`
	Message any    `json:"message"`
} //	@name	ErrorDto

// NewErrorDto creates a new [ErrorDto].
func NewErrorDto(status int, message any) ErrorDto {
	return ErrorDto{
		Status:  status,
		Error:   http.StatusText(status),
		Message: message,
	}
}
