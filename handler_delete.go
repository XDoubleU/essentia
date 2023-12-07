package essentia

import (
	"net/http"

	"github.com/XDoubleU/essentia/pkg/router"
)

type Delete struct {
	Generic
}

func (essentia Essentia) Delete(path string, handlerFunc router.HandlerFunc) {
	essentia.Generic(http.MethodDelete, path, handlerFunc)
}
