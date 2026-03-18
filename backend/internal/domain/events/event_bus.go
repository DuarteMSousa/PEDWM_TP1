package events

// Observer contract for domain events.
type Observer interface {
	Update(event Event)
}

// EventBus is intentionally minimal in this branch.
// TODO(team-eventbus): replace with full publish/subscribe implementation.
type EventBus struct{}

func NewEventBus() *EventBus {
	return &EventBus{}
}

func (b *EventBus) Subscribe(observer Observer) int {
	_ = observer
	return 0
}

func (b *EventBus) Unsubscribe(subscriptionID int) {
	_ = subscriptionID
}

func (b *EventBus) Publish(event Event) {
	_ = event
}
