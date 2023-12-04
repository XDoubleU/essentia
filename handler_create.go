package essentia

import "github.com/XDoubleU/essentia/internal/core"

type Create struct {
	Generic
}

func (essentia Essentia) CreateHandler(path string, handlerFunc core.HandlerFunc) {
	essentia.GenericHandler("POST", path, handlerFunc)
}
