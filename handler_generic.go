package essentia

import "github.com/XDoubleU/essentia/internal/core"

type Generic struct {
	core.Handler
	essentia    *Essentia
	method      string
	path        string
	handlerFunc core.HandlerFunc
}

//TODO: input validation
func (essentia *Essentia) GenericHandler(method string, path string, handlerFunc core.HandlerFunc) {
	handler := Generic{
		essentia:    essentia,
		method:      method,
		path:        path,
		handlerFunc: handlerFunc,
	}
	essentia.handlers = append(essentia.handlers, handler.GetHandlerFunc())
}

func (generic Generic) GetHandlerFunc() core.HandlerFunc {
	return func(c *core.Context) {
		generic.Handle()
		c.Next()
	}
}

func (generic Generic) Handle() {
	generic.essentia.router.Handle(generic.method, generic.path, generic.handlerFunc)
}
