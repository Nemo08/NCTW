package user

import (
		"github.com/google/uuid"

		"github.com/Nemo08/NCTW/services/api"
)

//UserRepository объект репозитория User
type UserRepository interface {
	Store(ctx api.Context, User User) (*User, error)
	GetUsers(ctx api.Context) ([]*User, int, error)
	FindByID(ctx api.Context, id uuid.UUID) (*User, error)
	Find(ctx api.Context, q string) ([]*User, int, error)
	UpdateUser(ctx api.Context, User User) (*User, error)
	DeleteUserByID(ctx api.Context, id uuid.UUID) error
	CheckPassword(login string, password string) (*User, error)
}
