package essentia

import (
	"net/http"

	"github.com/XDoubleU/essentia/internal/helpers"
	"github.com/XDoubleU/essentia/pkg/router"
)

type GetPagedRepository[TData any, TId any] interface {
	GetPaged(pageIndex int, pageSize int) []TData
}

type getPaged interface {
	GetPaged(pageIndex int, pageSize int) []any
}

type GetPaged[TData any, TId any] struct {
	Generic
	Repo GetPagedRepository[TData, TId]
}

func (g GetPaged[TData, TId]) GetPaged(pageIndex int, pageSize int) []any {
	return helpers.CastToAnyArray(g.Repo.GetPaged(pageIndex, pageSize))
}

func (essentia *Essentia) GetPaged(path string, g getPaged) {
	//TODO configure validator
	essentia.Generic(http.MethodGet, path, func(ctx *router.Context) {
		// TODO do something with data
		// TODO parse and use pageIndex and pageSize
		//data := g.GetPaged(-1, -1)
	})
}
