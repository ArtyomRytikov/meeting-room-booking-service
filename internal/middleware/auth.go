package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"test-backend-1-ArtyomRytikov/internal/auth"
)

type contextKey string

const (
	UserIDKey contextKey = "user_id"
	RoleKey   contextKey = "role"
)

type errorItem struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type errorBody struct {
	Error errorItem `json:"error"`
}

func writeJSONError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(errorBody{
		Error: errorItem{
			Code:    code,
			Message: message,
		},
	})
}

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			writeJSONError(w, http.StatusUnauthorized, "UNAUTHORIZED", "missing authorization header")
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" || parts[1] == "" {
			writeJSONError(w, http.StatusUnauthorized, "UNAUTHORIZED", "invalid authorization header")
			return
		}

		claims, err := auth.ParseToken(parts[1])
		if err != nil {
			writeJSONError(w, http.StatusUnauthorized, "UNAUTHORIZED", "invalid token")
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
			writeJSONError(w, http.StatusForbidden, "FORBIDDEN", "forbidden")
			return
		}

		next.ServeHTTP(w, r)
	})
}
