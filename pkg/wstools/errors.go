package wstools

import (
	"context"
	"net/http"

	"github.com/xdoubleu/essentia/pkg/contexttools"
	"github.com/xdoubleu/essentia/pkg/httptools"
	"github.com/xdoubleu/essentia/pkg/sentrytools"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

// ErrorResponse is used to handle any kind of error that occured on a WebSocket.
func ErrorResponse(
	ctx context.Context,
	conn *websocket.Conn,
	status int,
	message any,
) {
	errorDto := httptools.ErrorDto{
		Status:  status,
		Error:   http.StatusText(status),
		Message: message,
	}
	err := wsjson.Write(ctx, conn, errorDto)
	if err != nil {
		sentrytools.SendErrorToSentry(ctx, err)
		contexttools.Logger(ctx).Print(err)
	}
}

// ServerErrorResponse is used to handle
// internal server errors that occured on a WebSocket.
func ServerErrorResponse(ctx context.Context, conn *websocket.Conn, err error) {
	sentrytools.SendErrorToSentry(ctx, err)

	message := httptools.MessageInternalServerError
	if contexttools.ShowErrors(ctx) {
		message = err.Error()
	}

	ErrorResponse(ctx, conn, http.StatusInternalServerError, message)
}

// UpgradeErrorResponse is used to handle an error that
// occured during the upgrade towards a WebSocket.
func UpgradeErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	sentrytools.SendErrorToSentry(r.Context(), err)
	w.WriteHeader(http.StatusInternalServerError)
}

// FailedValidationResponse is used to handle
// an error of a [validate.Validator] that occured on a WebSocket.
func FailedValidationResponse(
	ctx context.Context,
	conn *websocket.Conn,
	errors map[string]string,
) {
	ErrorResponse(ctx, conn, http.StatusUnprocessableEntity, errors)
}
