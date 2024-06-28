package httptools

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/XDoubleU/essentia/pkg/validate"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type ISubjectMessageDto interface {
	validate.IValidatedType
	GetSubject() string
}

type WebsocketHandler[T ISubjectMessageDto] struct {
	url               string
	onCloseCallBack   OnCloseCallback
	subjectHandlerMap map[string]func(
		w http.ResponseWriter,
		r *http.Request,
		conn *websocket.Conn,
		msg T,
	)
}

type OnCloseCallback = func(conn *websocket.Conn)

func CreateWebsocketHandler[T ISubjectMessageDto](url string) WebsocketHandler[T] {
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
	}
}

func (h *WebsocketHandler[T]) SetOnCloseCallback(callback OnCloseCallback) {
	h.onCloseCallBack = callback
}

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

func (h WebsocketHandler[T]) GetHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
			OriginPatterns: []string{h.url},
		})
		if err != nil {
			WSUpgradeErrorResponse(w, r, err)
			return
		}

		defer func() {
			h.onCloseCallBack(conn)
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
