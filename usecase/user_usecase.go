package usecase

import (
	"github.com/google/uuid"

	ent "github.com/Nemo08/NCTW/entity"
	log "github.com/Nemo08/NCTW/infrastructure/logger"
	int "github.com/Nemo08/NCTW/interfaces"
)

//UserUsecase основная структура usecase
type UserUsecase interface {
	GetAllUsers() ([]*ent.User, error)
	AddUser(User ent.User) (*ent.User, error)
	FindByID(id uuid.UUID) (*ent.User, error)
	Find(q string) ([]*ent.User, error)
	UpdateUser(User ent.User) (*ent.User, error)
	DeleteUserByID(id uuid.UUID) error
	CheckPassword(login string, password string) (*ent.User, error)
}

type userUsecase struct {
	repo int.UserRepository
	log  log.LogInterface
}

//NewUserUsecase создание объекта usecase для User
func NewUserUsecase(l log.LogInterface, r int.UserRepository) *userUsecase {
	return &userUsecase{
		repo: r,
		log:  l,
	}
}

func (uc *userUsecase) GetAllUsers() ([]*ent.User, error) {
	uc.log.LogMessage("Get all users")

	users, err := uc.repo.GetAllUsers()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (uc *userUsecase) AddUser(u ent.User) (*ent.User, error) {
	uc.log.LogMessage("Add user", u)
	return uc.repo.Store(u)
}

func (uc *userUsecase) FindByID(id uuid.UUID) (*ent.User, error) {
	uc.log.LogMessage("Find user by id ", id)
	User, err := uc.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return User, nil
}

func (uc *userUsecase) Find(q string) ([]*ent.User, error) {
	uc.log.LogMessage("Find string info in users:", q)

	users, err := uc.repo.Find(q)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (uc *userUsecase) UpdateUser(u ent.User) (*ent.User, error) {
	uc.log.LogMessage("Update user", u)
	return uc.repo.UpdateUser(u)
}

func (uc *userUsecase) DeleteUserByID(id uuid.UUID) error {
	uc.log.LogMessage("Delete user by id ", id)
	return uc.repo.DeleteUserByID(id)
}

func (uc *userUsecase) CheckPassword(login string, password string) (*ent.User, error) {
	uc.log.LogMessage("Check password of ", login)
	return uc.repo.CheckPassword(login, password)
}
