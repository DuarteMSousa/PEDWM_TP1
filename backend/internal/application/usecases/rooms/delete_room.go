package rooms

import (
	"backend/internal/application/ports"
	domainevents "backend/internal/domain/events"
)

func (s *Service) DeleteRoom(roomID string, requesterID string) error {
	room, err := s.roomRepo.DeleteRoom(roomID, requesterID)
	if err != nil {
		return err
	}

	deletedPayload := map[string]any{
		"roomId": roomID,
	}
	s.publish(domainevents.NewRoomDeletedEvent(roomID, deletedPayload))
	s.publish(
		domainevents.NewRoomDeletedEvent(
			ports.LobbyRoomID,
			map[string]any{
				"roomId": roomID,
				"room":   roomPayload(room),
			},
		),
	)

	return nil
}
