package essentia

import "github.com/XDoubleU/essentia/internal/core"

type Delete struct {
	Generic
}

func (essentia Essentia) DeleteHandler(path string, handlerFunc core.HandlerFunc) {
	essentia.GenericHandler("DELETE", path, handlerFunc)
}
