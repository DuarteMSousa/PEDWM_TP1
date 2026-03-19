package games

import (
	"backend/internal/application/ports"
	"errors"
)

var ErrNotImplemented = errors.New("game use case not implemented yet")

type Service struct {
	gameRepo ports.GameRepository
}

func NewService(gameRepo ports.GameRepository) *Service {
	return &Service{gameRepo: gameRepo}
}
