package httptools

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/XDoubleU/essentia/pkg/contexttools"
	"github.com/XDoubleU/essentia/pkg/logger"
	"github.com/XDoubleU/essentia/pkg/tools"
	"github.com/getsentry/sentry-go"
	"nhooyr.io/websocket"
)

var (
	ErrRecordNotFound    = errors.New("record not found")
	ErrRecordUniqueValue = errors.New("record unique value already used")
)

//nolint:lll // can't make these lines shorter
const (
	MessageInternalServerError = "the server encountered a problem and could not process your request"
	MessageTooManyRequests     = "rate limit exceeded"
	MessageForbidden           = "user has no access to this resource"
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
		logger.GetLogger().Print(err)
	}
}

func ServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	sendErrorToSentry(r.Context(), err)

	showErrors := contexttools.GetContextValue[bool](
		r,
		contexttools.ShowErrorsContextKey,
	)

	message := MessageInternalServerError
	if showErrors != nil && *showErrors {
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
	ErrorResponse(w, r, http.StatusTooManyRequests, MessageTooManyRequests)
}

func UnauthorizedResponse(w http.ResponseWriter,
	r *http.Request, message string) {
	ErrorResponse(w, r, http.StatusUnauthorized, message)
}

func ForbiddenResponse(w http.ResponseWriter, r *http.Request) {
	ErrorResponse(w, r, http.StatusForbidden, MessageForbidden)
}

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

	if err == nil || errors.Is(err, ErrRecordUniqueValue) {
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

	if err == nil || errors.Is(err, ErrRecordNotFound) {
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

func FailedValidationResponse(
	w http.ResponseWriter,
	r *http.Request,
	errors map[string]string,
) {
	ErrorResponse(w, r, http.StatusUnprocessableEntity, errors)
}

func WSErrorResponse(
	_ http.ResponseWriter,
	r *http.Request,
	conn *websocket.Conn,
	beforeClosingFunc func(conn *websocket.Conn),
	err error,
) {
	if websocket.CloseStatus(err) == websocket.StatusNormalClosure ||
		websocket.CloseStatus(err) == websocket.StatusGoingAway ||
		websocket.CloseStatus(err) == websocket.StatusNoStatusRcvd {
		return
	}

	sendErrorToSentry(r.Context(), err)

	beforeClosingFunc(conn)

	conn.Close(websocket.StatusInternalError, MessageInternalServerError)
}

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
