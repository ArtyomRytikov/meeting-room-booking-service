package handler

import (
	"net/http"

	"test-backend-1-ArtyomRytikov/internal/middleware"
	"test-backend-1-ArtyomRytikov/internal/service"
)

func NewRouter(
	roomService *service.RoomService,
	scheduleService *service.ScheduleService,
	bookingService *service.BookingService,
) http.Handler {
	mux := http.NewServeMux()

	roomHandler := NewRoomHandler(roomService)
	scheduleHandler := NewScheduleHandler(scheduleService)
	bookingHandler := NewBookingHandler(bookingService)

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
			writeJSON(w, http.StatusOK, map[string]string{"status": "admin ok"})
		})),
	))

	mux.Handle("/rooms/create", middleware.RequireAuth(
		middleware.RequireRole("admin", http.HandlerFunc(roomHandler.Create)),
	))

	mux.Handle("/rooms/list", middleware.RequireAuth(
		http.HandlerFunc(roomHandler.List),
	))

	mux.Handle("/bookings/create", middleware.RequireAuth(
		http.HandlerFunc(bookingHandler.Create),
	))

	mux.Handle("/bookings/my", middleware.RequireAuth(
		http.HandlerFunc(bookingHandler.My),
	))

	mux.Handle("/bookings/list", middleware.RequireAuth(
		middleware.RequireRole("admin", http.HandlerFunc(bookingHandler.List)),
	))

	mux.Handle("/rooms/", middleware.RequireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && hasSuffix(r.URL.Path, "/schedule/create"):
			middleware.RequireRole("admin", http.HandlerFunc(scheduleHandler.Create)).ServeHTTP(w, r)
			return
		case r.Method == http.MethodGet && hasSuffix(r.URL.Path, "/slots/list"):
			scheduleHandler.ListSlots(w, r)
			return
		default:
			writeAPIError(w, http.StatusNotFound, "NOT_FOUND", "not found")
			return
		}
	})))

	mux.Handle("/bookings/", middleware.RequireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && hasSuffix(r.URL.Path, "/cancel"):
			bookingHandler.Cancel(w, r)
			return
		default:
			writeAPIError(w, http.StatusNotFound, "NOT_FOUND", "not found")
			return
		}
	})))

	return mux
}

func hasSuffix(path, suffix string) bool {
	if len(path) < len(suffix) {
		return false
	}
	return path[len(path)-len(suffix):] == suffix
}
