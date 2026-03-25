package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"test-backend-1-ArtyomRytikov/internal/domain"
)

type scheduleRepoMock struct {
	roomExistsFn         func(ctx context.Context, roomID string) (bool, error)
	scheduleExistsFn     func(ctx context.Context, roomID string) (bool, error)
	createScheduleWithFn func(ctx context.Context, schedule domain.Schedule, days int) (domain.Schedule, error)
	listSlotsByDateFn    func(ctx context.Context, roomID string, date time.Time) ([]domain.Slot, error)
	initFn               func(ctx context.Context) error
}

func (m *scheduleRepoMock) Init(ctx context.Context) error {
	if m.initFn != nil {
		return m.initFn(ctx)
	}
	return nil
}

func (m *scheduleRepoMock) RoomExists(ctx context.Context, roomID string) (bool, error) {
	if m.roomExistsFn != nil {
		return m.roomExistsFn(ctx, roomID)
	}
	return false, nil
}

func (m *scheduleRepoMock) ScheduleExists(ctx context.Context, roomID string) (bool, error) {
	if m.scheduleExistsFn != nil {
		return m.scheduleExistsFn(ctx, roomID)
	}
	return false, nil
}

func (m *scheduleRepoMock) CreateScheduleWithSlots(ctx context.Context, schedule domain.Schedule, days int) (domain.Schedule, error) {
	if m.createScheduleWithFn != nil {
		return m.createScheduleWithFn(ctx, schedule, days)
	}
	return domain.Schedule{}, nil
}

func (m *scheduleRepoMock) ListSlotsByDate(ctx context.Context, roomID string, date time.Time) ([]domain.Slot, error) {
	if m.listSlotsByDateFn != nil {
		return m.listSlotsByDateFn(ctx, roomID, date)
	}
	return []domain.Slot{}, nil
}

func TestScheduleService_Create_InvalidRoomID(t *testing.T) {
	svc := NewScheduleService(&scheduleRepoMock{})

	_, err := svc.Create(context.Background(), domain.Schedule{
		RoomID:     "bad-id",
		DaysOfWeek: []int{1, 2},
		StartTime:  "09:00",
		EndTime:    "10:00",
	})
	if err == nil {
		t.Fatal("expected error")
	}
	if err.Error() != "invalid room id" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestScheduleService_Create_RoomNotFound(t *testing.T) {
	svc := NewScheduleService(&scheduleRepoMock{
		roomExistsFn: func(ctx context.Context, roomID string) (bool, error) {
			return false, nil
		},
	})

	_, err := svc.Create(context.Background(), domain.Schedule{
		RoomID:     "11111111-1111-1111-1111-111111111111",
		DaysOfWeek: []int{1},
		StartTime:  "09:00",
		EndTime:    "10:00",
	})
	if err == nil {
		t.Fatal("expected error")
	}
	if err.Error() != "room not found" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestScheduleService_Create_ScheduleExists(t *testing.T) {
	svc := NewScheduleService(&scheduleRepoMock{
		roomExistsFn: func(ctx context.Context, roomID string) (bool, error) {
			return true, nil
		},
		scheduleExistsFn: func(ctx context.Context, roomID string) (bool, error) {
			return true, nil
		},
	})

	_, err := svc.Create(context.Background(), domain.Schedule{
		RoomID:     "11111111-1111-1111-1111-111111111111",
		DaysOfWeek: []int{1},
		StartTime:  "09:00",
		EndTime:    "10:00",
	})
	if err == nil {
		t.Fatal("expected error")
	}
	if err.Error() != "schedule exists" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestScheduleService_Create_InvalidDayOfWeek(t *testing.T) {
	svc := NewScheduleService(&scheduleRepoMock{
		roomExistsFn: func(ctx context.Context, roomID string) (bool, error) {
			return true, nil
		},
		scheduleExistsFn: func(ctx context.Context, roomID string) (bool, error) {
			return false, nil
		},
	})

	_, err := svc.Create(context.Background(), domain.Schedule{
		RoomID:     "11111111-1111-1111-1111-111111111111",
		DaysOfWeek: []int{8},
		StartTime:  "09:00",
		EndTime:    "10:00",
	})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestScheduleService_Create_Success(t *testing.T) {
	svc := NewScheduleService(&scheduleRepoMock{
		roomExistsFn: func(ctx context.Context, roomID string) (bool, error) {
			return true, nil
		},
		scheduleExistsFn: func(ctx context.Context, roomID string) (bool, error) {
			return false, nil
		},
		createScheduleWithFn: func(ctx context.Context, schedule domain.Schedule, days int) (domain.Schedule, error) {
			schedule.ID = "schedule-1"
			return schedule, nil
		},
	})

	result, err := svc.Create(context.Background(), domain.Schedule{
		RoomID:     "11111111-1111-1111-1111-111111111111",
		DaysOfWeek: []int{1, 2, 3},
		StartTime:  "09:00",
		EndTime:    "10:00",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.ID != "schedule-1" {
		t.Fatalf("expected schedule id")
	}
}

func TestScheduleService_ListSlots_InvalidRoomID(t *testing.T) {
	svc := NewScheduleService(&scheduleRepoMock{})

	_, err := svc.ListSlotsByDate(context.Background(), "bad-id", "2026-03-26")
	if err == nil {
		t.Fatal("expected error")
	}
	if err.Error() != "invalid room id" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestScheduleService_ListSlots_InvalidDate(t *testing.T) {
	svc := NewScheduleService(&scheduleRepoMock{
		roomExistsFn: func(ctx context.Context, roomID string) (bool, error) {
			return true, nil
		},
	})

	_, err := svc.ListSlotsByDate(context.Background(), "11111111-1111-1111-1111-111111111111", "bad-date")
	if err == nil {
		t.Fatal("expected error")
	}
	if err.Error() != "invalid date" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestScheduleService_Init_Error(t *testing.T) {
	svc := NewScheduleService(&scheduleRepoMock{
		initFn: func(ctx context.Context) error {
			return errors.New("init failed")
		},
	})

	err := svc.Init(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
}
