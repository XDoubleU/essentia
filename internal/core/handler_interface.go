package core

type Handler interface {
	GetHandlerFunc()
	Handle()
}
