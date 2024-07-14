package httptools

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/XDoubleU/essentia/pkg/contexttools"
	"github.com/XDoubleU/essentia/pkg/tools"
	"github.com/getsentry/sentry-go"
	"nhooyr.io/websocket"
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
	env := ErrorDto{
		Status:  status,
		Error:   http.StatusText(status),
		Message: message,
	}
	err := WriteJSON(w, status, env, nil)
	if err != nil {
		sendErrorToSentry(r.Context(), err)
		contexttools.GetLogger(r).Print(err)
	}
}

// ServerErrorResponse is used to handle internal server errors.
func ServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	sendErrorToSentry(r.Context(), err)

	message := MessageInternalServerError
	if contexttools.GetShowErrors(r) {
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

// WSErrorResponse is used to handle an error that occured on a Websocket.
func WSErrorResponse(
	_ http.ResponseWriter,
	r *http.Request,
	conn *websocket.Conn,
	beforeClosingFunc func(conn *websocket.Conn),
	err error,
) {
	// don't want to capture close errors
	if websocket.CloseStatus(err) != -1 {
		return
	}

	sendErrorToSentry(r.Context(), err)

	if beforeClosingFunc != nil {
		beforeClosingFunc(conn)
	}

	conn.Close(websocket.StatusInternalError, MessageInternalServerError)
}

// WSUpgradeErrorResponse is used to handle an error that
// occured during the upgrade towards a Websocket.
func WSUpgradeErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	sendErrorToSentry(r.Context(), err)
	w.WriteHeader(http.StatusInternalServerError)
}

func sendErrorToSentry(ctx context.Context, err error) {
	if hub := sentry.GetHubFromContext(ctx); hub != nil {
		hub.WithScope(func(scope *sentry.Scope) {
			scope.SetLevel(sentry.LevelError)
			hub.CaptureException(err)
		})
	}
}
