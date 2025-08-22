package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"taller_challenge/internal"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// EventController handles HTTP requests for events
type EventController struct {
	eventRepo internal.EventRepositoryInterface
}

// NewEventController creates a new event controller
func NewEventController(eventRepo internal.EventRepositoryInterface) *EventController {
	return &EventController{
		eventRepo: eventRepo,
	}
}

type createEventInput struct {
	Title       string    `json:"title"`
	Description *string   `json:"description"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
}

// CreateEvent handles POST /events
func (ec *EventController) CreateEvent(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	var in createEventInput
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&in); err != nil {
		http.Error(w, fmt.Sprintf("invalid JSON: %v", err), http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(in.Title) == "" {
		http.Error(w, "title is required", http.StatusBadRequest)
		return
	}
	if len(in.Title) > 100 {
		http.Error(w, "title must be <= 100 characters", http.StatusBadRequest)
		return
	}
	if in.StartTime.IsZero() || in.EndTime.IsZero() {
		http.Error(w, "start_time and end_time are required (RFC3339)", http.StatusBadRequest)
		return
	}
	if !in.StartTime.Before(in.EndTime) {
		http.Error(w, "start_time must be before end_time", http.StatusBadRequest)
		return
	}

	id := uuid.New()
	createdAt := time.Now().UTC()

	event := internal.EventDB{
		ID:          id,
		Title:       in.Title,
		Description: in.Description,
		StartTime:   in.StartTime.UTC(),
		EndTime:     in.EndTime.UTC(),
		CreatedAt:   createdAt,
		UpdatedAt:   createdAt,
	}

	createdEvent, err := ec.eventRepo.CreateEvent(ctx, event)
	if err != nil {
		log.Printf("Error creating event: %v", err)
		if ctx.Err() == context.DeadlineExceeded {
			http.Error(w, "Request timeout", http.StatusRequestTimeout)
			return
		}
		http.Error(w, "Failed to create event", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdEvent)
}

// GetEvents handles GET /events
func (ec *EventController) GetEvents(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	events, err := ec.eventRepo.GetEvents(ctx)
	if err != nil {
		log.Printf("Error getting events: %v", err)
		if ctx.Err() == context.DeadlineExceeded {
			http.Error(w, "Request timeout", http.StatusRequestTimeout)
			return
		}
		http.Error(w, "Failed to get events", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}

// GetEventByID handles GET /events/{id}
func (ec *EventController) GetEventByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid UUID format", http.StatusBadRequest)
		return
	}

	event, err := ec.eventRepo.GetEventByID(ctx, id)
	if err != nil {
		log.Printf("Error getting event by ID: %v", err)
		http.Error(w, "Event not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(event)
}

// SetupRoutes configures the HTTP routes
func (ec *EventController) SetupRoutes() *mux.Router {
	router := mux.NewRouter()

	// Events endpoints
	router.HandleFunc("/events", ec.CreateEvent).Methods("POST")
	router.HandleFunc("/events", ec.GetEvents).Methods("GET")
	router.HandleFunc("/events/{id}", ec.GetEventByID).Methods("GET")

	return router
}

// StartServer starts the HTTP server with graceful shutdown
func StartServer(eventRepo internal.EventRepositoryInterface, port string) {
	controller := NewEventController(eventRepo)
	router := controller.SetupRoutes()

	router.Use(loggingMiddleware)

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on port %s", port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Server is shutting down...")

	// Create a context with timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

// loggingMiddleware logs incoming HTTP requests
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %v", r.Method, r.RequestURI, time.Since(start))
	})
}
