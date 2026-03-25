package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDummyLogin_Success(t *testing.T) {
	body := []byte(`{"role":"user"}`)
	req := httptest.NewRequest(http.MethodPost, "/dummyLogin", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	DummyLogin(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestDummyLogin_InvalidBody(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/dummyLogin", bytes.NewReader([]byte(`bad json`)))
	w := httptest.NewRecorder()

	DummyLogin(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestDummyLogin_InvalidRole(t *testing.T) {
	body := []byte(`{"role":"manager"}`)
	req := httptest.NewRequest(http.MethodPost, "/dummyLogin", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	DummyLogin(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}
