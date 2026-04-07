package room

import (
	"backend/internal/domain/events"
	"errors"
	"testing"
)

type roomObserverStub struct {
	events []events.Event
}

func (o *roomObserverStub) Update(event events.Event) {
	o.events = append(o.events, event)
}

func TestNewRoomValidationAndHost(t *testing.T) {
	t.Parallel()

	if _, err := NewRoom("", "host1", "Host"); !errors.Is(err, ErrInvalidRoomID) {
		t.Fatalf("expected ErrInvalidRoomID, got %v", err)
	}
	if _, err := NewRoom("room1", "", "Host"); !errors.Is(err, ErrInvalidHost) {
		t.Fatalf("expected ErrInvalidHost for empty host ID, got %v", err)
	}
	if _, err := NewRoom("room1", "host1", ""); !errors.Is(err, ErrInvalidHost) {
		t.Fatalf("expected ErrInvalidHost for empty username, got %v", err)
	}

	r, err := NewRoom("room1", "host1", "Host")
	if err != nil {
		t.Fatalf("expected valid room, got error %v", err)
	}
	if r.Status != OPEN {
		t.Fatalf("expected room status OPEN, got %s", r.Status)
	}
	if len(r.Players) != 1 {
		t.Fatalf("expected exactly 1 player (host), got %d", len(r.Players))
	}
	if _, ok := r.Players["host1"]; !ok {
		t.Fatal("expected host player to exist in room")
	}
}

func TestAddPlayerPublishesEventAndPreventsDuplicate(t *testing.T) {
	t.Parallel()

	r, err := NewRoom("room1", "host1", "Host")
	if err != nil {
		t.Fatalf("failed to create room: %v", err)
	}

	bus := events.NewEventBus()
	obs := &roomObserverStub{}
	bus.Subscribe(obs)
	r.SetEventBus(bus)

	if err := r.AddPlayer("p2", "Player2"); err != nil {
		t.Fatalf("expected AddPlayer to succeed, got %v", err)
	}

	if len(obs.events) != 1 {
		t.Fatalf("expected 1 published event, got %d", len(obs.events))
	}
	if obs.events[0].Type != events.EventPlayerJoined {
		t.Fatalf("expected EventPlayerJoined, got %s", obs.events[0].Type)
	}
	if obs.events[0].RoomID != r.ID {
		t.Fatalf("expected event roomID %q, got %q", r.ID, obs.events[0].RoomID)
	}

	err = r.AddPlayer("p2", "Player2")
	if !errors.Is(err, ErrPlayerAlreadyInRoom) {
		t.Fatalf("expected ErrPlayerAlreadyInRoom, got %v", err)
	}
}

func TestRemovePlayerClosesEmptyRoomAndPublishesCloseEvent(t *testing.T) {
	t.Parallel()

	r, err := NewRoom("room1", "host1", "Host")
	if err != nil {
		t.Fatalf("failed to create room: %v", err)
	}

	bus := events.NewEventBus()
	obs := &roomObserverStub{}
	bus.Subscribe(obs)
	r.SetEventBus(bus)

	if err := r.RemovePlayer("host1"); err != nil {
		t.Fatalf("expected RemovePlayer to succeed, got %v", err)
	}

	if r.Status != CLOSED {
		t.Fatalf("expected room status CLOSED, got %s", r.Status)
	}
	if len(r.Players) != 0 {
		t.Fatalf("expected no players after removing host, got %d", len(r.Players))
	}

	if len(obs.events) < 2 {
		t.Fatalf("expected at least 2 events (PLAYER_LEFT + ROOM_CLOSED), got %d", len(obs.events))
	}
	last := obs.events[len(obs.events)-1]
	if last.Type != events.EventRoomClosed {
		t.Fatalf("expected last event EventRoomClosed, got %s", last.Type)
	}
}

func TestCreateGameSetsInGameAndCreatesFourPlayers(t *testing.T) {
	t.Parallel()

	r, err := NewRoom("room1", "host1", "Host")
	if err != nil {
		t.Fatalf("failed to create room: %v", err)
	}

	if err := r.CreateGame(); err != nil {
		t.Fatalf("expected CreateGame to succeed, got %v", err)
	}

	if r.Status != IN_GAME {
		t.Fatalf("expected room status IN_GAME, got %s", r.Status)
	}
	if r.Game == nil {
		t.Fatal("expected room.Game to be initialized")
	}
	if len(r.Game.Players) != 4 {
		t.Fatalf("expected game to have 4 players (with bots), got %d", len(r.Game.Players))
	}
}
