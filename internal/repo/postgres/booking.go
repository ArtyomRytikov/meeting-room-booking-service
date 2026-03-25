package postgres

import (
	"context"
	"errors"
	"time"

	"test-backend-1-ArtyomRytikov/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BookingRepository struct {
	pool *pgxpool.Pool
}

func NewBookingRepository(pool *pgxpool.Pool) *BookingRepository {
	return &BookingRepository{pool: pool}
}

func (r *BookingRepository) Init(ctx context.Context) error {
	queries := []string{
		`
		CREATE TABLE IF NOT EXISTS bookings (
			id UUID PRIMARY KEY,
			slot_id UUID NOT NULL REFERENCES slots(id) ON DELETE CASCADE,
			room_id UUID NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
			user_id UUID NOT NULL,
			status TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			cancelled_at TIMESTAMPTZ NULL
		);
		`,
		`
		CREATE UNIQUE INDEX IF NOT EXISTS ux_bookings_slot_active
		ON bookings(slot_id)
		WHERE status = 'active';
		`,
		`
		CREATE INDEX IF NOT EXISTS ix_bookings_user_status
		ON bookings(user_id, status);
		`,
		`
		CREATE INDEX IF NOT EXISTS ix_bookings_room_status
		ON bookings(room_id, status);
		`,
	}
	for _, q := range queries {
		if _, err := r.pool.Exec(ctx, q); err != nil {
			return err
		}
	}
	return nil
}

func (r *BookingRepository) FindSlot(ctx context.Context, slotID string) (*domain.SlotDetails, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, room_id, start_at, end_at
		FROM slots
		WHERE id = $1
	`, slotID)

	var s domain.SlotDetails
	if err := row.Scan(&s.SlotID, &s.RoomID, &s.StartAt, &s.EndAt); err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *BookingRepository) Create(ctx context.Context, userID, slotID string) (domain.Booking, error) {
	slot, err := r.FindSlot(ctx, slotID)
	if err != nil {
		return domain.Booking{}, err
	}

	id := uuid.NewString()
	createdAt := time.Now().UTC()

	_, err = r.pool.Exec(ctx, `
		INSERT INTO bookings (id, slot_id, room_id, user_id, status, created_at)
		VALUES ($1, $2, $3, $4, 'active', $5)
	`, id, slot.SlotID, slot.RoomID, userID, createdAt)
	if err != nil {
		return domain.Booking{}, err
	}

	return domain.Booking{
		ID:        id,
		SlotID:    slot.SlotID,
		RoomID:    slot.RoomID,
		UserID:    userID,
		Status:    "active",
		Start:     slot.StartAt.UTC().Format(time.RFC3339),
		End:       slot.EndAt.UTC().Format(time.RFC3339),
		CreatedAt: createdAt.Format(time.RFC3339),
	}, nil
}

func (r *BookingRepository) FindByID(ctx context.Context, bookingID string) (*domain.Booking, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT b.id, b.slot_id, b.room_id, b.user_id, b.status, s.start_at, s.end_at, b.created_at
		FROM bookings b
		JOIN slots s ON s.id = b.slot_id
		WHERE b.id = $1
	`, bookingID)

	var booking domain.Booking
	var startAt time.Time
	var endAt time.Time
	var createdAt time.Time

	if err := row.Scan(
		&booking.ID,
		&booking.SlotID,
		&booking.RoomID,
		&booking.UserID,
		&booking.Status,
		&startAt,
		&endAt,
		&createdAt,
	); err != nil {
		return nil, err
	}

	booking.Start = startAt.UTC().Format(time.RFC3339)
	booking.End = endAt.UTC().Format(time.RFC3339)
	booking.CreatedAt = createdAt.UTC().Format(time.RFC3339)

	return &booking, nil
}

func (r *BookingRepository) Cancel(ctx context.Context, bookingID string) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE bookings
		SET status = 'cancelled', cancelled_at = NOW()
		WHERE id = $1 AND status = 'active'
	`, bookingID)
	return err
}

func (r *BookingRepository) ListMy(ctx context.Context, userID string) ([]domain.Booking, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT b.id, b.slot_id, b.room_id, b.user_id, b.status, s.start_at, s.end_at, b.created_at
		FROM bookings b
		JOIN slots s ON s.id = b.slot_id
		WHERE b.user_id = $1
		  AND s.start_at >= NOW()
		ORDER BY s.start_at ASC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanBookings(rows)
}

func (r *BookingRepository) CountAll(ctx context.Context) (int, error) {
	var total int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM bookings`).Scan(&total)
	return total, err
}

func (r *BookingRepository) ListAllPaginated(ctx context.Context, limit, offset int) ([]domain.Booking, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT b.id, b.slot_id, b.room_id, b.user_id, b.status, s.start_at, s.end_at, b.created_at
		FROM bookings b
		JOIN slots s ON s.id = b.slot_id
		ORDER BY s.start_at ASC
		LIMIT $1 OFFSET $2
	`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanBookings(rows)
}

type bookingScanner interface {
	Next() bool
	Scan(dest ...any) error
	Err() error
}

func scanBookings(rows bookingScanner) ([]domain.Booking, error) {
	result := make([]domain.Booking, 0)

	for rows.Next() {
		var booking domain.Booking
		var startAt time.Time
		var endAt time.Time
		var createdAt time.Time

		if err := rows.Scan(
			&booking.ID,
			&booking.SlotID,
			&booking.RoomID,
			&booking.UserID,
			&booking.Status,
			&startAt,
			&endAt,
			&createdAt,
		); err != nil {
			return nil, err
		}

		booking.Start = startAt.UTC().Format(time.RFC3339)
		booking.End = endAt.UTC().Format(time.RFC3339)
		booking.CreatedAt = createdAt.UTC().Format(time.RFC3339)

		result = append(result, booking)
	}

	return result, rows.Err()
}

var ErrNotFound = errors.New("not found")
