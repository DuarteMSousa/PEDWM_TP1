package interfaces

import domainevents "backend/internal/domain/events"

type EventPublisher interface {
	Publish(event domainevents.Event)
}
