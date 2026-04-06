package friendship

import (
	"errors"
	"time"
)

// FriendshipStatus represents the current state of a friendship.
type FriendshipStatus string

const (
	StatusPending  FriendshipStatus = "PENDING"
	StatusAccepted FriendshipStatus = "ACCEPTED"
	StatusRejected FriendshipStatus = "REJECTED"
)

var (
	ErrRequesterRequired     = errors.New("requester id is required")
	ErrAddresseeRequired     = errors.New("addressee id is required")
	ErrSameUser              = errors.New("cannot send friend request to yourself")
	ErrAlreadyAccepted       = errors.New("friendship already accepted")
	ErrAlreadyRejected       = errors.New("friendship already rejected")
	ErrNotPending            = errors.New("friendship is not pending")
	ErrFriendshipNotFound    = errors.New("friendship not found")
	ErrFriendshipAlreadySent = errors.New("friend request already sent")
)

// Friendship represents a friendship relationship between two users.
type Friendship struct {
	RequesterID string           `json:"requester_id"`
	AddresseeID string           `json:"addressee_id"`
	Status      FriendshipStatus `json:"status"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}

// NewFriendship creates a new friendship request in the PENDING state.
// Returns an error if the IDs are empty or identical.
func NewFriendship(requesterID, addresseeID string) (*Friendship, error) {
	if requesterID == "" {
		return nil, ErrRequesterRequired
	}
	if addresseeID == "" {
		return nil, ErrAddresseeRequired
	}
	if requesterID == addresseeID {
		return nil, ErrSameUser
	}

	now := time.Now()
	return &Friendship{
		RequesterID: requesterID,
		AddresseeID: addresseeID,
		Status:      StatusPending,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

// Accept changes the friendship status to ACCEPTED. It can only be called when the status is PENDING.
func (f *Friendship) Accept() error {
	if f.Status != StatusPending {
		return ErrNotPending
	}
	f.Status = StatusAccepted
	f.UpdatedAt = time.Now()
	return nil
}

// Reject changes the friendship status to REJECTED. It can only be called when the status is PENDING.
func (f *Friendship) Reject() error {
	if f.Status != StatusPending {
		return ErrNotPending
	}
	f.Status = StatusRejected
	f.UpdatedAt = time.Now()
	return nil
}
