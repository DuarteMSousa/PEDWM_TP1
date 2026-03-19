package rooms

import (
	"backend/internal/application/ports"
	domainevents "backend/internal/domain/events"
)

func (s *Service) publish(event domainevents.Event) {
	if s == nil || s.publisher == nil {
		return
	}
	s.publisher.Publish(event)
}

func roomPayload(room ports.Room) map[string]any {
	return map[string]any{
		"id":           room.ID,
		"name":         room.Name,
		"hostPlayerId": room.HostID,
		"status":       room.Status,
		"maxPlayers":   room.MaxPlayers,
		"isPrivate":    room.IsPrivate,
		"playersCount": len(room.PlayerIDs),
		"playerIds":    room.PlayerIDs,
	}
}
