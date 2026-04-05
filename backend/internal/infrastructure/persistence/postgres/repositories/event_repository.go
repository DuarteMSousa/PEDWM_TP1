package repositories

import (
	"backend/internal/domain/events"
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5/pgxpool"
)

type EventPostgresRepository struct {
	pool *pgxpool.Pool
}

func NewEventPostgresRepository(pool *pgxpool.Pool) *EventPostgresRepository {
	return &EventPostgresRepository{pool: pool}
}

func (r *EventPostgresRepository) Save(event events.Event) error {
	ctx := context.Background()

	_, err := r.pool.Exec(ctx, `
        INSERT INTO events (id, game_id, event_type, sequence, timestamp, payload)
        VALUES ($1, $2, $3, $4, $5, $6)
        ON CONFLICT (id) DO UPDATE
        SET event_type = EXCLUDED.event_type,
            sequence   = EXCLUDED.sequence,
            timestamp  = EXCLUDED.timestamp,
            payload    = EXCLUDED.payload
    `,
		event.ID,
		event.GameID,
		event.Type,
		event.Sequence,
		event.Timestamp,
		event.Payload,
	)

	return err
}

func (r *EventPostgresRepository) FindByRoomID(roomID string) ([]events.Event, error) {
	ctx := context.Background()
	rows, err := r.pool.Query(ctx, `
		SELECT e.id, e.game_id, e.event_type, e.sequence, e.timestamp, e.payload
		FROM events e
		JOIN games g ON e.game_id = g.id
		WHERE g.room_id = $1
		ORDER BY e.sequence ASC
	`, roomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var eventsList []events.Event
	for rows.Next() {
		var event events.Event
		var rawPayload []byte
		if err := rows.Scan(&event.ID, &event.GameID, &event.Type, &event.Sequence, &event.Timestamp, &rawPayload); err != nil {
			return nil, err
		}
		event.Payload = UnmarshalPayload(event.Type, rawPayload)
		eventsList = append(eventsList, event)
	}
	return eventsList, nil
}

func (r *EventPostgresRepository) FindByGameID(gameID string) ([]events.Event, error) {
	ctx := context.Background()
	rows, err := r.pool.Query(ctx, `
		SELECT id, game_id, event_type, sequence, timestamp, payload
		FROM events
		WHERE game_id = $1
		ORDER BY sequence ASC
	`, gameID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var eventsList []events.Event
	for rows.Next() {
		var event events.Event
		var rawPayload []byte
		if err := rows.Scan(&event.ID, &event.GameID, &event.Type, &event.Sequence, &event.Timestamp, &rawPayload); err != nil {
			return nil, err
		}
		event.Payload = UnmarshalPayload(event.Type, rawPayload)
		eventsList = append(eventsList, event)
	}
	return eventsList, nil
}

func UnmarshalPayload(eventType events.EventType, raw []byte) any {
	if len(raw) == 0 {
		return nil
	}
	var target any
	switch eventType {
	case events.EventPlayerJoined:
		target = &events.PlayerJoinedPayload{}
	case events.EventPlayerLeft:
		target = &events.PlayerLeftPayload{}
	case events.EventRoundStarted:
		target = &events.RoundStartedPayload{}
	case events.EventTrickStarted:
		target = &events.TrickStartedPayload{}
	case events.EventTrumpRevealed:
		target = &events.TrumpRevealedPayload{}
	case events.EventCardDealt:
		target = &events.CardDealtPayload{}
	case events.EventCardPlayed:
		target = &events.CardPlayedPayload{}
	case events.EventTurnChanged:
		target = &events.TurnChangedPayload{}
	case events.EventTrickEnded:
		target = &events.TrickEndedPayload{}
	case events.EventRoundEnded:
		target = &events.RoundEndedPayload{}
	case events.EventGameScore:
		target = &events.GameScorePayload{}
	case events.EventGameStarted:
		target = &events.GameStartedPayload{}
	case events.EventGameEnded:
		target = &events.GameEndedPayload{}
	case events.EventRoomClosed:
		target = &events.RoomClosedPayload{}
	case events.EventBotStrategyChanged:
		target = &events.BotStrategyChangedPayload{}
	default:
		return nil
	}
	if err := json.Unmarshal(raw, target); err != nil {
		return nil
	}
	return target
}
