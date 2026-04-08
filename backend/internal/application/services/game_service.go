package services

import (
	"backend/internal/application/interfaces"
	"backend/internal/domain/game"
	"log/slog"
)

// GameService manages the retrieval and updating of game states.
type GameService struct {
	gameRepo interfaces.GameRepository
}

// NewGameService creates a new GameService.
func NewGameService(gameRepo interfaces.GameRepository) *GameService {
	return &GameService{gameRepo: gameRepo}
}

// GetUserGames returns the list of games for a user.
func (s *GameService) GetUserGames(userID string) ([]*game.Game, error) {
	slog.Debug("retrieving user games", "userID", userID)
	games, err := s.gameRepo.GetByUserID(userID)
	if err != nil {
		slog.Error("error retrieving user games", "userID", userID, "error", err)
		return nil, err
	}

	return games, nil
}

// GetGame returns a game by ID.
func (s *GameService) GetGame(id string) (*game.Game, error) {
	return s.gameRepo.FindByID(id)
}

// SetGameStatus updates the status of a game.
func (s *GameService) SetGameStatus(gameID string, status game.GameStatus) (*game.Game, error) {
	slog.Info("updating game status", "gameID", gameID, "status", status)

	g, err := s.gameRepo.FindByID(gameID)
	if err != nil {
		slog.Error("error finding game", "gameID", gameID, "error", err)
		return nil, err
	}

	g.Status = status
	if err := s.gameRepo.Save(g); err != nil {
		slog.Error("error persisting game status", "gameID", gameID, "error", err)
		return nil, err
	}

	slog.Info("game status updated", "gameID", gameID, "status", status)
	return g, nil
}
