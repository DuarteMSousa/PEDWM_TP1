package events_infrastructure

import (
	"strings"

	"backend/internal/application/interfaces"
	domain "backend/internal/domain/events"
)

// Minimal bridge EventBus -> WebSocket room broadcast.
type EventPersistanceObserver struct {
	eventService    interfaces.EventService
	eventDispatcher *EventDispatcher
}

func NewEventPersistanceObserver(eventService interfaces.EventService, eventDispatcher *EventDispatcher) *EventPersistanceObserver {
	return &EventPersistanceObserver{eventService: eventService, eventDispatcher: eventDispatcher}
}

func (o *EventPersistanceObserver) Update(event domain.Event) {
	if o == nil || o.eventService == nil {
		return
	}

	gameId := strings.TrimSpace(event.GameID)

	if gameId != "" {
		o.eventService.SaveEvent(event)
	}

	if o.eventDispatcher != nil {
		o.eventDispatcher.Dispatch(event)
	}
}
