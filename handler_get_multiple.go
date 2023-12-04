package essentia

import "github.com/XDoubleU/essentia/internal/core"

type GetMultiple struct {
	Generic
}

//TODO: pagination
func (essentia *Essentia) GetMultipleHandler(path string, handlerFunc core.HandlerFunc) {
	essentia.GenericHandler("GET", path, handlerFunc)
}
