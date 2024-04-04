package essentia

import (
	"net/http"

	"github.com/XDoubleU/essentia/pkg/router"
)

type Update struct {
	Generic
}

func (essentia Essentia) Update(path string, handlerFunc router.HandlerFunc) {
	essentia.Generic(http.MethodPatch, path, handlerFunc)
}