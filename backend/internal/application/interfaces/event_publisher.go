package interfaces

import domainevents "backend/internal/domain/events"

// EventPublisher defines the contract for publishing domain events.
type EventPublisher interface {
	Publish(event domainevents.Event)
}
