package main

import (
	"net/http"

	"github.com/XDoubleU/essentia/pkg/middleware"
	"github.com/justinas/alice"
)

func (app application) Routes() (*http.Handler, error) {
	mux := http.NewServeMux()

	app.healthRoutes(mux)

	middleware, err := middleware.Default(
		app.logger,
		app.config.Env == TestEnv,
		app.config.AllowedOrigins,
		nil,
	)
	if err != nil {
		return nil, err
	}

	standard := alice.New(middleware...)
	handler := standard.Then(mux)

	return &handler, nil
}
