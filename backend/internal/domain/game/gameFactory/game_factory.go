package game_factory

import (
	"backend/internal/domain/game"
	game_strategy "backend/internal/domain/game/gameStrategy"
	"backend/internal/domain/player"
	bot_strategy "backend/internal/domain/player/botStrategy"
	"backend/internal/domain/room"
	"backend/internal/domain/team"
)

func CreateSuecaGame(room room.Room) *game.Game {
	roomPlayers := make([]*player.Player, 0, len(room.Players))
	for _, p := range room.Players {
		roomPlayers = append(roomPlayers, p)
	}
	team1 := team.Team{ID: "team1"}
	team2 := team.Team{ID: "team2"}
	for i, p := range roomPlayers {
		if i%2 == 0 {
			team1.Players = append(team1.Players, p)
		} else {
			team2.Players = append(team2.Players, p)
		}
	}

	game := game.NewGame([]team.Team{team1, team2}, game_strategy.NewSuecaGameScoringStrategy(), bot_strategy.NewEasyBotStrategy())
	return game
}
