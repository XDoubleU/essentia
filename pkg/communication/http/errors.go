package http

import (
	"errors"
	"net/http"

	"github.com/xdoubleu/essentia/pkg/context"
	errortools "github.com/xdoubleu/essentia/pkg/errors"
	"github.com/xdoubleu/essentia/pkg/logging"
)

func HandleError(w http.ResponseWriter, r *http.Request, err error) {
	notFoundError := errortools.NotFoundError{}
	conflictError := errortools.ConflictError{}
	badRequestError := errortools.BadRequestError{}

	switch {
	case errors.As(err, &badRequestError):
		BadRequestResponse(w, r, badRequestError)
	case errors.As(err, &notFoundError):
		NotFoundResponse(w, r, notFoundError)
	case errors.As(err, &conflictError):
		ConflictResponse(w, r, conflictError)
	default:
		ServerErrorResponse(w, r, err)
	}
}

// ErrorResponse is used to handle any kind of error.
func ErrorResponse(w http.ResponseWriter, r *http.Request,
	status int, message any) {
	errorDto := errortools.NewErrorDto(status, message)
	err := WriteJSON(w, status, errorDto, nil)
	if err != nil {
		context.Logger(r.Context()).
			ErrorContext(r.Context(), "failed to write JSON", logging.ErrAttr(err))
	}
}

// ServerErrorResponse is used to handle internal server errors.
func ServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	context.Logger(r.Context()).
		ErrorContext(r.Context(), "server error occurred", logging.ErrAttr(err))

	message := errortools.MessageInternalServerError
	if context.ShowErrors(r.Context()) {
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
	ErrorResponse(w, r, http.StatusTooManyRequests, errortools.MessageTooManyRequests)
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
	ErrorResponse(w, r, http.StatusForbidden, errortools.MessageForbidden)
}

// ConflictResponse is used to handle an error when a resource already exists.
func ConflictResponse(
	w http.ResponseWriter,
	r *http.Request,
	err errortools.ConflictError,
) {
	outputErr := make(map[string]string)
	outputErr[err.JsonField] = err.Error()
	ErrorResponse(w, r, http.StatusNotFound, err)
}

// NotFoundResponse is used to handle an error when a resource wasn't found.
func NotFoundResponse(
	w http.ResponseWriter,
	r *http.Request,
	err errortools.NotFoundError,
) {
	outputErr := make(map[string]string)
	outputErr[err.JsonField] = err.Error()
	ErrorResponse(w, r, http.StatusNotFound, err)
}

// FailedValidationResponse is used to handle an error of a [validate.Validator].
func FailedValidationResponse(
	w http.ResponseWriter,
	r *http.Request,
	errors map[string]string,
) {
	ErrorResponse(w, r, http.StatusUnprocessableEntity, errors)
}
