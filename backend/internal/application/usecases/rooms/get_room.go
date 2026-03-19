package rooms

import "backend/internal/application/ports"

func (s *Service) GetRoom(roomID string) (ports.Room, bool) {
	return s.roomRepo.GetRoom(roomID)
}
