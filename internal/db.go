package internal

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

// Event: database struct from postgres
type EventDB struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Description *string   `json:"description" db:"description"`
	StartTime   time.Time `json:"start_time" db:"start_time"`
	EndTime     time.Time `json:"end_time" db:"end_time"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type EventRepository struct {
	db *sql.DB
}

// NewEventRepository creates a new event repository
func NewEventRepository(db *sql.DB) *EventRepository {
	return &EventRepository{db: db}
}

// CreateEvent inserts a new event into the database
func (r *EventRepository) CreateEvent(ctx context.Context, event EventDB) (*EventDB, error) {
	query := `
		INSERT INTO events (title, description, start_time, end_time) 
		VALUES ($1, $2, $3, $4) 
		RETURNING id, title, description, start_time, end_time, created_at, updated_at`

	row := r.db.QueryRowContext(ctx, query, event.Title, event.Description, event.StartTime, event.EndTime)

	var createdEvent EventDB
	err := row.Scan(
		&createdEvent.ID,
		&createdEvent.Title,
		&createdEvent.Description,
		&createdEvent.StartTime,
		&createdEvent.EndTime,
		&createdEvent.CreatedAt,
		&createdEvent.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create event: %w", err)
	}

	log.Printf("Event created successfully with ID: %s", createdEvent.ID)
	return &createdEvent, nil
}

// GetEvents retrieves all events from the database
func (r *EventRepository) GetEvents(ctx context.Context) ([]EventDB, error) {
	query := `
		SELECT id, title, description, start_time, end_time, created_at, updated_at 
		FROM events 
		ORDER BY start_time ASC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query events: %w", err)
	}
	defer rows.Close()

	var events []EventDB
	for rows.Next() {
		var event EventDB
		err := rows.Scan(
			&event.ID,
			&event.Title,
			&event.Description,
			&event.StartTime,
			&event.EndTime,
			&event.CreatedAt,
			&event.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan event: %w", err)
		}
		events = append(events, event)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating events: %w", err)
	}

	log.Printf("Retrieved %d events", len(events))
	return events, nil
}

// GetEventByID retrieves a specific event by ID
func (r *EventRepository) GetEventByID(ctx context.Context, id uuid.UUID) (*EventDB, error) {
	query := `
		SELECT id, title, description, start_time, end_time, created_at, updated_at 
		FROM events 
		WHERE id = $1`

	row := r.db.QueryRowContext(ctx, query, id)

	var event EventDB
	err := row.Scan(
		&event.ID,
		&event.Title,
		&event.Description,
		&event.StartTime,
		&event.EndTime,
		&event.CreatedAt,
		&event.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("event not found")
		}
		return nil, fmt.Errorf("failed to get event by ID: %w", err)
	}

	return &event, nil
}
