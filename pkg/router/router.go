package router

import (
	"fmt"
	"net/http"
)

type Router struct {
	mux        *http.ServeMux
	middleware []HandlerFunc
}

func NewRouter() *Router {
	return &Router{
		mux: http.DefaultServeMux,
	}
}

func (router Router) Handle(method string, path string, handlers ...HandlerFunc) {
	handlers = append(router.middleware, handlers...)

	router.mux.HandleFunc(
		fmt.Sprintf("%s %s", method, path),
		func(w http.ResponseWriter, r *http.Request) {
			NewContext(w, r, handlers).Next()
		},
	)
}

func (router *Router) AddMiddleware(middleware ...HandlerFunc) {
	router.middleware = append(router.middleware, middleware...)
}

func (router Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router.mux.ServeHTTP(w, r)
}
