package essentia

import (
	"github.com/XDoubleU/essentia/pkg/parser"
	"github.com/XDoubleU/essentia/pkg/router"
)

type Generic struct {
	router.Handler
	essentia    *Essentia
	method      string
	path        string
	parser      *parser.Parser
	handlerFunc router.HandlerFunc
}

func (essentia *Essentia) Generic(
	method string,
	path string,
	handlerFunc router.HandlerFunc,
) {
	handler := Generic{
		essentia:    essentia,
		method:      method,
		path:        path,
		handlerFunc: handlerFunc,
	}
	essentia.handlers = append(essentia.handlers, handler.GetHandlerFunc())
}

func (generic *Generic) SetParser(parser *parser.Parser) {
	generic.parser = parser
}

func (generic *Generic) GetHandlerFunc() router.HandlerFunc {
	return func(c *router.Context) {
		if generic.parser != nil && !generic.parser.Parse(c) {
			//TODO: throw error or smth
		}

		generic.Handle()
		c.Next()
	}
}

func (generic Generic) Handle() {
	generic.essentia.router.Handle(generic.method, generic.path, generic.handlerFunc)
}
