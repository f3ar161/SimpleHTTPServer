package internal

import (
	"context"

	"github.com/google/uuid"
)

// EventRepositoryInterface defines the contract for event repository operations.
// This interface abstracts the database operations, allowing for easier testing
type EventRepositoryInterface interface {
	CreateEvent(ctx context.Context, event EventDB) (*EventDB, error)
	GetEvents(ctx context.Context) ([]EventDB, error)
	GetEventByID(ctx context.Context, id uuid.UUID) (*EventDB, error)
}
