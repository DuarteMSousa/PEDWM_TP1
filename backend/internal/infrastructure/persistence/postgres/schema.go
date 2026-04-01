package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func EnsureSchema(ctx context.Context, pool *pgxpool.Pool) error {
	statements := []string{

		// =========================
		// ENUM TYPES
		// =========================
		`DO $$
		BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'friendship_status') THEN
				CREATE TYPE friendship_status AS ENUM ('PENDING', 'ACCEPTED', 'REJECTED');
			END IF;
		END$$;`,

		`DO $$
		BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'room_status') THEN
				CREATE TYPE room_status AS ENUM ('OPEN', 'IN_GAME', 'CLOSED');
			END IF;
		END$$;`,

		`DO $$
		BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'game_status') THEN
				CREATE TYPE game_status AS ENUM ('PENDING', 'IN_PROGRESS', 'FINISHED');
			END IF;
		END$$;`,

		// =========================
		// USERS
		// =========================
		`CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY,
			username VARCHAR NOT NULL UNIQUE,
			password VARCHAR NOT NULL
		)`,

		// =========================
		// FRIENDSHIPS
		// =========================
		`CREATE TABLE IF NOT EXISTS friendships (
			requester_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			addressee_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			status friendship_status NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			PRIMARY KEY (requester_id, addressee_id)
		)`,

		// =========================
		// USER STATS
		// =========================
		`CREATE TABLE IF NOT EXISTS user_stats (
			user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
			games INT NOT NULL DEFAULT 0,
			wins INT NOT NULL DEFAULT 0,
			elo INT NOT NULL DEFAULT 1000
		)`,

		// =========================
		// ROOMS
		// =========================
		`CREATE TABLE IF NOT EXISTS rooms (
			id UUID PRIMARY KEY,
			host_id UUID NOT NULL REFERENCES users(id),
			status room_status NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,

		// =========================
		// ROOM PLAYERS
		// =========================
		`CREATE TABLE IF NOT EXISTS room_players (
			room_id UUID NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			seat INT NOT NULL DEFAULT 0,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			PRIMARY KEY (room_id, user_id)
		)`,

		`ALTER TABLE room_players
			ADD COLUMN IF NOT EXISTS seat INT NOT NULL DEFAULT 0`,

		// =========================
		// GAMES
		// =========================
		`CREATE TABLE IF NOT EXISTS games (
			id UUID PRIMARY KEY,
			room_id UUID NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
			status game_status NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,

		// =========================
		// EVENTS
		// =========================
		`CREATE TABLE IF NOT EXISTS events (
			id UUID PRIMARY KEY,
			game_id UUID NOT NULL REFERENCES games(id) ON DELETE CASCADE, 
			event_type TEXT NOT NULL,
			sequence INT NOT NULL,
			timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			payload JSONB
		)`,

		// =========================
		// INDEXES
		// =========================
		`CREATE INDEX IF NOT EXISTS idx_friendships_requester ON friendships(requester_id)`,
		`CREATE INDEX IF NOT EXISTS idx_friendships_addressee ON friendships(addressee_id)`,

		`CREATE INDEX IF NOT EXISTS idx_rooms_host ON rooms(host_id)`,

		`CREATE INDEX IF NOT EXISTS idx_games_room_id ON games(room_id)`,

		`CREATE INDEX IF NOT EXISTS idx_room_players_room ON room_players(room_id)`,
		`CREATE INDEX IF NOT EXISTS idx_room_players_user ON room_players(user_id)`,

		`CREATE INDEX IF NOT EXISTS idx_events_game ON events(game_id)`,
		`CREATE INDEX IF NOT EXISTS idx_events_user ON events(user_id)`,

		`CREATE UNIQUE INDEX IF NOT EXISTS one_active_game_per_room
			ON games(room_id)
			WHERE status = 'IN_PROGRESS'`,
	}

	for _, stmt := range statements {
		if _, err := pool.Exec(ctx, stmt); err != nil {
			return err
		}
	}

	return nil
}
