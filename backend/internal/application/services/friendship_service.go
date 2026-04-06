package services

import (
	"backend/internal/application/interfaces"
	"backend/internal/domain/friendship"
	"log/slog"
)

// FriendshipService manages friend requests and relationships between users.
type FriendshipService struct {
	repo     interfaces.FriendshipRepository
	userRepo interfaces.UserRepository
}

// NewFriendshipService creates a new FriendshipService.
func NewFriendshipService(repo interfaces.FriendshipRepository, userRepo interfaces.UserRepository) *FriendshipService {
	return &FriendshipService{repo: repo, userRepo: userRepo}
}

// SendRequest creates a friend request between two users.
func (s *FriendshipService) SendRequest(requesterID, addresseeID string) (*friendship.Friendship, error) {
	slog.Info("sending friend request", "requesterID", requesterID, "addresseeID", addresseeID)

	// Validate both users exist
	if _, err := s.userRepo.FindByID(requesterID); err != nil {
		slog.Warn("friend request failed: requester not found", "requesterID", requesterID)
		return nil, ErrUserNotFound
	}
	if _, err := s.userRepo.FindByID(addresseeID); err != nil {
		slog.Warn("friend request failed: addressee not found", "addresseeID", addresseeID)
		return nil, ErrUserNotFound
	}

	// Check if friendship already exists in either direction
	existing, _ := s.repo.Find(requesterID, addresseeID)
	if existing != nil {
		slog.Warn("friend request failed: already exists", "requesterID", requesterID, "addresseeID", addresseeID)
		return nil, friendship.ErrFriendshipAlreadySent
	}
	existing, _ = s.repo.Find(addresseeID, requesterID)
	if existing != nil {
		slog.Warn("friend request failed: already exists (inverse)", "requesterID", requesterID, "addresseeID", addresseeID)
		return nil, friendship.ErrFriendshipAlreadySent
	}

	f, err := friendship.NewFriendship(requesterID, addresseeID)
	if err != nil {
		slog.Error("error creating friendship", "error", err)
		return nil, err
	}

	if err := s.repo.Save(f); err != nil {
		slog.Error("error persisting friendship", "error", err)
		return nil, err
	}

	slog.Info("friend request sent", "requesterID", requesterID, "addresseeID", addresseeID)
	return f, nil
}

// AcceptRequest accepts a pending friend request.
func (s *FriendshipService) AcceptRequest(requesterID, addresseeID string) (*friendship.Friendship, error) {
	slog.Info("accepting friend request", "requesterID", requesterID, "addresseeID", addresseeID)

	f, err := s.repo.Find(requesterID, addresseeID)
	if err != nil || f == nil {
		slog.Warn("friendship not found to accept", "requesterID", requesterID, "addresseeID", addresseeID)
		return nil, friendship.ErrFriendshipNotFound
	}

	if err := f.Accept(); err != nil {
		slog.Warn("error accepting friendship", "error", err)
		return nil, err
	}

	if err := s.repo.Save(f); err != nil {
		slog.Error("error persisting acceptance", "error", err)
		return nil, err
	}

	slog.Info("friend request accepted", "requesterID", requesterID, "addresseeID", addresseeID)
	return f, nil
}

// RejectRequest rejects a pending friend request.
func (s *FriendshipService) RejectRequest(requesterID, addresseeID string) (*friendship.Friendship, error) {
	slog.Info("rejecting friend request", "requesterID", requesterID, "addresseeID", addresseeID)

	f, err := s.repo.Find(requesterID, addresseeID)
	if err != nil || f == nil {
		slog.Warn("friendship not found to reject", "requesterID", requesterID, "addresseeID", addresseeID)
		return nil, friendship.ErrFriendshipNotFound
	}

	if err := f.Reject(); err != nil {
		slog.Warn("error rejecting friendship", "error", err)
		return nil, err
	}

	if err := s.repo.Save(f); err != nil {
		slog.Error("error persisting rejection", "error", err)
		return nil, err
	}

	slog.Info("friend request rejected", "requesterID", requesterID, "addresseeID", addresseeID)
	return f, nil
}

// RemoveFriend removes an existing friendship (both directions).
func (s *FriendshipService) RemoveFriend(requesterID, addresseeID string) error {
	slog.Info("removing friendship", "requesterID", requesterID, "addresseeID", addresseeID)

	// Try both directions
	f, _ := s.repo.Find(requesterID, addresseeID)
	if f != nil {
		return s.repo.Delete(requesterID, addresseeID)
	}

	f, _ = s.repo.Find(addresseeID, requesterID)
	if f != nil {
		return s.repo.Delete(addresseeID, requesterID)
	}

	slog.Warn("friendship not found to remove", "requesterID", requesterID, "addresseeID", addresseeID)
	return friendship.ErrFriendshipNotFound
}

// GetFriends returns the list of accepted friendships for a user.
func (s *FriendshipService) GetFriends(userID string) ([]*friendship.Friendship, error) {
	return s.repo.FindByUser(userID)
}

// GetPendingRequests returns the pending friend requests for a user.
func (s *FriendshipService) GetPendingRequests(userID string) ([]*friendship.Friendship, error) {
	return s.repo.FindPendingForUser(userID)
}
