package services

import (
	"backend/internal/domain/user"
	"errors"
	"testing"
)

func TestUserStatsServiceGetByUserID(t *testing.T) {
	t.Parallel()

	statsRepo := &fakeUserStatsRepo{
		byUserID: map[string]*user.UserStats{
			"u1": {UserId: "u1", Games: 3, Wins: 2, Elo: 1020},
		},
	}
	service := NewUserStatsService(statsRepo, &fakeUserRepo{})

	stats, err := service.GetByUserID("u1")
	if err != nil {
		t.Fatalf("expected stats fetch success, got %v", err)
	}
	if stats.Games != 3 || stats.Wins != 2 || stats.Elo != 1020 {
		t.Fatalf("unexpected stats: %+v", stats)
	}

	if _, err := service.GetByUserID("missing"); !errors.Is(err, ErrUserStatsNotFound) {
		t.Fatalf("expected ErrUserStatsNotFound, got %v", err)
	}
}

func TestUserStatsServiceRecordGameCreatesWhenMissing(t *testing.T) {
	t.Parallel()

	statsRepo := &fakeUserStatsRepo{byUserID: map[string]*user.UserStats{}}
	service := NewUserStatsService(statsRepo, &fakeUserRepo{})

	stats, err := service.RecordGame("u1", true)
	if err != nil {
		t.Fatalf("expected record game success, got %v", err)
	}
	if stats.Games != 1 || stats.Wins != 1 || stats.Elo != 1010 {
		t.Fatalf("unexpected stats after win: %+v", stats)
	}

	// Existing user stats should be updated in place.
	stats, err = service.RecordGame("u1", false)
	if err != nil {
		t.Fatalf("expected second record game success, got %v", err)
	}
	if stats.Games != 2 || stats.Wins != 1 || stats.Elo != 1000 {
		t.Fatalf("unexpected stats after loss: %+v", stats)
	}
}

func TestUserStatsServiceRecordGameSaveError(t *testing.T) {
	t.Parallel()

	statsRepo := &fakeUserStatsRepo{
		byUserID: map[string]*user.UserStats{
			"u1": {UserId: "u1", Games: 0, Wins: 0, Elo: 1000},
		},
		saveErr: errors.New("save failed"),
	}
	service := NewUserStatsService(statsRepo, &fakeUserRepo{})

	if _, err := service.RecordGame("u1", true); err == nil {
		t.Fatal("expected save error")
	}
}
