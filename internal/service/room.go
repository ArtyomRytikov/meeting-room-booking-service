package service

import (
	"context"
	"errors"
	"strings"

	"test-backend-1-ArtyomRytikov/internal/domain"
)

type RoomRepository interface {
	Init(ctx context.Context) error
	Create(ctx context.Context, room domain.Room) (domain.Room, error)
	List(ctx context.Context) ([]domain.Room, error)
}

type RoomService struct {
	repo RoomRepository
}

func NewRoomService(repo RoomRepository) *RoomService {
	return &RoomService{repo: repo}
}

func (s *RoomService) Init(ctx context.Context) error {
	return s.repo.Init(ctx)
}

func (s *RoomService) Create(ctx context.Context, room domain.Room) (domain.Room, error) {
	if strings.TrimSpace(room.Name) == "" {
		return domain.Room{}, errors.New("name is required")
	}

	if room.Capacity != nil && *room.Capacity < 0 {
		return domain.Room{}, errors.New("capacity must be >= 0")
	}

	return s.repo.Create(ctx, room)
}

func (s *RoomService) List(ctx context.Context) ([]domain.Room, error) {
	return s.repo.List(ctx)
}
