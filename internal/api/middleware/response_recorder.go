package middleware

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"net/http"
)

type ResponseRecorder struct {
	http.ResponseWriter
	Status      int
	Written     int64
	WroteHeader bool
	Headers     http.Header
	Body        *bytes.Buffer
}

func NewResponseRecorder(w http.ResponseWriter) *ResponseRecorder {
	// Copie os headers existentes do ResponseWriter original
	return &ResponseRecorder{
		ResponseWriter: w,
		Status:         http.StatusOK, // Status padrão é 200 OK
		Headers:        w.Header(),    // Use os headers do writer original
		Body:           &bytes.Buffer{},
	}
}

func (r *ResponseRecorder) WriteHeader(status int) {
	if !r.WroteHeader {
		r.Status = status
		r.WroteHeader = true
		r.ResponseWriter.WriteHeader(status)
	}
}

func (r *ResponseRecorder) Write(b []byte) (int, error) {
	if !r.WroteHeader {
		r.WriteHeader(r.Status)
	}

	r.Body.Write(b)
	n, err := r.ResponseWriter.Write(b)
	r.Written += int64(n)
	return n, err
}

func (r *ResponseRecorder) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hijacker, ok := r.ResponseWriter.(http.Hijacker); ok {
		return hijacker.Hijack()
	}
	return nil, nil, fmt.Errorf("hijacking not supported")
}

func (r *ResponseRecorder) Flush() {
	if flusher, ok := r.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

// Push implementa http.Pusher
func (r *ResponseRecorder) Push(target string, opts *http.PushOptions) error {
	if pusher, ok := r.ResponseWriter.(http.Pusher); ok {
		return pusher.Push(target, opts)
	}
	return http.ErrNotSupported
}
