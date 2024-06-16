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

type SubjectHandler = func(
	w http.ResponseWriter,
	r *http.Request,
	conn *websocket.Conn,
	msg ISubjectMessageDto,
)

type WebsocketHandler struct {
	url               string
	onCloseCallBack   OnCloseCallback
	subjectHandlerMap map[string]SubjectHandler
}

type OnCloseCallback = func(conn *websocket.Conn)

func CreateWebsocketHandler(url string) WebsocketHandler {
	if strings.Contains(url, "://") {
		url = strings.Split(url, "://")[1]
	}

	return WebsocketHandler{
		url: url,
		subjectHandlerMap: make(
			map[string]SubjectHandler,
		),
	}
}

func (h *WebsocketHandler) SetOnCloseCallback(callback OnCloseCallback) {
	h.onCloseCallBack = callback
}

func (h *WebsocketHandler) AddSubjectHandler(
	subject string,
	handler SubjectHandler,
) {
	_, ok := h.subjectHandlerMap[subject]
	if ok {
		panic("subject and handler already in map")
	}

	h.subjectHandlerMap[subject] = handler
}

func (h WebsocketHandler) GetHandler() http.HandlerFunc {
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

		var msg ISubjectMessageDto
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
