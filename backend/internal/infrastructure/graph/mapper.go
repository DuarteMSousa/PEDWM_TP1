package graph

import (
	"backend/internal/application/services"
	"backend/internal/domain/card"
	domainevents "backend/internal/domain/events"
	"backend/internal/domain/game"
	"backend/internal/domain/player"
	"backend/internal/domain/room"
	"backend/internal/domain/team"
	"backend/internal/domain/trick"
	"backend/internal/domain/user"
	"backend/internal/infrastructure/graph/model"
	"sort"
)

func mapUser(u *user.User) *model.User {
	return &model.User{
		ID:       u.ID,
		Username: u.Username,
	}
}

func mapUserStats(s *user.UserStats) *model.UserStats {
	return &model.UserStats{
		UserID: s.UserId,
		Games:  int32(s.Games),
		Wins:   int32(s.Wins),
		Elo:    int32(s.Elo),
	}
}

func mapRoom(r *room.Room) *model.Room {
	orderedPlayers := make([]*player.Player, 0, len(r.Players))
	for _, p := range r.Players {
		orderedPlayers = append(orderedPlayers, p)
	}
	sort.Slice(orderedPlayers, func(i, j int) bool {
		if orderedPlayers[i].Sequence == orderedPlayers[j].Sequence {
			return orderedPlayers[i].ID < orderedPlayers[j].ID
		}
		return orderedPlayers[i].Sequence < orderedPlayers[j].Sequence
	})

	players := make([]*model.RoomPlayer, 0, len(orderedPlayers))
	for _, p := range orderedPlayers {
		players = append(players, &model.RoomPlayer{
			ID:       p.ID,
			Username: p.Name,
		})
	}

	botStrategy := model.BotStrategyTypeEasy
	if r.BotStrategy != nil {
		botStrategy = model.BotStrategyType(r.BotStrategy.GetType())
	}

	return &model.Room{
		ID:          r.ID,
		HostID:      r.HostID,
		Players:     players,
		Status:      model.RoomStatus(r.Status),
		CreatedAt:   r.CreatedAt,
		BotStrategy: botStrategy,
	}
}

func mapGameCard(c card.Card) *model.GameCard {
	return &model.GameCard{
		ID:   c.ID,
		Suit: string(c.Suit),
		Rank: string(c.Rank),
	}
}

func mapGameTablePlay(p trick.Play) *model.GameTablePlay {
	return &model.GameTablePlay{
		PlayerID: p.PlayerID,
		Card:     mapGameCard(p.Card),
	}
}

func mapEvent(e domainevents.Event) *model.Event {
	return &model.Event{
		ID:        e.ID,
		Type:      model.EventType(e.Type),
		GameID:    &e.GameID,
		RoomID:    &e.RoomID,
		Timestamp: e.Timestamp,
		Sequence:  int32(e.Sequence),
		Payload:   mapEventPayload(e.Payload),
	}
}

func mapTeamPlayers(players []*player.Player) []*model.Player {
	result := make([]*model.Player, 0, len(players))
	for _, p := range players {
		result = append(result, &model.Player{
			ID:       p.ID,
			Name:     p.Name,
			Type:     model.PlayerType(p.Type),
			Sequence: int32(p.Sequence),
		})
	}
	return result
}

func mapDomainTeams(teams []team.Team) []*model.Team {
	result := make([]*model.Team, 0, len(teams))
	for _, t := range teams {
		teamPlayers := make([]*player.Player, 0, len(t.Players))
		for i := range t.Players {
			teamPlayers = append(teamPlayers, t.Players[i])
		}

		result = append(result, &model.Team{
			ID:      t.ID,
			Players: mapTeamPlayers(teamPlayers),
		})
	}
	return result
}

func mapTeamScores(score map[string]int) []*model.TeamScore {
	result := make([]*model.TeamScore, 0, len(score))
	for teamID, points := range score {
		result = append(result, &model.TeamScore{TeamID: teamID, Points: int32(points)})
	}
	return result
}

func mapEventPayload(payload any) model.EventPayload {
	switch p := payload.(type) {
	case domainevents.PlayerJoinedPayload:
		return model.PlayerJoinedEventPayload{PlayerID: p.PlayerID, Name: p.Name, Slot: int32(p.Slot)}
	case *domainevents.PlayerJoinedPayload:
		return model.PlayerJoinedEventPayload{PlayerID: p.PlayerID, Name: p.Name, Slot: int32(p.Slot)}
	case domainevents.PlayerLeftPayload:
		return model.PlayerLeftEventPayload{PlayerID: p.PlayerID, RoomID: p.RoomID}
	case *domainevents.PlayerLeftPayload:
		return model.PlayerLeftEventPayload{PlayerID: p.PlayerID, RoomID: p.RoomID}
	case domainevents.RoundStartedPayload:
		return model.RoundStartedEventPayload{RoundNumber: int32(p.RoundNumber), DealerID: p.DealerID}
	case *domainevents.RoundStartedPayload:
		return model.RoundStartedEventPayload{RoundNumber: int32(p.RoundNumber), DealerID: p.DealerID}
	case domainevents.TrickStartedPayload:
		return model.TrickStartedEventPayload{LeaderID: p.LeaderID}
	case *domainevents.TrickStartedPayload:
		return model.TrickStartedEventPayload{LeaderID: p.LeaderID}
	case domainevents.TrumpRevealedPayload:
		return model.TrumpRevealedEventPayload{Card: mapGameCard(p.Card), Suit: string(p.Suit)}
	case *domainevents.TrumpRevealedPayload:
		return model.TrumpRevealedEventPayload{Card: mapGameCard(p.Card), Suit: string(p.Suit)}
	case domainevents.CardDealtPayload:
		return model.CardDealtEventPayload{PlayerID: p.PlayerID, Card: mapGameCard(p.Card)}
	case *domainevents.CardDealtPayload:
		return model.CardDealtEventPayload{PlayerID: p.PlayerID, Card: mapGameCard(p.Card)}
	case domainevents.CardPlayedPayload:
		return model.CardPlayedEventPayload{PlayerID: p.PlayerID, Card: mapGameCard(p.Card)}
	case *domainevents.CardPlayedPayload:
		return model.CardPlayedEventPayload{PlayerID: p.PlayerID, Card: mapGameCard(p.Card)}
	case domainevents.TurnChangedPayload:
		return model.TurnChangedEventPayload{PlayerID: p.PlayerID}
	case *domainevents.TurnChangedPayload:
		return model.TurnChangedEventPayload{PlayerID: p.PlayerID}
	case domainevents.TrickEndedPayload:
		return model.TrickEndedEventPayload{WinnerID: p.WinnerID, Points: int32(p.Points)}
	case *domainevents.TrickEndedPayload:
		return model.TrickEndedEventPayload{WinnerID: p.WinnerID, Points: int32(p.Points)}
	case domainevents.RoundEndedPayload:
		return model.RoundEndedEventPayload{WinnerTeam: p.WinnerTeam, Score: mapTeamScores(p.Score)}
	case *domainevents.RoundEndedPayload:
		return model.RoundEndedEventPayload{WinnerTeam: p.WinnerTeam, Score: mapTeamScores(p.Score)}
	case domainevents.GameScorePayload:
		return model.GameScoreUpdatedEventPayload{Score: mapTeamScores(p.Score)}
	case *domainevents.GameScorePayload:
		return model.GameScoreUpdatedEventPayload{Score: mapTeamScores(p.Score)}
	case domainevents.GameStartedPayload:
		return model.GameStartedEventPayload{Teams: mapDomainTeams(p.Teams)}
	case *domainevents.GameStartedPayload:
		return model.GameStartedEventPayload{Teams: mapDomainTeams(p.Teams)}
	case domainevents.GameEndedPayload:
		return model.GameEndedEventPayload{Winner: p.Winner, FinalScores: mapTeamScores(p.FinalScores), Teams: mapDomainTeams(p.Teams)}
	case *domainevents.GameEndedPayload:
		return model.GameEndedEventPayload{Winner: p.Winner, FinalScores: mapTeamScores(p.FinalScores), Teams: mapDomainTeams(p.Teams)}
	case domainevents.RoomClosedPayload:
		return model.RoomClosedEventPayload{RoomID: p.RoomID}
	case *domainevents.RoomClosedPayload:
		return model.RoomClosedEventPayload{RoomID: p.RoomID}
	default:
		return nil
	}
}

func mapGameSnapshot(snapshot *services.GameSnapshot) *model.GameSnapshot {
	if snapshot == nil {
		return nil
	}

	hand := make([]*model.GameCard, 0, len(snapshot.MyHand))
	for _, c := range snapshot.MyHand {
		hand = append(hand, mapGameCard(c))
	}

	tablePlays := make([]*model.GameTablePlay, 0, len(snapshot.TablePlays))
	for _, p := range snapshot.TablePlays {
		tablePlays = append(tablePlays, mapGameTablePlay(p))
	}

	scores := make([]*model.TeamScore, 0, len(snapshot.Scores))
	for teamID, points := range snapshot.Scores {
		scores = append(scores, &model.TeamScore{
			TeamID: teamID,
			Points: int32(points),
		})
	}

	var trumpSuit *string
	if snapshot.TrumpSuit != "" {
		value := snapshot.TrumpSuit
		trumpSuit = &value
	}

	var currentPlayerID *string
	if snapshot.CurrentPlayerID != "" {
		value := snapshot.CurrentPlayerID
		currentPlayerID = &value
	}

	var teams []*model.Team
	if snapshot.Teams != nil {
		for _, t := range snapshot.Teams {
			teams = append(teams, &model.Team{
				ID: t.ID,
				Players: func() []*model.Player {
					result := make([]*model.Player, 0, len(t.Players))
					for _, p := range t.Players {
						result = append(result, &model.Player{
							ID:       p.ID,
							Name:     p.Name,
							Sequence: int32(p.Sequence),
							Type:     model.PlayerType(p.Type),
						})
					}
					return result
				}(),
			})
		}
	}

	snap := &model.GameSnapshot{
		RoomID:          snapshot.RoomID,
		GameID:          snapshot.GameID,
		Status:          snapshot.Status,
		TrumpSuit:       trumpSuit,
		CurrentPlayerID: currentPlayerID,
		MyHand:          hand,
		Teams:           teams,
		TablePlays:      tablePlays,
		Scores:          scores,
	}

	return snap
}

func mapGame(g *game.Game) *model.Game {

	var players []*model.GamePlayer
	for _, team := range g.Teams {
		for _, player := range team.Players {

			players = append(players, &model.GamePlayer{
				ID:       player.ID,
				Username: player.Name,
			})

		}
	}

	var resultEvents []*model.Event

	for _, e := range g.Events {
		resultEvents = append(resultEvents, &model.Event{
			ID:        e.ID,
			Type:      model.EventType(e.Type),
			GameID:    &e.GameID,
			RoomID:    &e.RoomID,
			Timestamp: e.Timestamp,
			Sequence:  int32(e.Sequence),
			Payload:   mapEventPayload(e.Payload),
		})
	}

	return &model.Game{
		ID:        g.ID.String(),
		RoomID:    &g.RoomID,
		Players:   players,
		Events:    resultEvents,
		CreatedAt: g.CreatedAt,
		UpdatedAt: g.UpdatedAt,
	}
}
