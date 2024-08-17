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

type Health struct {
	IsDatabaseActive bool
}

func (app *application) getHealthHandler(w http.ResponseWriter,
	r *http.Request) {

	data := Health{
		IsDatabaseActive: app.db.Ping(r.Context()) == nil,
	}

	err := httptools.WriteJSON(w, http.StatusOK, data, nil)
	if err != nil {
		httptools.ServerErrorResponse(w, r, err)
	}
}
