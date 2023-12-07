package main

import (
	"net/http"

	"github.com/XDoubleU/essentia"
)

func App() http.Handler {
	app := essentia.Default()

	app.SetRepository(Data{}, DataRepository{})

	app.Get("/multiple", Data{}, true)
	app.Get("/single/{id}", Data{}, false)
	app.Create("/create", nil)
	app.Update("/update/{id}", nil)
	app.Delete("/delete/{id}", nil)

	return app
}

func main() {
	http.ListenAndServe(":8000", App())
}
