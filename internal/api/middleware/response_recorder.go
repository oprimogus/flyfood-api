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
	Body        *bytes.Buffer // Se precisar capturar o body
}

func NewResponseRecorder(w http.ResponseWriter) *ResponseRecorder {
	return &ResponseRecorder{
		ResponseWriter: w,
		Status:         http.StatusOK,
		Headers:        make(http.Header),
		Body:           &bytes.Buffer{},
	}
}

func (r *ResponseRecorder) Write(b []byte) (int, error) {
	if !r.WroteHeader {
		r.WriteHeader(http.StatusOK)
	}

	// Captura o body se necessário
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
