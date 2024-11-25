package middleware

import (
	"context"
	"github.com/google/uuid"
	"net/http"
)

type ContextKey string

const (
	RequestKey ContextKey = "request_data"
)

type RequestData struct {
	TraceID  string `json:"trace_id"`
	Method   string `json:"method"`
	Path     string `json:"path"`
	ClientIP string `json:"client_ip"`
}

func GetRequestData(ctx context.Context) *RequestData {
	return ctx.Value(string(RequestKey)).(*RequestData)
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceID := uuid.New().String()
		reqData := &RequestData{
			TraceID:  traceID,
			Method:   r.Method,
			Path:     r.URL.Path,
			ClientIP: r.RemoteAddr,
		}
		ctx := context.WithValue(r.Context(), string(RequestKey), reqData)
		nr := r.WithContext(ctx)
		next.ServeHTTP(w, nr)
	})
}
