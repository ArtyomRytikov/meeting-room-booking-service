package middleware

import (
	"context"
	"net/http"
	"strings"

	"test-backend-1-ArtyomRytikov/internal/auth"
)

type contextKey string

const (
	UserIDKey contextKey = "user_id"
	RoleKey   contextKey = "role"
)

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, `{"error":{"code":"UNAUTHORIZED","message":"missing authorization header"}}`, http.StatusUnauthorized)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" || parts[1] == "" {
			http.Error(w, `{"error":{"code":"UNAUTHORIZED","message":"invalid authorization header"}}`, http.StatusUnauthorized)
			return
		}

		claims, err := auth.ParseToken(parts[1])
		if err != nil {
			http.Error(w, `{"error":{"code":"UNAUTHORIZED","message":"invalid token"}}`, http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, RoleKey, claims.Role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RequireRole(role string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got, _ := r.Context().Value(RoleKey).(string)
		if got != role {
			http.Error(w, `{"error":{"code":"FORBIDDEN","message":"forbidden"}}`, http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
