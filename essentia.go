package essentia

import (
	"github.com/XDoubleU/essentia/internal/core"
	"github.com/XDoubleU/essentia/internal/middleware"
)

type Essentia struct {
	router   *core.Router
	handlers []core.HandlerFunc
}

func New() *Essentia {
	essentia := &Essentia{}
	essentia.router = core.NewRouter()
	return essentia
}

func (essentia *Essentia) Use(middleware ...core.HandlerFunc) {
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

func (essentia Essentia) Serve(address string) {
	//http.ListenAndServe(address, essentia)
}
