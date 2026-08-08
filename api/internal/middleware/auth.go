package middleware

import (
	"cafenetchi-api/internal/contextx"
	"cafenetchi-api/internal/helpers"
	"cafenetchi-api/internal/types"
	"cafenetchi-api/internal/utils"
	"net/http"
	"strings"
)

func Auth(secret string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if !strings.HasPrefix(auth, "Bearer ") {
				helpers.Error(w, types.ErrNotAuthenticated)
				return
			}

			token := strings.TrimPrefix(auth, "Bearer ")
			claims, err := utils.ParseJWT(token, secret)
			if err != nil {
				helpers.Error(w, types.ErrNotAuthenticated)
				return
			}

			ctx := contextx.WithUserID(r.Context(), claims.UserID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
