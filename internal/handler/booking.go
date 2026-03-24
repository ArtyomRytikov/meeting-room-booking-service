package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"test-backend-1-ArtyomRytikov/internal/middleware"
	"test-backend-1-ArtyomRytikov/internal/service"
)

type BookingHandler struct {
	service *service.BookingService
}

func NewBookingHandler(service *service.BookingService) *BookingHandler {
	return &BookingHandler{service: service}
}

type createBookingRequest struct {
	SlotID string `json:"slotId"`
}

func (h *BookingHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeAPIError(w, http.StatusMethodNotAllowed, "INVALID_REQUEST", "method not allowed")
		return
	}

	userID, _ := r.Context().Value(middleware.UserIDKey).(string)

	var req createBookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeAPIError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid request body")
		return
	}

	booking, err := h.service.Create(r.Context(), userID, req.SlotID)
	if err != nil {
		switch err.Error() {
		case "slotId is required":
			writeAPIError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		case "slot not found":
			writeAPIError(w, http.StatusNotFound, "SLOT_NOT_FOUND", "slot not found")
		case "slot already booked":
			writeAPIError(w, http.StatusBadRequest, "SLOT_ALREADY_BOOKED", "slot already booked")
		default:
			writeAPIError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to create booking")
		}
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"booking": booking,
	})
}

func (h *BookingHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeAPIError(w, http.StatusMethodNotAllowed, "INVALID_REQUEST", "method not allowed")
		return
	}

	bookingID, ok := bookingIDFromCancelPath(r.URL.Path)
	if !ok {
		writeAPIError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid bookingId path")
		return
	}

	userID, _ := r.Context().Value(middleware.UserIDKey).(string)
	role, _ := r.Context().Value(middleware.RoleKey).(string)

	err := h.service.Cancel(r.Context(), bookingID, userID, role)
	if err != nil {
		switch err.Error() {
		case "booking not found":
			writeAPIError(w, http.StatusNotFound, "BOOKING_NOT_FOUND", "booking not found")
		case "forbidden":
			writeAPIError(w, http.StatusForbidden, "FORBIDDEN", "forbidden")
		case "booking not found or already cancelled":
			writeAPIError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		default:
			writeAPIError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to cancel booking")
		}
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"status": "cancelled",
	})
}

func (h *BookingHandler) My(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(middleware.UserIDKey).(string)

	bookings, err := h.service.ListMy(r.Context(), userID)
	if err != nil {
		writeAPIError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to list bookings")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"bookings": bookings,
	})
}

func (h *BookingHandler) List(w http.ResponseWriter, r *http.Request) {
	bookings, err := h.service.ListAll(r.Context())
	if err != nil {
		writeAPIError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to list bookings")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"bookings": bookings,
	})
}

func bookingIDFromCancelPath(path string) (string, bool) {
	const prefix = "/bookings/"
	const suffix = "/cancel"

	if !strings.HasPrefix(path, prefix) || !strings.HasSuffix(path, suffix) {
		return "", false
	}

	value := strings.TrimSuffix(strings.TrimPrefix(path, prefix), suffix)
	value = strings.Trim(value, "/")
	if value == "" {
		return "", false
	}

	return value, true
}
