package main

import (
	"context"
	"net/http"

	wstools "github.com/XDoubleU/essentia/pkg/communication/ws"
	"github.com/XDoubleU/essentia/pkg/validate"
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
	)
	wsHandler.AddTopic(
		"topic",
		app.config.AllowedOrigins,
		func(_ context.Context, _ *wstools.Topic) (any, error) {
			return ResponseMessageDto{
				Message: "Hello, World!",
			}, nil
		})

	return wsHandler.Handler()
}
