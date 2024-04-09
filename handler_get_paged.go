package essentia

import (
	"net/http"

	"github.com/XDoubleU/essentia/internal/helpers"
	"github.com/XDoubleU/essentia/pkg/router"
)

type GenericPagedGetter[TData any, TId any] interface {
	PagedGet(pageIndex int, pageSize int) []TData
}

type PagedGetter interface {
	PagedGet(pageIndex int, pageSize int) []any
}

type GetPaged[TData any, TId any] struct {
	Generic
	Repo GenericPagedGetter[TData, TId]
}

func (g GetPaged[TData, TId]) PagedGet(pageIndex int, pageSize int) []any {
	return helpers.CastToAnyArray(g.Repo.PagedGet(pageIndex, pageSize))
}

func (essentia *Engine) GetPaged(path string, g PagedGetter) {
	// TODO configure validator
	essentia.Generic(http.MethodGet, path, nil, func(ctx *router.Context) {
		// TODO do something with data
		// TODO parse and use pageIndex and pageSize
		// data := g.GetPaged(-1, -1)
	})
}
