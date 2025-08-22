package internal

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// MockEventRepository
type MockEventRepository struct {
	createEventFunc  func(ctx context.Context, event EventDB) (*EventDB, error)
	getEventsFunc    func(ctx context.Context) ([]EventDB, error)
	getEventByIDFunc func(ctx context.Context, id uuid.UUID) (*EventDB, error)
}

func NewMockEventRepository() *MockEventRepository {
	return &MockEventRepository{}
}

func (m *MockEventRepository) CreateEvent(ctx context.Context, event EventDB) (*EventDB, error) {
	if m.createEventFunc != nil {
		return m.createEventFunc(ctx, event)
	}
	return nil, errors.New("mock not configured")
}

func (m *MockEventRepository) GetEvents(ctx context.Context) ([]EventDB, error) {
	if m.getEventsFunc != nil {
		return m.getEventsFunc(ctx)
	}
	return nil, errors.New("mock not configured")
}

func (m *MockEventRepository) GetEventByID(ctx context.Context, id uuid.UUID) (*EventDB, error) {
	if m.getEventByIDFunc != nil {
		return m.getEventByIDFunc(ctx, id)
	}
	return nil, errors.New("mock not configured")
}

func TestCreateEvent(t *testing.T) {
	tests := []struct {
		name     string
		event    EventDB
		mockFunc func(ctx context.Context, event EventDB) (*EventDB, error)
		wantErr  bool
		errMsg   string
	}{
		{
			name: "successful creation",
			event: EventDB{
				Title:       "Test Event",
				Description: stringPtr("Test Description"),
				StartTime:   time.Now().Add(time.Hour),
				EndTime:     time.Now().Add(2 * time.Hour),
			},
			mockFunc: func(ctx context.Context, event EventDB) (*EventDB, error) {
				return &EventDB{
					ID:          uuid.New(),
					Title:       event.Title,
					Description: event.Description,
					StartTime:   event.StartTime,
					EndTime:     event.EndTime,
					CreatedAt:   time.Now().UTC(),
					UpdatedAt:   time.Now().UTC(),
				}, nil
			},
			wantErr: false,
		},
		{
			name: "database error",
			event: EventDB{
				Title:     "Test Event",
				StartTime: time.Now().Add(time.Hour),
				EndTime:   time.Now().Add(2 * time.Hour),
			},
			mockFunc: func(ctx context.Context, event EventDB) (*EventDB, error) {
				return nil, errors.New("database connection failed")
			},
			wantErr: true,
			errMsg:  "database connection failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := NewMockEventRepository()
			mockRepo.createEventFunc = tt.mockFunc

			result, err := mockRepo.CreateEvent(context.Background(), tt.event)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.NotEqual(t, uuid.Nil, result.ID)
			assert.Equal(t, tt.event.Title, result.Title)
		})
	}
}

func TestGetEvents(t *testing.T) {
	tests := []struct {
		name      string
		mockFunc  func(ctx context.Context) ([]EventDB, error)
		wantCount int
		wantErr   bool
	}{
		{
			name: "empty result",
			mockFunc: func(ctx context.Context) ([]EventDB, error) {
				return []EventDB{}, nil
			},
			wantCount: 0,
			wantErr:   false,
		},
		{
			name: "multiple events",
			mockFunc: func(ctx context.Context) ([]EventDB, error) {
				now := time.Now()
				return []EventDB{
					{ID: uuid.New(), Title: "Event 1", StartTime: now.Add(1 * time.Hour), EndTime: now.Add(2 * time.Hour)},
					{ID: uuid.New(), Title: "Event 2", StartTime: now.Add(2 * time.Hour), EndTime: now.Add(3 * time.Hour)},
				}, nil
			},
			wantCount: 2,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := NewMockEventRepository()
			mockRepo.getEventsFunc = tt.mockFunc

			result, err := mockRepo.GetEvents(context.Background())

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Len(t, result, tt.wantCount)
		})
	}
}

func TestGetEventByID(t *testing.T) {
	testID := uuid.New()

	tests := []struct {
		name     string
		id       uuid.UUID
		mockFunc func(ctx context.Context, id uuid.UUID) (*EventDB, error)
		wantErr  bool
	}{
		{
			name: "existing event",
			id:   testID,
			mockFunc: func(ctx context.Context, id uuid.UUID) (*EventDB, error) {
				if id == testID {
					return &EventDB{
						ID:        testID,
						Title:     "Test Event",
						StartTime: time.Now().Add(1 * time.Hour),
						EndTime:   time.Now().Add(2 * time.Hour),
					}, nil
				}
				return nil, errors.New("event not found")
			},
			wantErr: false,
		},
		{
			name: "non-existing event",
			id:   uuid.New(),
			mockFunc: func(ctx context.Context, id uuid.UUID) (*EventDB, error) {
				return nil, fmt.Errorf("event not found")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := NewMockEventRepository()
			mockRepo.getEventByIDFunc = tt.mockFunc

			result, err := mockRepo.GetEventByID(context.Background(), tt.id)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, tt.id, result.ID)
		})
	}
}

// Helper function to create string pointers
func stringPtr(s string) *string {
	return &s
}
