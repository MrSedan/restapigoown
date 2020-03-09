package apiserver

import (
	"bufio"
	"errors"
	"net"
	"net/http"

	"github.com/felixge/httpsnoop"
)

//ResponseWriter is a new respinseWriter
type ResponseWriter struct {
	http.ResponseWriter
	code     int
	Hijacker http.Hijacker
}

// WriteHeader writing header to ResponseWriter struct
func (w *ResponseWriter) WriteHeader(statusCode int) {
	w.code = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func newResponseWriter(w http.ResponseWriter) *ResponseWriter {
	hijacker, _ := w.(http.Hijacker)
	return &ResponseWriter{
		ResponseWriter: httpsnoop.Wrap(w, httpsnoop.Hooks{}),
		Hijacker:       hijacker,
	}
}

// Hijack hiajcking http
func (w *ResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if w.Hijacker == nil {
		return nil, nil, errors.New("http.Hijacker not implemented by underlying http.ResponseWriter")
	}
	return w.Hijacker.Hijack()
}
