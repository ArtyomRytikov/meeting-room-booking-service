package handler

import (
	"encoding/json"
	"net/http"

	"test-backend-1-ArtyomRytikov/internal/auth"
)

type dummyLoginRequest struct {
	Role string `json:"role"`
}

type tokenResponse struct {
	Token string `json:"token"`
}

type errorBody struct {
	Error errorItem `json:"error"`
}

type errorItem struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func writeAPIError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, errorBody{
		Error: errorItem{
			Code:    code,
			Message: message,
		},
	})
}

func DummyLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeAPIError(w, http.StatusMethodNotAllowed, "INVALID_REQUEST", "method not allowed")
		return
	}

	var req dummyLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeAPIError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid request body")
		return
	}

	token, err := auth.GenerateToken(req.Role)
	if err != nil {
		writeAPIError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid role")
		return
	}

	writeJSON(w, http.StatusOK, tokenResponse{Token: token})
}
