package httptools

import (
	"bufio"
	"errors"
	"net"
	"net/http"
)

type ResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{w, -1}
}

func (w *ResponseWriter) WriteHeader(statusCode int) {
	if w.statusCode != -1 {
		return
	}

	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w ResponseWriter) StatusCode() int {
	return w.statusCode
}

func (w *ResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h, ok := w.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, errors.New("hijack not supported")
	}
	return h.Hijack()
}
