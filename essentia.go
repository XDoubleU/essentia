package essentia

import (
	"net/http"

	"github.com/XDoubleU/essentia/pkg/middleware"
	"github.com/XDoubleU/essentia/pkg/repositories"
	"github.com/XDoubleU/essentia/pkg/router"
)

type Essentia struct {
	router   *router.Router
	handlers []router.HandlerFunc
}

func New() *Essentia {
	essentia := &Essentia{}
	essentia.router = router.NewRouter()
	return essentia
}

func (essentia *Essentia) Use(middleware ...router.HandlerFunc) {
	essentia.router.AddMiddleware(middleware...)
}

func Minimal() *Essentia {
	essentia := New()
	essentia.Use(middleware.Logger(), middleware.Recover())
	return essentia
}

func Default() *Essentia {
	essentia := Minimal()
	essentia.Use(middleware.Helmet(), middleware.Cors(), middleware.RateLimit())
	return essentia
}

func (essentia *Essentia) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	essentia.router.ServeHTTP(w, r)
}

func SetRepository[TData any, TId any](
	e *Essentia,
	repo repositories.Repository[TData, TId],
) {
	router.SetRepository[TData, TId](e.router, repo)
}
