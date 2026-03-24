package postgres

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	"test-backend-1-ArtyomRytikov/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ScheduleRepository struct {
	pool *pgxpool.Pool
}

func NewScheduleRepository(pool *pgxpool.Pool) *ScheduleRepository {
	return &ScheduleRepository{pool: pool}
}

func (r *ScheduleRepository) Init(ctx context.Context) error {
	queries := []string{
		`
		CREATE TABLE IF NOT EXISTS schedules (
			id UUID PRIMARY KEY,
			room_id UUID NOT NULL UNIQUE REFERENCES rooms(id) ON DELETE CASCADE,
			days_of_week INT[] NOT NULL,
			start_time TEXT NOT NULL,
			end_time TEXT NOT NULL
		);
		`,
		`
		CREATE TABLE IF NOT EXISTS slots (
			id UUID PRIMARY KEY,
			room_id UUID NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
			schedule_id UUID NOT NULL REFERENCES schedules(id) ON DELETE CASCADE,
			start_at TIMESTAMPTZ NOT NULL,
			end_at TIMESTAMPTZ NOT NULL
		);
		`,
		`
		CREATE UNIQUE INDEX IF NOT EXISTS ux_slots_room_time
		ON slots(room_id, start_at, end_at);
		`,
		`
		CREATE INDEX IF NOT EXISTS ix_slots_room_start
		ON slots(room_id, start_at);
		`,
	}

	for _, q := range queries {
		if _, err := r.pool.Exec(ctx, q); err != nil {
			return err
		}
	}

	return nil
}

func (r *ScheduleRepository) RoomExists(ctx context.Context, roomID string) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM rooms WHERE id = $1)`, roomID).Scan(&exists)
	return exists, err
}

func (r *ScheduleRepository) ScheduleExists(ctx context.Context, roomID string) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM schedules WHERE room_id = $1)`, roomID).Scan(&exists)
	return exists, err
}

func (r *ScheduleRepository) CreateScheduleWithSlots(
	ctx context.Context,
	schedule domain.Schedule,
	days int,
) (domain.Schedule, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return domain.Schedule{}, err
	}
	defer tx.Rollback(ctx)

	schedule.ID = uuid.NewString()

	_, err = tx.Exec(ctx, `
		INSERT INTO schedules (id, room_id, days_of_week, start_time, end_time)
		VALUES ($1, $2, $3, $4, $5)
	`,
		schedule.ID,
		schedule.RoomID,
		schedule.DaysOfWeek,
		schedule.StartTime,
		schedule.EndTime,
	)
	if err != nil {
		return domain.Schedule{}, err
	}

	startHour, startMinute, err := parseHHMM(schedule.StartTime)
	if err != nil {
		return domain.Schedule{}, err
	}
	endHour, endMinute, err := parseHHMM(schedule.EndTime)
	if err != nil {
		return domain.Schedule{}, err
	}

	now := time.Now().UTC()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	for i := 0; i < days; i++ {
		day := today.AddDate(0, 0, i)
		apiWeekday := toAPIWeekday(day.Weekday())
		if !slices.Contains(schedule.DaysOfWeek, apiWeekday) {
			continue
		}

		slotStart := time.Date(day.Year(), day.Month(), day.Day(), startHour, startMinute, 0, 0, time.UTC)
		slotEndBoundary := time.Date(day.Year(), day.Month(), day.Day(), endHour, endMinute, 0, 0, time.UTC)

		for current := slotStart; current.Before(slotEndBoundary); current = current.Add(30 * time.Minute) {
			next := current.Add(30 * time.Minute)
			if next.After(slotEndBoundary) {
				break
			}

			_, err := tx.Exec(ctx, `
				INSERT INTO slots (id, room_id, schedule_id, start_at, end_at)
				VALUES ($1, $2, $3, $4, $5)
				ON CONFLICT (room_id, start_at, end_at) DO NOTHING
			`,
				uuid.NewString(),
				schedule.RoomID,
				schedule.ID,
				current,
				next,
			)
			if err != nil {
				return domain.Schedule{}, err
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return domain.Schedule{}, err
	}

	return schedule, nil
}

func (r *ScheduleRepository) ListSlotsByDate(ctx context.Context, roomID string, date time.Time) ([]domain.Slot, error) {
	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	end := start.Add(24 * time.Hour)

	rows, err := r.pool.Query(ctx, `
		SELECT s.id, s.room_id, s.start_at, s.end_at
		FROM slots s
		WHERE s.room_id = $1
		  AND s.start_at >= $2
		  AND s.start_at < $3
		ORDER BY s.start_at ASC
	`, roomID, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.Slot
	for rows.Next() {
		var slot domain.Slot
		var startAt time.Time
		var endAt time.Time

		if err := rows.Scan(&slot.ID, &slot.RoomID, &startAt, &endAt); err != nil {
			return nil, err
		}

		slot.Start = startAt.UTC().Format(time.RFC3339)
		slot.End = endAt.UTC().Format(time.RFC3339)

		result = append(result, slot)
	}

	return result, rows.Err()
}

func parseHHMM(value string) (int, int, error) {
	t, err := time.Parse("15:04", value)
	if err != nil {
		return 0, 0, err
	}
	return t.Hour(), t.Minute(), nil
}

func toAPIWeekday(w time.Weekday) int {
	if w == time.Sunday {
		return 7
	}
	return int(w)
}

var ErrInvalidTimeRange = errors.New("invalid time range")

func ValidateTimeRange(startTime, endTime string) error {
	start, err := time.Parse("15:04", startTime)
	if err != nil {
		return err
	}
	end, err := time.Parse("15:04", endTime)
	if err != nil {
		return err
	}

	if !start.Before(end) {
		return ErrInvalidTimeRange
	}

	diff := end.Sub(start)
	if diff%(30*time.Minute) != 0 {
		return fmt.Errorf("time range must be divisible by 30 minutes")
	}

	return nil
}
