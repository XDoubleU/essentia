package essentia

import (
	"net/http"

	"github.com/XDoubleU/essentia/pkg/router"
)

type Create struct {
	Generic
}

func (essentia Essentia) Create(path string, handlerFunc router.HandlerFunc) {
	essentia.Generic(http.MethodPost, path, handlerFunc)
}
