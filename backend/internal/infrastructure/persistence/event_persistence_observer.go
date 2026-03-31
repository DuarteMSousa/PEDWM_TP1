package persistence

import (
	"strings"

	"backend/internal/application/interfaces"
	domain "backend/internal/domain/events"
)

// Minimal bridge EventBus -> WebSocket room broadcast.
type EventPersistanceObserver struct {
	eventService interfaces.EventService
}

func NewEventPersistanceObserver(eventService interfaces.EventService) *EventPersistanceObserver {
	return &EventPersistanceObserver{eventService: eventService}
}

func (o *EventPersistanceObserver) Update(event domain.Event) {
	if o == nil || o.eventService == nil {
		return
	}

	gameId := strings.TrimSpace(event.GameID)

	if gameId == "" {
		return
	}

	o.eventService.SaveEvent(event)
}
