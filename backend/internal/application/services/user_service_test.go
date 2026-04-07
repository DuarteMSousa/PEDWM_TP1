package services

import (
	"backend/internal/domain/user"
	"errors"
	"testing"
)

type fakeUserRepo struct {
	byID              map[string]*user.User
	byUsername        map[string]*user.User
	saveErr           error
	findByIDErr       error
	findByUsernameErr error
}

func (f *fakeUserRepo) Save(u *user.User) error {
	if f.saveErr != nil {
		return f.saveErr
	}
	if f.byID == nil {
		f.byID = map[string]*user.User{}
	}
	if f.byUsername == nil {
		f.byUsername = map[string]*user.User{}
	}
	f.byID[u.ID] = u
	f.byUsername[u.Username] = u
	return nil
}

func (f *fakeUserRepo) FindByID(id string) (*user.User, error) {
	if f.findByIDErr != nil {
		return nil, f.findByIDErr
	}
	return f.byID[id], nil
}

func (f *fakeUserRepo) FindByUsername(username string) (*user.User, error) {
	if f.findByUsernameErr != nil {
		return nil, f.findByUsernameErr
	}
	return f.byUsername[username], nil
}

type fakeUserStatsRepo struct {
	byUserID map[string]*user.UserStats
	saveErr  error
}

func (f *fakeUserStatsRepo) Save(stats *user.UserStats) error {
	if f.saveErr != nil {
		return f.saveErr
	}
	if f.byUserID == nil {
		f.byUserID = map[string]*user.UserStats{}
	}
	f.byUserID[stats.UserId] = stats
	return nil
}

func (f *fakeUserStatsRepo) FindByUserID(userID string) (*user.UserStats, error) {
	return f.byUserID[userID], nil
}

func TestUserServiceRegisterRejectsExistingUsername(t *testing.T) {
	t.Parallel()

	existing := &user.User{ID: "u1", Username: "ana", Password: "hash"}
	repo := &fakeUserRepo{
		byID:       map[string]*user.User{existing.ID: existing},
		byUsername: map[string]*user.User{existing.Username: existing},
	}
	statsRepo := &fakeUserStatsRepo{}
	service := NewUserService(repo, statsRepo)

	_, err := service.Register("ana", "123456")
	if !errors.Is(err, ErrUsernameExists) {
		t.Fatalf("expected ErrUsernameExists, got %v", err)
	}
}

func TestUserServiceRegisterSuccessCreatesUserAndStats(t *testing.T) {
	t.Parallel()

	repo := &fakeUserRepo{
		byID:       map[string]*user.User{},
		byUsername: map[string]*user.User{},
	}
	statsRepo := &fakeUserStatsRepo{byUserID: map[string]*user.UserStats{}}
	service := NewUserService(repo, statsRepo)

	created, err := service.Register("maria", "123456")
	if err != nil {
		t.Fatalf("expected register success, got %v", err)
	}
	if created == nil || created.ID == "" {
		t.Fatal("expected created user with ID")
	}
	if _, ok := repo.byID[created.ID]; !ok {
		t.Fatal("expected user persisted in repo")
	}

	stats, ok := statsRepo.byUserID[created.ID]
	if !ok {
		t.Fatal("expected initial user stats to be saved")
	}
	if stats.Elo != 1000 || stats.Games != 0 || stats.Wins != 0 {
		t.Fatalf("unexpected initial stats: %+v", stats)
	}
}

func TestUserServiceLogin(t *testing.T) {
	t.Parallel()

	u, err := user.NewUser("joao", "123456")
	if err != nil {
		t.Fatalf("failed to create user fixture: %v", err)
	}

	repo := &fakeUserRepo{
		byID:       map[string]*user.User{u.ID: u},
		byUsername: map[string]*user.User{u.Username: u},
	}
	service := NewUserService(repo, &fakeUserStatsRepo{})

	if _, err := service.Login("missing", "123456"); !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials for unknown user, got %v", err)
	}
	if _, err := service.Login("joao", "wrong"); !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials for wrong password, got %v", err)
	}

	logged, err := service.Login("joao", "123456")
	if err != nil {
		t.Fatalf("expected login success, got %v", err)
	}
	if logged.ID != u.ID {
		t.Fatalf("expected logged user %q, got %q", u.ID, logged.ID)
	}
}

func TestUserServiceGetUserNotFound(t *testing.T) {
	t.Parallel()

	service := NewUserService(&fakeUserRepo{
		byID:       map[string]*user.User{},
		byUsername: map[string]*user.User{},
	}, &fakeUserStatsRepo{})

	if _, err := service.GetUser("missing"); !errors.Is(err, ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
	if _, err := service.GetUserByUsername("missing"); !errors.Is(err, ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}
