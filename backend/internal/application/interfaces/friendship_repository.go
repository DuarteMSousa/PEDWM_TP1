package interfaces

import "backend/internal/domain/friendship"

// FriendshipRepository defines the contract for friendship persistence.
type FriendshipRepository interface {
	Save(f *friendship.Friendship) error
	Find(requesterID, addresseeID string) (*friendship.Friendship, error)
	FindByUser(userID string) ([]*friendship.Friendship, error)
	FindPendingForUser(userID string) ([]*friendship.Friendship, error)
	Delete(requesterID, addresseeID string) error
}
