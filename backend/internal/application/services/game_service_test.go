package services

import (
	"backend/internal/domain/game"
	"errors"
	"testing"
)

type fakeGameRepo struct {
	byID      map[string]*game.Game
	byUserID  map[string][]*game.Game
	findErr   error
	saveErr   error
	byUserErr error
}

func (f *fakeGameRepo) Save(g *game.Game) error {
	if f.saveErr != nil {
		return f.saveErr
	}
	if f.byID == nil {
		f.byID = map[string]*game.Game{}
	}
	f.byID[g.ID.String()] = g
	return nil
}

func (f *fakeGameRepo) FindByID(id string) (*game.Game, error) {
	if f.findErr != nil {
		return nil, f.findErr
	}
	return f.byID[id], nil
}

func (f *fakeGameRepo) FindByRoomID(roomID string) ([]*game.Game, error) {
	return nil, nil
}

func (f *fakeGameRepo) GetByUserID(userID string) ([]*game.Game, error) {
	if f.byUserErr != nil {
		return nil, f.byUserErr
	}
	return f.byUserID[userID], nil
}

func TestGameServiceGetUserGames(t *testing.T) {
	t.Parallel()

	g1 := &game.Game{}
	repo := &fakeGameRepo{
		byUserID: map[string][]*game.Game{
			"u1": []*game.Game{g1},
		},
	}
	service := NewGameService(repo)

	games, err := service.GetUserGames("u1")
	if err != nil {
		t.Fatalf("expected GetUserGames success, got %v", err)
	}
	if len(games) != 1 || games[0] != g1 {
		t.Fatalf("unexpected games result: %+v", games)
	}

	repo.byUserErr = errors.New("repo error")
	if _, err := service.GetUserGames("u1"); err == nil {
		t.Fatal("expected error when repository fails")
	}
}

func TestGameServiceSetGameStatus(t *testing.T) {
	t.Parallel()

	g := &game.Game{Status: game.PENDING}
	repo := &fakeGameRepo{
		byID: map[string]*game.Game{
			"g1": g,
		},
	}
	service := NewGameService(repo)

	updated, err := service.SetGameStatus("g1", game.FINISHED)
	if err != nil {
		t.Fatalf("expected SetGameStatus success, got %v", err)
	}
	if updated.Status != game.FINISHED {
		t.Fatalf("expected updated status FINISHED, got %s", updated.Status)
	}

	repo.saveErr = errors.New("save failed")
	if _, err := service.SetGameStatus("g1", game.IN_PROGRESS); err == nil {
		t.Fatal("expected error when save fails")
	}
}
