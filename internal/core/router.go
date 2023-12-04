package core

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Router struct {
	router     *httprouter.Router
	middleware []HandlerFunc
}

func NewRouter() *Router {
	return &Router{
		router: httprouter.New(),
	}
}

func (router Router) Handle(method string, path string, handlers ...HandlerFunc) {
	handlers = append(router.middleware, handlers...)

	router.router.Handle(method, path, func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
		NewContext(w, req, params, handlers).Next()
	})
}

func (router *Router) AddMiddleware(middleware ...HandlerFunc) {
	router.middleware = append(router.middleware, middleware...)
}
