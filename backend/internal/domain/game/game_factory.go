package game

import (
	"backend/internal/domain/game"
	"backend/internal/domain/room"
)

func CreateSuecaGame(room room.Room) game.Game {
	roomPlayers := make([]game.Player, 0, len(room.Players))
	for _, p := range room.Players {
		roomPlayers = append(roomPlayers, p)
	}
	team1 := game.Team{ID: "team1"}
	team2 := game.Team{ID: "team2"}
	for i, p := range roomPlayers {
		if i%2 == 0 {
			team1.Players = append(team1.Players, p)
		} else {
			team2.Players = append(team2.Players, p)
		}
	}

	game := game.NewGame([]game.Team{team1, team2}, SuecaScoringStrategy, EasyBotStrategy)
	return game
}
