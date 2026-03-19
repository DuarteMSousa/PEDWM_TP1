package players

import "backend/internal/application/ports"

type Service struct {
	playerRepo ports.PlayerRepository
}

func NewService(playerRepo ports.PlayerRepository) *Service {
	return &Service{playerRepo: playerRepo}
}
