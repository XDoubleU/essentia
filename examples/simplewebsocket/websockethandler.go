package main

import (
	"net/http"

	wstools "github.com/xdoubleu/essentia/pkg/communication/ws"
	"github.com/xdoubleu/essentia/pkg/validate"
)

type SubscribeMessageDto struct {
	TopicName string `json:"topicName"`
}

type ResponseMessageDto struct {
	Message string `json:"message"`
}

func (msg SubscribeMessageDto) Validate() *validate.Validator {
	return validate.New()
}

func (msg SubscribeMessageDto) Topic() string {
	return msg.TopicName
}

func (app *application) websocketRoutes(mux *http.ServeMux) {
	mux.HandleFunc(
		"GET /",
		app.getWebSocketHandler(),
	)
}

func (app *application) getWebSocketHandler() http.HandlerFunc {

	wsHandler := wstools.CreateWebSocketHandler[SubscribeMessageDto](
		1,
		10,
		app.config.AllowedOrigins,
	)
	wsHandler.AddTopic(
		"topic",
		func(_ *wstools.Topic) any {
			return ResponseMessageDto{
				Message: "Hello, World!",
			}
		})

	return wsHandler.Handler()
}
