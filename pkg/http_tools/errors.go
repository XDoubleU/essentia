package http_tools

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/XDoubleU/essentia/pkg/tools"
	"github.com/getsentry/sentry-go"
	"nhooyr.io/websocket"
)

var (
	ErrRecordNotFound    = errors.New("record not found")
	ErrRecordUniqueValue = errors.New("record unique value already used")
)

type ErrorDto struct {
	Status  int    `json:"status"`
	Error   string `json:"error"`
	Message any    `json:"message"`
} //	@name	ErrorDto

func ErrorResponse(w http.ResponseWriter,
	_ *http.Request, status int, message any) {
	env := ErrorDto{
		Status:  status,
		Error:   http.StatusText(status),
		Message: message,
	}
	err := WriteJSON(w, status, env, nil)
	if err != nil {
		GetLogger().Print(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func ServerErrorResponse(w http.ResponseWriter,
	r *http.Request, err error, hideError bool) {
	if hub := sentry.GetHubFromContext(r.Context()); hub != nil {
		hub.WithScope(func(scope *sentry.Scope) {
			scope.SetLevel(sentry.LevelError)
			hub.CaptureException(err)
		})
	}

	message := "the server encountered a problem and could not process your request"
	if !hideError {
		message = err.Error()
	}

	ErrorResponse(w, r, http.StatusInternalServerError, message)
}

func BadRequestResponse(w http.ResponseWriter,
	r *http.Request, err error) {
	ErrorResponse(w, r, http.StatusBadRequest, err.Error())
}

func RateLimitExceededResponse(w http.ResponseWriter,
	r *http.Request) {
	message := "rate limit exceeded"
	ErrorResponse(w, r, http.StatusTooManyRequests, message)
}

func UnauthorizedResponse(w http.ResponseWriter,
	r *http.Request, message string) {
	ErrorResponse(w, r, http.StatusUnauthorized, message)
}

func ForbiddenResponse(w http.ResponseWriter,
	r *http.Request) {
	message := "user has no access to this resource"
	ErrorResponse(w, r, http.StatusForbidden, message)
}

func ConflictResponse(
	w http.ResponseWriter,
	r *http.Request,
	err error,
	resourceName string,
	identifier string,
	identifierValue string,
	jsonField string,
	hideError bool,
) {
	value := tools.AnyToString(identifierValue)

	if err == nil || errors.Is(err, ErrRecordUniqueValue) {
		message := fmt.Sprintf(
			"%s with %s '%s' already exists",
			resourceName,
			identifier,
			value,
		)
		err := make(map[string]string)
		err[jsonField] = message
		ErrorResponse(w, r, http.StatusConflict, err)
	} else {
		ServerErrorResponse(w, r, err, hideError)
	}
}

func NotFoundResponse(
	w http.ResponseWriter,
	r *http.Request,
	err error,
	resourceName string,
	identifier string, //nolint:unparam //should keep param
	identifierValue any,
	jsonField string,
	hideError bool,
) {
	value := tools.AnyToString(identifierValue)

	if err == nil || errors.Is(err, ErrRecordNotFound) {
		message := fmt.Sprintf(
			"%s with %s '%s' doesn't exist",
			resourceName,
			identifier,
			value,
		)

		err := make(map[string]string)
		err[jsonField] = message

		ErrorResponse(w, r, http.StatusNotFound, err)
	} else {
		ServerErrorResponse(w, r, err, hideError)
	}
}

func FailedValidationResponse(
	w http.ResponseWriter,
	r *http.Request,
	errors map[string]string,
) {
	ErrorResponse(w, r, http.StatusUnprocessableEntity, errors)
}

func WSErrorResponse(err error) {
	if websocket.CloseStatus(err) == websocket.StatusNormalClosure ||
		websocket.CloseStatus(err) == websocket.StatusGoingAway {
		return
	}

	//todo
}
