package router

type Handler interface {
	GetHandlerFunc()
	Handle()
}

type HandlerFunc func(*Context)
