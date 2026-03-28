package application

import (
	"backend/internal/application/interfaces"
	"backend/internal/domain/friendship"
)

type FriendshipService struct {
	repo     interfaces.FriendshipRepository
	userRepo interfaces.UserRepository
}

func NewFriendshipService(repo interfaces.FriendshipRepository, userRepo interfaces.UserRepository) *FriendshipService {
	return &FriendshipService{repo: repo, userRepo: userRepo}
}

func (s *FriendshipService) SendRequest(requesterID, addresseeID string) (*friendship.Friendship, error) {
	// Validate both users exist
	if _, err := s.userRepo.FindByID(requesterID); err != nil {
		return nil, ErrUserNotFound
	}
	if _, err := s.userRepo.FindByID(addresseeID); err != nil {
		return nil, ErrUserNotFound
	}

	// Check if friendship already exists in either direction
	existing, _ := s.repo.Find(requesterID, addresseeID)
	if existing != nil {
		return nil, friendship.ErrFriendshipAlreadySent
	}
	existing, _ = s.repo.Find(addresseeID, requesterID)
	if existing != nil {
		return nil, friendship.ErrFriendshipAlreadySent
	}

	f, err := friendship.NewFriendship(requesterID, addresseeID)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Save(f); err != nil {
		return nil, err
	}

	return f, nil
}

func (s *FriendshipService) AcceptRequest(requesterID, addresseeID string) (*friendship.Friendship, error) {
	f, err := s.repo.Find(requesterID, addresseeID)
	if err != nil || f == nil {
		return nil, friendship.ErrFriendshipNotFound
	}

	if err := f.Accept(); err != nil {
		return nil, err
	}

	if err := s.repo.Save(f); err != nil {
		return nil, err
	}

	return f, nil
}

func (s *FriendshipService) RejectRequest(requesterID, addresseeID string) (*friendship.Friendship, error) {
	f, err := s.repo.Find(requesterID, addresseeID)
	if err != nil || f == nil {
		return nil, friendship.ErrFriendshipNotFound
	}

	if err := f.Reject(); err != nil {
		return nil, err
	}

	if err := s.repo.Save(f); err != nil {
		return nil, err
	}

	return f, nil
}

func (s *FriendshipService) RemoveFriend(requesterID, addresseeID string) error {
	// Try both directions
	f, _ := s.repo.Find(requesterID, addresseeID)
	if f != nil {
		return s.repo.Delete(requesterID, addresseeID)
	}

	f, _ = s.repo.Find(addresseeID, requesterID)
	if f != nil {
		return s.repo.Delete(addresseeID, requesterID)
	}

	return friendship.ErrFriendshipNotFound
}

func (s *FriendshipService) GetFriends(userID string) ([]*friendship.Friendship, error) {
	return s.repo.FindByUser(userID)
}

func (s *FriendshipService) GetPendingRequests(userID string) ([]*friendship.Friendship, error) {
	return s.repo.FindPendingForUser(userID)
}
