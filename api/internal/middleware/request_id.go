package middleware

import (
	"cafenetchi-api/internal/contextx"
	"net/http"

	"github.com/google/uuid"
)

// archived for now
// using chi request ID middleware
const RequestIDHeader = "X-Request-ID"

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get(RequestIDHeader)

		if requestID == "" {
			requestID = uuid.NewString()
		}

		ctx := contextx.WithRequestID(r.Context(), requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})

}
