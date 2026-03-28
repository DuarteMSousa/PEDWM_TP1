package graph

import (
	"backend/internal/domain/room"
	"backend/internal/infrastructure/graph/model"
)

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
