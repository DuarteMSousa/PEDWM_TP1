package graph

import (
	"backend/internal/domain/friendship"
	"backend/internal/domain/room"
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
			ID:       p.UserID,
			Username: p.Username,
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
