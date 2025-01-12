package middleware

import (
	"context"
	"github.com/google/uuid"
	logger "github.com/oprimogus/cardapiogo/pkg/log"
	"log/slog"
	"net/http"
)

func GetRequestData(ctx context.Context) *logger.RequestData {
	return ctx.Value(string(logger.RequestKey)).(*logger.RequestData)
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceID := uuid.New().String()

		reqData := &logger.RequestData{
			TraceID:  traceID,
			Method:   r.Method,
			Path:     r.URL.Path,
			ClientIP: r.RemoteAddr,
		}

		ctx := context.WithValue(r.Context(), string(logger.RequestKey), reqData)
		nr := r.WithContext(ctx)
		next.ServeHTTP(w, nr)
		slog.InfoContext(nr.Context(), "request handled")
	})
}
