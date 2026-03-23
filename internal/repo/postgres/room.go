package postgres

import (
	"context"

	"test-backend-1-ArtyomRytikov/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RoomRepository struct {
	pool *pgxpool.Pool
}

func NewRoomRepository(pool *pgxpool.Pool) *RoomRepository {
	return &RoomRepository{pool: pool}
}

func (r *RoomRepository) Init(ctx context.Context) error {
	query := `
	CREATE TABLE IF NOT EXISTS rooms (
		id UUID PRIMARY KEY,
		name TEXT NOT NULL,
		description TEXT NULL,
		capacity INT NULL
	);
	`
	_, err := r.pool.Exec(ctx, query)
	return err
}

func (r *RoomRepository) Create(ctx context.Context, room domain.Room) (domain.Room, error) {
	id := uuid.NewString()

	_, err := r.pool.Exec(
		ctx,
		`INSERT INTO rooms (id, name, description, capacity) VALUES ($1, $2, $3, $4)`,
		id,
		room.Name,
		room.Description,
		room.Capacity,
	)
	if err != nil {
		return domain.Room{}, err
	}

	room.ID = id
	return room, nil
}

func (r *RoomRepository) List(ctx context.Context) ([]domain.Room, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, name, description, capacity
		FROM rooms
		ORDER BY name ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []domain.Room

	for rows.Next() {
		var room domain.Room
		var description *string
		var capacity *int

		if err := rows.Scan(&room.ID, &room.Name, &description, &capacity); err != nil {
			return nil, err
		}

		room.Description = description
		room.Capacity = capacity

		rooms = append(rooms, room)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return rooms, nil
}
