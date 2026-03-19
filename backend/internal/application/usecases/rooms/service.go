package rooms

import "backend/internal/application/ports"

type Service struct {
	roomRepo  ports.RoomRepository
	publisher ports.EventPublisher
}

func NewService(roomRepo ports.RoomRepository, publisher ports.EventPublisher) *Service {
	return &Service{
		roomRepo:  roomRepo,
		publisher: publisher,
	}
}
