package router

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/XDoubleU/essentia/pkg/repositories"
)

type Router struct {
	mux          *http.ServeMux
	middleware   []HandlerFunc
	repositories map[string]repositories.Repository[any, any]
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
			NewContext(w, r, handlers, router.repositories).Next()
		},
	)
}

func (router *Router) AddMiddleware(middleware ...HandlerFunc) {
	router.middleware = append(router.middleware, middleware...)
}

func (router Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router.mux.ServeHTTP(w, r)
}

func SetRepository[TData any, TId any](
	r *Router,
	repo repositories.Repository[TData, TId],
) {
	r.repositories[reflect.TypeFor[TData]().String()] = repo.(repositories.Repository[any, any])
}
