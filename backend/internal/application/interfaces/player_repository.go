package interfaces

import "backend/internal/domain/player"

// PlayerRepository defines the contract for player persistence.
type PlayerRepository interface {
	CreateOrGetByNickname(nickname string) (*player.Player, error)
	GetPlayer(playerID string) (*player.Player, bool)
	ListPlayers() []*player.Player
	PlayersByIDs(ids []string) []*player.Player
}
