package events_infrastructure

import (
	"strings"

	"backend/internal/application/interfaces"
	domain "backend/internal/domain/events"
)

// EventPersistanceObserver implements Observer and persists domain events
// in the database, forwarding them to the EventDispatcher.
type EventPersistanceObserver struct {
	eventService    interfaces.EventService
	eventDispatcher *EventDispatcher
}

// NewEventPersistanceObserver creates a new persistence observer.
func NewEventPersistanceObserver(eventService interfaces.EventService, eventDispatcher *EventDispatcher) *EventPersistanceObserver {
	return &EventPersistanceObserver{eventService: eventService, eventDispatcher: eventDispatcher}
}

// Update persists the event and forwards it to the dispatcher.
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
