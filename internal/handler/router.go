package handler

import (
	"encoding/json"
	"net/http"

	"test-backend-1-ArtyomRytikov/internal/middleware"
)

func NewRouter() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/_info", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	mux.HandleFunc("/dummyLogin", DummyLogin)

	mux.Handle("/me", middleware.RequireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(middleware.UserIDKey).(string)
		role, _ := r.Context().Value(middleware.RoleKey).(string)

		writeJSON(w, http.StatusOK, map[string]string{
			"userId": userID,
			"role":   role,
		})
	})))

	mux.Handle("/admin/ping", middleware.RequireAuth(
		middleware.RequireRole("admin", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_ = json.NewEncoder(w).Encode(map[string]string{"status": "admin ok"})
		})),
	))

	return mux
}
