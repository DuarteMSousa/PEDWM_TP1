package interfaces

import (
	"backend/internal/domain/events"
)

type EventService interface {
	SaveEvent(event events.Event) error
}
