package user

import (
	"github.com/google/uuid"

	"github.com/Nemo08/NCTW/services/api"
)

//Usecase основная структура usecase
type Usecase interface {
	Get(ctx api.Context) ([]*User, int64, error)
	Add(ctx api.Context, User User) (*User, error)
	FindByID(ctx api.Context, id uuid.UUID) (*User, error)
	Find(ctx api.Context, q string) ([]*User, int64, error)
	Update(ctx api.Context, User User) (*User, error)
	DeleteByID(ctx api.Context, id uuid.UUID) error
	CheckPassword(login string, password string) (*User, error)
}

type usecaseStruct struct {
	repo Repository
}

//NewUserUsecase создание объекта usecase для User
func NewUsecase(r Repository) *usecaseStruct {
	return &usecaseStruct{
		repo: r,
	}
}

func (uc *usecaseStruct) Get(ctx api.Context) ([]*User, int64, error) {
	ctx.Log.Info("Get users")

	users, count, err := uc.repo.Get(ctx)
	if err != nil {
		return nil, 0, err
	}
	return users, count, nil
}

func (uc *usecaseStruct) Add(ctx api.Context, u User) (*User, error) {
	//ctx.Log.Info("Add user", u)
	return uc.repo.Store(ctx, u)
}

func (uc *usecaseStruct) FindByID(ctx api.Context, id uuid.UUID) (*User, error) {
	ctx.Log.Info("Find user by id ", id)
	User, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return User, nil
}

func (uc *usecaseStruct) Find(ctx api.Context, q string) ([]*User, int64, error) {
	//ctx.Log.Info("Find string info in users:", q, "limit:", limit, "offset:", offset)

	users, count, err := uc.repo.Find(ctx, q)
	if err != nil {
		return nil, 0, err
	}
	return users, count, nil
}

func (uc *usecaseStruct) Update(ctx api.Context, u User) (*User, error) {
	ctx.Log.Info("Update user", u)
	return uc.repo.Update(ctx, u)
}

func (uc *usecaseStruct) DeleteByID(ctx api.Context, id uuid.UUID) error {
	ctx.Log.Info("Delete user by id ", id)
	return uc.repo.DeleteByID(ctx, id)
}

func (uc *usecaseStruct) CheckPassword(login string, password string) (*User, error) {
	return uc.repo.CheckPassword(login, password)
}
