package interfaces

import "backend/internal/domain/user"

type UserRepository interface {
	Save(u *user.User) error
	FindByID(id string) (*user.User, error)
	FindByUsername(username string) (*user.User, error)
}
