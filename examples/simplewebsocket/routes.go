package main

import (
	"net/http"

	"github.com/XDoubleU/essentia/pkg/middleware"
	"github.com/getsentry/sentry-go"
	"github.com/justinas/alice"
)

func (app application) Routes() http.Handler {
	mux := http.NewServeMux()

	app.websocketRoutes(mux)

	middleware, err := middleware.DefaultWithSentry(
		app.logger,
		app.config.AllowedOrigins,
		app.config.Env,
		sentry.ClientOptions{},
	)
	if err != nil {
		panic(err)
	}

	standard := alice.New(middleware...)
	return standard.Then(mux)
}
