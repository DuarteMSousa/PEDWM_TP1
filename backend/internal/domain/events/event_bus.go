package events

import (
	"log/slog"
	"sync"
)

// Subject represents the contract for domain events.
type Subject interface {
	Subscribe(observer Observer)
	Unsubscribe(observer Observer)
	Publish(event Event)
}

// Observer represents the contract for domain events.
type Observer interface {
	Update(event Event)
}

// EventBus distributes domain events to registered observers.
// It is thread-safe through sync.RWMutex.
type EventBus struct {
	mu          sync.RWMutex
	subscribers []Observer
}

var (
	defaultBusMu sync.RWMutex
	defaultBus   = NewEventBus()
)

// NewEventBus creates a new event bus without observers.
func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: make([]Observer, 0),
	}
}

// DefaultBus returns the global instance of the event bus.
func DefaultBus() *EventBus {
	defaultBusMu.RLock()
	defer defaultBusMu.RUnlock()
	return defaultBus
}

// SetDefaultBus replaces the global instance of the event bus.
func SetDefaultBus(bus *EventBus) {
	if bus == nil {
		return
	}

	defaultBusMu.Lock()
	defaultBus = bus
	defaultBusMu.Unlock()
}

// Subscribe registers an observer to receive events.
func (b *EventBus) Subscribe(observer Observer) {
	if b == nil || observer == nil {
		return
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	b.subscribers = append(b.subscribers, observer)
}

// Unsubscribe removes an observer from the bus.
func (b *EventBus) Unsubscribe(observer Observer) {
	if b == nil || observer == nil {
		return
	}

	b.mu.Lock()
	defer b.mu.Unlock()
	for i, o := range b.subscribers {
		if o == observer {
			b.subscribers = append(b.subscribers[:i], b.subscribers[i+1:]...)
			break
		}
	}
}

// Publish notifies all registered observers with the given event.
func (b *EventBus) Publish(event Event) {
	if b == nil {
		return
	}

	slog.Debug("publishing event", "type", event.Type, "gameID", event.GameID, "roomID", event.RoomID)

	b.mu.RLock()
	observers := make([]Observer, 0, len(b.subscribers))
	for _, observer := range b.subscribers {
		observers = append(observers, observer)
	}
	b.mu.RUnlock()

	for _, observer := range observers {
		observer.Update(event)
	}
}
