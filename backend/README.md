# Backend

## Run with Docker (recommended)

From repository root:

```bash
docker compose up --build
```

This starts:

- PostgreSQL on `localhost:5432`
- Backend on `localhost:4000`
- pgAdmin on `localhost:5050`

`DATABASE_URL` is already defined in `docker-compose.yml`, so you do not need
to pass anything in the command line.

pgAdmin login:

- Email: `admin@pedwm.com`
- Password: `admin123`

Postgres server setup in pgAdmin:

- Host: `postgres`
- Port: `5432`
- Database: `sueca`
- User: `sueca`
- Password: `sueca123`

## Run

```bash
go run ./cmd/api
```

Default address: `:4000` (override with `API_ADDR`).

### Persistence mode

- PostgreSQL only.
- `DATABASE_URL` is required and schema is auto-created on startup.

Example:

```bash
$env:DATABASE_URL="postgres://postgres:postgres@localhost:5432/sueca?sslmode=disable"
go run ./cmd/api
```

## Endpoints

- `GET|POST /graphql`
- `GET /ws` (WebSocket upgrade)

## GraphQL

### Create player (unique nickname)

```graphql
mutation {
  createPlayer(nickname: "Alex") {
    id
    nickname
    createdAt
  }
}
```

### Create room

```graphql
mutation {
  createRoom(input: { name: "Mesa Alex", hostPlayerId: "player_1" }) {
    id
    name
    hostPlayerId
    playersCount
    status
  }
}
```

### Create private room (with password)

```graphql
mutation {
  createRoom(
    input: {
      name: "Mesa Privada"
      hostPlayerId: "player_1"
      isPrivate: true
      password: "1234"
    }
  ) {
    id
    isPrivate
    playersCount
  }
}
```

### Join room

```graphql
mutation {
  joinRoom(roomId: "room_1", playerId: "player_2", password: "1234") {
    id
    playersCount
  }
}
```

### Delete room

```graphql
mutation {
  deleteRoom(roomId: "room_1", requesterId: "player_1")
}
```

### Queries

```graphql
query {
  players {
    id
    nickname
  }
  rooms {
    id
    name
    status
    playersCount
    players {
      id
      nickname
    }
  }
}
```

## Notes

- Nickname uniqueness is case-insensitive (`Alex` and `alex` map to the same player).
- If the same nickname logs in again, the same player is returned (idempotent login).
- Data is persisted in PostgreSQL.
- `deleteRoom` can only be executed by the room host.
- Server starts with no pre-created rooms.
- Private rooms require a password on join.
