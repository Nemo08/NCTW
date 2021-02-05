package usecase

import (
	"github.com/google/uuid"

	ent "github.com/Nemo08/NCTW/entity"
	log "github.com/Nemo08/NCTW/infrastructure/logger"
)

//UserUsecase основная структура usecase
type UserUsecase interface {
	GetUsers(limit, offset int) ([]*ent.User, int, error)
	AddUser(User ent.User) (*ent.User, error)
	FindByID(id uuid.UUID) (*ent.User, error)
	Find(q string, limit, offset int) ([]*ent.User, int, error)
	UpdateUser(User ent.User) (*ent.User, error)
	DeleteUserByID(id uuid.UUID) error
	CheckPassword(login string, password string) (*ent.User, error)
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

func (uc *UserUsecaseStruct) GetUsers(limit, offset int) ([]*ent.User, int, error) {
	log.LogMessage("Get users, limit:", limit, "offset:", offset)

	users, count, err := uc.repo.GetUsers(limit, offset)
	if err != nil {
		return nil, 0, err
	}
	return users, count, nil
}

func (uc *UserUsecaseStruct) AddUser(u ent.User) (*ent.User, error) {
	log.LogMessage("Add user", u)
	return uc.repo.Store(u)
}

func (uc *UserUsecaseStruct) FindByID(id uuid.UUID) (*ent.User, error) {
	log.LogMessage("Find user by id ", id)
	User, err := uc.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return User, nil
}

func (uc *UserUsecaseStruct) Find(q string, limit, offset int) ([]*ent.User, int, error) {
	log.LogMessage("Find string info in users:", q, "limit:", limit, "offset:", offset)

	users, count, err := uc.repo.Find(q, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	return users, count, nil
}

func (uc *UserUsecaseStruct) UpdateUser(u ent.User) (*ent.User, error) {
	log.LogMessage("Update user", u)
	return uc.repo.UpdateUser(u)
}

func (uc *UserUsecaseStruct) DeleteUserByID(id uuid.UUID) error {
	log.LogMessage("Delete user by id ", id)
	return uc.repo.DeleteUserByID(id)
}

func (uc *UserUsecaseStruct) CheckPassword(login string, password string) (*ent.User, error) {
	log.LogMessage("Check password of ", login)
	return uc.repo.CheckPassword(login, password)
}
