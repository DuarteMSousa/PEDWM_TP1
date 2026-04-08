package events_infrastructure

import (
	"backend/internal/domain/events"
	"backend/internal/domain/game"
	"backend/internal/domain/player"
	"backend/internal/domain/room"
	"errors"
	"log/slog"
)

var (
	ErrMissingPlayer = errors.New("playerId is required")
)

// NewPlayerLeftEventHandler creates a handler for the PLAYER_LEFT event.
// It removes the player from the room through the RoomService.
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

// NewGameEndedEventHandler creates a handler for the GAME_ENDED event.
// It records the game result for each player and updates the game status.
func NewGameEndedEventHandler(userStatsService UserStatsService, gameService GameService, roomService RoomService) EventHandler {
	return func(event events.Event) error {
		payload := event.Payload.(events.GameEndedPayload)

		for _, team := range payload.Teams {
			for _, p := range team.Players {
				if p.Type == player.BOT {
					continue
				}
				won := payload.Winner == team.ID
				_, err := userStatsService.RecordGame(p.ID, won)
				if err != nil {
					slog.Error("error recording game result for player", "playerID", p.ID, "won", won, "error", err)
				}
			}
		}

		_, err := gameService.SetGameStatus(event.GameID, game.FINISHED)
		if err != nil {
			return err
		}

		_, error := roomService.SetRoomStatus(event.RoomID, room.OPEN)
		if error != nil {
			return error
		}

		return nil
	}
}

// NewRoomClosedEventHandler creates a handler for the ROOM_CLOSED event.
// It deletes the room from persistence.
func NewRoomClosedEventHandler(roomService RoomService) EventHandler {
	return func(event events.Event) error {
		p := event.Payload.(events.RoomClosedPayload)
		slog.Info("processing ROOM_CLOSED event", "roomID", p.RoomID)

		err := roomService.DeleteRoom(p.RoomID)
		if err != nil {
			slog.Error("error deleting room", "roomID", p.RoomID, "error", err)
			return err
		}

		slog.Info("room deleted successfully", "roomID", p.RoomID)
		return nil
	}
}
