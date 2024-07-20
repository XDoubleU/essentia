package main

import (
	"net/http"

	"github.com/XDoubleU/essentia/pkg/validate"
	"github.com/XDoubleU/essentia/pkg/wstools"
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

func (msg SubscribeMessageDto) GetTopicName() string {
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
		app.config.AllowedOrigins,
	)
	wsHandler.AddTopic(
		"topic",
		ResponseMessageDto{
			Message: "Hello, World!",
		})

	return wsHandler.Handler()
}
