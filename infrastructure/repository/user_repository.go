package repository

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"

	ent "github.com/Nemo08/NCTW/entity"
	cfg "github.com/Nemo08/NCTW/infrastructure/config"
	log "github.com/Nemo08/NCTW/infrastructure/logger"
)

//DbUser стуктура для хранения User в базе
type DbUser struct {
	ID           uuid.UUID  `gorm:"type:uuid;primary_key;"`
	CreatedAt    time.Time  `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time  `gorm:"default:CURRENT_TIMESTAMP"`
	DeletedAt    *time.Time `sql:"index"`
	Login        string     `gorm:"index;unique;not null"`
	PasswordHash string     `gorm:"not null"`
}

func db2user(i DbUser) ent.User {
	return ent.User{
		ID:           i.ID,
		Login:        i.Login,
		PasswordHash: i.PasswordHash,
	}
}

func user2db(i ent.User) DbUser {
	return DbUser{
		ID:           i.ID,
		Login:        i.Login,
		PasswordHash: i.PasswordHash,
	}
}

type userRepositorySqlite struct {
	db  *gorm.DB
	log log.LogInterface
}

//NewUserRepositorySqlite создание объекта репозитория для User
func NewUserRepositorySqlite(l log.LogInterface, c cfg.ConfigInterface, db *gorm.DB) *userRepositorySqlite {
	return &userRepositorySqlite{
		db:  db,
		log: l,
	}
}

func (urs *userRepositorySqlite) Store(user ent.User) (*ent.User, error) {
	var d DbUser
	d = user2db(user)
	d.ID = uuid.New()

	errSlice := urs.db.Create(&d).GetErrors()
	var estr string
	if len(errSlice) != 0 {

		for _, err := range errSlice {
			urs.log.LogError("Error while user create", err)
			estr = estr + err.Error()
		}
		return &user, errors.New("Error while user create:" + estr)
	}
	u := db2user(d)
	return &u, nil
}

func (urs *userRepositorySqlite) GetAllUsers() ([]*ent.User, error) {
	var users []*ent.User
	var DbUsers []*DbUser
	urs.db.Find(&DbUsers)
	for _, d := range DbUsers {
		e := db2user(*d)
		users = append(users, &e)
	}

	return users, nil
}

func (urs *userRepositorySqlite) FindByID(id uuid.UUID) (*ent.User, error) {
	var d DbUser
	urs.db.Where("id = ?", id).First(&d)
	e := db2user(d)
	return &e, nil
}

func (urs *userRepositorySqlite) Find(q string) ([]*ent.User, error) {
	var users []*ent.User
	var DbUsers []*DbUser

	urs.db.Where("utflower(login) LIKE ?", strings.ToLower("%"+q+"%")).Find(&DbUsers)
	for _, d := range DbUsers {
		e := db2user(*d)
		users = append(users, &e)
	}
	return users, nil
}

func (urs *userRepositorySqlite) UpdateUser(u ent.User) (*ent.User, error) {
	d := user2db(u)
	urs.db.Where("id = ?", d.ID).Save(&d)
	r := db2user(d)
	return &r, nil
}

func (urs *userRepositorySqlite) DeleteUserByID(id uuid.UUID) error {
	urs.db.Where("id = ?", id).Delete(&DbUser{})
	return nil
}

func (urs *userRepositorySqlite) CheckPassword(login string, password string) (*ent.User, error) {
	var d DbUser

	urs.db.Where("login = ? AND password = ?", login, password).Take(&d)
	u := db2user(d)
	return &u, nil
}
