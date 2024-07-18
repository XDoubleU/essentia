package main

import (
	"net/http"

	"github.com/XDoubleU/essentia/pkg/httptools"
	"github.com/XDoubleU/essentia/pkg/validate"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type SubjectMessageDto struct {
	Subject string
}

type ResponseMessageDto struct {
	Message string
}

func (msg SubjectMessageDto) Validate() *validate.Validator {
	return validate.New()
}

func (msg SubjectMessageDto) GetSubject() string {
	return msg.Subject
}

func (app *application) websocketRoutes(mux *http.ServeMux) {
	mux.HandleFunc(
		"GET /",
		app.getWebSocketHandler(),
	)
}

func (app *application) getWebSocketHandler() http.HandlerFunc {

	wsHandler := httptools.CreateWebsocketHandler[SubjectMessageDto](app.config.AllowedOrigins)
	wsHandler.AddSubjectHandler("subject", subjectHandler)

	return wsHandler.GetHandler()
}

func subjectHandler(w http.ResponseWriter, r *http.Request, conn *websocket.Conn, msg SubjectMessageDto) {
	err := wsjson.Write(r.Context(), conn, ResponseMessageDto{
		Message: "Hello, World!",
	})
	if err != nil {
		httptools.WSErrorResponse(
			w,
			r,
			conn,
			nil,
			err,
		)
		return
	}
}
