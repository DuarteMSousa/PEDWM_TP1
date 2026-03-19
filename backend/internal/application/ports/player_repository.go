package ports

type PlayerRepository interface {
	CreateOrGetByNickname(nickname string) (Player, error)
	GetPlayer(playerID string) (Player, bool)
	ListPlayers() []Player
	PlayersByIDs(ids []string) []Player
}
