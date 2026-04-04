package services

import (
	"backend/internal/application/interfaces"
	"backend/internal/domain/game"
)

type GameService struct {
	gameRepo interfaces.GameRepository
}

func NewGameService(gameRepo interfaces.GameRepository) *GameService {
	return &GameService{gameRepo: gameRepo}
}

func (s *GameService) GetUserGames(userID string) ([]*game.Game, error) {
	games, err := s.gameRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	return games, nil
}
