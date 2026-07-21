package middleware

import (
	"cafenetchi-api/internal/helpers"
	"cafenetchi-api/internal/limiter"
	"cafenetchi-api/internal/types"
	"net/http"
	"strconv"
)

func RateLimit(limiter limiter.Limiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			// TODO: Better key generation.
			// For now use the client IP. For OTP you may later use:
			// ip + ":" + phone
			key := r.RemoteAddr
			allowed, retryAfter, err := limiter.Allow(ctx, key)
			if err != nil {
				helpers.Error(w, types.ErrInternalServer)
				return
			}

			if !allowed {
				w.Header().Set("Retry-After", strconv.FormatInt(int64(retryAfter.Seconds()), 10))
				helpers.Error(w, types.ErrTooManyRequest)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
