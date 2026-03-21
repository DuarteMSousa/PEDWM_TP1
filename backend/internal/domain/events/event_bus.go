package events

import "sync"

// Subject contract for domain events.
type Subject interface {
	Subscribe(observer Observer)
	Unsubscribe(observer Observer)
	Publish(event Event)
}

// Observer contract for domain events.
type Observer interface {
	Update(event Event)
}

// EventBus dispatches domain events to subscribed observers.
type EventBus struct {
	mu          sync.RWMutex
	subscribers []Observer
}

func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: make([]Observer, 0),
	}
}

func (b *EventBus) Subscribe(observer Observer) {
	if b == nil || observer == nil {
		return
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	b.subscribers = append(b.subscribers, observer)
}

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

func (b *EventBus) Publish(event Event) {
	if b == nil {
		return
	}

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
