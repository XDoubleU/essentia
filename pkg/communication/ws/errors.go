package ws

import (
	"context"
	"net/http"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	contexttools "github.com/xdoubleu/essentia/pkg/context"
	errortools "github.com/xdoubleu/essentia/pkg/errors"
	"github.com/xdoubleu/essentia/pkg/logging"
)

// ErrorResponse is used to handle any kind of error that occured on a WebSocket.
func ErrorResponse(
	ctx context.Context,
	conn *websocket.Conn,
	status int,
	message any,
) {
	errorDto := errortools.NewErrorDto(status, message)
	err := wsjson.Write(ctx, conn, errorDto)
	if err != nil {
		contexttools.Logger(ctx).
			ErrorContext(ctx, "failed to write JSON", logging.ErrAttr(err))
	}
}

// ServerErrorResponse is used to handle
// internal server errors that occured on a WebSocket.
func ServerErrorResponse(ctx context.Context, conn *websocket.Conn, err error) {
	contexttools.Logger(ctx).
		ErrorContext(ctx, "server error occurred", logging.ErrAttr(err))

	message := errortools.MessageInternalServerError
	if contexttools.ShowErrors(ctx) {
		message = err.Error()
	}

	ErrorResponse(ctx, conn, http.StatusInternalServerError, message)
}

// UpgradeErrorResponse is used to handle an error that
// occured during the upgrade towards a WebSocket.
func UpgradeErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	contexttools.Logger(r.Context()).
		ErrorContext(r.Context(), "WS upgrade error occurred", logging.ErrAttr(err))
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

// ForbiddenResponse is used to handle an error when a user
// isn't authorized to access a certain resource.
func ForbiddenResponse(ctx context.Context, conn *websocket.Conn) {
	ErrorResponse(ctx, conn, http.StatusForbidden, errortools.MessageForbidden)
}
