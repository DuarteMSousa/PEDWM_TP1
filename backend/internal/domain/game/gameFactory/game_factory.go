package game_factory

import (
	"backend/internal/domain/game"
	"backend/internal/domain/player"
	bot_strategy "backend/internal/domain/player/botStrategy"
	"backend/internal/domain/team"
	"sort"
)

func CreateSuecaGame(players map[string]*player.Player) *game.Game {
	roomPlayers := make([]*player.Player, 0, len(players))
	for _, p := range players {
		roomPlayers = append(roomPlayers, p)
	}

	sort.Slice(roomPlayers, func(i, j int) bool {
		return roomPlayers[i].Sequence < roomPlayers[j].Sequence
	})

	team1 := team.Team{ID: "team1"}
	team2 := team.Team{ID: "team2"}
	for i, p := range roomPlayers {
		if i%2 == 0 {
			team1.Players = append(team1.Players, p)
		} else {
			team2.Players = append(team2.Players, p)
		}
	}

	game := game.NewGame([]*team.Team{&team1, &team2}, game.NewSuecaGameScoringStrategy(), bot_strategy.NewEasyBotStrategy())
	return game
}
