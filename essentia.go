package essentia

import (
	"github.com/julienschmidt/httprouter"
)

type Essentia struct {
	router     *httprouter.Router
	middleware []HandlerFunc
}

func New() *Essentia {
	essentia := &Essentia{}
	essentia.router = httprouter.New()
	return essentia
}

func (essentia *Essentia) Use(middleware ...HandlerFunc) {
	essentia.middleware = append(essentia.middleware, middleware...)
}

func Minimal() *Essentia {
	essentia := New()
	essentia.Use(Logger(), Recover())
	return essentia
}

func Default() *Essentia {
	essentia := Minimal()
	essentia.Use(Helmet(), Cors(), RateLimit())
	return essentia
}
