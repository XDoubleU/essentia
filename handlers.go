package essentia

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type HandlerFunc func(*Context)

func (essentia *Essentia) Handle(method string, path string) {
	essentia.router.Handle(method, path, func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
		essentia.createContext(w, req, params).Next()
	})
}
