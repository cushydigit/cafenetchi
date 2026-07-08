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
	UserID ContextKey = "user_id"
)

func Auth(secret string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")
			claims, err := utils.ParseJWT(token, secret)
			if err != nil {
				helpers.Error(w, types.ErrInvalidOTP)
				return
			}

			ctx := context.WithValue(r.Context(), UserID, claims)

			next.ServeHTTP(
				w,
				r.WithContext(ctx),
			)
		})
	}
}
