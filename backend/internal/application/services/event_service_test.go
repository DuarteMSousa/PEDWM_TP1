package services

import (
	"backend/internal/domain/events"
	"errors"
	"testing"
)

type fakeEventRepo struct {
	saveErr      error
	findRoomErr  error
	findGameErr  error
	lastSaved    events.Event
	eventsByRoom map[string][]events.Event
	eventsByGame map[string][]events.Event
}

func (f *fakeEventRepo) Save(event events.Event) error {
	if f.saveErr != nil {
		return f.saveErr
	}
	f.lastSaved = event
	return nil
}

func (f *fakeEventRepo) FindByRoomID(roomID string) ([]events.Event, error) {
	if f.findRoomErr != nil {
		return nil, f.findRoomErr
	}
	return f.eventsByRoom[roomID], nil
}

func (f *fakeEventRepo) FindByGameID(gameID string) ([]events.Event, error) {
	if f.findGameErr != nil {
		return nil, f.findGameErr
	}
	return f.eventsByGame[gameID], nil
}

func TestEventServiceSaveEvent(t *testing.T) {
	t.Parallel()

	repo := &fakeEventRepo{}
	service := NewEventService(repo)
	ev := events.NewRoundStartedEvent("game_1")

	if err := service.SaveEvent(ev); err != nil {
		t.Fatalf("expected save success, got %v", err)
	}
	if repo.lastSaved.ID != ev.ID {
		t.Fatalf("expected saved event %q, got %q", ev.ID, repo.lastSaved.ID)
	}

	repo.saveErr = errors.New("save failed")
	if err := service.SaveEvent(ev); err == nil {
		t.Fatal("expected save error")
	}
}

func TestEventServiceGetEventsByRoomAndGame(t *testing.T) {
	t.Parallel()

	roomEvent := events.NewTurnChangedEvent("g1", "p1")
	gameEvent := events.NewRoundStartedEvent("g2")

	repo := &fakeEventRepo{
		eventsByRoom: map[string][]events.Event{
			"room_1": {roomEvent},
		},
		eventsByGame: map[string][]events.Event{
			"game_2": {gameEvent},
		},
	}
	service := NewEventService(repo)

	byRoom, err := service.GetEventsByRoomID("room_1")
	if err != nil {
		t.Fatalf("expected room query success, got %v", err)
	}
	if len(byRoom) != 1 || byRoom[0].ID != roomEvent.ID {
		t.Fatalf("unexpected room events result: %+v", byRoom)
	}

	byGame, err := service.GetEventsByGameID("game_2")
	if err != nil {
		t.Fatalf("expected game query success, got %v", err)
	}
	if len(byGame) != 1 || byGame[0].ID != gameEvent.ID {
		t.Fatalf("unexpected game events result: %+v", byGame)
	}

	repo.findRoomErr = errors.New("room query failed")
	if _, err := service.GetEventsByRoomID("room_1"); err == nil {
		t.Fatal("expected room query error")
	}

	repo.findGameErr = errors.New("game query failed")
	if _, err := service.GetEventsByGameID("game_2"); err == nil {
		t.Fatal("expected game query error")
	}
}
