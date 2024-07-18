package main

import (
	"net/http"

	"github.com/XDoubleU/essentia/pkg/httptools"
)

func (app *application) healthRoutes(mux *http.ServeMux) {
	mux.HandleFunc(
		"GET /health",
		app.getHealthHandler,
	)
}

func (app *application) getHealthHandler(w http.ResponseWriter,
	r *http.Request) {

	err := httptools.WriteJSON(w, http.StatusOK, nil, nil)
	if err != nil {
		httptools.ServerErrorResponse(w, r, err)
	}
}
