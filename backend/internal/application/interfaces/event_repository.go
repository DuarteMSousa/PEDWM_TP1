package interfaces

import (
	"backend/internal/domain/events"
)

type EventRepository interface {
	Save(event events.Event) error
	FindByGameID(gameID string) ([]events.Event, error)
	FindByType(eventType events.EventType) ([]events.Event, error)
}
