package main

import (
	"net/http"

	httptools "github.com/xdoubleu/essentia/pkg/communication/http"
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
