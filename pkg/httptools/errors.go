package httptools

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/xdoubleu/essentia/pkg/contexttools"
	"github.com/xdoubleu/essentia/pkg/logging"
	"github.com/xdoubleu/essentia/pkg/tools"
)

var (
	// ErrResourceNotFound is an error with value "resource not found".
	ErrResourceNotFound = errors.New("resource not found")
	// ErrResourceUniqueValue is an error with value "resource unique value already used".
	ErrResourceUniqueValue = errors.New("resource unique value already used")
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

// ErrorResponse is used to handle any kind of error.
func ErrorResponse(w http.ResponseWriter, r *http.Request,
	status int, message any) {
	errorDto := ErrorDto{
		Status:  status,
		Error:   http.StatusText(status),
		Message: message,
	}
	err := WriteJSON(w, status, errorDto, nil)
	if err != nil {
		contexttools.Logger(r.Context()).
			ErrorContext(r.Context(), "failed to write JSON", logging.ErrAttr(err))
	}
}

// ServerErrorResponse is used to handle internal server errors.
func ServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	contexttools.Logger(r.Context()).
		ErrorContext(r.Context(), "server error occurred", logging.ErrAttr(err))

	message := MessageInternalServerError
	if contexttools.ShowErrors(r.Context()) {
		message = err.Error()
	}

	ErrorResponse(w, r, http.StatusInternalServerError, message)
}

// BadRequestResponse is used to handle an error when a request is incorrect.
func BadRequestResponse(w http.ResponseWriter,
	r *http.Request, err error) {
	ErrorResponse(w, r, http.StatusBadRequest, err.Error())
}

// RateLimitExceededResponse is used to handle an error when the rate limit is exceeded.
func RateLimitExceededResponse(w http.ResponseWriter,
	r *http.Request) {
	ErrorResponse(w, r, http.StatusTooManyRequests, MessageTooManyRequests)
}

// UnauthorizedResponse is used to handle an error when a user
// isn't authorized.
func UnauthorizedResponse(w http.ResponseWriter,
	r *http.Request, message string) {
	ErrorResponse(w, r, http.StatusUnauthorized, message)
}

// ForbiddenResponse is used to handle an error when a user
// isn't authorized to access a certain resource.
func ForbiddenResponse(w http.ResponseWriter, r *http.Request) {
	ErrorResponse(w, r, http.StatusForbidden, MessageForbidden)
}

// ConflictResponse is used to handle an error when a resource already exists.
func ConflictResponse(
	w http.ResponseWriter,
	r *http.Request,
	err error,
	resourceName string,
	identifierValue any,
	jsonField string,
) {
	value, err2 := tools.AnyToString(identifierValue)
	if err2 != nil {
		ServerErrorResponse(w, r, err2)
		return
	}

	if err == nil || errors.Is(err, ErrResourceUniqueValue) {
		message := fmt.Sprintf(
			"%s with %s '%s' already exists",
			resourceName,
			jsonField,
			value,
		)
		err := make(map[string]string)
		err[jsonField] = message
		ErrorResponse(w, r, http.StatusConflict, err)
	} else {
		ServerErrorResponse(w, r, err)
	}
}

// NotFoundResponse is used to handle an error when a resource wasn't found.
func NotFoundResponse(
	w http.ResponseWriter,
	r *http.Request,
	err error,
	resourceName string,
	identifierValue any,
	jsonField string,
) {
	value, err2 := tools.AnyToString(identifierValue)
	if err2 != nil {
		ServerErrorResponse(w, r, err2)
		return
	}

	if err == nil || errors.Is(err, ErrResourceNotFound) {
		message := fmt.Sprintf(
			"%s with %s '%s' doesn't exist",
			resourceName,
			jsonField,
			value,
		)

		err := make(map[string]string)
		err[jsonField] = message

		ErrorResponse(w, r, http.StatusNotFound, err)
	} else {
		ServerErrorResponse(w, r, err)
	}
}

// FailedValidationResponse is used to handle an error of a [validate.Validator].
func FailedValidationResponse(
	w http.ResponseWriter,
	r *http.Request,
	errors map[string]string,
) {
	ErrorResponse(w, r, http.StatusUnprocessableEntity, errors)
}
