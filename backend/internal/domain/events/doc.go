// Package events defines the domain event system used for
// asynchronous communication between aggregates and system layers. It follows the
// Observer pattern (Subject/Observer) through the EventBus.
//
// Supported event types: PlayerJoined, PlayerLeft, GameStarted,
// GameEnded, RoundStarted, RoundEnded, TrickStarted, TrickEnded,
// CardPlayed, CardDealt, TrumpRevealed, TurnChanged, GameScoreUpdated,
// RoomClosed, and BotStrategyChanged.
//
// Each event carries a typed payload with relevant data and
// metadata such as ID, type, timestamp, and sequence.
package events
