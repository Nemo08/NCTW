package user

import (
	"github.com/google/uuid"

	"github.com/Nemo08/NCTW/infrastructure/logger"
	"github.com/Nemo08/NCTW/services/api"
)

//UserUsecase основная структура usecase
type Usecase interface {
	Get(ctx api.Context) ([]*User, int, error)
	Add(ctx api.Context, User User) (*User, error)
	FindByID(ctx api.Context, id uuid.UUID) (*User, error)
	Find(ctx api.Context, q string) ([]*User, int, error)
	Update(ctx api.Context, User User) (*User, error)
	DeleteByID(ctx api.Context, id uuid.UUID) error
	CheckPassword(login string, password string) (*User, error)
}

type UsecaseStruct struct {
	repo Repository
	log  logger.Logr
}

//NewUserUsecase создание объекта usecase для User
func NewUsecase(log logger.Logr, r Repository) *UsecaseStruct {
	return &UsecaseStruct{
		repo: r,
		log:  log,
	}
}

func (uc *UsecaseStruct) Get(ctx api.Context) ([]*User, int, error) {
	uc.log.Info("Get users")

	users, count, err := uc.repo.Get(ctx)
	if err != nil {
		return nil, 0, err
	}
	return users, count, nil
}

func (uc *UsecaseStruct) Add(ctx api.Context, u User) (*User, error) {
	uc.log.Info("Add user", u)
	return uc.repo.Store(ctx, u)
}

func (uc *UsecaseStruct) FindByID(ctx api.Context, id uuid.UUID) (*User, error) {
	uc.log.Info("Find user by id ", id)
	User, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return User, nil
}

func (uc *UsecaseStruct) Find(ctx api.Context, q string) ([]*User, int, error) {
	//uc.log.Info("Find string info in users:", q, "limit:", limit, "offset:", offset)

	users, count, err := uc.repo.Find(ctx, q)
	if err != nil {
		return nil, 0, err
	}
	return users, count, nil
}

func (uc *UsecaseStruct) Update(ctx api.Context, u User) (*User, error) {
	uc.log.Info("Update user", u)
	return uc.repo.Update(ctx, u)
}

func (uc *UsecaseStruct) DeleteByID(ctx api.Context, id uuid.UUID) error {
	uc.log.Info("Delete user by id ", id)
	return uc.repo.DeleteByID(ctx, id)
}

func (uc *UsecaseStruct) CheckPassword(login string, password string) (*User, error) {
	uc.log.Info("Check password of ", login)
	return uc.repo.CheckPassword(login, password)
}
