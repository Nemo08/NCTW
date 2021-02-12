package user

import (
	"github.com/google/uuid"

	log "github.com/Nemo08/NCTW/infrastructure/logger"
	"github.com/Nemo08/NCTW/infrastructure/router"
)

//UserUsecase основная структура usecase
type UserUsecase interface {
	GetUsers(ctx router.ApiContext) ([]*User, int, error)
	AddUser(ctx router.ApiContext, User User) (*User, error)
	FindByID(ctx router.ApiContext, id uuid.UUID) (*User, error)
	Find(ctx router.ApiContext, q string) ([]*User, int, error)
	UpdateUser(ctx router.ApiContext, User User) (*User, error)
	DeleteUserByID(ctx router.ApiContext, id uuid.UUID) error
	CheckPassword(login string, password string) (*User, error)
}

type UserUsecaseStruct struct {
	repo UserRepository
}

//NewUserUsecase создание объекта usecase для User
func NewUserUsecase(r UserRepository) *UserUsecaseStruct {
	return &UserUsecaseStruct{
		repo: r,
	}
}

func (uc *UserUsecaseStruct) GetUsers(ctx router.ApiContext) ([]*User, int, error) {
	//log.LogMessage("Get users, limit:", limit, "offset:", offset)

	users, count, err := uc.repo.GetUsers(ctx)
	if err != nil {
		return nil, 0, err
	}
	return users, count, nil
}

func (uc *UserUsecaseStruct) AddUser(ctx router.ApiContext, u User) (*User, error) {
	log.LogMessage("Add user", u)
	return uc.repo.Store(ctx, u)
}

func (uc *UserUsecaseStruct) FindByID(ctx router.ApiContext, id uuid.UUID) (*User, error) {
	log.LogMessage("Find user by id ", id)
	User, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return User, nil
}

func (uc *UserUsecaseStruct) Find(ctx router.ApiContext, q string) ([]*User, int, error) {
	//log.LogMessage("Find string info in users:", q, "limit:", limit, "offset:", offset)

	users, count, err := uc.repo.Find(ctx, q)
	if err != nil {
		return nil, 0, err
	}
	return users, count, nil
}

func (uc *UserUsecaseStruct) UpdateUser(ctx router.ApiContext, u User) (*User, error) {
	log.LogMessage("Update user", u)
	return uc.repo.UpdateUser(ctx, u)
}

func (uc *UserUsecaseStruct) DeleteUserByID(ctx router.ApiContext, id uuid.UUID) error {
	log.LogMessage("Delete user by id ", id)
	return uc.repo.DeleteUserByID(ctx, id)
}

func (uc *UserUsecaseStruct) CheckPassword(login string, password string) (*User, error) {
	log.LogMessage("Check password of ", login)
	return uc.repo.CheckPassword(login, password)
}
