package players

import "backend/internal/application/ports"

func (s *Service) GetPlayer(playerID string) (ports.Player, bool) {
	return s.playerRepo.GetPlayer(playerID)
}
