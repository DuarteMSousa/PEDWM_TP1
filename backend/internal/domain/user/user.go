package user

import (
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUsernameRequired = errors.New("username is required")
	ErrPasswordRequired = errors.New("password is required")
	ErrPasswordTooShort = errors.New("password must be at least 6 characters")
	ErrInvalidPassword  = errors.New("invalid password")
)

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// NewUser creates a new user with a hashed password.
func NewUser(username, plainPassword string) (*User, error) {
	if username == "" {
		return nil, ErrUsernameRequired
	}
	if plainPassword == "" {
		return nil, ErrPasswordRequired
	}
	if len(plainPassword) < 6 {
		return nil, ErrPasswordTooShort
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:       uuid.New().String(),
		Username: username,
		Password: string(hashed),
	}, nil
}

// CheckPassword verifies a plain password against the stored hash.
func (u *User) CheckPassword(plainPassword string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainPassword)) == nil
}
