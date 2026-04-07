package events

import "testing"

type observerStub struct {
	received []Event
}

func (o *observerStub) Update(event Event) {
	o.received = append(o.received, event)
}

func TestEventBusPublishToSubscribers(t *testing.T) {
	t.Parallel()

	bus := NewEventBus()
	obs := &observerStub{}
	bus.Subscribe(obs)

	event := NewTurnChangedEvent("game_1", "p1")
	bus.Publish(event)

	if len(obs.received) != 1 {
		t.Fatalf("expected 1 received event, got %d", len(obs.received))
	}
	if obs.received[0].Type != EventTurnChanged {
		t.Fatalf("expected event type %s, got %s", EventTurnChanged, obs.received[0].Type)
	}
}

func TestEventBusUnsubscribeStopsReceiving(t *testing.T) {
	t.Parallel()

	bus := NewEventBus()
	obs := &observerStub{}
	bus.Subscribe(obs)
	bus.Unsubscribe(obs)

	event := NewRoundStartedEvent("game_1")
	bus.Publish(event)

	if len(obs.received) != 0 {
		t.Fatalf("expected 0 received events after unsubscribe, got %d", len(obs.received))
	}
}
