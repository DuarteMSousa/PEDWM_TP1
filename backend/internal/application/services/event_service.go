package services

import (
	"backend/internal/application/interfaces"
	"backend/internal/domain/events"
	"log/slog"
)

// EventService manages the persistence and querying of domain events.
type EventService struct {
	repo interfaces.EventRepository
}

// NewEventService creates a new EventService.
func NewEventService(repo interfaces.EventRepository) *EventService {
	return &EventService{repo: repo}
}

// SaveEvent persists a domain event.
func (s *EventService) SaveEvent(event events.Event) error {
	slog.Debug("persisting event", "eventID", event.ID, "type", event.Type, "gameID", event.GameID)
	return s.repo.Save(event)
}

// GetEventsByRoomID returns the events associated with a room.
func (s *EventService) GetEventsByRoomID(roomID string) ([]events.Event, error) {
	slog.Debug("getting events by room", "roomID", roomID)
	return s.repo.FindByRoomID(roomID)
}

// GetEventsByGameID returns the events associated with a game.
func (s *EventService) GetEventsByGameID(gameID string) ([]events.Event, error) {
	slog.Debug("getting events by game", "gameID", gameID)
	return s.repo.FindByGameID(gameID)
}
