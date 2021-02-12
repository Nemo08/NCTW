package user

import (
	"github.com/Nemo08/NCTW/infrastructure/router"

	"github.com/google/uuid"
)

//UserRepository объект репозитория User
type UserRepository interface {
	Store(ctx router.ApiContext, User User) (*User, error)
	GetUsers(ctx router.ApiContext) ([]*User, int, error)
	FindByID(ctx router.ApiContext, id uuid.UUID) (*User, error)
	Find(ctx router.ApiContext, q string) ([]*User, int, error)
	UpdateUser(ctx router.ApiContext, User User) (*User, error)
	DeleteUserByID(ctx router.ApiContext, id uuid.UUID) error
	CheckPassword(login string, password string) (*User, error)
}
