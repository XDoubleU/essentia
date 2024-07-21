package main

import (
	"net/http"

	"github.com/getsentry/sentry-go"
	"github.com/justinas/alice"
	"github.com/xdoubleu/essentia/pkg/middleware"
)

func (app application) Routes() (*http.Handler, error) {
	mux := http.NewServeMux()

	app.websocketRoutes(mux)

	middleware, err := middleware.DefaultWithSentry(
		app.logger,
		app.config.AllowedOrigins,
		app.config.Env,
		sentry.ClientOptions{},
	)
	if err != nil {
		return nil, err
	}

	standard := alice.New(middleware...)
	handler := standard.Then(mux)

	return &handler, nil
}
