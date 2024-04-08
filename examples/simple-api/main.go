package main

import (
	"net/http"

	"github.com/XDoubleU/essentia"
	"github.com/XDoubleU/essentia/pkg/router"
)

func setupRouter(dataRepo *DataRepository) *essentia.Engine {
	r := essentia.Default()

	r.Generic(http.MethodGet, "/generic", func(ctx *router.Context) {
		rsData := map[string]string{
			"message": "ok",
		}
		ctx.Writer.WriteJSON(http.StatusOK, rsData, nil)
	})
	r.GetSingle("/single/{id}", essentia.GetSingle[Data, string]{
		Repo: dataRepo,
	})
	r.GetPaged("/paged", essentia.GetPaged[Data, string]{
		Repo: dataRepo,
	})
	//r.Create("/create", nil)
	//r.Update("/update/{id}", nil)
	//r.Delete("/delete/{id}", nil)

	return r
}

func main() {
	dataRepo := NewDataRepository()
	r := setupRouter(&dataRepo)
	r.Run(":8000")
}
