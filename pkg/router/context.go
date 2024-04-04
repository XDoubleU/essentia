package router

import (
	"encoding/json"
	"log"
	"net/http"
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
	r *http.Request,
	middleware []HandlerFunc,
) *Context {
	var body map[string]string
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		//TODO: handle error
	}

	return &Context{
		Writer:  &ResponseWriter{w, 0},
		Request: r,
		Middleware: &Middleware{
			index:    -1,
			handlers: middleware,
		},
		body: body,
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

	if !ok {
		log.Panicf("No key %s in context data", key)
	}

	return value
}

func (c Context) GetQueryValue(name string) []string {
	value, ok := c.Request.URL.Query()[name]

	if !ok {
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

func (c Context) GetBodyValue(name string) string {
	value, ok := c.body[name]

	if !ok {
		log.Panicf("No param %s in query params", name)
	}

	return value
}
