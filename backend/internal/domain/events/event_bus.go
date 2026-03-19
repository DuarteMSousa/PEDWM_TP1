package events

import "sync"

// Observer contract for domain events.
type Observer interface {
	Update(event Event)
}

// EventBus dispatches domain events to subscribed observers.
type EventBus struct {
	mu           sync.RWMutex
	nextID       int
	subscription map[int]Observer
}

func NewEventBus() *EventBus {
	return &EventBus{
		subscription: make(map[int]Observer),
	}
}

func (b *EventBus) Subscribe(observer Observer) int {
	if b == nil || observer == nil {
		return 0
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	b.nextID++
	id := b.nextID
	b.subscription[id] = observer
	return id
}

func (b *EventBus) Unsubscribe(subscriptionID int) {
	if b == nil || subscriptionID <= 0 {
		return
	}

	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.subscription, subscriptionID)
}

func (b *EventBus) Publish(event Event) {
	if b == nil {
		return
	}

	b.mu.RLock()
	observers := make([]Observer, 0, len(b.subscription))
	for _, observer := range b.subscription {
		observers = append(observers, observer)
	}
	b.mu.RUnlock()

	for _, observer := range observers {
		observer.Update(event)
	}
}
