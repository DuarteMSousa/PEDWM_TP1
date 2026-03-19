package events

import (
	"backend/internal/application/ports"
	domainevents "backend/internal/domain/events"
)

type EventBusPublisher struct {
	bus *domainevents.EventBus
}

func NewEventBusPublisher(bus *domainevents.EventBus) ports.EventPublisher {
	return &EventBusPublisher{bus: bus}
}

func (p *EventBusPublisher) Publish(event domainevents.Event) {
	if p == nil || p.bus == nil {
		return
	}
	p.bus.Publish(event)
}
