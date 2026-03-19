package players

import "backend/internal/application/ports"

func (s *Service) ListPlayers() []ports.Player {
	return s.playerRepo.ListPlayers()
}

func (s *Service) PlayersByIDs(ids []string) []ports.Player {
	return s.playerRepo.PlayersByIDs(ids)
}
