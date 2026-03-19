package ports

import "time"

const (
	DefaultMaxPlayers = 4
	LobbyRoomID       = "lobby"
)

type RoomStatus string

const (
	RoomStatusOpen   RoomStatus = "OPEN"
	RoomStatusInGame RoomStatus = "IN_GAME"
	RoomStatusClosed RoomStatus = "CLOSED"
)

type Player struct {
	ID        string    `json:"id"`
	Nickname  string    `json:"nickname"`
	CreatedAt time.Time `json:"createdAt"`
}

type Room struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	HostID     string     `json:"hostId"`
	Status     RoomStatus `json:"status"`
	MaxPlayers int        `json:"maxPlayers"`
	IsPrivate  bool       `json:"isPrivate"`
	PlayerIDs  []string   `json:"playerIds"`
	CreatedAt  time.Time  `json:"createdAt"`
}

type RoomView struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	PlayersCount int    `json:"playersCount"`
	MaxPlayers   int    `json:"maxPlayers"`
	IsPrivate    bool   `json:"isPrivate"`
}
