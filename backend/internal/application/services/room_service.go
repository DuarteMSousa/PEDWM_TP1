package services

import (
	"backend/internal/application/interfaces"
	"backend/internal/domain/card"
	"backend/internal/domain/events"
	"backend/internal/domain/room"
	"backend/internal/domain/team"
	"backend/internal/domain/trick"
	events_infrastructure "backend/internal/infrastructure/events"
	"backend/internal/infrastructure/websocket"
	"errors"
	"log/slog"

	"github.com/google/uuid"
)

var (
	ErrRoomNotFound = errors.New("room not found")
	ErrGameNotFound = errors.New("game not found")
)

// GameSnapshot contains the current state of a game visible to a player.
type GameSnapshot struct {
	RoomID          string
	GameID          string
	Status          string
	TrumpSuit       string
	CurrentPlayerID string
	MyHand          []card.Card
	TablePlays      []trick.Play
	Teams           map[string]*team.Team
	Scores          map[string]int
}

// RoomService orchestrates room operations: create, join, leave, start game.
type RoomService struct {
	repo            interfaces.RoomRepository
	gameRepo        interfaces.GameRepository
	userRepo        interfaces.UserRepository
	eventService    *EventService
	eventDispatcher *events_infrastructure.EventDispatcher
	hub             *websocket.Hub
}

// NewRoomService creates a new RoomService.
func NewRoomService(repo interfaces.RoomRepository, gameRepo interfaces.GameRepository, userRepo interfaces.UserRepository, eventService *EventService, hub *websocket.Hub) *RoomService {
	return &RoomService{repo: repo, gameRepo: gameRepo, userRepo: userRepo, eventService: eventService, hub: hub}
}

// CreateRoom creates a new room and configures the event bus and observers.
func (s *RoomService) CreateRoom(hostID string) (*room.Room, error) {
	slog.Info("creating room", "hostID", hostID)

	user, err := s.userRepo.FindByID(hostID)
	if err != nil || user == nil {
		slog.Warn("room creation failed: host not found", "hostID", hostID)
		return nil, err
	}

	r, err := room.NewRoom(uuid.New().String(), hostID, user.Username)
	if err != nil {
		slog.Error("failed to create room", "hostID", hostID, "error", err)
		return nil, err
	}

	eventBus := events.NewEventBus()
	r.SetEventBus(eventBus)

	wsObserver := websocket.NewWebSocketObserver(s.hub)
	eventBus.Subscribe(wsObserver)

	if s.eventDispatcher == nil {
		s.eventDispatcher = events_infrastructure.GetEventDispatcherInstance()
	}

	persistenceObserver := events_infrastructure.NewEventPersistanceObserver(s.eventService, s.eventDispatcher)
	eventBus.Subscribe(persistenceObserver)

	s.hub.CreateRoomHub(r)

	slog.Info("room created successfully", "roomID", r.ID, "hostID", hostID)
	return r, s.repo.Save(r)
}

// JoinRoom adds a player to an existing room.
func (s *RoomService) JoinRoom(roomID, userID string) (*room.Room, error) {
	slog.Info("joining room", "roomID", roomID, "userID", userID)

	room := s.hub.GetRoom(roomID)
	if room == nil {
		slog.Warn("join room failed: room not found", "roomID", roomID)
		return nil, ErrRoomNotFound
	}

	user, err := s.userRepo.FindByID(userID)
	if err != nil || user == nil {
		slog.Warn("join room failed: user not found", "userID", userID)
		return nil, err
	}

	if err := room.AddPlayer(userID, user.Username); err != nil {
		slog.Warn("failed to add player to room", "roomID", roomID, "userID", userID, "error", err)
		return nil, err
	}

	if err := s.repo.Save(room); err != nil {
		slog.Error("failed to persist room after join", "roomID", roomID, "error", err)
		return nil, err
	}

	slog.Info("player joined room", "roomID", roomID, "userID", userID, "totalPlayers", len(room.Players))
	return room, nil
}

// LeaveRoom removes a player from a room. Closes the room if it becomes empty.
func (s *RoomService) LeaveRoom(roomID, userID string) (*room.Room, error) {
	slog.Info("leaving room", "roomID", roomID, "userID", userID)

	rm := s.hub.GetRoom(roomID)
	if rm == nil {
		slog.Warn("leave room failed: room not found", "roomID", roomID)
		return nil, ErrRoomNotFound
	}

	if err := rm.RemovePlayer(userID); err != nil {
		slog.Error("failed to remove player from room", "roomID", roomID, "userID", userID, "error", err)
		return nil, err
	}

	if rm.Status != room.CLOSED {
		if err := s.repo.Save(rm); err != nil {
			slog.Error("failed to persist room after leave", "roomID", roomID, "error", err)
			return nil, err
		}
	}

	slog.Info("player left room", "roomID", roomID, "userID", userID, "remainingPlayers", len(rm.Players))
	return rm, nil
}

// DeleteRoom removes a room from persistence.
func (s *RoomService) DeleteRoom(roomID string) error {
	slog.Info("deleting room", "roomID", roomID)
	s.hub.DeleteRoom(roomID)
	if err := s.repo.Delete(roomID); err != nil {
		slog.Error("failed to delete room", "roomID", roomID, "error", err)
		return err
	}
	slog.Info("room deleted", "roomID", roomID)
	return nil
}

// StartGame starts a Sueca game in the specified room.
func (s *RoomService) StartGame(roomID string) (*room.Room, error) {
	slog.Info("starting game", "roomID", roomID)

	room := s.hub.GetRoom(roomID)
	if room == nil {
		slog.Warn("start game failed: room not found", "roomID", roomID)
		return nil, ErrRoomNotFound
	}

	if err := room.CreateGame(); err != nil {
		slog.Error("failed to create game", "roomID", roomID, "error", err)
		return nil, err
	}

	if err := s.repo.Save(room); err != nil {
		slog.Error("failed to persist room after game creation", "roomID", roomID, "error", err)
		return nil, err
	}
	if room.Game != nil {
		if err := s.gameRepo.Save(room.Game); err != nil {
			slog.Error("failed to persist game", "roomID", roomID, "gameID", room.Game.ID, "error", err)
			return nil, err
		}
	}

	room.Game.State.Enter()

	slog.Info("game started successfully", "roomID", roomID, "gameID", room.Game.ID)
	return room, nil
}

// GetRoom returns a room by ID (in-memory hub or database).
func (s *RoomService) GetRoom(id string) (*room.Room, error) {
	if r := s.hub.GetRoom(id); r != nil {
		return r, nil
	}
	return s.repo.FindByID(id)
}

// GetRooms returns all rooms.
func (s *RoomService) GetRooms() ([]*room.Room, error) {
	return s.repo.FindAll()
}

// GetGameSnapshot returns the current game state visible to a player.
func (s *RoomService) GetGameSnapshot(roomID, playerID string) (*GameSnapshot, error) {
	slog.Debug("fetching game snapshot", "roomID", roomID, "playerID", playerID)
	room := s.hub.GetRoom(roomID)
	if room == nil {
		return nil, ErrRoomNotFound
	}
	if room.Game == nil {
		return nil, ErrGameNotFound
	}

	g := room.Game
	currentRound := g.CurrentRound()
	if currentRound == nil {
		return nil, ErrGameNotFound
	}

	requestedPlayer, err := g.GetPlayer(playerID)
	if err != nil {
		return nil, err
	}

	scores := make(map[string]int, len(g.Score))
	for teamID, value := range g.Score {
		scores[teamID] = value
	}

	myHand := make([]card.Card, 0)
	if requestedPlayer.Hand != nil {
		myHand = append(myHand, requestedPlayer.Hand.Cards...)
	}

	tablePlays := make([]trick.Play, 0)
	currentPlayerID := ""
	if currentRound.CurrentTrick != nil {
		tablePlays = append(tablePlays, currentRound.CurrentTrick.Plays...)
		nextPlayerID, nextErr := currentRound.CurrentTrick.TurnOrder.Next()
		if nextErr == nil {
			currentPlayerID = nextPlayerID
		}
	}

	return &GameSnapshot{
		RoomID:          room.ID,
		GameID:          g.ID.String(),
		Status:          string(g.Status),
		TrumpSuit:       string(currentRound.TrumpSuit),
		CurrentPlayerID: currentPlayerID,
		MyHand:          myHand,
		TablePlays:      tablePlays,
		Scores:          scores,
		Teams:           g.Teams,
	}, nil
}
