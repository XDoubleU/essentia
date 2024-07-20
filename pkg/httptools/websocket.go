package httptools

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/XDoubleU/essentia/pkg/validate"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

// SubscribeMessageDto is implemented by all messages
// used to subscribe to a certain handler of a [WebsocketHandler].
type SubscribeMessageDto interface {
	validate.ValidatedType
	GetTopic() string
}

// A WebsocketHandler handles incoming requests to a
// websocket and makes sure that the right handler is called for each subject.
type WebsocketHandler[T SubscribeMessageDto] struct {
	allowedOrigins    []string
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
func CreateWebsocketHandler[T SubscribeMessageDto](allowedOrigins []string) WebsocketHandler[T] {
	for i, url := range allowedOrigins {
		if strings.Contains(url, "://") {
			allowedOrigins[i] = strings.Split(url, "://")[1]
		}
	}

	return WebsocketHandler[T]{
		allowedOrigins: allowedOrigins,
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

// AddTopicHandler adds a handler for a
// specific topic provided by a [SubscribeMessageDto].
func (h *WebsocketHandler[T]) AddTopicHandler(
	topic string,
	handler func(
		w http.ResponseWriter,
		r *http.Request,
		conn *websocket.Conn,
		msg T,
	),
) {
	_, ok := h.topicHandlerMap[topic]
	if ok {
		panic(fmt.Sprintf("topic '%s' already has a handler", topic))
	}

	h.topicHandlerMap[topic] = handler
}

// Handler returns the [http.HandlerFunc] of a [WebsocketHandler].
func (h WebsocketHandler[T]) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//nolint:exhaustruct //other fields are optional
		conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
			OriginPatterns: h.allowedOrigins,
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

		handler, ok := h.topicHandlerMap[msg.GetTopic()]
		if !ok {
			errorDto := ErrorDto{
				Status:  0,
				Error:   "unknown topic",
				Message: fmt.Sprintf("no handler found for '%s'", msg.GetTopic()),
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
