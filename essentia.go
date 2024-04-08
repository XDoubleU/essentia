package essentia

import (
	"net/http"

	"github.com/XDoubleU/essentia/pkg/middleware"
	"github.com/XDoubleU/essentia/pkg/router"
)

type Engine struct {
	router router.Router
}

func New() *Engine {
	return &Engine{
		router: router.NewRouter(),
	}
}

func (essentia *Engine) Use(middleware ...router.HandlerFunc) {
	essentia.router.AddMiddleware(middleware...)
}

func Minimal() *Engine {
	essentia := New()
	essentia.Use(middleware.Logger(), middleware.Recover())
	return essentia
}

func Default() *Engine {
	essentia := Minimal()
	essentia.Use(middleware.Helmet(), middleware.Cors(), middleware.RateLimit())
	return essentia
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	e.router.ServeHTTP(w, r)
}

func (e *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, e)
}
