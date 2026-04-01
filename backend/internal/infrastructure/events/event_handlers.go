package events_infrastructure

// Este ficheiro contém as implementações concretas dos CommandHandler
// para cada tipo de comando suportado pelo sistema de WebSocket.
// Cada handler é uma função que recebe o contexto do comando e o payload
// JSON, e interage com os repositórios e comandos de domínio adequados.

import (
	"backend/internal/domain/events"
	"encoding/json"
	"errors"
)

var (
	ErrMissingPlayer = errors.New("playerId is required")
)

// NewPlayCardHandler cria um handler para o comando "play_card".
// O handler localiza o jogo ativo na sala do cliente, cria um
// PlayCardCommand e executa-o sobre o jogo.
func NewPlayerLeftEventHandler(roomService RoomService) EventHandler {
	return func(payload json.RawMessage) error {
		var p events.PlayerLeftPayload
		if err := json.Unmarshal(payload, &p); err != nil {
			return err
		}
		if p.PlayerID == "" {
			return ErrMissingPlayer
		}

		_, err := roomService.LeaveRoom(p.RoomID, p.PlayerID)
		if err != nil {
			return err
		}

		return nil
	}
}

func NewGameEndedEventHandler(userStatsService UserStatsService) EventHandler {
	return func(payload json.RawMessage) error {
		var p events.GameEndedPayload
		if err := json.Unmarshal(payload, &p); err != nil {
			return err
		}

		for team := range p.Teams {
			for _, player := range p.Teams[team].Players {
				won := p.Winner == team
				_, err := userStatsService.RecordGame(player.ID, won)
				if err != nil {
					return err
				}
			}
		}

		return nil
	}
}
