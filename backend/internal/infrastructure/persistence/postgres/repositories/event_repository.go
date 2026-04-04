package repositories

import (
	"backend/internal/domain/events"
	"context"

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
		if err := rows.Scan(&event.ID, &event.GameID, &event.Type, &event.Sequence, &event.Timestamp, &event.Payload); err != nil {
			return nil, err
		}
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
		if err := rows.Scan(&event.ID, &event.GameID, &event.Type, &event.Sequence, &event.Timestamp, &event.Payload); err != nil {
			return nil, err
		}
		eventsList = append(eventsList, event)
	}
	return eventsList, nil
}
