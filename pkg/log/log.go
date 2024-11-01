package logger

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ContextKey string

const (
	TraceIDKey     ContextKey = "trace_id"
	UserIDKey      ContextKey = "user_id"
	RequestDataKey ContextKey = "request_data"
)

type RequestData struct {
	UserID   string `json:"user_id"`
	TraceID  string `json:"trace_id"`
	Method   string `json:"method"`
	Path     string `json:"path"`
	ClientIP string `json:"client_ip"`
}

type ContextualHandler struct {
	out  io.Writer
	opts slog.HandlerOptions
	ctx  context.Context
}

func NewContextualHandler(out io.Writer, opts *slog.HandlerOptions) *ContextualHandler {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}
	return &ContextualHandler{
		out:  out,
		opts: *opts,
	}
}

func (h *ContextualHandler) Handle(ctx context.Context, r slog.Record) error {
	m := make(map[string]interface{})

	// Adiciona campos padrão
	m["time"] = r.Time.Format(time.RFC3339)
	m["level"] = r.Level.String()
	m["message"] = r.Message

	if reqData, ok := ctx.Value(RequestDataKey).(*RequestData); ok {
		m["request"] = reqData
		// m["transaction_id"] = reqData.TraceID
		// m["method"] = reqData.Method
		// m["path"] = reqData.Path
		// m["client_ip"] = reqData.ClientIP
	}

	attrs := make(map[string]interface{})
	r.Attrs(func(a slog.Attr) bool {
		attrs[a.Key] = a.Value.Any()
		return true
	})

	if len(attrs) > 0 {
		m["attributes"] = attrs
	}

	encoded, err := json.Marshal(m)
	if err != nil {
		return err
	}

	_, err = h.out.Write(append(encoded, '\n'))
	return err
}

func (h *ContextualHandler) Enabled(context.Context, slog.Level) bool {
	return true
}

func (h *ContextualHandler) WithAttrs([]slog.Attr) slog.Handler {
	return h
}

func (h *ContextualHandler) WithGroup(string) slog.Handler {
	return h
}

func GinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString(string(UserIDKey))
		traceID := uuid.New().String()
		reqData := &RequestData{
			UserID:   userID,
			TraceID:  traceID,
			Method:   c.Request.Method,
			Path:     c.Request.URL.Path,
			ClientIP: c.ClientIP(),
		}
		c.Set(string(TraceIDKey), traceID)
		ctx := context.WithValue(c.Request.Context(), RequestDataKey, reqData)
		c.Request = c.Request.WithContext(ctx)

		c.Next()

		slog.InfoContext(ctx, "request completed",
			slog.Int("status", c.Writer.Status()),
		)
	}
}

func InitLogger(out io.Writer, level slog.Level) {
	handler := NewContextualHandler(out, &slog.HandlerOptions{
		Level:     level,
		AddSource: true,
	})
	logger := slog.New(handler)
	slog.SetDefault(logger)
}
