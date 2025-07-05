package middleware

import (
	"context"
	"net/http"
	"strings"

	"merchshop/internal/api/http/auth"
)

type contextKey string

const UserIDKey contextKey = "user_id"

func AuthMiddleware(tokenManager auth.TokenManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Неавторизован", http.StatusUnauthorized)
				return
			}

			headerParts := strings.Split(authHeader, " ")
			if len(headerParts) != 2 || headerParts[0] != "Bearer" {
				http.Error(w, "Неверный header", http.StatusUnauthorized)
				return
			}

			userID, err := tokenManager.Parse(headerParts[1])
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
