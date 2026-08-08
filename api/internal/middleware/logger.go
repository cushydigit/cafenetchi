package middleware

import (
	"cafenetchi-api/internal/contextx"
	"cafenetchi-api/internal/logger"
	"log/slog"
	"net/http"
	"time"

	chi_middleware "github.com/go-chi/chi/v5/middleware"
)

func RequestLogger(log *logger.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			ww := chi_middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			next.ServeHTTP(ww, r)

			attrs := []any{
				slog.String("request_id", contextx.RequestID(r.Context())),
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Int("status", ww.Status()),
				slog.Duration("duration", time.Since(start)),
				slog.String("remote_ip", r.RemoteAddr),
				slog.String("user_agent", r.UserAgent()),
				slog.Int("bytes", ww.BytesWritten()),
			}

			log.LogRequest("request completed", ww.Status(), attrs...)
		})
	}
}
