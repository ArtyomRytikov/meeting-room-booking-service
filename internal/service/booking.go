package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"test-backend-1-ArtyomRytikov/internal/domain"
)

type BookingRepository interface {
	Init(ctx context.Context) error
	FindSlot(ctx context.Context, slotID string) (*domain.SlotDetails, error)
	Create(ctx context.Context, userID, slotID string) (domain.Booking, error)
	FindByID(ctx context.Context, bookingID string) (*domain.Booking, error)
	Cancel(ctx context.Context, bookingID string) error
	ListMy(ctx context.Context, userID string) ([]domain.Booking, error)
	CountAll(ctx context.Context) (int, error)
	ListAllPaginated(ctx context.Context, limit, offset int) ([]domain.Booking, error)
}

type Pagination struct {
	Page       int `json:"page"`
	PageSize   int `json:"pageSize"`
	Total      int `json:"total"`
	TotalPages int `json:"totalPages"`
}

type BookingService struct {
	repo BookingRepository
}

func NewBookingService(repo BookingRepository) *BookingService {
	return &BookingService{repo: repo}
}

func (s *BookingService) Init(ctx context.Context) error {
	return s.repo.Init(ctx)
}

func (s *BookingService) Create(ctx context.Context, userID, slotID string) (domain.Booking, error) {
	if strings.TrimSpace(slotID) == "" {
		return domain.Booking{}, errors.New("slotId is required")
	}

	slot, err := s.repo.FindSlot(ctx, slotID)
	if err != nil {
		msg := err.Error()
		if strings.Contains(msg, "no rows in result set") {
			return domain.Booking{}, errors.New("slot not found")
		}
		return domain.Booking{}, err
	}

	if slot.StartAt.Before(time.Now().UTC()) {
		return domain.Booking{}, errors.New("slot is in the past")
	}

	booking, err := s.repo.Create(ctx, userID, slotID)
	if err != nil {
		msg := err.Error()
		switch {
		case strings.Contains(msg, "ux_bookings_slot_active"):
			return domain.Booking{}, errors.New("slot already booked")
		default:
			return domain.Booking{}, err
		}
	}

	return booking, nil
}

func (s *BookingService) Cancel(ctx context.Context, bookingID, requesterID, requesterRole string) (*domain.Booking, error) {
	booking, err := s.repo.FindByID(ctx, bookingID)
	if err != nil {
		return nil, errors.New("booking not found")
	}

	if booking.UserID != requesterID {
		return nil, errors.New("forbidden")
	}

	if booking.Status == "cancelled" {
		return booking, nil
	}

	if err := s.repo.Cancel(ctx, bookingID); err != nil {
		return nil, err
	}

	return s.repo.FindByID(ctx, bookingID)
}

func (s *BookingService) ListMy(ctx context.Context, userID string) ([]domain.Booking, error) {
	return s.repo.ListMy(ctx, userID)
}

func (s *BookingService) ListAll(ctx context.Context, page, pageSize int) ([]domain.Booking, Pagination, error) {
	if page < 1 {
		return nil, Pagination{}, errors.New("invalid page")
	}
	if pageSize < 1 || pageSize > 100 {
		return nil, Pagination{}, errors.New("invalid pageSize")
	}

	total, err := s.repo.CountAll(ctx)
	if err != nil {
		return nil, Pagination{}, err
	}

	offset := (page - 1) * pageSize
	bookings, err := s.repo.ListAllPaginated(ctx, pageSize, offset)
	if err != nil {
		return nil, Pagination{}, err
	}

	totalPages := 0
	if total > 0 {
		totalPages = (total + pageSize - 1) / pageSize
	}

	return bookings, Pagination{
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}
