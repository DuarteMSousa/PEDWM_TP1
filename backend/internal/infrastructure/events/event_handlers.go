package events_infrastructure

// Este ficheiro contém as implementações concretas dos CommandHandler
// para cada tipo de comando suportado pelo sistema de WebSocket.
// Cada handler é uma função que recebe o contexto do comando e o payload
// JSON, e interage com os repositórios e comandos de domínio adequados.

import (
	"backend/internal/domain/events"
	"backend/internal/domain/game"
	"errors"
	"log"
)

var (
	ErrMissingPlayer = errors.New("playerId is required")
)

// NewPlayCardHandler cria um handler para o comando "play_card".
// O handler localiza o jogo ativo na sala do cliente, cria um
// PlayCardCommand e executa-o sobre o jogo.
func NewPlayerLeftEventHandler(roomService RoomService) EventHandler {
	return func(event events.Event) error {
		p := event.Payload.(events.PlayerLeftPayload)

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

func NewGameEndedEventHandler(userStatsService UserStatsService, gameService GameService) EventHandler {
	return func(event events.Event) error {
		p := event.Payload.(events.GameEndedPayload)

		for _, team := range p.Teams {
			for _, player := range team.Players {
				won := p.Winner == team.ID
				_, err := userStatsService.RecordGame(player.ID, won)
				if err != nil {
					return err
				}
			}
		}

		_, err := gameService.SetGameStatus(event.GameID, game.FINISHED)
		if err != nil {
			return err
		}

		return nil
	}
}

func NewRoomClosedEventHandler(roomService RoomService) EventHandler {
	return func(event events.Event) error {
		p := event.Payload.(events.RoomClosedPayload)
		log.Printf("handling RoomClosed event for room %s", p.RoomID)

		err := roomService.DeleteRoom(p.RoomID)
		if err != nil {
			log.Printf("Error deleting room %s: %v", p.RoomID, err)
			return err
		}

		log.Printf("Room %s deleted successfully", p.RoomID)
		return nil
	}
}
