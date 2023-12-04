package essentia

import "github.com/XDoubleU/essentia/internal/core"

type Update struct {
	Generic
}

func (essentia Essentia) UpdateHandler(path string, handlerFunc core.HandlerFunc) {
	essentia.GenericHandler("PATCH", path, handlerFunc)
}
