package essentia

import (
	"github.com/XDoubleU/essentia/pkg/router"
)

type Generic struct {
	router.Handler
}

func (e *Engine) Generic(
	method string,
	path string,
	handlerFunc router.HandlerFunc,
) {
	e.router.Handle(method, path, handlerFunc)
}

/*
func (generic *Generic) GetHandlerFunc() router.HandlerFunc {
	return func(c *router.Context) {
		if generic.parser != nil && !generic.parser.Parse(c) {
			//TODO: throw error or smth
		}

		generic.Handle()
		c.Next()
	}
}
*/
