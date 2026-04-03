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

func (r *EventPostgresRepository) FindByGameID(gameID string) ([]events.Event, error) {
	ctx := context.Background()

	rows, err := r.pool.Query(ctx, `
		SELECT e.id, e.game_id, g.room_id, e.event_type, e.sequence, e.timestamp, e.payload
		FROM events e
		JOIN games g ON g.id = e.game_id
		WHERE e.game_id = $1
		ORDER BY e.sequence ASC, e.timestamp ASC
	`, gameID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanEvents(rows)
}

func (r *EventPostgresRepository) FindByType(eventType events.EventType) ([]events.Event, error) {
	ctx := context.Background()

	rows, err := r.pool.Query(ctx, `
		SELECT e.id, e.game_id, g.room_id, e.event_type, e.sequence, e.timestamp, e.payload
		FROM events e
		JOIN games g ON g.id = e.game_id
		WHERE e.event_type = $1
		ORDER BY e.timestamp DESC
	`, string(eventType))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanEvents(rows)
}

type scannableEventRows interface {
	Next() bool
	Scan(dest ...any) error
	Err() error
}

func scanEvents(rows scannableEventRows) ([]events.Event, error) {
	result := make([]events.Event, 0)

	for rows.Next() {
		var (
			e          events.Event
			eventType  string
			payloadRaw []byte
		)

		if err := rows.Scan(
			&e.ID,
			&e.GameID,
			&e.RoomID,
			&eventType,
			&e.Sequence,
			&e.Timestamp,
			&payloadRaw,
		); err != nil {
			return nil, err
		}

		e.Type = events.EventType(eventType)
		if len(payloadRaw) > 0 {
			var payload any
			if err := json.Unmarshal(payloadRaw, &payload); err != nil {
				return nil, err
			}
			e.Payload = payload
		}

		result = append(result, e)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
