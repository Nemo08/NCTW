package usecase

import (
	"github.com/google/uuid"

	ent "github.com/Nemo08/NCTW/entity"
)

//UserRepository объект репозитория User
type UserRepository interface {
	Store(User ent.User) (*ent.User, error)
	GetAllUsers() ([]*ent.User, error)
	FindByID(id uuid.UUID) (*ent.User, error)
	Find(q string) ([]*ent.User, error)
	UpdateUser(User ent.User) (*ent.User, error)
	DeleteUserByID(id uuid.UUID) error
	CheckPassword(login string, password string) (*ent.User, error)
}
