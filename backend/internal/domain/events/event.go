package events

import (
	"backend/internal/domain/card"
	bot_strategy "backend/internal/domain/player/botStrategy"
	"backend/internal/domain/team"
	"time"

	"github.com/google/uuid"
)

type EventType string

const (
	// EventGameCreated   EventType = "GAME_CREATED"
	EventRoundStarted  EventType = "ROUND_STARTED"
	EventTrickStarted  EventType = "TRICK_STARTED"
	EventTrumpRevealed EventType = "TRUMP_REVEALED"
	EventCardDealt     EventType = "CARD_DEALT"
	EventRoundEnded    EventType = "ROUND_ENDED"
	EventGameScore     EventType = "GAME_SCORE_UPDATED"

	EventPlayerJoined EventType = "PLAYER_JOINED"
	EventPlayerLeft   EventType = "PLAYER_LEFT"
	EventGameStarted  EventType = "GAME_STARTED"
	EventTurnChanged  EventType = "TURN_CHANGED"
	EventCardPlayed   EventType = "CARD_PLAYED"
	EventTrickEnded   EventType = "TRICK_ENDED"
	EventGameEnded    EventType = "GAME_ENDED"
	EventRoomClosed   EventType = "ROOM_CLOSED"

	EventBotStrategyChanged EventType = "BOT_STRATEGY_CHANGED"
)

type Event struct {
	ID        string    `json:"id"`
	Type      EventType `json:"type"`
	GameID    string    `json:"gameId,omitempty"`
	RoomID    string    `json:"roomId,omitempty"`
	Timestamp time.Time `json:"timestamp"`
	Sequence  int       `json:"sequence,omitempty"`
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
	RoomID   string `json:"roomId"`
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

type TurnChangedPayload struct {
	PlayerID string `json:"playerId"`
}

type TrickEndedPayload struct {
	WinnerID string `json:"winnerId"`
	Points   int    `json:"points"`
}

type RoundEndedPayload struct {
	Score      map[string]int `json:"score"`
	WinnerTeam string         `json:"winnerTeam"`
}

type GameScorePayload struct {
	Score map[string]int `json:"score"`
}

type GameStartedPayload struct {
	Teams []team.Team `json:"teams,omitempty"`
}

type GameEndedPayload struct {
	FinalScores map[string]int `json:"finalScores"`
	Winner      string         `json:"winner"`
	Teams       []team.Team    `json:"teams,omitempty"`
}

type RoomClosedPayload struct {
	RoomID string `json:"roomId"`
}

type BotStrategyChangedPayload struct {
	RoomID      string                       `json:"roomId"`
	BotStrategy bot_strategy.BotStrategyType `json:"botStrategy"`
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

func NewPlayerLeftEvent(gameID string, playerID string, roomID string) Event {
	return newEvent(EventPlayerLeft, gameID, PlayerLeftPayload{PlayerID: playerID, RoomID: roomID})
}

func NewGameStartedEvent(gameID string, teams []team.Team) Event {
	return newEvent(EventGameStarted, gameID, GameStartedPayload{Teams: teams})
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

func NewTurnChangedEvent(gameID string, playerID string) Event {
	return newEvent(EventTurnChanged, gameID, TurnChangedPayload{PlayerID: playerID})
}

func NewTrickEndedEvent(gameID string, winnerID string, points int) Event {
	return newEvent(EventTrickEnded, gameID, TrickEndedPayload{WinnerID: winnerID, Points: points})
}

func NewRoundEndedEvent(gameID string, score map[string]int, winnerTeam string) Event {
	return newEvent(EventRoundEnded, gameID, RoundEndedPayload{Score: score, WinnerTeam: winnerTeam})
}

func NewGameScoreUpdatedEvent(gameID string, score map[string]int) Event {
	return newEvent(EventGameScore, gameID, GameScorePayload{Score: score})
}

func NewGameEndedEvent(gameID string, finalScores map[string]int, winner string, teams []team.Team) Event {
	return newEvent(EventGameEnded, gameID, GameEndedPayload{FinalScores: finalScores, Winner: winner, Teams: teams})
}

func NewRoomClosedEvent(roomID string) Event {
	event := newEvent(EventRoomClosed, "", RoomClosedPayload{RoomID: roomID})
	event.RoomID = roomID
	return event
}

func NewBotStrategyChangedEvent(roomID string, strategy bot_strategy.BotStrategyType) Event {
	event := newEvent(EventBotStrategyChanged, "", BotStrategyChangedPayload{RoomID: roomID, BotStrategy: strategy})
	event.RoomID = roomID
	return event
}
