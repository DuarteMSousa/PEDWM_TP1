package graph

import (
	"backend/internal/domain"
	"backend/internal/graph/model"
)

func mapRoom(r *domain.Room) *model.Room {

	players := []*model.Player{}

	for _, p := range r.Players {
		players = append(players, &model.Player{
			ID:     p.ID,
			Name:   p.Name,
			Type:   model.PlayerType(p.Type),
			TeamID: &p.TeamID,
		})
	}

	var gameID *string
	if r.GameID != "" {
		gameID = &r.GameID
	}

	return &model.Room{
		ID:        r.ID,
		HostID:    r.HostID,
		Status:    model.RoomStatus(r.Status),
		GameID:    gameID,
		Players:   players,
		CreatedAt: r.CreatedAt,
	}
}
