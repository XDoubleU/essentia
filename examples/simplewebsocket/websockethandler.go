package main

import (
	"net/http"

	"github.com/xdoubleu/essentia/pkg/validate"
	"github.com/xdoubleu/essentia/pkg/wstools"
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
		ResponseMessageDto{
			Message: "Hello, World!",
		})

	return wsHandler.Handler()
}
