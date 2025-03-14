package middleware

import (
	"encoding/json"
	"fmt"
	"github.com/oprimogus/flyfood-api/internal/xerrors"
	logger "github.com/oprimogus/flyfood-api/pkg/log"
	"log/slog"
	"net/http"
	"runtime/debug"
	"strings"
)

// Recovery é um middleware que recupera e lida com panics durante o processamento de uma requisição
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				// Converte o recover para um erro, se possível
				var err error
				switch x := rec.(type) {
				case error:
					err = x
				case string:
					err = fmt.Errorf("%s", x)
				default:
					err = fmt.Errorf("%v", x)
				}

				// Obtém o stack trace completo
				stackTrace := debug.Stack()

				// Prepara os detalhes do erro
				errorDetails := prepareErrorDetails(r, err, stackTrace)

				// Registra o erro de forma detalhada
				slog.ErrorContext(r.Context(), "Panic recuperado",
					"error", errorDetails.ErrorMessage,
					"trace_id", errorDetails.TraceID)

				// Imprime o erro no console para facilitar debugging local
				fmt.Println(formatConsoleOutput(errorDetails))

				// Responde com erro interno se o header ainda não foi escrito
				if !isHeaderWritten(w) {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)

					// Obtém o trace ID do contexto
					reqData, ok := r.Context().Value(string(logger.RequestKey)).(*logger.RequestData)
					traceID := ""
					if ok && reqData != nil {
						traceID = reqData.TraceID
					}

					// Codifica o erro
					_ = json.NewEncoder(w).Encode(xerrors.InternalServer(traceID, "Erro interno do servidor"))
				}
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// errorDetails contém informações detalhadas sobre o erro
type errorDetails struct {
	ErrorMessage string
	StackTrace   string
	TraceID      string
	RequestURL   string
	Method       string
	ClientIP     string
}

// prepareErrorDetails prepara os detalhes do erro para log e debug
func prepareErrorDetails(r *http.Request, err error, stackTrace []byte) errorDetails {
	// Obtém o trace ID do contexto, se disponível
	traceID := ""
	if reqData, ok := r.Context().Value(string(logger.RequestKey)).(*logger.RequestData); ok && reqData != nil {
		traceID = reqData.TraceID
	}

	return errorDetails{
		ErrorMessage: err.Error(),
		StackTrace:   sanitizeStackTrace(string(stackTrace)),
		TraceID:      traceID,
		RequestURL:   r.URL.String(),
		Method:       r.Method,
		ClientIP:     r.RemoteAddr,
	}
}

// sanitizeStackTrace remove linhas redundantes e formata o stack trace
func sanitizeStackTrace(trace string) string {
	// Divide o stack trace em linhas
	lines := strings.Split(trace, "\n")

	// Remove linhas vazias e mantém uma formatação limpa
	var cleanLines []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.Contains(line, "runtime/panic.go") {
			cleanLines = append(cleanLines, line)
		}
	}

	return strings.Join(cleanLines, "\n")
}

// formatConsoleOutput cria uma saída formatada para console
func formatConsoleOutput(details errorDetails) string {
	return fmt.Sprintf(`
========== PANIC RECOVERED ==========
Error: %s
Trace ID: %s
URL: %s
Method: %s
Client IP: %s

Stack Trace:
%s

==========================================
`,
		details.ErrorMessage,
		details.TraceID,
		details.RequestURL,
		details.Method,
		details.ClientIP,
		details.StackTrace,
	)
}

// isHeaderWritten verifica se o header já foi escrito
func isHeaderWritten(w http.ResponseWriter) bool {
	if rw, ok := w.(*ResponseRecorder); ok {
		return rw.WroteHeader
	}
	return false
}
