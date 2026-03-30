package websocket

// Este ficheiro contém as implementações concretas dos CommandHandler
// para cada tipo de comando suportado pelo sistema de WebSocket.
// Cada handler é uma função que recebe o contexto do comando e o payload
// JSON, e interage com os repositórios e comandos de domínio adequados.

import (
	"backend/internal/application/interfaces"
	command "backend/internal/domain/commands"
	"backend/internal/domain/game"
	"encoding/json"
	"errors"
)

var (
	ErrNoActiveGame = errors.New("no active game in this room")
	ErrMissingCard  = errors.New("cardId is required")
)

// PlayCardPayload representa o payload esperado para o comando play_card.
type PlayCardPayload struct {
	CardID string `json:"cardId"`
}

// NewPlayCardHandler cria um handler para o comando "play_card".
// O handler localiza o jogo ativo na sala do cliente, cria um
// PlayCardCommand e executa-o sobre o jogo.
func NewPlayCardHandler(gameRepo interfaces.GameRepository) CommandHandler {
	return func(ctx *CommandContext, payload json.RawMessage) error {
		var p PlayCardPayload
		if err := json.Unmarshal(payload, &p); err != nil {
			return err
		}
		if p.CardID == "" {
			return ErrMissingCard
		}

		activeGame, err := findActiveGame(gameRepo, ctx.RoomID)
		if err != nil {
			return err
		}

		cmd := command.NewPlayCardCommand(ctx.PlayerID, p.CardID)
		cmd.Execute(activeGame)

		return gameRepo.Save(activeGame)
	}
}

// findActiveGame procura o jogo com estado IN_PROGRESS associado a uma sala.
func findActiveGame(gameRepo interfaces.GameRepository, roomID string) (*game.Game, error) {
	games, err := gameRepo.FindByRoomID(roomID)
	if err != nil {
		return nil, err
	}

	for _, g := range games {
		if g.Status == game.IN_PROGRESS {
			return g, nil
		}
	}

	return nil, ErrNoActiveGame
}
