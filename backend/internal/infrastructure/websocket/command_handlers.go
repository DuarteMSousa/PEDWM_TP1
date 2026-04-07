package websocket

// Este ficheiro contém as implementações concretas dos CommandHandler
// para cada tipo de comando suportado pelo sistema de WebSocket.
// Cada handler é uma função que recebe o contexto do comando e o payload
// JSON, e interage com os repositórios e comandos de domínio adequados.

import (
	command "backend/internal/domain/commands"
	"backend/internal/domain/game"
	bot_strategy "backend/internal/domain/player/botStrategy"
	"backend/internal/domain/room"
	"encoding/json"
	"errors"
)

var (
	ErrNoActiveGame             = errors.New("no active game in this room")
	ErrMissingCard              = errors.New("cardId is required")
	ErrMissingBotStrategyOption = errors.New("bot strategy option is required")
)

// PlayCardPayload representa o payload esperado para o comando play_card.
type PlayCardPayload struct {
	CardID string `json:"cardId"`
}

type BotStrategyChangePayload struct {
	Option bot_strategy.BotStrategyType `json:"option"`
}

// NewPlayCardHandler cria um handler para o comando "play_card".
// O handler localiza o jogo ativo na sala do cliente, cria um
// PlayCardCommand e executa-o sobre o jogo.
func NewPlayCardHandler(hub *Hub) CommandHandler {
	return func(ctx *CommandContext, payload json.RawMessage) error {
		var p PlayCardPayload
		if err := json.Unmarshal(payload, &p); err != nil {
			return err
		}
		if p.CardID == "" {
			return ErrMissingCard
		}

		room, err := findRoom(hub, ctx.RoomID)
		if err != nil {
			return err
		}

		cmd := command.NewPlayCardCommand(ctx.PlayerID, p.CardID)
		if err := cmd.Execute(room); err != nil {
			return err
		}

		return nil
	}
}

// NewBotStrategyHandler cria um handler para o comando "bot_strategy".
// O handler localiza o jogo ativo na sala do cliente, cria um
// BotStrategyCommand e executa-o sobre o jogo.
func NewChangeBotStrategyHandler(hub *Hub) CommandHandler {
	return func(ctx *CommandContext, payload json.RawMessage) error {
		var p BotStrategyChangePayload
		if err := json.Unmarshal(payload, &p); err != nil {
			return err
		}
		if p.Option == "" {
			return ErrMissingBotStrategyOption
		}

		room, err := findRoom(hub, ctx.RoomID)
		if err != nil {
			return err
		}

		cmd := command.NewChangeBotStrategyCommand(p.Option)
		if err := cmd.Execute(room); err != nil {
			return err
		}

		return nil
	}
}

// findGame procura o jogo.
func findGame(hub *Hub, roomID string) (*game.Game, error) {
	roomHub := hub.GetRoomHub(roomID)
	if roomHub == nil {
		return nil, ErrNoActiveGame
	}

	game := roomHub.room.Game
	if game != nil {
		return game, nil
	}

	return nil, ErrNoActiveGame
}

// findRoom procura a sala.
func findRoom(hub *Hub, roomID string) (*room.Room, error) {
	roomHub := hub.GetRoomHub(roomID)
	if roomHub == nil {
		return nil, ErrNoActiveGame
	}
	return roomHub.room, nil
}
