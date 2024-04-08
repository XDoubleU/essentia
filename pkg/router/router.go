package router

import (
	"fmt"
	"net/http"
)

type Router struct {
	mux        *http.ServeMux
	middleware []HandlerFunc
}

func NewRouter() Router {
	return Router{
		mux:        http.DefaultServeMux,
		middleware: make([]HandlerFunc, 0),
	}
}

func (router *Router) Handle(method string, path string, handler HandlerFunc) {
	var handlers []HandlerFunc
	handlers = append(handlers, router.middleware...)
	handlers = append(handlers, handler)

	router.mux.HandleFunc(
		fmt.Sprintf("%s %s", method, path),
		func(w http.ResponseWriter, r *http.Request) {
			c := NewContext(w, r, handlers)
			if err := c.readBody(); err != nil {
				// todo handle error
				panic(err)
			}

			c.Next()
		},
	)
}

func (router *Router) AddMiddleware(middleware ...HandlerFunc) {
	router.middleware = append(router.middleware, middleware...)
}

func (router Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router.mux.ServeHTTP(w, r)
}
