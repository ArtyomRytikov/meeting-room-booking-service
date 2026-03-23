package handler

import (
	"encoding/json"
	"net/http"

	"test-backend-1-ArtyomRytikov/internal/domain"
	"test-backend-1-ArtyomRytikov/internal/service"
)

type RoomHandler struct {
	service *service.RoomService
}

func NewRoomHandler(service *service.RoomService) *RoomHandler {
	return &RoomHandler{service: service}
}

type createRoomRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
	Capacity    *int    `json:"capacity"`
}

func (h *RoomHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeAPIError(w, http.StatusMethodNotAllowed, "INVALID_REQUEST", "method not allowed")
		return
	}

	var req createRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeAPIError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid request body")
		return
	}

	room, err := h.service.Create(r.Context(), domain.Room{
		Name:        req.Name,
		Description: req.Description,
		Capacity:    req.Capacity,
	})
	if err != nil {
		writeAPIError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"room": room,
	})
}

func (h *RoomHandler) List(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeAPIError(w, http.StatusMethodNotAllowed, "INVALID_REQUEST", "method not allowed")
		return
	}

	rooms, err := h.service.List(r.Context())
	if err != nil {
		writeAPIError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to list rooms")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"rooms": rooms,
	})
}
