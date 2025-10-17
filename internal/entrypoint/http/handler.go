package http

import (
	"context"
	"encoding/json"
	"event-api/internal/domain/model"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Service interface {
	CreateEvent(context context.Context, event *model.Event) error
	GetEvent(context context.Context, id string) (*model.Event, error)
	GetAllEvents(context context.Context) ([]model.Event, error)
}

type EventHandler struct {
	svc Service
}

func NewEventHandler(svc Service) *EventHandler {
	return &EventHandler{svc: svc}
}

// CreateEvent godoc
// @Summary Create a new event
// @Description Create a new event with the provided details
// @Tags events
// @Accept json
// @Produce json
// @Param event body CreateEventRequest true "Event creation data"
// @Success 201 {object} model.Event "Event created successfully"
// @Failure 400 {object} ErrorResponse "Invalid input data"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /events [post]
func (h *EventHandler) createEvent(w http.ResponseWriter, r *http.Request) {
	var req CreateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handlerErrorResponse(w, err)
		return
	}

	if errors := req.Validate(); len(errors) > 0 {
		respondWithJSON(w, http.StatusBadRequest, errors)
		return
	}

	event := &model.Event{
		Title:       req.Title,
		Description: req.Description,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		CreatedAt:   time.Now().UTC(),
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	if err := h.svc.CreateEvent(ctx, event); err != nil {
		handlerErrorResponse(w, err)
		return
	}

	respondWithJSON(w, http.StatusCreated, event)
}

// getEvent godoc
// @Summary Get event by ID
// @Description Get a specific event by its UUID
// @Tags events
// @Accept json
// @Produce json
// @Param id path string true "Event ID" Format(uuid)
// @Success 200 {object} model.Event "Event found"
// @Failure 400 {object} ErrorResponse "Invalid UUID format"
// @Failure 404 {object} ErrorResponse "Event not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /events/{id} [get]
func (h *EventHandler) getEvent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if _, err := uuid.Parse(id); err != nil {
		handlerErrorResponse(w, model.NewValidationErrorWithTime(fmt.Sprintf("Invalid UUID format: %s", err.Error())))
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	event, err := h.svc.GetEvent(ctx, id)
	if err != nil {
		handlerErrorResponse(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, event)
}

// getAllEvents godoc
// @Summary List all events
// @Description Get all events ordered by start time in ascending order. MISSING PAGINATION!!!
// @Tags events
// @Accept json
// @Produce json
// @Success 200 {array} []model.Event "List of events"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /events [get]
func (h *EventHandler) getAllEvents(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	events, err := h.svc.GetAllEvents(ctx)
	if err != nil {
		handlerErrorResponse(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, events)
}

// HealthCheck godoc
// @Summary Health check endpoint
// @Description Check if the service is healthy
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string "Service status"
// @Failure 503 {object} ErrorResponse "Service unavailable"
// @Router /health [get]
func (h *EventHandler) healthCheck(w http.ResponseWriter, _ *http.Request) {
	respondWithJSON(w, http.StatusOK, map[string]string{"status": "healthy"})
}
