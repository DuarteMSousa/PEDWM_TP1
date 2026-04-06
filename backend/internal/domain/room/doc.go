// Package room defines the Room aggregate (game room). A room groups
// players (up to 4), manages the lifecycle from OPEN → IN_GAME → CLOSED,
// and is responsible for creating the game (Game) when all players are
// ready. Publishes domain events through the EventBus.
package room
