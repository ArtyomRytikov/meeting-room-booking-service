package handler

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"test-backend-1-ArtyomRytikov/internal/domain"
	"test-backend-1-ArtyomRytikov/internal/service"
)

type roomRepoForHandlerMock struct {
	createFn func(ctx context.Context, room domain.Room) (domain.Room, error)
	listFn   func(ctx context.Context) ([]domain.Room, error)
	initFn   func(ctx context.Context) error
}

func (m *roomRepoForHandlerMock) Init(ctx context.Context) error {
	if m.initFn != nil {
		return m.initFn(ctx)
	}
	return nil
}

func (m *roomRepoForHandlerMock) Create(ctx context.Context, room domain.Room) (domain.Room, error) {
	if m.createFn != nil {
		return m.createFn(ctx, room)
	}
	return domain.Room{}, nil
}

func (m *roomRepoForHandlerMock) List(ctx context.Context) ([]domain.Room, error) {
	if m.listFn != nil {
		return m.listFn(ctx)
	}
	return []domain.Room{}, nil
}

func TestRoomHandler_Create_Success(t *testing.T) {
	repo := &roomRepoForHandlerMock{
		createFn: func(ctx context.Context, room domain.Room) (domain.Room, error) {
			room.ID = "room-1"
			return room, nil
		},
	}
	h := NewRoomHandler(service.NewRoomService(repo))

	req := httptest.NewRequest(http.MethodPost, "/rooms/create", bytes.NewBufferString(`{"name":"Room A","capacity":6}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Create(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}
}

func TestRoomHandler_Create_InvalidBody(t *testing.T) {
	h := NewRoomHandler(service.NewRoomService(&roomRepoForHandlerMock{}))

	req := httptest.NewRequest(http.MethodPost, "/rooms/create", bytes.NewBufferString(`bad json`))
	w := httptest.NewRecorder()

	h.Create(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestRoomHandler_Create_InvalidMethod(t *testing.T) {
	h := NewRoomHandler(service.NewRoomService(&roomRepoForHandlerMock{}))

	req := httptest.NewRequest(http.MethodGet, "/rooms/create", nil)
	w := httptest.NewRecorder()

	h.Create(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", w.Code)
	}
}

func TestRoomHandler_List_Success(t *testing.T) {
	repo := &roomRepoForHandlerMock{
		listFn: func(ctx context.Context) ([]domain.Room, error) {
			return []domain.Room{{ID: "1", Name: "Room A"}}, nil
		},
	}
	h := NewRoomHandler(service.NewRoomService(repo))

	req := httptest.NewRequest(http.MethodGet, "/rooms/list", nil)
	w := httptest.NewRecorder()

	h.List(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestRoomHandler_List_InternalError(t *testing.T) {
	repo := &roomRepoForHandlerMock{
		listFn: func(ctx context.Context) ([]domain.Room, error) {
			return nil, errors.New("db error")
		},
	}
	h := NewRoomHandler(service.NewRoomService(repo))

	req := httptest.NewRequest(http.MethodGet, "/rooms/list", nil)
	w := httptest.NewRecorder()

	h.List(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}
