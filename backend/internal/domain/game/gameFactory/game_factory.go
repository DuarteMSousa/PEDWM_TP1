package game_factory

import (
	"backend/internal/domain/events"
	"backend/internal/domain/game"
	"backend/internal/domain/player"
	bot_strategy "backend/internal/domain/player/botStrategy"
	"backend/internal/domain/team"
	"sort"
	"strconv"
)

func CreateSuecaGame(players map[string]*player.Player, botStrategy bot_strategy.IBotStrategy, eventBus *events.EventBus) *game.Game {
	roomPlayers := make([]*player.Player, 0, len(players))
	for _, p := range players {
		roomPlayers = append(roomPlayers, p)
	}

	if len(roomPlayers) < 4 {
		for i := len(roomPlayers); i < 4; i++ {
			currentSequence := 1
			for _, p := range roomPlayers {
				if p.Sequence == currentSequence {
					currentSequence = p.Sequence + 1
				}
			}
			botPlayer := player.NewPlayer("bot"+strconv.Itoa(i), "Bot "+strconv.Itoa(i), currentSequence)
			botPlayer.Type = player.BOT
			roomPlayers = append(roomPlayers, botPlayer)
		}
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

	game := game.NewGame([]*team.Team{&team1, &team2}, game.NewSuecaGameScoringStrategy(), botStrategy)
	if eventBus != nil {
		game.SetEventBus(eventBus)
	}

	return game
}
