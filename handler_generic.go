package essentia

import (
	"github.com/XDoubleU/essentia/pkg/parser"
	"github.com/XDoubleU/essentia/pkg/router"
)

type Generic struct {
	router.Handler
}

func (e *Engine) Generic(
	method string,
	path string,
	parser *parser.Parser,
	handlerFunc router.HandlerFunc,
) {
	e.router.Handle(method, path, func(c *router.Context) {
		if parser != nil && !parser.Parse(c) {
			//TODO: handle errors
			panic("Couldn't parse")
		}

		handlerFunc(c)
	})
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
