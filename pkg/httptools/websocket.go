package httptools

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/XDoubleU/essentia/pkg/validate"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

// SubjectMessageDto is implemented by all messages
// used to subscribe to a certain handler of a [WebsocketHandler].
type SubjectMessageDto interface {
	validate.ValidatedType
	GetSubject() string
}

// A WebsocketHandler handles incoming requests to a
// websocket and makes sure that the right handler is called for each subject.
type WebsocketHandler[T SubjectMessageDto] struct {
	url               string
	onCloseCallBack   OnCloseCallback
	subjectHandlerMap map[string]func(
		w http.ResponseWriter,
		r *http.Request,
		conn *websocket.Conn,
		msg T,
	)
}

// OnCloseCallback is called when the websocket is closed.
type OnCloseCallback = func(conn *websocket.Conn)

// CreateWebsocketHandler creates a new [WebsocketHandler].
func CreateWebsocketHandler[T SubjectMessageDto](url string) WebsocketHandler[T] {
	if strings.Contains(url, "://") {
		url = strings.Split(url, "://")[1]
	}

	return WebsocketHandler[T]{
		url: url,
		subjectHandlerMap: make(
			map[string]func(
				w http.ResponseWriter,
				r *http.Request,
				conn *websocket.Conn,
				msg T,
			),
		),
		onCloseCallBack: nil,
	}
}

// SetOnCloseCallback sets the function to call when closing a [websocket.Conn].
func (h *WebsocketHandler[T]) SetOnCloseCallback(callback OnCloseCallback) {
	h.onCloseCallBack = callback
}

// AddSubjectHandler adds a handler for a
// specific subject provided by a [SubjectMessageDto].
func (h *WebsocketHandler[T]) AddSubjectHandler(
	subject string,
	handler func(
		w http.ResponseWriter,
		r *http.Request,
		conn *websocket.Conn,
		msg T,
	),
) {
	_, ok := h.subjectHandlerMap[subject]
	if ok {
		panic("subject and handler already in map")
	}

	h.subjectHandlerMap[subject] = handler
}

// GetHandler returns the [http.HandlerFunc] of a [WebsocketHandler].
func (h WebsocketHandler[T]) GetHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//nolint:exhaustruct //other fields are optional
		conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
			OriginPatterns: []string{h.url},
		})
		if err != nil {
			WSUpgradeErrorResponse(w, r, err)
			return
		}

		defer func() {
			if h.onCloseCallBack != nil {
				h.onCloseCallBack(conn)
			}

			conn.Close(websocket.StatusInternalError, "")
		}()

		var msg T
		err = wsjson.Read(r.Context(), conn, &msg)
		if err != nil {
			WSErrorResponse(w, r, conn, h.onCloseCallBack, err)
			return
		}

		if v := msg.Validate(); !v.Valid() {
			FailedValidationResponse(w, r, v.Errors)
			return
		}

		handler, ok := h.subjectHandlerMap[msg.GetSubject()]
		if !ok {
			errorDto := ErrorDto{
				Status:  0,
				Error:   "unknown subject",
				Message: fmt.Sprintf("no handler found for '%s'", msg.GetSubject()),
			}
			err = wsjson.Write(r.Context(), conn, errorDto)
			if err != nil {
				WSErrorResponse(w, r, conn, h.onCloseCallBack, err)
				return
			}
			return
		}

		handler(w, r, conn, msg)
	}
}
