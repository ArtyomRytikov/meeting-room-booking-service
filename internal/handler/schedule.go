package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"test-backend-1-ArtyomRytikov/internal/domain"
	"test-backend-1-ArtyomRytikov/internal/service"
)

type ScheduleHandler struct {
	service *service.ScheduleService
}

func NewScheduleHandler(service *service.ScheduleService) *ScheduleHandler {
	return &ScheduleHandler{service: service}
}

type createScheduleRequest struct {
	DaysOfWeek []int  `json:"daysOfWeek"`
	StartTime  string `json:"startTime"`
	EndTime    string `json:"endTime"`
}

func (h *ScheduleHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeAPIError(w, http.StatusMethodNotAllowed, "INVALID_REQUEST", "method not allowed")
		return
	}

	roomID, ok := roomIDFromSchedulePath(r.URL.Path)
	if !ok {
		writeAPIError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid roomId path")
		return
	}

	var req createScheduleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeAPIError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid request body")
		return
	}

	schedule, err := h.service.Create(r.Context(), domain.Schedule{
		RoomID:     roomID,
		DaysOfWeek: req.DaysOfWeek,
		StartTime:  req.StartTime,
		EndTime:    req.EndTime,
	})
	if err != nil {
		switch err.Error() {
		case "room not found":
			writeAPIError(w, http.StatusNotFound, "ROOM_NOT_FOUND", "room not found")
		case "schedule exists":
			writeAPIError(w, http.StatusBadRequest, "SCHEDULE_EXISTS", "schedule already exists")
		default:
			writeAPIError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		}
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"schedule": schedule,
	})
}

func (h *ScheduleHandler) ListSlots(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeAPIError(w, http.StatusMethodNotAllowed, "INVALID_REQUEST", "method not allowed")
		return
	}

	roomID, ok := roomIDFromSlotsPath(r.URL.Path)
	if !ok {
		writeAPIError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid roomId path")
		return
	}

	date := r.URL.Query().Get("date")
	if strings.TrimSpace(date) == "" {
		writeAPIError(w, http.StatusBadRequest, "INVALID_REQUEST", "date is required")
		return
	}

	slots, err := h.service.ListSlotsByDate(r.Context(), roomID, date)
	if err != nil {
		switch err.Error() {
		case "room not found":
			writeAPIError(w, http.StatusNotFound, "ROOM_NOT_FOUND", "room not found")
		case "invalid date":
			writeAPIError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid date")
		default:
			writeAPIError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to list slots")
		}
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"slots": slots,
	})
}

func roomIDFromSchedulePath(path string) (string, bool) {
	const prefix = "/rooms/"
	const suffix = "/schedule/create"

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

func roomIDFromSlotsPath(path string) (string, bool) {
	const prefix = "/rooms/"
	const suffix = "/slots/list"

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
