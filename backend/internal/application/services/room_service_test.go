package services

import (
	"backend/internal/domain/game"
	"backend/internal/domain/room"
	"backend/internal/domain/user"
	"backend/internal/infrastructure/websocket"
	"errors"
	"fmt"
	"testing"
	"time"
)

type fakeRoomRepo struct {
	byID      map[string]*room.Room
	saveErr   error
	deleteErr error
	findErr   error
	allErr    error
	saveCalls int
	delCalls  int
}

func (f *fakeRoomRepo) Save(r *room.Room) error {
	if f.saveErr != nil {
		return f.saveErr
	}
	if f.byID == nil {
		f.byID = map[string]*room.Room{}
	}
	f.byID[r.ID] = r
	f.saveCalls++
	return nil
}

func (f *fakeRoomRepo) FindByID(id string) (*room.Room, error) {
	if f.findErr != nil {
		return nil, f.findErr
	}
	return f.byID[id], nil
}

func (f *fakeRoomRepo) FindAll() ([]*room.Room, error) {
	if f.allErr != nil {
		return nil, f.allErr
	}
	result := make([]*room.Room, 0, len(f.byID))
	for _, r := range f.byID {
		result = append(result, r)
	}
	return result, nil
}

func (f *fakeRoomRepo) Delete(id string) error {
	if f.deleteErr != nil {
		return f.deleteErr
	}
	delete(f.byID, id)
	f.delCalls++
	return nil
}

type fakeRoomSvcGameRepo struct {
	saveErr  error
	findErr  error
	byID     map[string]*game.Game
	saveHits int
}

func (f *fakeRoomSvcGameRepo) Save(g *game.Game) error {
	if f.saveErr != nil {
		return f.saveErr
	}
	if f.byID == nil {
		f.byID = map[string]*game.Game{}
	}
	f.byID[g.ID.String()] = g
	f.saveHits++
	return nil
}

func (f *fakeRoomSvcGameRepo) FindByID(id string) (*game.Game, error) {
	if f.findErr != nil {
		return nil, f.findErr
	}
	return f.byID[id], nil
}

func (f *fakeRoomSvcGameRepo) FindByRoomID(roomID string) ([]*game.Game, error) {
	return nil, nil
}

func (f *fakeRoomSvcGameRepo) GetByUserID(userID string) ([]*game.Game, error) {
	return nil, nil
}

func uniqueID(prefix string) string {
	return fmt.Sprintf("%s_%d", prefix, time.Now().UnixNano())
}

func makeServiceFixtures() (*RoomService, *fakeRoomRepo, *fakeRoomSvcGameRepo, *fakeUserRepo, *websocket.Hub) {
	roomRepo := &fakeRoomRepo{byID: map[string]*room.Room{}}
	gameRepo := &fakeRoomSvcGameRepo{byID: map[string]*game.Game{}}
	userRepo := &fakeUserRepo{
		byID:       map[string]*user.User{},
		byUsername: map[string]*user.User{},
	}
	hub := websocket.GetHubInstance()
	service := NewRoomService(roomRepo, gameRepo, userRepo, nil, hub)
	return service, roomRepo, gameRepo, userRepo, hub
}

func TestRoomServiceCreateRoom(t *testing.T) {
	service, roomRepo, _, userRepo, hub := makeServiceFixtures()

	hostID := uniqueID("host")
	host := &user.User{ID: hostID, Username: "HostUser"}
	userRepo.byID[hostID] = host

	created, err := service.CreateRoom(hostID)
	if err != nil {
		t.Fatalf("expected create room success, got %v", err)
	}
	if created == nil {
		t.Fatal("expected created room")
	}
	defer hub.DeleteRoom(created.ID)

	if roomRepo.saveCalls != 1 {
		t.Fatalf("expected room saveCalls 1, got %d", roomRepo.saveCalls)
	}
	if got := hub.GetRoom(created.ID); got == nil {
		t.Fatal("expected room to be registered in hub")
	}

	userRepo.findByIDErr = errors.New("user find failed")
	if _, err := service.CreateRoom(hostID); err == nil {
		t.Fatal("expected error when user lookup fails")
	}
}

func TestRoomServiceJoinRoom(t *testing.T) {
	service, roomRepo, _, userRepo, hub := makeServiceFixtures()

	hostID := uniqueID("host")
	r, err := room.NewRoom(uniqueID("room"), hostID, "Host")
	if err != nil {
		t.Fatalf("failed to create room fixture: %v", err)
	}
	hub.CreateRoomHub(r)
	defer hub.DeleteRoom(r.ID)

	joinerID := uniqueID("joiner")
	userRepo.byID[joinerID] = &user.User{ID: joinerID, Username: "Joiner"}

	joined, err := service.JoinRoom(r.ID, joinerID)
	if err != nil {
		t.Fatalf("expected join success, got %v", err)
	}
	if joined == nil || joined.Players[joinerID] == nil {
		t.Fatal("expected joined player in room")
	}
	if roomRepo.saveCalls != 1 {
		t.Fatalf("expected room saveCalls 1, got %d", roomRepo.saveCalls)
	}

	if _, err := service.JoinRoom("missing_room", joinerID); !errors.Is(err, ErrRoomNotFound) {
		t.Fatalf("expected ErrRoomNotFound, got %v", err)
	}

	userRepo.findByIDErr = errors.New("user find failed")
	if _, err := service.JoinRoom(r.ID, uniqueID("u")); err == nil {
		t.Fatal("expected error when user lookup fails")
	}
	userRepo.findByIDErr = nil

	// Duplicate player path (AddPlayer error).
	if _, err := service.JoinRoom(r.ID, joinerID); !errors.Is(err, room.ErrPlayerAlreadyInRoom) {
		t.Fatalf("expected ErrPlayerAlreadyInRoom, got %v", err)
	}
}

func TestRoomServiceLeaveRoom(t *testing.T) {
	service, roomRepo, _, _, hub := makeServiceFixtures()

	if _, err := service.LeaveRoom("missing_room", "p1"); !errors.Is(err, ErrRoomNotFound) {
		t.Fatalf("expected ErrRoomNotFound, got %v", err)
	}

	hostID := uniqueID("host")
	r, err := room.NewRoom(uniqueID("room"), hostID, "Host")
	if err != nil {
		t.Fatalf("failed to create room fixture: %v", err)
	}
	if err := r.AddPlayer(uniqueID("p2"), "P2"); err != nil {
		t.Fatalf("failed to add second player: %v", err)
	}
	hub.CreateRoomHub(r)
	defer hub.DeleteRoom(r.ID)

	// OPEN room -> Save should be called.
	_, err = service.LeaveRoom(r.ID, hostID)
	if err != nil {
		t.Fatalf("expected leave success, got %v", err)
	}
	if roomRepo.saveCalls != 1 {
		t.Fatalf("expected saveCalls 1 for OPEN room, got %d", roomRepo.saveCalls)
	}

	// CLOSED room -> Save should be skipped.
	r2, err := room.NewRoom(uniqueID("room"), uniqueID("host"), "Host2")
	if err != nil {
		t.Fatalf("failed to create second room fixture: %v", err)
	}
	hub.CreateRoomHub(r2)
	defer hub.DeleteRoom(r2.ID)

	before := roomRepo.saveCalls
	if _, err := service.LeaveRoom(r2.ID, r2.HostID); err != nil {
		t.Fatalf("expected leave success on single-player room, got %v", err)
	}
	if roomRepo.saveCalls != before {
		t.Fatalf("expected saveCalls unchanged for CLOSED room, got before=%d after=%d", before, roomRepo.saveCalls)
	}
}

func TestRoomServiceDeleteRoom(t *testing.T) {
	service, roomRepo, _, _, hub := makeServiceFixtures()

	r, err := room.NewRoom(uniqueID("room"), uniqueID("host"), "Host")
	if err != nil {
		t.Fatalf("failed to create room fixture: %v", err)
	}
	hub.CreateRoomHub(r)

	if err := service.DeleteRoom(r.ID); err != nil {
		t.Fatalf("expected delete success, got %v", err)
	}
	if roomRepo.delCalls != 1 {
		t.Fatalf("expected delete calls 1, got %d", roomRepo.delCalls)
	}
	if got := hub.GetRoom(r.ID); got != nil {
		t.Fatal("expected room removed from hub")
	}

	roomRepo.deleteErr = errors.New("delete failed")
	if err := service.DeleteRoom(uniqueID("room")); err == nil {
		t.Fatal("expected delete error")
	}
}

func TestRoomServiceGetRoomAndGetRooms(t *testing.T) {
	service, roomRepo, _, _, hub := makeServiceFixtures()

	fromHub, err := room.NewRoom(uniqueID("room"), uniqueID("host"), "Host")
	if err != nil {
		t.Fatalf("failed to create hub room fixture: %v", err)
	}
	hub.CreateRoomHub(fromHub)
	defer hub.DeleteRoom(fromHub.ID)

	got, err := service.GetRoom(fromHub.ID)
	if err != nil {
		t.Fatalf("expected GetRoom success, got %v", err)
	}
	if got == nil || got.ID != fromHub.ID {
		t.Fatalf("expected room from hub %q, got %+v", fromHub.ID, got)
	}

	fromRepo, err := room.NewRoom(uniqueID("room"), uniqueID("host"), "RepoHost")
	if err != nil {
		t.Fatalf("failed to create repo room fixture: %v", err)
	}
	roomRepo.byID[fromRepo.ID] = fromRepo
	got, err = service.GetRoom(fromRepo.ID)
	if err != nil {
		t.Fatalf("expected GetRoom from repo success, got %v", err)
	}
	if got == nil || got.ID != fromRepo.ID {
		t.Fatalf("expected room from repo %q, got %+v", fromRepo.ID, got)
	}

	rooms, err := service.GetRooms()
	if err != nil {
		t.Fatalf("expected GetRooms success, got %v", err)
	}
	if len(rooms) == 0 {
		t.Fatal("expected at least one room from repository")
	}
}

func TestRoomServiceStartGame(t *testing.T) {
	service, roomRepo, gameRepo, _, hub := makeServiceFixtures()

	if _, err := service.StartGame("missing_room"); !errors.Is(err, ErrRoomNotFound) {
		t.Fatalf("expected ErrRoomNotFound, got %v", err)
	}

	r, err := room.NewRoom(uniqueID("room"), uniqueID("host"), "Host")
	if err != nil {
		t.Fatalf("failed to create room fixture: %v", err)
	}
	hub.CreateRoomHub(r)
	defer hub.DeleteRoom(r.ID)

	started, err := service.StartGame(r.ID)
	if err != nil {
		t.Fatalf("expected start game success, got %v", err)
	}
	if started.Status != room.IN_GAME {
		t.Fatalf("expected room status IN_GAME, got %s", started.Status)
	}
	if started.Game == nil {
		t.Fatal("expected game to be created")
	}
	if roomRepo.saveCalls == 0 {
		t.Fatal("expected room to be persisted on StartGame")
	}
	if gameRepo.saveHits == 0 {
		t.Fatal("expected game to be persisted on StartGame")
	}
}

func TestRoomServiceGetGameSnapshot(t *testing.T) {
	service, _, _, _, hub := makeServiceFixtures()

	if _, err := service.GetGameSnapshot("missing_room", "p1"); !errors.Is(err, ErrRoomNotFound) {
		t.Fatalf("expected ErrRoomNotFound, got %v", err)
	}

	noGameRoom, err := room.NewRoom(uniqueID("room"), uniqueID("host"), "Host")
	if err != nil {
		t.Fatalf("failed to create room fixture: %v", err)
	}
	hub.CreateRoomHub(noGameRoom)
	defer hub.DeleteRoom(noGameRoom.ID)
	if _, err := service.GetGameSnapshot(noGameRoom.ID, noGameRoom.HostID); !errors.Is(err, ErrGameNotFound) {
		t.Fatalf("expected ErrGameNotFound when room has no game, got %v", err)
	}

	gameRoom, err := room.NewRoom(uniqueID("room"), uniqueID("host"), "Host")
	if err != nil {
		t.Fatalf("failed to create game room fixture: %v", err)
	}
	if err := gameRoom.CreateGame(); err != nil {
		t.Fatalf("failed to create game in room: %v", err)
	}
	gameRoom.Game.State.Enter()
	hub.CreateRoomHub(gameRoom)
	defer hub.DeleteRoom(gameRoom.ID)

	snapshot, err := service.GetGameSnapshot(gameRoom.ID, gameRoom.HostID)
	if err != nil {
		t.Fatalf("expected snapshot success, got %v", err)
	}
	if snapshot == nil {
		t.Fatal("expected non-nil snapshot")
	}
	if snapshot.RoomID != gameRoom.ID {
		t.Fatalf("expected snapshot roomID %q, got %q", gameRoom.ID, snapshot.RoomID)
	}
	if snapshot.GameID == "" {
		t.Fatal("expected snapshot to include GameID")
	}
	if len(snapshot.MyHand) == 0 {
		t.Fatal("expected snapshot to include player's hand")
	}
}
