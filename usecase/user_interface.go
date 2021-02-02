package usecase

import (
	"github.com/google/uuid"

	ent "github.com/Nemo08/NCTW/entity"
)

//UserRepository объект репозитория User
type UserRepository interface {
	Store(User ent.User) (*ent.User, error)
	GetUsers(limit, offset int) ([]*ent.User, int, error)
	FindByID(id uuid.UUID) (*ent.User, error)
	Find(q string, limit, offset int) ([]*ent.User, int, error)
	UpdateUser(User ent.User) (*ent.User, error)
	DeleteUserByID(id uuid.UUID) error
	CheckPassword(login string, password string) (*ent.User, error)
}
