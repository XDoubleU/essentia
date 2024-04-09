package essentia

import (
	"net/http"

	"github.com/XDoubleU/essentia/pkg/router"
)

type Create struct {
	Generic
}

func (e Engine) Create(path string, handlerFunc router.HandlerFunc) {
	e.Generic(http.MethodPost, path, nil, handlerFunc)
}
