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

func (r *ResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if r.Hijacker == nil {
		return nil, nil, errors.New("http.Hijacker not implemented by underlying http.ResponseWriter")
	}
	return r.Hijacker.Hijack()
}
