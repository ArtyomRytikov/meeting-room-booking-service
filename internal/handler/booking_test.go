package handler

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"test-backend-1-ArtyomRytikov/internal/domain"
	"test-backend-1-ArtyomRytikov/internal/middleware"
	"test-backend-1-ArtyomRytikov/internal/service"
)

type bookingRepoForHandlerMock struct {
	initFn             func(ctx context.Context) error
	findSlotFn         func(ctx context.Context, slotID string) (*domain.SlotDetails, error)
	createFn           func(ctx context.Context, userID, slotID string) (domain.Booking, error)
	findByIDFn         func(ctx context.Context, bookingID string) (*domain.Booking, error)
	cancelFn           func(ctx context.Context, bookingID string) error
	listMyFn           func(ctx context.Context, userID string) ([]domain.Booking, error)
	countAllFn         func(ctx context.Context) (int, error)
	listAllPaginatedFn func(ctx context.Context, limit, offset int) ([]domain.Booking, error)
}

func (m *bookingRepoForHandlerMock) Init(ctx context.Context) error {
	if m.initFn != nil {
		return m.initFn(ctx)
	}
	return nil
}

func (m *bookingRepoForHandlerMock) FindSlot(ctx context.Context, slotID string) (*domain.SlotDetails, error) {
	if m.findSlotFn != nil {
		return m.findSlotFn(ctx, slotID)
	}
	return &domain.SlotDetails{
		SlotID:  slotID,
		RoomID:  "room-1",
		StartAt: time.Now().UTC().Add(1 * time.Hour),
		EndAt:   time.Now().UTC().Add(2 * time.Hour),
	}, nil
}

func (m *bookingRepoForHandlerMock) Create(ctx context.Context, userID, slotID string) (domain.Booking, error) {
	if m.createFn != nil {
		return m.createFn(ctx, userID, slotID)
	}
	return domain.Booking{
		ID:     "booking-1",
		SlotID: slotID,
		UserID: userID,
		Status: "active",
	}, nil
}

func (m *bookingRepoForHandlerMock) FindByID(ctx context.Context, bookingID string) (*domain.Booking, error) {
	if m.findByIDFn != nil {
		return m.findByIDFn(ctx, bookingID)
	}
	return &domain.Booking{
		ID:     bookingID,
		UserID: "user-1",
		Status: "active",
	}, nil
}

func (m *bookingRepoForHandlerMock) Cancel(ctx context.Context, bookingID string) error {
	if m.cancelFn != nil {
		return m.cancelFn(ctx, bookingID)
	}
	return nil
}

func (m *bookingRepoForHandlerMock) ListMy(ctx context.Context, userID string) ([]domain.Booking, error) {
	if m.listMyFn != nil {
		return m.listMyFn(ctx, userID)
	}
	return []domain.Booking{{ID: "booking-1"}}, nil
}

func (m *bookingRepoForHandlerMock) CountAll(ctx context.Context) (int, error) {
	if m.countAllFn != nil {
		return m.countAllFn(ctx)
	}
	return 1, nil
}

func (m *bookingRepoForHandlerMock) ListAllPaginated(ctx context.Context, limit, offset int) ([]domain.Booking, error) {
	if m.listAllPaginatedFn != nil {
		return m.listAllPaginatedFn(ctx, limit, offset)
	}
	return []domain.Booking{{ID: "booking-1"}}, nil
}

func TestBookingHandler_Create_Success(t *testing.T) {
	h := NewBookingHandler(service.NewBookingService(&bookingRepoForHandlerMock{}))

	req := httptest.NewRequest(http.MethodPost, "/bookings/create", bytes.NewBufferString(`{"slotId":"slot-1"}`))
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, "user-1"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Create(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}
}

func TestBookingHandler_Create_InvalidBody(t *testing.T) {
	h := NewBookingHandler(service.NewBookingService(&bookingRepoForHandlerMock{}))

	req := httptest.NewRequest(http.MethodPost, "/bookings/create", bytes.NewBufferString(`bad json`))
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, "user-1"))
	w := httptest.NewRecorder()

	h.Create(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestBookingHandler_Create_Conflict(t *testing.T) {
	repo := &bookingRepoForHandlerMock{
		createFn: func(ctx context.Context, userID, slotID string) (domain.Booking, error) {
			return domain.Booking{}, errors.New("ux_bookings_slot_active")
		},
	}
	h := NewBookingHandler(service.NewBookingService(repo))

	req := httptest.NewRequest(http.MethodPost, "/bookings/create", bytes.NewBufferString(`{"slotId":"slot-1"}`))
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, "user-1"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Create(w, req)

	if w.Code != http.StatusConflict {
		t.Fatalf("expected 409, got %d", w.Code)
	}
}

func TestBookingHandler_Cancel_Success(t *testing.T) {
	h := NewBookingHandler(service.NewBookingService(&bookingRepoForHandlerMock{
		findByIDFn: func(ctx context.Context, bookingID string) (*domain.Booking, error) {
			return &domain.Booking{
				ID:     bookingID,
				UserID: "user-1",
				Status: "active",
			}, nil
		},
	}))

	req := httptest.NewRequest(http.MethodPost, "/bookings/booking-1/cancel", nil)
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, "user-1")
	ctx = context.WithValue(ctx, middleware.RoleKey, "user")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	h.Cancel(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestBookingHandler_My_Success(t *testing.T) {
	h := NewBookingHandler(service.NewBookingService(&bookingRepoForHandlerMock{}))

	req := httptest.NewRequest(http.MethodGet, "/bookings/my", nil)
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, "user-1"))
	w := httptest.NewRecorder()

	h.My(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestBookingHandler_List_Success(t *testing.T) {
	h := NewBookingHandler(service.NewBookingService(&bookingRepoForHandlerMock{}))

	req := httptest.NewRequest(http.MethodGet, "/bookings/list?page=1&pageSize=20", nil)
	w := httptest.NewRecorder()

	h.List(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestBookingHandler_List_InvalidPage(t *testing.T) {
	h := NewBookingHandler(service.NewBookingService(&bookingRepoForHandlerMock{}))

	req := httptest.NewRequest(http.MethodGet, "/bookings/list?page=bad&pageSize=20", nil)
	w := httptest.NewRecorder()

	h.List(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestBookingIDFromCancelPath(t *testing.T) {
	id, ok := bookingIDFromCancelPath("/bookings/booking-1/cancel")
	if !ok || id != "booking-1" {
		t.Fatal("expected valid booking id parse")
	}
}
