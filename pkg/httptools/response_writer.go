package httptools

import (
	"bufio"
	"errors"
	"net"
	"net/http"
)

// A ResponseWriter is used to capture set status codes.
type ResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

// NewResponseWriter returns a new [ResponseWriter].
func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{w, -1}
}

// WriteHeader sets the internal statusCode value of a [ResponseWriter].
func (w *ResponseWriter) WriteHeader(statusCode int) {
	if w.statusCode != -1 {
		return
	}

	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

// StatusCode returns the status code of a [ResponseWriter].
func (w ResponseWriter) StatusCode() int {
	return w.statusCode
}

// Hijack lets the caller take over the connection.
// After a call to Hijack the HTTP server library
// will not do anything else with the connection.
func (w *ResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h, ok := w.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, errors.New("hijack not supported")
	}
	return h.Hijack()
}
