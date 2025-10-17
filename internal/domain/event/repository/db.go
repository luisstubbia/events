package repository

import (
	"context"
	"database/sql"
	"errors"
	"event-api/internal/domain/model"
	"fmt"
	"time"

	"github.com/lib/pq"
)

const (
	InsertEvent    = `INSERT INTO events (id, title, description, start_time, end_time) VALUES ($1, $2, $3, $4, $5) RETURNING id, title, description, start_time, end_time, created_at`
	SelectEvent    = `SELECT id, title, description, start_time, end_time, created_at FROM events WHERE id = $1`
	QueryAllEvents = `SELECT id, title, description, start_time, end_time, created_at FROM events ORDER BY start_time ASC`
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) RetrieveById(ctx context.Context, id string) (*model.Event, error) {
	var event model.Event
	err := r.db.QueryRowContext(ctx, SelectEvent, id).Scan(
		&event.ID,
		&event.Title,
		&event.Description,
		&event.StartTime,
		&event.EndTime,
		&event.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.NewNotFoundError(fmt.Sprintf("event by id: %s was not found", id), time.Now())
		}
		return nil, fmt.Errorf("failed to get event: %w", err)
	}
	return &event, nil
}

func (r *Repository) Save(ctx context.Context, event *model.Event) error {
	_, err := r.db.ExecContext(ctx, InsertEvent,
		event.ID,
		event.Title,
		event.Description,
		event.StartTime,
		event.EndTime,
	)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23505" {
				return model.NewValidationErrorWithTime("Event already exists. Duplicated event id")
			}
			return fmt.Errorf("postgres error: %s (code %s)", pqErr.Message, pqErr.Code)
		}
		return fmt.Errorf("failed to create event: %w", err)
	}
	return nil
}

func (r *Repository) RetrieveAll(ctx context.Context) ([]model.Event, error) {
	rows, err := r.db.QueryContext(ctx, QueryAllEvents)
	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			fmt.Printf("failed to close rows: %s", err)
		}
	}(rows)

	var events []model.Event
	for rows.Next() {
		var event model.Event
		var description sql.NullString

		err := rows.Scan(
			&event.ID,
			&event.Title,
			&description,
			&event.StartTime,
			&event.EndTime,
			&event.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	if len(events) == 0 {
		return nil, model.NewNotFoundError("no events found", time.Now())
	}
	return events, nil
}
