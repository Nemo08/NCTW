package user2

import (
	"github.com/Nemo08/NCTW/infrastructure/core"
	"github.com/google/uuid"
)

//Usecase основная структура usecase
type Usecase interface {
	Get(ctx core.ServiceContext) error
	Store(ctx core.ServiceContext, User User) error
	FindByID(ctx core.ServiceContext, id uuid.UUID) error
	Find(ctx core.ServiceContext, q string) error
	Update(ctx core.ServiceContext, User User) error
	DeleteByID(ctx core.ServiceContext, id uuid.UUID) error
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

func (uc *usecaseStruct) Get(ctx core.ServiceContext) error {
	ctx.Log.Info("Get users")
	return uc.repo.Get(ctx)
}

func (uc *usecaseStruct) Store(ctx core.ServiceContext, u User) error {
	//ctx.Log.Info("Add user", u)
	return uc.repo.Store(ctx, u)
}

func (uc *usecaseStruct) FindByID(ctx core.ServiceContext, id uuid.UUID) error {
	ctx.Log.Info("Find user by id ", id)
	return uc.repo.FindByID(ctx, id)
}

func (uc *usecaseStruct) Find(ctx core.ServiceContext, q string) error {
	return uc.repo.Find(ctx, q)
}

func (uc *usecaseStruct) Update(ctx core.ServiceContext, u User) error {
	ctx.Log.Info("Update user", u)
	return uc.repo.Update(ctx, u)
}

func (uc *usecaseStruct) DeleteByID(ctx core.ServiceContext, id uuid.UUID) error {
	ctx.Log.Info("Delete user by id ", id)
	return uc.repo.DeleteByID(ctx, id)
}