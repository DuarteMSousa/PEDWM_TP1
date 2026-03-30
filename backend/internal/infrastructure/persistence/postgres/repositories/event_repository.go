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
