package essentia

import "github.com/XDoubleU/essentia/internal/core"

type GetSingle struct {
	Generic
}

func (essentia Essentia) GetSingleHandler(path string, handlerFunc core.HandlerFunc) {
	essentia.GenericHandler("GET", path, handlerFunc)
}
