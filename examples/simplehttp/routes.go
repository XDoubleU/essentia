package main

import (
	"net/http"

	"github.com/justinas/alice"
	"github.com/xdoubleu/essentia/pkg/middleware"
)

func (app application) Routes() (*http.Handler, error) {
	mux := http.NewServeMux()

	app.healthRoutes(mux)

	middleware, err := middleware.Default(
		app.logger,
		app.config.AllowedOrigins,
	)
	if err != nil {
		return nil, err
	}

	standard := alice.New(middleware...)
	handler := standard.Then(mux)

	return &handler, nil
}
