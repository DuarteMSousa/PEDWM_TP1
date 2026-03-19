package rooms

import (
	"backend/internal/application/ports"
	domainevents "backend/internal/domain/events"
)

func (s *Service) JoinRoom(roomID string, playerID string, password string) (ports.RoomView, error) {
	view, snapshot, player, err := s.roomRepo.JoinRoom(roomID, playerID, password)
	if err != nil {
		return ports.RoomView{}, err
	}

	s.publish(domainevents.NewPlayerJoinedEvent(roomID, player.ID, player.Nickname))
	s.publish(
		domainevents.NewRoomUpdatedEvent(
			roomID,
			map[string]any{
				"room": roomPayload(snapshot),
			},
		),
	)
	s.publish(
		domainevents.NewRoomUpdatedEvent(
			ports.LobbyRoomID,
			map[string]any{
				"roomId": roomID,
				"room":   roomPayload(snapshot),
			},
		),
	)

	return view, nil
}
