package players

import "backend/internal/application/ports"

func (s *Service) CreatePlayer(nickname string) (ports.Player, error) {
	return s.playerRepo.CreateOrGetByNickname(nickname)
}
