package events

import (
	"backend/internal/domain/card"
	"fmt"
	"time"
)

// NOTE: Event model kept intentionally minimal so websocket transport can be
// developed independently. Event payload design and complete event catalog are
// owned by another teammate.

type EventType string

const (
	EventRoomCreated  EventType = "ROOM_CREATED"
	EventRoomUpdated  EventType = "ROOM_UPDATED"
	EventRoomDeleted  EventType = "ROOM_DELETED"
	EventPlayerJoined EventType = "PLAYER_JOINED"
	EventPlayerLeft   EventType = "PLAYER_LEFT"
	EventGameStarted  EventType = "GAME_STARTED"
	EventCardPlayed   EventType = "CARD_PLAYED"
	EventTrickEnded   EventType = "TRICK_ENDED"
	EventGameEnded    EventType = "GAME_ENDED"
)

type Event struct {
	ID        string    `json:"id"`
	Type      EventType `json:"type"`
	RoomID    string    `json:"roomId,omitempty"`
	GameID    string    `json:"gameId,omitempty"`
	Timestamp time.Time `json:"timestamp"`
	Payload   any       `json:"payload,omitempty"`
}

func (e Event) WithPayload(payload any) Event {
	e.Payload = payload
	return e
}

func newEvent(typ EventType, roomID string, gameID string, payload any) Event {
	now := time.Now().UTC()
	return Event{
		ID:        fmt.Sprintf("evt_%d", now.UnixNano()),
		Type:      typ,
		RoomID:    roomID,
		GameID:    gameID,
		Timestamp: now,
		Payload:   payload,
	}
}

func NewRoomCreatedEvent(roomID string, hostID string) Event {
	return newEvent(EventRoomCreated, roomID, "", map[string]any{"hostId": hostID})
}

func NewRoomUpdatedEvent(roomID string, payload any) Event {
	return newEvent(EventRoomUpdated, roomID, "", payload)
}

func NewRoomDeletedEvent(roomID string, payload any) Event {
	return newEvent(EventRoomDeleted, roomID, "", payload)
}

func NewPlayerJoinedEvent(roomID string, playerID string, nickname string) Event {
	return newEvent(
		EventPlayerJoined,
		roomID,
		"",
		map[string]any{
			"playerId": playerID,
			"nickname": nickname,
		},
	)
}

func NewPlayerLeftEvent(roomID string, playerID string) Event {
	return newEvent(EventPlayerLeft, roomID, "", map[string]any{"playerId": playerID})
}

func NewGameStartedEvent(roomID string, gameID string) Event {
	return newEvent(EventGameStarted, roomID, gameID, nil)
}

func NewCardPlayedEvent(gameID string, playerID string, card card.Card) Event {
	return newEvent(EventCardPlayed, "", gameID, map[string]any{"playerId": playerID, "card": card})
}

func NewTrickEndedEvent(gameID string, winnerID string, points int) Event {
	return newEvent(EventTrickEnded, "", gameID, map[string]any{"winnerId": winnerID, "points": points})
}

func NewGameEndedEvent(gameID string) Event {
	return newEvent(EventGameEnded, "", gameID, nil)
}
