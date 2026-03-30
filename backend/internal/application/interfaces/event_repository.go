package interfaces

import (
	"backend/internal/domain/events"
)

type EventRepository interface {
	Save(event events.Event) error
}
