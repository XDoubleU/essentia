package essentia

import (
	"net/http"

	"github.com/XDoubleU/essentia/pkg/router"
)

type GetSingleRepository[TData any, TId any] interface {
	GetSingle(id TId) *TData
}

type getSingle interface {
	GetSingle(id any) any
}

type GetSingle[TData any, TId any] struct {
	Generic
	Repo GetSingleRepository[TData, TId]
}

func (g GetSingle[TData, TId]) GetSingle(id any) any {
	v, ok := id.(TId)
	if !ok {
		//TODO error
		return nil
	}

	return g.Repo.GetSingle(v)
}

func (essentia *Essentia) GetSingle(path string, g getSingle) {
	//TODO configure validator
	essentia.Generic(http.MethodGet, path, func(ctx *router.Context) {
		// TODO do something with data
		//data := g.GetSingle(-1)
	})
}
