package rooms

import "backend/internal/application/ports"

func (s *Service) ListRoomsDetailed() []ports.Room {
	return s.roomRepo.ListRoomsDetailed()
}

func (s *Service) ListRoomViews() []ports.RoomView {
	return s.roomRepo.ListRoomViews()
}
