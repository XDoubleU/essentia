package router

import (
	"log"
	"net/http"
)

type Middleware struct {
	index    int8
	handlers []HandlerFunc
}

type Context struct {
	Request    *http.Request
	Writer     *ResponseWriter
	Middleware *Middleware
	data       map[string]any
}

func NewContext(
	w http.ResponseWriter,
	r *http.Request,
	middleware []HandlerFunc,
) *Context {
	return &Context{
		Writer:  &ResponseWriter{w, 0},
		Request: r,
		Middleware: &Middleware{
			index:    -1,
			handlers: middleware,
		},
		data: make(map[string]any),
	}
}

func (c *Context) Next() {
	c.Middleware.index++
	s := int8(len(c.Middleware.handlers))
	for ; c.Middleware.index < s; c.Middleware.index++ {
		c.Middleware.handlers[c.Middleware.index](c)
	}
}

func (c *Context) SetData(key string, value any) {
	c.data[key] = value
}

func (c Context) GetData(key string) any {
	value, ok := c.data[key]

	if !ok || value == nil {
		log.Panicf("No key %s in context data", key)
	}

	return value
}

func (c Context) GetQueryValue(name string) []string {
	value, ok := c.Request.URL.Query()[name]

	if !ok || value == nil {
		log.Panicf("No param %s in query params", name)
	}

	return value
}

func (c Context) GetPathValue(name string) string {
	value := c.Request.PathValue(name)

	if len(value) == 0 {
		log.Panicf("No param %s in url params", name)
	}

	return value
}
