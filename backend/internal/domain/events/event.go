package events

import (
	"backend/internal/domain/card"
	"time"

	"github.com/google/uuid"
)

// NOTE: Event model kept intentionally minimal so websocket transport can be
// developed independently. Event payload design and complete event catalog are
// owned by another teammate.

type EventType string

const (
	// EventGameCreated   EventType = "GAME_CREATED"
	EventRoundStarted  EventType = "ROUND_STARTED"
	EventTrickStarted  EventType = "TRICK_STARTED"
	EventTrumpRevealed EventType = "TRUMP_REVEALED"
	EventCardDealt     EventType = "CARD_DEALT"
	EventRoundEnded    EventType = "ROUND_ENDED"

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
	GameID    string    `json:"gameId,omitempty"`
	Timestamp time.Time `json:"timestamp"`
	Payload   any       `json:"payload,omitempty"`
}

type GameCreatedPayload struct {
	CreatorID string `json:"creatorId"`
	Settings  any    `json:"settings,omitempty"`
}

type PlayerJoinedPayload struct {
	PlayerID string `json:"playerId"`
	Name     string `json:"name"`
	Slot     int    `json:"slot"`
}

type PlayerLeftPayload struct {
	PlayerID string `json:"playerId"`
}

type RoundStartedPayload struct {
	RoundNumber int    `json:"roundNumber"`
	DealerID    string `json:"dealerId"`
}

type TrickStartedPayload struct {
	LeaderID string `json:"leaderId"`
}

type TrumpRevealedPayload struct {
	Card card.Card `json:"card"`
	Suit card.Suit `json:"suit"`
}

type CardDealtPayload struct {
	PlayerID string    `json:"playerId"`
	Card     card.Card `json:"card"`
}

type CardPlayedPayload struct {
	PlayerID string    `json:"playerId"`
	Card     card.Card `json:"card"`
}

type TrickEndedPayload struct {
	WinnerID string `json:"winnerId"`
	Points   int    `json:"points"`
}

type RoundEndedPayload struct {
	TeamAScore int    `json:"teamAScore"`
	TeamBScore int    `json:"teamBScore"`
	WinnerTeam string `json:"winnerTeam"`
}

type GameEndedPayload struct {
	FinalScores map[string]int `json:"finalScores"`
	WinnerTeam  string         `json:"winnerTeam"`
}

func (e Event) WithPayload(payload any) Event {
	e.Payload = payload
	return e
}

func newEvent(typ EventType, gameID string, payload any) Event {
	now := time.Now().UTC()
	return Event{
		ID:        uuid.NewString(),
		Type:      typ,
		GameID:    gameID,
		Timestamp: now,
		Payload:   payload,
	}
}

// func NewGameCreatedEvent(gameID string, creatorID string, settings any) Event {
// 	return newEvent(EventGameCreated, gameID, GameCreatedPayload{CreatorID: creatorID, Settings: settings})
// }

func NewPlayerJoinedEvent(gameID string, playerID string, name string, slot int) Event {
	return newEvent(EventPlayerJoined, gameID, PlayerJoinedPayload{PlayerID: playerID, Name: name, Slot: slot})
}

func NewPlayerLeftEvent(gameID string, playerID string) Event {
	return newEvent(EventPlayerLeft, gameID, PlayerLeftPayload{PlayerID: playerID})
}

func NewGameStartedEvent(gameID string) Event {
	return newEvent(EventGameStarted, gameID, nil)
}

func NewRoundStartedEvent(gameID string) Event {
	return newEvent(EventRoundStarted, gameID, nil)
}

func NewTrickStartedEvent(gameID string, leaderID string) Event {
	return newEvent(EventTrickStarted, gameID, TrickStartedPayload{LeaderID: leaderID})
}

func NewTrumpRevealedEvent(gameID string, trumpCard card.Card) Event {
	return newEvent(EventTrumpRevealed, gameID, TrumpRevealedPayload{Card: trumpCard, Suit: trumpCard.Suit})
}

func NewCardDealtEvent(gameID string, playerID string, card card.Card) Event {
	return newEvent(EventCardDealt, gameID, CardDealtPayload{PlayerID: playerID, Card: card})
}

func NewCardPlayedEvent(gameID string, playerID string, playedCard card.Card) Event {
	return newEvent(EventCardPlayed, gameID, CardPlayedPayload{PlayerID: playerID, Card: playedCard})
}

func NewTrickEndedEvent(gameID string, winnerID string, points int) Event {
	return newEvent(EventTrickEnded, gameID, TrickEndedPayload{WinnerID: winnerID, Points: points})
}

func NewRoundEndedEvent(gameID string, teamAScore int, teamBScore int, winnerTeam string) Event {
	return newEvent(EventRoundEnded, gameID, RoundEndedPayload{TeamAScore: teamAScore, TeamBScore: teamBScore, WinnerTeam: winnerTeam})
}

func NewGameEndedEvent(gameID string, finalScores map[string]int, winnerTeam string) Event {
	return newEvent(EventGameEnded, gameID, GameEndedPayload{FinalScores: finalScores, WinnerTeam: winnerTeam})
}
