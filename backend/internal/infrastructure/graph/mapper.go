package graph

import (
	application "backend/internal/application/services"
	"backend/internal/domain/card"
	"backend/internal/domain/friendship"
	"backend/internal/domain/room"
	"backend/internal/domain/trick"
	"backend/internal/domain/user"
	"backend/internal/infrastructure/graph/model"
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
	players := make([]*model.RoomPlayer, 0, len(r.Players))

	for _, p := range r.Players {
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

func mapGameSnapshot(snapshot *application.GameSnapshot) *model.GameSnapshot {
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

	return &model.GameSnapshot{
		RoomID:          snapshot.RoomID,
		GameID:          snapshot.GameID,
		Status:          snapshot.Status,
		TrumpSuit:       trumpSuit,
		CurrentPlayerID: currentPlayerID,
		MyHand:          hand,
		TablePlays:      tablePlays,
		Scores:          scores,
	}
}
