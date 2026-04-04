package graph

import (
	"backend/internal/application/services"
	"backend/internal/domain/card"
	"backend/internal/domain/friendship"
	"backend/internal/domain/player"
	"backend/internal/domain/room"
	"backend/internal/domain/trick"
	"backend/internal/domain/user"
	"backend/internal/infrastructure/graph/model"
	"encoding/json"
	"log"
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

func mapFriendship(f *friendship.Friendship) *model.Friendship {
	return &model.Friendship{
		RequesterID: f.RequesterID,
		AddresseeID: f.AddresseeID,
		Status:      model.FriendshipStatus(f.Status),
		CreatedAt:   f.CreatedAt,
		UpdatedAt:   f.UpdatedAt,
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

	return &model.Room{
		ID:        r.ID,
		HostID:    r.HostID,
		Players:   players,
		Status:    model.RoomStatus(r.Status),
		CreatedAt: r.CreatedAt,
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

	b, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		log.Println("Erro ao mapear:", err)
	}
	log.Printf("Mapped GameSnapshot:\n%s", string(b))
	return snap
}
