package router

import (
	"fmt"
	"net/http"

	"github.com/XDoubleU/essentia/internal/logging"
)

type Router struct {
	logger     logging.Logger
	mux        *http.ServeMux
	middleware []HandlerFunc
}

func NewRouter() Router {
	return Router{
		logger:     logging.NewLogger(),
		mux:        http.NewServeMux(),
		middleware: make([]HandlerFunc, 0),
	}
}

func (r *Router) Handle(method string, path string, handler HandlerFunc) {
	var handlers []HandlerFunc
	handlers = append(handlers, r.middleware...)
	handlers = append(handlers, handler)

	pattern := fmt.Sprintf("%s %s", method, path)

	r.logger.Infof("Adding %s to mux", pattern)
	r.mux.HandleFunc(
		pattern,
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
