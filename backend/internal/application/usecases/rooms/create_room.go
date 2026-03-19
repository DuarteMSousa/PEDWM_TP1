package rooms

import (
	"backend/internal/application/ports"
	domainevents "backend/internal/domain/events"
)

func (s *Service) CreateRoom(
	name string,
	hostPlayerID string,
	maxPlayers int,
	isPrivate bool,
	password string,
) (ports.Room, error) {
	room, err := s.roomRepo.CreateRoom(name, hostPlayerID, maxPlayers, isPrivate, password)
	if err != nil {
		return ports.Room{}, err
	}

	s.publish(
		domainevents.NewRoomCreatedEvent(
			ports.LobbyRoomID,
			hostPlayerID,
		).WithPayload(roomPayload(room)),
	)
	s.publish(domainevents.NewRoomCreatedEvent(room.ID, hostPlayerID))

	return room, nil
}
