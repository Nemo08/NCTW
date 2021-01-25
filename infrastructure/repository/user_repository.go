package repository

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"

	ent "github.com/Nemo08/NCTW/entity"
	cfg "github.com/Nemo08/NCTW/infrastructure/config"
	log "github.com/Nemo08/NCTW/infrastructure/logger"
)

type DbUser struct {
	ent.User `gorm:"embedded"`
}

type userRepositorySqlite struct {
	db  *gorm.DB
	log log.LogInterface
}

func NewUserRepositorySqlite(l log.LogInterface, c cfg.ConfigInterface, db *gorm.DB) *userRepositorySqlite {
	return &userRepositorySqlite{
		db:  db,
		log: l,
	}
}

func (urs *userRepositorySqlite) Store(User ent.User) (ent.User, error) {
	var u DbUser
	u.User = User
	a := uuid.New()

	u.ID = a

	err_slice := urs.db.Create(&u).GetErrors()
	if len(err_slice) != 0 {
		for _, err := range err_slice {
			urs.log.LogError("Error while user create", err)
		}
		return u.User, errors.New("Error while user create")
	}
	return u.User, nil
}

func (urs *userRepositorySqlite) GetAllUsers() ([]*ent.User, error) {
	var users []*ent.User
	var dbusers []*DbUser
	urs.db.Find(&dbusers)
	for _, u := range dbusers {
		users = append(users, &u.User)
	}

	return users, nil
}

func (urs *userRepositorySqlite) FindById(id uuid.UUID) (*ent.User, error) {
	var DUser DbUser
	urs.db.Where("id = ?", id).First(&DUser)
	return &DUser.User, nil
}

func (urs *userRepositorySqlite) Find(q string) ([]*ent.User, error) {
	var users []*ent.User
	var dbusers []*DbUser
	urs.db.Where("search_string LIKE ?", strings.ToLower("%"+q+"%")).Find(&dbusers)
	for _, u := range dbusers {
		users = append(users, &u.User)
	}

	return users, nil
}

func (urs *userRepositorySqlite) UpdateUser(User ent.User) (ent.User, error) {
	var DUser DbUser
	DUser.User = User

	urs.db.Where("id = ?", DUser.ID).Save(&DUser)
	return DUser.User, nil
}

func (urs *userRepositorySqlite) DeleteUserById(id uuid.UUID) error {
	urs.db.Where("id = ?", id).Delete(&DbUser{})
	return nil
}

func (urs *userRepositorySqlite) CheckPassword(login string, password string) (ent.User, error) {
	var user DbUser

	urs.db.Where("login = ? AND password = ?", login, password).Take(&user)
	return user.User, nil
}
