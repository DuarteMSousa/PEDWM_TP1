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
	players := make([]*model.Player, 0, len(r.Players))

	for _, p := range r.Players {
		players = append(players, &model.Player{
			ID:       p.ID,
			Name:     p.Name,
			Type:     model.PlayerType(p.Type),
			Sequence: int32(p.Sequence),
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
