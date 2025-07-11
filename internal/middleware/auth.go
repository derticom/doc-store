package middleware

import (
	"context"
	"net/http"

	"github.com/derticom/doc-store/internal/usecase/auth"
)

type ctxKey string

const UserIDKey ctxKey = "userID"

func AuthMiddleware(store auth.SessionStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.URL.Query().Get("token")
			if token == "" {
				http.Error(w, "unauthorized: missing token", http.StatusUnauthorized)
				return
			}

			userID, err := store.GetUserID(r.Context(), token)
			if err != nil {
				http.Error(w, "unauthorized: invalid token", http.StatusForbidden)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
