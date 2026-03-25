package service

import (
	"context"
	"errors"
	"testing"

	"test-backend-1-ArtyomRytikov/internal/domain"
)

type roomRepoMock struct {
	createFn func(ctx context.Context, room domain.Room) (domain.Room, error)
	listFn   func(ctx context.Context) ([]domain.Room, error)
	initFn   func(ctx context.Context) error
}

func (m *roomRepoMock) Init(ctx context.Context) error {
	if m.initFn != nil {
		return m.initFn(ctx)
	}
	return nil
}

func (m *roomRepoMock) Create(ctx context.Context, room domain.Room) (domain.Room, error) {
	if m.createFn != nil {
		return m.createFn(ctx, room)
	}
	return domain.Room{}, nil
}

func (m *roomRepoMock) List(ctx context.Context) ([]domain.Room, error) {
	if m.listFn != nil {
		return m.listFn(ctx)
	}
	return []domain.Room{}, nil
}

func TestRoomService_Create_Success(t *testing.T) {
	repo := &roomRepoMock{
		createFn: func(ctx context.Context, room domain.Room) (domain.Room, error) {
			room.ID = "room-1"
			return room, nil
		},
	}

	svc := NewRoomService(repo)

	capacity := 6
	room, err := svc.Create(context.Background(), domain.Room{
		Name:     "Room A",
		Capacity: &capacity,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if room.ID != "room-1" {
		t.Fatalf("expected room id to be set")
	}
	if room.Name != "Room A" {
		t.Fatalf("unexpected room name: %s", room.Name)
	}
}

func TestRoomService_Create_EmptyName(t *testing.T) {
	svc := NewRoomService(&roomRepoMock{})

	_, err := svc.Create(context.Background(), domain.Room{Name: ""})
	if err == nil {
		t.Fatal("expected error for empty name")
	}
	if err.Error() != "name is required" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRoomService_Create_NegativeCapacity(t *testing.T) {
	svc := NewRoomService(&roomRepoMock{})

	capacity := -1
	_, err := svc.Create(context.Background(), domain.Room{
		Name:     "Room A",
		Capacity: &capacity,
	})
	if err == nil {
		t.Fatal("expected error for negative capacity")
	}
	if err.Error() != "capacity must be >= 0" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRoomService_List(t *testing.T) {
	repo := &roomRepoMock{
		listFn: func(ctx context.Context) ([]domain.Room, error) {
			return []domain.Room{
				{ID: "1", Name: "Room A"},
				{ID: "2", Name: "Room B"},
			}, nil
		},
	}
	svc := NewRoomService(repo)

	rooms, err := svc.List(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(rooms) != 2 {
		t.Fatalf("expected 2 rooms, got %d", len(rooms))
	}
}

func TestRoomService_Init_Error(t *testing.T) {
	svc := NewRoomService(&roomRepoMock{
		initFn: func(ctx context.Context) error {
			return errors.New("init failed")
		},
	})

	err := svc.Init(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
}
