package core

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type HandlerFunc func(*Context)

type Middleware struct {
	index    int8
	handlers []HandlerFunc
}

type Context struct {
	Request    *http.Request
	Writer     *ResponseWriter
	Middleware *Middleware
	Params     httprouter.Params
	Data       map[string]interface{}
}

func NewContext(w http.ResponseWriter, req *http.Request, params httprouter.Params, middleware []HandlerFunc) *Context {
	return &Context{
		Writer:  &ResponseWriter{w, 0},
		Request: req,
		Middleware: &Middleware{
			index:    -1,
			handlers: middleware,
		},
		Params: params,
		Data:   make(map[string]interface{}),
	}
}

func (c *Context) Next() {
	c.Middleware.index++
	s := int8(len(c.Middleware.handlers))
	for ; c.Middleware.index < s; c.Middleware.index++ {
		c.Middleware.handlers[c.Middleware.index](c)
	}
}

func (c *Context) Set(key string, value interface{}) {
	c.Data[key] = value
}

func (c *Context) Get(key string) interface{} {
	value, ok := c.Data[key]

	if !ok || value == nil {
		log.Panicf("No key %s in context data", key)
	}

	return value
}
