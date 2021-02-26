package user

import (
	"github.com/google/uuid"

	"github.com/Nemo08/NCTW/services/api"
)

//Repository объект репозитория User
type Repository interface {
	Store(ctx api.Context, User User) (*User, error)
	Get(ctx api.Context) ([]*User, error)
	FindByID(ctx api.Context, id uuid.UUID) (*User, error)
	Find(ctx api.Context, q string) ([]*User, error)
	Update(ctx api.Context, User User) (*User, error)
	DeleteByID(ctx api.Context, id uuid.UUID) error
	CheckPassword(login string, password string) (*User, error)
}
