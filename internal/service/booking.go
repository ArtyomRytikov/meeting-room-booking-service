package service

import (
	"context"
	"errors"
	"strings"

	"test-backend-1-ArtyomRytikov/internal/domain"
)

type BookingRepository interface {
	Init(ctx context.Context) error
	Create(ctx context.Context, userID, slotID string) (domain.Booking, error)
	FindByID(ctx context.Context, bookingID string) (*domain.Booking, error)
	Cancel(ctx context.Context, bookingID string) error
	ListMy(ctx context.Context, userID string) ([]domain.Booking, error)
	ListAll(ctx context.Context) ([]domain.Booking, error)
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

	booking, err := s.repo.Create(ctx, userID, slotID)
	if err != nil {
		msg := err.Error()
		switch {
		case strings.Contains(msg, "no rows in result set"):
			return domain.Booking{}, errors.New("slot not found")
		case strings.Contains(msg, "ux_bookings_slot_active"):
			return domain.Booking{}, errors.New("slot already booked")
		default:
			return domain.Booking{}, err
		}
	}

	return booking, nil
}

func (s *BookingService) Cancel(ctx context.Context, bookingID, requesterID, requesterRole string) error {
	booking, err := s.repo.FindByID(ctx, bookingID)
	if err != nil {
		return errors.New("booking not found")
	}

	if requesterRole != "admin" && booking.UserID != requesterID {
		return errors.New("forbidden")
	}

	if booking.Status != "active" {
		return errors.New("booking not found or already cancelled")
	}

	return s.repo.Cancel(ctx, bookingID)
}

func (s *BookingService) ListMy(ctx context.Context, userID string) ([]domain.Booking, error) {
	return s.repo.ListMy(ctx, userID)
}

func (s *BookingService) ListAll(ctx context.Context) ([]domain.Booking, error) {
	return s.repo.ListAll(ctx)
}
