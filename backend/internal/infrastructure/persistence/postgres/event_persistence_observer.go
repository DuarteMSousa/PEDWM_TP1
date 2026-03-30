package postgres

import (
	"strings"

	application "backend/internal/application/services"
	domain "backend/internal/domain/events"
)

// Minimal bridge EventBus -> WebSocket room broadcast.
type EventPersistanceObserver struct {
	eventService *application.EventService
}

func NewEventPersistanceObserver(eventService *application.EventService) *EventPersistanceObserver {
	return &EventPersistanceObserver{eventService: eventService}
}

func (o *EventPersistanceObserver) Update(event domain.Event) {
	if o == nil || o.eventService == nil {
		return
	}

	roomID := strings.TrimSpace(event.RoomID)

	if roomID == "" {
		return
	}

	o.eventService.SaveEvent(event)
}
