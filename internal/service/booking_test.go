package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"test-backend-1-ArtyomRytikov/internal/domain"
)

type bookingRepoMock struct {
	findSlotFn         func(ctx context.Context, slotID string) (*domain.SlotDetails, error)
	createFn           func(ctx context.Context, userID, slotID string) (domain.Booking, error)
	findByIDFn         func(ctx context.Context, bookingID string) (*domain.Booking, error)
	cancelFn           func(ctx context.Context, bookingID string) error
	listMyFn           func(ctx context.Context, userID string) ([]domain.Booking, error)
	countAllFn         func(ctx context.Context) (int, error)
	listAllPaginatedFn func(ctx context.Context, limit, offset int) ([]domain.Booking, error)
	initFn             func(ctx context.Context) error
}

func (m *bookingRepoMock) Init(ctx context.Context) error {
	if m.initFn != nil {
		return m.initFn(ctx)
	}
	return nil
}

func (m *bookingRepoMock) FindSlot(ctx context.Context, slotID string) (*domain.SlotDetails, error) {
	if m.findSlotFn != nil {
		return m.findSlotFn(ctx, slotID)
	}
	return nil, errors.New("no rows in result set")
}

func (m *bookingRepoMock) Create(ctx context.Context, userID, slotID string) (domain.Booking, error) {
	if m.createFn != nil {
		return m.createFn(ctx, userID, slotID)
	}
	return domain.Booking{}, nil
}

func (m *bookingRepoMock) FindByID(ctx context.Context, bookingID string) (*domain.Booking, error) {
	if m.findByIDFn != nil {
		return m.findByIDFn(ctx, bookingID)
	}
	return nil, errors.New("not found")
}

func (m *bookingRepoMock) Cancel(ctx context.Context, bookingID string) error {
	if m.cancelFn != nil {
		return m.cancelFn(ctx, bookingID)
	}
	return nil
}

func (m *bookingRepoMock) ListMy(ctx context.Context, userID string) ([]domain.Booking, error) {
	if m.listMyFn != nil {
		return m.listMyFn(ctx, userID)
	}
	return []domain.Booking{}, nil
}

func (m *bookingRepoMock) CountAll(ctx context.Context) (int, error) {
	if m.countAllFn != nil {
		return m.countAllFn(ctx)
	}
	return 0, nil
}

func (m *bookingRepoMock) ListAllPaginated(ctx context.Context, limit, offset int) ([]domain.Booking, error) {
	if m.listAllPaginatedFn != nil {
		return m.listAllPaginatedFn(ctx, limit, offset)
	}
	return []domain.Booking{}, nil
}

func TestBookingService_Create_EmptySlotID(t *testing.T) {
	svc := NewBookingService(&bookingRepoMock{})

	_, err := svc.Create(context.Background(), "user-1", "")
	if err == nil {
		t.Fatal("expected error")
	}
	if err.Error() != "slotId is required" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBookingService_Create_SlotNotFound(t *testing.T) {
	svc := NewBookingService(&bookingRepoMock{
		findSlotFn: func(ctx context.Context, slotID string) (*domain.SlotDetails, error) {
			return nil, errors.New("no rows in result set")
		},
	})

	_, err := svc.Create(context.Background(), "user-1", "slot-1")
	if err == nil {
		t.Fatal("expected error")
	}
	if err.Error() != "slot not found" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBookingService_Create_SlotInPast(t *testing.T) {
	svc := NewBookingService(&bookingRepoMock{
		findSlotFn: func(ctx context.Context, slotID string) (*domain.SlotDetails, error) {
			return &domain.SlotDetails{
				SlotID:  slotID,
				RoomID:  "room-1",
				StartAt: time.Now().UTC().Add(-1 * time.Hour),
				EndAt:   time.Now().UTC(),
			}, nil
		},
	})

	_, err := svc.Create(context.Background(), "user-1", "slot-1")
	if err == nil {
		t.Fatal("expected error")
	}
	if err.Error() != "slot is in the past" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBookingService_Create_SlotAlreadyBooked(t *testing.T) {
	svc := NewBookingService(&bookingRepoMock{
		findSlotFn: func(ctx context.Context, slotID string) (*domain.SlotDetails, error) {
			return &domain.SlotDetails{
				SlotID:  slotID,
				RoomID:  "room-1",
				StartAt: time.Now().UTC().Add(1 * time.Hour),
				EndAt:   time.Now().UTC().Add(2 * time.Hour),
			}, nil
		},
		createFn: func(ctx context.Context, userID, slotID string) (domain.Booking, error) {
			return domain.Booking{}, errors.New("ux_bookings_slot_active")
		},
	})

	_, err := svc.Create(context.Background(), "user-1", "slot-1")
	if err == nil {
		t.Fatal("expected error")
	}
	if err.Error() != "slot already booked" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBookingService_Create_Success(t *testing.T) {
	svc := NewBookingService(&bookingRepoMock{
		findSlotFn: func(ctx context.Context, slotID string) (*domain.SlotDetails, error) {
			return &domain.SlotDetails{
				SlotID:  slotID,
				RoomID:  "room-1",
				StartAt: time.Now().UTC().Add(1 * time.Hour),
				EndAt:   time.Now().UTC().Add(2 * time.Hour),
			}, nil
		},
		createFn: func(ctx context.Context, userID, slotID string) (domain.Booking, error) {
			return domain.Booking{
				ID:     "booking-1",
				SlotID: slotID,
				UserID: userID,
				Status: "active",
			}, nil
		},
	})

	booking, err := svc.Create(context.Background(), "user-1", "slot-1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if booking.ID != "booking-1" {
		t.Fatalf("expected booking id")
	}
}

func TestBookingService_Cancel_NotFound(t *testing.T) {
	svc := NewBookingService(&bookingRepoMock{
		findByIDFn: func(ctx context.Context, bookingID string) (*domain.Booking, error) {
			return nil, errors.New("not found")
		},
	})

	_, err := svc.Cancel(context.Background(), "booking-1", "user-1", "user")
	if err == nil {
		t.Fatal("expected error")
	}
	if err.Error() != "booking not found" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBookingService_Cancel_Forbidden(t *testing.T) {
	svc := NewBookingService(&bookingRepoMock{
		findByIDFn: func(ctx context.Context, bookingID string) (*domain.Booking, error) {
			return &domain.Booking{
				ID:     bookingID,
				UserID: "owner-1",
				Status: "active",
			}, nil
		},
	})

	_, err := svc.Cancel(context.Background(), "booking-1", "user-1", "user")
	if err == nil {
		t.Fatal("expected error")
	}
	if err.Error() != "forbidden" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBookingService_Cancel_Idempotent(t *testing.T) {
	svc := NewBookingService(&bookingRepoMock{
		findByIDFn: func(ctx context.Context, bookingID string) (*domain.Booking, error) {
			return &domain.Booking{
				ID:     bookingID,
				UserID: "user-1",
				Status: "cancelled",
			}, nil
		},
	})

	booking, err := svc.Cancel(context.Background(), "booking-1", "user-1", "user")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if booking.Status != "cancelled" {
		t.Fatalf("expected cancelled status")
	}
}

func TestBookingService_ListAll_InvalidPage(t *testing.T) {
	svc := NewBookingService(&bookingRepoMock{})

	_, _, err := svc.ListAll(context.Background(), 0, 20)
	if err == nil {
		t.Fatal("expected error")
	}
	if err.Error() != "invalid page" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBookingService_ListAll_InvalidPageSize(t *testing.T) {
	svc := NewBookingService(&bookingRepoMock{})

	_, _, err := svc.ListAll(context.Background(), 1, 0)
	if err == nil {
		t.Fatal("expected error")
	}
	if err.Error() != "invalid pageSize" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBookingService_ListAll_Success(t *testing.T) {
	svc := NewBookingService(&bookingRepoMock{
		countAllFn: func(ctx context.Context) (int, error) {
			return 3, nil
		},
		listAllPaginatedFn: func(ctx context.Context, limit, offset int) ([]domain.Booking, error) {
			return []domain.Booking{
				{ID: "b1"},
				{ID: "b2"},
			}, nil
		},
	})

	bookings, pagination, err := svc.ListAll(context.Background(), 1, 2)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(bookings) != 2 {
		t.Fatalf("expected 2 bookings")
	}
	if pagination.Total != 3 {
		t.Fatalf("expected total=3, got %d", pagination.Total)
	}
	if pagination.TotalPages != 2 {
		t.Fatalf("expected totalPages=2, got %d", pagination.TotalPages)
	}
}
