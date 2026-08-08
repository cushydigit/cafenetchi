package logger

import (
	"cafenetchi-api/internal/contextx"
	"context"
	"log/slog"
)

// enrich logger with context values (request_id, user_id)
func (l *Logger) WithContext(ctx context.Context) *slog.Logger {
	attrs := []any{}

	if id := contextx.RequestID(ctx); id != "" {
		attrs = append(attrs, "request_id", id)
	}

	if userID := contextx.UserID(ctx); userID != 0 {
		attrs = append(attrs, slog.Int64("user_id", userID))
	}

	return l.With(attrs...)
}
