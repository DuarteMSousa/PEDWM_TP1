package friendship

import (
	"errors"
	"time"
)

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

type Friendship struct {
	RequesterID string           `json:"requester_id"`
	AddresseeID string           `json:"addressee_id"`
	Status      FriendshipStatus `json:"status"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}

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

func (f *Friendship) Accept() error {
	if f.Status != StatusPending {
		return ErrNotPending
	}
	f.Status = StatusAccepted
	f.UpdatedAt = time.Now()
	return nil
}

func (f *Friendship) Reject() error {
	if f.Status != StatusPending {
		return ErrNotPending
	}
	f.Status = StatusRejected
	f.UpdatedAt = time.Now()
	return nil
}
