package application

import (
	"backend/internal/application/interfaces"
	"backend/internal/domain/card"
	"backend/internal/domain/events"
	"backend/internal/domain/room"
	"backend/internal/domain/trick"
	"backend/internal/infrastructure/persistence"
	"backend/internal/infrastructure/websocket"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrRoomNotFound = errors.New("room not found")
	ErrGameNotFound = errors.New("game not found")
)

type GameSnapshot struct {
	RoomID          string
	GameID          string
	Status          string
	TrumpSuit       string
	CurrentPlayerID string
	MyHand          []card.Card
	TablePlays      []trick.Play
	Scores          map[string]int
}

type RoomService struct {
	repo         interfaces.RoomRepository
	gameRepo     interfaces.GameRepository
	userRepo     interfaces.UserRepository
	eventService *EventService
	hub          *websocket.Hub
}

func NewRoomService(repo interfaces.RoomRepository, gameRepo interfaces.GameRepository, userRepo interfaces.UserRepository, eventService *EventService, hub *websocket.Hub) *RoomService {
	return &RoomService{repo: repo, gameRepo: gameRepo, userRepo: userRepo, eventService: eventService, hub: hub}
}

func (s *RoomService) CreateRoom(hostID string) (*room.Room, error) {
	user, err := s.userRepo.FindByID(hostID)
	if err != nil || user == nil {
		return nil, err
	}

	r, err := room.NewRoom(uuid.New().String(), hostID, user.Username)
	if err != nil {
		return nil, err
	}

	eventBus := events.NewEventBus()
	r.SetEventBus(eventBus)

	wsObserver := websocket.NewWebSocketObserver(s.hub)
	eventBus.Subscribe(wsObserver)

	persistenceObserver := persistence.NewEventPersistanceObserver(s.eventService)
	eventBus.Subscribe(persistenceObserver)

	s.hub.CreateRoomHub(r)

	return r, s.repo.Save(r)
}

func (s *RoomService) JoinRoom(roomID, userID string) (*room.Room, error) {
	rm, err := s.ensureRealtimeRoom(roomID)
	if err != nil || rm == nil {
		return nil, ErrRoomNotFound
	}

	user, err := s.userRepo.FindByID(userID)
	if err != nil || user == nil {
		return nil, err
	}

	if err := rm.AddPlayer(userID, user.Username); err != nil {
		if !errors.Is(err, room.ErrPlayerAlreadyInRoom) {
			return nil, err
		}
	}

	if err := s.repo.Save(rm); err != nil {
		return nil, err
	}
	s.publishRoomSyncedEvent(rm)

	return rm, nil
}

func (s *RoomService) LeaveRoom(roomID, userID string) (*room.Room, error) {
	rm, err := s.ensureRealtimeRoom(roomID)
	if err != nil || rm == nil {
		return nil, ErrRoomNotFound
	}

	if err := rm.RemovePlayer(userID); err != nil {
		return nil, err
	}
	s.publishRoomSyncedEvent(rm)

	if len(rm.Players) == 0 {
		if err := s.repo.Delete(rm.ID); err != nil {
			return nil, err
		}
		return rm, nil
	}

	if err := s.repo.Save(rm); err != nil {
		return nil, err
	}

	return rm, nil
}

func (s *RoomService) StartGame(roomID string) (*room.Room, error) {
	rm, err := s.ensureRealtimeRoom(roomID)
	if err != nil || rm == nil {
		return nil, ErrRoomNotFound
	}

	if err := rm.StartGame(); err != nil {
		return nil, err
	}

	if err := s.repo.Save(rm); err != nil {
		return nil, err
	}
	if rm.Game != nil {
		if err := s.gameRepo.Save(rm.Game); err != nil {
			return nil, err
		}
	}
	s.publishRoomSyncedEvent(rm)

	return rm, nil
}

func (s *RoomService) publishRoomSyncedEvent(room *room.Room) {
	if room == nil || room.EventBus == nil {
		return
	}

	room.EventBus.Publish(events.Event{
		ID:        uuid.NewString(),
		Type:      events.EventType("ROOM_SYNCED"),
		RoomID:    room.ID,
		Timestamp: time.Now().UTC(),
		Payload: map[string]any{
			"status":       room.Status,
			"playersCount": len(room.Players),
		},
	})
}

func (s *RoomService) GetRoom(id string) (*room.Room, error) {
	if liveRoom := s.hub.GetRoom(id); liveRoom != nil {
		return liveRoom, nil
	}
	return s.repo.FindByID(id)
}

func (s *RoomService) GetRooms() ([]*room.Room, error) {
	return s.repo.FindAll()
}

func (s *RoomService) GetGameSnapshot(roomID, playerID string) (*GameSnapshot, error) {
	rm, err := s.ensureRealtimeRoom(roomID)
	if err != nil || rm == nil {
		return nil, ErrRoomNotFound
	}
	if rm.Game == nil {
		return nil, ErrGameNotFound
	}

	g := rm.Game
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
		RoomID:          rm.ID,
		GameID:          g.ID.String(),
		Status:          string(g.Status),
		TrumpSuit:       string(currentRound.TrumpSuit),
		CurrentPlayerID: currentPlayerID,
		MyHand:          myHand,
		TablePlays:      tablePlays,
		Scores:          scores,
	}, nil
}

func (s *RoomService) ensureRealtimeRoom(roomID string) (*room.Room, error) {
	if liveRoom := s.hub.GetRoom(roomID); liveRoom != nil {
		if liveRoom.EventBus == nil {
			eventBus := events.NewEventBus()
			eventBus.Subscribe(websocket.NewWebSocketObserver(s.hub))
			eventBus.Subscribe(persistence.NewEventPersistanceObserver(s.eventService))
			liveRoom.SetEventBus(eventBus)
		}
		return liveRoom, nil
	}

	rm, err := s.repo.FindByID(roomID)
	if err != nil || rm == nil {
		return nil, ErrRoomNotFound
	}

	eventBus := events.NewEventBus()
	eventBus.Subscribe(websocket.NewWebSocketObserver(s.hub))
	eventBus.Subscribe(persistence.NewEventPersistanceObserver(s.eventService))
	rm.SetEventBus(eventBus)

	s.hub.CreateRoomHub(rm)
	return rm, nil
}
