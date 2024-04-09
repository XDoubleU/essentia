package essentia

import (
	"fmt"
	"net/http"

	"github.com/XDoubleU/essentia/pkg/router"
)

type GenericSingleGetter[TData any, TId any] interface {
	SingleGet(id TId) *TData
}

type SingleGetter interface {
	SingleGet(id any) any
}

type GetSingle[TData any, TId any] struct {
	Generic
	Repo GenericSingleGetter[TData, TId]
}

func (g GetSingle[TData, TId]) SingleGet(id any) any {
	v, ok := id.(TId)
	if !ok {
		// TODO error
		return nil
	}

	return g.Repo.SingleGet(v)
}

func (e *Engine) GetSingle(path string, g SingleGetter) {
	// TODO configure validator

	e.Generic(http.MethodGet, path, nil, func(c *router.Context) {
		id, ok := c.PathValues["id"]
		if !ok {
			//todo: handle error
			return
		}

		fmt.Printf("id: %s\n", id)

		// TODO do something with data
		// data := g.GetSingle(id)
	})
}
