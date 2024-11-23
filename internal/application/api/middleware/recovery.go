package middleware

import (
	"encoding/json"
	"github.com/oprimogus/cardapiogo/internal/xerrors"
	"log/slog"
	"net/http"
	"runtime"
)

func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				stack := make([]byte, 4<<10) // 4KB
				length := runtime.Stack(stack, false)

				defer func() {
					slog.ErrorContext(r.Context(), "PANIC RECOVERED",
						"error", err,
						"stack", string(stack[:length]),
						"url", r.URL.String(),
						"method", r.Method,
						"client_ip", r.RemoteAddr,
					)
				}()

				if !isHeaderWritten(w) {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)
					reqData := r.Context().Value(string(RequestKey)).(*RequestData)

					_ = json.NewEncoder(w).Encode(xerrors.InternalServer(reqData.TraceID, ""))
				}
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func isHeaderWritten(w http.ResponseWriter) bool {
	if rw, ok := w.(*ResponseRecorder); ok {
		return rw.WroteHeader
	}
	return false
}
