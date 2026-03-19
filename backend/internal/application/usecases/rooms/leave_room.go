package rooms

import (
	"backend/internal/application/ports"
	domainevents "backend/internal/domain/events"
)

func (s *Service) LeaveRoom(roomID string, playerID string) (ports.RoomView, error) {
	view, snapshot, err := s.roomRepo.LeaveRoom(roomID, playerID)
	if err != nil {
		return ports.RoomView{}, err
	}

	s.publish(domainevents.NewPlayerLeftEvent(roomID, playerID))
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
