package essentia

import (
	"github.com/XDoubleU/essentia/pkg/input"
	"github.com/XDoubleU/essentia/pkg/router"
)

type Generic struct {
	router.Handler
	essentia    *Essentia
	method      string
	path        string
	validator   *input.Validator
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

func (generic *Generic) SetValidator(validator *input.Validator) {
	generic.validator = validator
}

func (generic *Generic) GetHandlerFunc() router.HandlerFunc {
	return func(c *router.Context) {
		if generic.validator != nil && !generic.validator.Validate(c) {
			//TODO: throw error or smth
		}

		generic.Handle()
		c.Next()
	}
}

func (generic Generic) Handle() {
	generic.essentia.router.Handle(generic.method, generic.path, generic.handlerFunc)
}
