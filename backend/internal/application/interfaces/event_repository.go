package interfaces

import (
	"backend/internal/domain/events"
)

// EventRepository defines the contract for event persistence.
type EventRepository interface {
	Save(event events.Event) error
	FindByRoomID(roomID string) ([]events.Event, error)
	FindByGameID(gameID string) ([]events.Event, error)
}
