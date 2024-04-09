package essentia

import (
	"net/http"

	"github.com/XDoubleU/essentia/pkg/router"
)

type Delete struct {
	Generic
}

func (e Engine) Delete(path string, handlerFunc router.HandlerFunc) {
	e.Generic(http.MethodDelete, path, nil, handlerFunc)
}
