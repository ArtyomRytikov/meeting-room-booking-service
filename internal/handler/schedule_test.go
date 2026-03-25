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
	"test-backend-1-ArtyomRytikov/internal/service"
)

type scheduleRepoForHandlerMock struct {
	initFn            func(ctx context.Context) error
	roomExistsFn      func(ctx context.Context, roomID string) (bool, error)
	scheduleExistsFn  func(ctx context.Context, roomID string) (bool, error)
	createWithSlotsFn func(ctx context.Context, schedule domain.Schedule, days int) (domain.Schedule, error)
	listSlotsByDateFn func(ctx context.Context, roomID string, date time.Time) ([]domain.Slot, error)
}

func (m *scheduleRepoForHandlerMock) Init(ctx context.Context) error {
	if m.initFn != nil {
		return m.initFn(ctx)
	}
	return nil
}

func (m *scheduleRepoForHandlerMock) RoomExists(ctx context.Context, roomID string) (bool, error) {
	if m.roomExistsFn != nil {
		return m.roomExistsFn(ctx, roomID)
	}
	return true, nil
}

func (m *scheduleRepoForHandlerMock) ScheduleExists(ctx context.Context, roomID string) (bool, error) {
	if m.scheduleExistsFn != nil {
		return m.scheduleExistsFn(ctx, roomID)
	}
	return false, nil
}

func (m *scheduleRepoForHandlerMock) CreateScheduleWithSlots(ctx context.Context, schedule domain.Schedule, days int) (domain.Schedule, error) {
	if m.createWithSlotsFn != nil {
		return m.createWithSlotsFn(ctx, schedule, days)
	}
	schedule.ID = "schedule-1"
	return schedule, nil
}

func (m *scheduleRepoForHandlerMock) ListSlotsByDate(ctx context.Context, roomID string, date time.Time) ([]domain.Slot, error) {
	if m.listSlotsByDateFn != nil {
		return m.listSlotsByDateFn(ctx, roomID, date)
	}
	return []domain.Slot{{ID: "slot-1", RoomID: roomID}}, nil
}

func TestScheduleHandler_Create_Success(t *testing.T) {
	h := NewScheduleHandler(service.NewScheduleService(&scheduleRepoForHandlerMock{}))

	req := httptest.NewRequest(http.MethodPost, "/rooms/11111111-1111-1111-1111-111111111111/schedule/create", bytes.NewBufferString(`{"daysOfWeek":[1,2],"startTime":"09:00","endTime":"10:00"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Create(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}
}

func TestScheduleHandler_Create_InvalidBody(t *testing.T) {
	h := NewScheduleHandler(service.NewScheduleService(&scheduleRepoForHandlerMock{}))

	req := httptest.NewRequest(http.MethodPost, "/rooms/11111111-1111-1111-1111-111111111111/schedule/create", bytes.NewBufferString(`bad json`))
	w := httptest.NewRecorder()

	h.Create(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestScheduleHandler_Create_ScheduleExists(t *testing.T) {
	repo := &scheduleRepoForHandlerMock{
		scheduleExistsFn: func(ctx context.Context, roomID string) (bool, error) {
			return true, nil
		},
	}
	h := NewScheduleHandler(service.NewScheduleService(repo))

	req := httptest.NewRequest(http.MethodPost, "/rooms/11111111-1111-1111-1111-111111111111/schedule/create", bytes.NewBufferString(`{"daysOfWeek":[1],"startTime":"09:00","endTime":"10:00"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Create(w, req)

	if w.Code != http.StatusConflict {
		t.Fatalf("expected 409, got %d", w.Code)
	}
}

func TestScheduleHandler_ListSlots_Success(t *testing.T) {
	h := NewScheduleHandler(service.NewScheduleService(&scheduleRepoForHandlerMock{}))

	req := httptest.NewRequest(http.MethodGet, "/rooms/11111111-1111-1111-1111-111111111111/slots/list?date=2026-03-26", nil)
	w := httptest.NewRecorder()

	h.ListSlots(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestScheduleHandler_ListSlots_InvalidDate(t *testing.T) {
	h := NewScheduleHandler(service.NewScheduleService(&scheduleRepoForHandlerMock{}))

	req := httptest.NewRequest(http.MethodGet, "/rooms/11111111-1111-1111-1111-111111111111/slots/list?date=bad-date", nil)
	w := httptest.NewRecorder()

	h.ListSlots(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestRoomIDHelpers(t *testing.T) {
	roomID, ok := roomIDFromSchedulePath("/rooms/11111111-1111-1111-1111-111111111111/schedule/create")
	if !ok || roomID == "" {
		t.Fatal("expected valid schedule path parse")
	}

	roomID, ok = roomIDFromSlotsPath("/rooms/11111111-1111-1111-1111-111111111111/slots/list")
	if !ok || roomID == "" {
		t.Fatal("expected valid slots path parse")
	}
}

func TestScheduleHandler_ListSlots_InternalError(t *testing.T) {
	repo := &scheduleRepoForHandlerMock{
		roomExistsFn: func(ctx context.Context, roomID string) (bool, error) {
			return false, errors.New("db error")
		},
	}
	h := NewScheduleHandler(service.NewScheduleService(repo))

	req := httptest.NewRequest(http.MethodGet, "/rooms/11111111-1111-1111-1111-111111111111/slots/list?date=2026-03-26", nil)
	w := httptest.NewRecorder()

	h.ListSlots(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}
