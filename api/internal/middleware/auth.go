package middleware

import (
	"cafenetchi-api/internal/helpers"
	"cafenetchi-api/internal/types"
	"cafenetchi-api/internal/utils"
	"context"
	"net/http"
)

type ContextKey string

const (
	UserIDKey ContextKey = "user_id"
)

func Auth(secret string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")
			claims, err := utils.ParseJWT(token, secret)
			if err != nil {
				helpers.Error(w, types.ErrNotAuthenticated)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)

			next.ServeHTTP(
				w,
				r.WithContext(ctx),
			)
		})
	}
}

func UserID(ctx context.Context) int64 {
	id, ok := ctx.Value(UserIDKey).(int64)
	if !ok {
		return 0
	}
	return id
}
