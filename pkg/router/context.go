package router

import (
	"log"
	"net/http"

	"github.com/XDoubleU/essentia/internal/helpers"
)

type Middleware struct {
	index    int8
	handlers []HandlerFunc
}

type Context struct {
	Request     *http.Request
	Writer      *ResponseWriter
	Middleware  *Middleware
	QueryValues map[string]any
	PathValues  map[string]any
	BodyValues  map[string]any
	body        map[string]string
	data        map[string]any
}

func NewContext(
	w http.ResponseWriter,
	req *http.Request,
	handlers []HandlerFunc,
) Context {
	return Context{
		Writer:  &ResponseWriter{w, 0},
		Request: req,
		Middleware: &Middleware{
			index:    -1,
			handlers: handlers,
		},
		body: make(map[string]string),
		data: make(map[string]any),
	}
}

func (c *Context) readBody() error {
	var body map[string]string
	if err := helpers.ReadJSON(c.Request.Body, &body, true); err != nil {
		return err
	}

	c.body = body

	return nil
}

func (c *Context) Next() {
	c.Middleware.index++

	for ; c.Middleware.index < int8(len(c.Middleware.handlers)); c.Middleware.index++ {
		c.Middleware.handlers[c.Middleware.index](c)
	}
}

func (c *Context) SetData(key string, value any) {
	c.data[key] = value
}

func (c Context) GetData(key string) any {
	value, ok := c.data[key]

	if !ok {
		log.Panicf("No key %s in context data", key)
	}

	return value
}

func (c Context) GetRawQueryValue(name string) []string {
	value, ok := c.Request.URL.Query()[name]

	if !ok {
		log.Panicf("No param %s in query params", name)
	}

	return value
}

func (c Context) GetRawPathValue(name string) string {
	value := c.Request.PathValue(name)

	if len(value) == 0 {
		log.Panicf("No param %s in url params", name)
	}

	return value
}

func (c Context) GetRawBodyValue(name string) string {
	value, ok := c.body[name]

	if !ok {
		log.Panicf("No param %s in query params", name)
	}

	return value
}
