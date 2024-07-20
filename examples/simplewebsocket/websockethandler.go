package main

import (
	"net/http"

	"github.com/XDoubleU/essentia/pkg/httptools"
	"github.com/XDoubleU/essentia/pkg/validate"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type SubscribeMessageDto struct {
	Topic string
}

type ResponseMessageDto struct {
	Message string
}

func (msg SubscribeMessageDto) Validate() *validate.Validator {
	return validate.New()
}

func (msg SubscribeMessageDto) GetTopic() string {
	return msg.Topic
}

func (app *application) websocketRoutes(mux *http.ServeMux) {
	mux.HandleFunc(
		"GET /",
		app.getWebSocketHandler(),
	)
}

func (app *application) getWebSocketHandler() http.HandlerFunc {

	wsHandler := httptools.CreateWebsocketHandler[SubscribeMessageDto](
		app.config.AllowedOrigins,
	)
	wsHandler.AddTopicHandler("topic", topicHandler)

	return wsHandler.Handler()
}

func topicHandler(
	w http.ResponseWriter,
	r *http.Request,
	conn *websocket.Conn,
	msg SubscribeMessageDto,
) {
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
