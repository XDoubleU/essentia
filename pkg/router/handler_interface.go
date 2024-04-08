package router

type Handler interface {
	GetHandlerFunc()
	Handle()
}

// TODO: return response struct & handle writing in essentia
type HandlerFunc func(*Context)
