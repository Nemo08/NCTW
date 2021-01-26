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

type DbUser struct {
	ID        uuid.UUID  `gorm:"type:uuid;primary_key;"`
	CreatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"-"`
	Login     string
	Password  string
}

func DB2Entity(i DbUser) ent.User {
	return ent.User{
		ID:       i.ID,
		Login:    i.Login,
		Password: i.Password,
	}
}

func Entity2DB(i ent.User) DbUser {
	return DbUser{
		ID:       i.ID,
		Login:    i.Login,
		Password: i.Password,
	}
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

func (urs *userRepositorySqlite) Store(user ent.User) (ent.User, error) {
	var u DbUser
	u = Entity2DB(user)
	u.ID = uuid.New()

	err_slice := urs.db.Create(&u).GetErrors()
	if len(err_slice) != 0 {
		for _, err := range err_slice {
			urs.log.LogError("Error while user create", err)
		}
		return DB2Entity(u), errors.New("Error while user create")
	}
	return DB2Entity(u), nil
}

func (urs *userRepositorySqlite) GetAllUsers() ([]*ent.User, error) {
	var users []*ent.User
	var dbusers []*DbUser
	urs.db.Find(&dbusers)
	for _, d := range dbusers {
		e := DB2Entity(*d)
		users = append(users, &e)
	}

	return users, nil
}

func (urs *userRepositorySqlite) FindById(id uuid.UUID) (*ent.User, error) {
	var d DbUser
	urs.db.Where("id = ?", id).First(&d)
	e := DB2Entity(d)
	return &e, nil
}

func (urs *userRepositorySqlite) Find(q string) ([]*ent.User, error) {
	var users []*ent.User
	var dbusers []*DbUser
	urs.db.Where("search_string LIKE ?", strings.ToLower("%"+q+"%")).Find(&dbusers)
	for _, d := range dbusers {
		e := DB2Entity(*d)
		users = append(users, &e)
	}

	return users, nil
}

func (urs *userRepositorySqlite) UpdateUser(u ent.User) (ent.User, error) {
	d := Entity2DB(u)
	urs.db.Where("id = ?", d.ID).Save(&d)
	return DB2Entity(d), nil
}

func (urs *userRepositorySqlite) DeleteUserById(id uuid.UUID) error {
	urs.db.Where("id = ?", id).Delete(&DbUser{})
	return nil
}

func (urs *userRepositorySqlite) CheckPassword(login string, password string) (ent.User, error) {
	var d DbUser

	urs.db.Where("login = ? AND password = ?", login, password).Take(&d)
	return DB2Entity(d), nil
}
