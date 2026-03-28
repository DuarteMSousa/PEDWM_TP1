package interfaces

import "backend/internal/domain/friendship"

type FriendshipRepository interface {
	Save(f *friendship.Friendship) error
	Find(requesterID, addresseeID string) (*friendship.Friendship, error)
	FindByUser(userID string) ([]*friendship.Friendship, error)
	FindPendingForUser(userID string) ([]*friendship.Friendship, error)
	Delete(requesterID, addresseeID string) error
}
