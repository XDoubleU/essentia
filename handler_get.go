package essentia

import (
	"net/http"

	"github.com/XDoubleU/essentia/pkg/router"
)

type Get struct {
	Generic
}

// TODO: pagination
func (essentia *Essentia) Get(path string, dataType any, hasPagination bool) {
	essentia.Generic(http.MethodGet, path, func(ctx *router.Context) {
		r := ctx.GetRepository(dataType)
		//todo: parse query
		r.GetPaged(-1, -1)
	})
}
