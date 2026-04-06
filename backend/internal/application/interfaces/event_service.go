package interfaces

import (
	"backend/internal/domain/events"
)

// EventService defines the contract for event persistence.
type EventService interface {
	SaveEvent(event events.Event) error
}
