package main

import (
	"net/http"

	"github.com/XDoubleU/essentia"
)

func App() http.Handler {
	app := essentia.Default()

	var dataRepo = DataRepository{}

	app.GetPaged("/paged", essentia.GetPaged[Data, string]{
		Repo: dataRepo,
	})
	app.GetSingle("/single/{id}", essentia.GetSingle[Data, string]{
		Repo: dataRepo,
	})
	app.Create("/create", nil)
	app.Update("/update/{id}", nil)
	app.Delete("/delete/{id}", nil)

	return app
}

func main() {
	http.ListenAndServe(":8000", App())
}
