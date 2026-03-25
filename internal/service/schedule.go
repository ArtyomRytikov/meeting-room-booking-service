package service

import (
	"context"
	"errors"
	"time"

	"test-backend-1-ArtyomRytikov/internal/domain"
	postgresrepo "test-backend-1-ArtyomRytikov/internal/repo/postgres"

	"github.com/google/uuid"
)

type ScheduleRepository interface {
	Init(ctx context.Context) error
	RoomExists(ctx context.Context, roomID string) (bool, error)
	ScheduleExists(ctx context.Context, roomID string) (bool, error)
	CreateScheduleWithSlots(ctx context.Context, schedule domain.Schedule, days int) (domain.Schedule, error)
	ListSlotsByDate(ctx context.Context, roomID string, date time.Time) ([]domain.Slot, error)
}

type ScheduleService struct {
	repo ScheduleRepository
}

func NewScheduleService(repo ScheduleRepository) *ScheduleService {
	return &ScheduleService{repo: repo}
}

func (s *ScheduleService) Init(ctx context.Context) error {
	return s.repo.Init(ctx)
}

func (s *ScheduleService) Create(ctx context.Context, schedule domain.Schedule) (domain.Schedule, error) {
	if _, err := uuid.Parse(schedule.RoomID); err != nil {
		return domain.Schedule{}, errors.New("invalid room id")
	}

	exists, err := s.repo.RoomExists(ctx, schedule.RoomID)
	if err != nil {
		return domain.Schedule{}, err
	}
	if !exists {
		return domain.Schedule{}, errors.New("room not found")
	}

	hasSchedule, err := s.repo.ScheduleExists(ctx, schedule.RoomID)
	if err != nil {
		return domain.Schedule{}, err
	}
	if hasSchedule {
		return domain.Schedule{}, errors.New("schedule exists")
	}

	if len(schedule.DaysOfWeek) == 0 {
		return domain.Schedule{}, errors.New("daysOfWeek is required")
	}
	for _, d := range schedule.DaysOfWeek {
		if d < 1 || d > 7 {
			return domain.Schedule{}, errors.New("daysOfWeek must contain values from 1 to 7")
		}
	}

	if err := postgresrepo.ValidateTimeRange(schedule.StartTime, schedule.EndTime); err != nil {
		return domain.Schedule{}, err
	}

	return s.repo.CreateScheduleWithSlots(ctx, schedule, 30)
}

func (s *ScheduleService) ListSlotsByDate(ctx context.Context, roomID, rawDate string) ([]domain.Slot, error) {
	if _, err := uuid.Parse(roomID); err != nil {
		return nil, errors.New("invalid room id")
	}

	exists, err := s.repo.RoomExists(ctx, roomID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("room not found")
	}

	date, err := time.Parse("2006-01-02", rawDate)
	if err != nil {
		return nil, errors.New("invalid date")
	}

	return s.repo.ListSlotsByDate(ctx, roomID, date.UTC())
}
