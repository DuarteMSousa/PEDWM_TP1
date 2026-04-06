package events_infrastructure

import (
	"log/slog"
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
		slog.Warn("event persistence observer failed: missing eventService")
		return
	}

	gameId := strings.TrimSpace(event.GameID)

	if gameId != "" {
		err := o.eventService.SaveEvent(event)
		if err != nil {
			slog.Error("failed to save event to database", "eventType", event.Type, "gameID", gameId, "error", err)
		}
	} else {
		slog.Debug("skipping event persistence: no gameID", "eventType", event.Type)
	}

	if o.eventDispatcher != nil {
		o.eventDispatcher.Dispatch(event)
	}
}
