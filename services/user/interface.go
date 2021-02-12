package user

import (
	"github.com/google/uuid"
)

//UserRepository объект репозитория User
type UserRepository interface {
	Store(User User) (*User, error)
	GetUsers(limit, offset int) ([]*User, int, error)
	FindByID(id uuid.UUID) (*User, error)
	Find(q string, limit, offset int) ([]*User, int, error)
	UpdateUser(User User) (*User, error)
	DeleteUserByID(id uuid.UUID) error
	CheckPassword(login string, password string) (*User, error)
}
