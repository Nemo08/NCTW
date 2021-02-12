package user

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"gopkg.in/guregu/null.v4"

	log "github.com/Nemo08/NCTW/infrastructure/logger"
	repo "github.com/Nemo08/NCTW/infrastructure/repository"
	"github.com/Nemo08/NCTW/infrastructure/router"
)

//DbUser стуктура для хранения User в базе
type DbUser struct {
	ID           uuid.UUID   `gorm:"type:uuid;primaryKey;PrioritizedPrimaryField"`
	CreatedAt    time.Time   `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time   `gorm:"default:CURRENT_TIMESTAMP"`
	DeletedAt    *time.Time  `sql:"index"`
	Login        null.String `gorm:"index;unique"`
	PasswordHash null.String
	Email        null.String
}

func db2user(i DbUser) User {
	return User{
		ID:           i.ID,
		Login:        i.Login,
		PasswordHash: i.PasswordHash,
		Email:        i.Email,
	}
}

func user2db(i User) DbUser {
	return DbUser{
		ID:           i.ID,
		Login:        i.Login,
		PasswordHash: i.PasswordHash,
		Email:        i.Email,
	}
}

type userRepositorySqlite struct {
	db *gorm.DB
}

//NewUserRepositorySqlite создание объекта репозитория для User
func NewUserRepositorySqlite(db *gorm.DB) *userRepositorySqlite {
	return &userRepositorySqlite{
		db: db,
	}
}

func (urs *userRepositorySqlite) Store(ctx router.ApiContext, user User) (*User, error) {
	var d DbUser = user2db(user)
	d.ID = uuid.New()

	errSlice := urs.db.Create(&d).GetErrors()
	var estr string
	if len(errSlice) != 0 {

		for _, err := range errSlice {
			log.LogError("Error while user create", err)
			estr = estr + err.Error()
		}
		return &user, errors.New("Error while user create:" + estr)
	}
	u := db2user(d)
	return &u, nil
}

func (urs *userRepositorySqlite) GetUsers(ctx router.ApiContext) ([]*User, int, error) {
	var users []*User
	var DbUsers []*DbUser
	var count int

	urs.db.Model(&DbUsers).Count(&count)
	g := urs.db.Scopes(repo.Paginate(ctx)).Find(&DbUsers)

	if g.Error != nil {
		return users, count, g.Error
	}

	for _, d := range DbUsers {
		e := db2user(*d)
		users = append(users, &e)
	}

	return users, count, nil
}

func (urs *userRepositorySqlite) FindByID(ctx router.ApiContext, id uuid.UUID) (*User, error) {
	var d DbUser
	var u User

	g := urs.db.Where("id = ?", id).First(&d)
	if g.Error != nil {
		return &u, g.Error
	}

	u = db2user(d)
	return &u, nil
}

func (urs *userRepositorySqlite) Find(ctx router.ApiContext, q string) ([]*User, int, error) {
	var users []*User
	var DbUsers []*DbUser
	var count int

	//считаем количество результатов в запросе
	urs.db.Where("utflower(login) LIKE ?", strings.ToLower("%"+q+"%")).Find(&DbUsers).Count(&count)
	g := urs.db.Scopes(repo.Paginate(ctx)).Where(
		"(utflower(login) LIKE ?) OR (utflower(email) LIKE ?)",
		strings.ToLower("%"+q+"%"), strings.ToLower("%"+q+"%")).Find(&DbUsers)
	if g.Error != nil {
		return users, 0, g.Error
	}
	for _, d := range DbUsers {
		e := db2user(*d)
		users = append(users, &e)
	}
	return users, count, nil
}

func (urs *userRepositorySqlite) UpdateUser(ctx router.ApiContext, u User) (*User, error) {
	d := user2db(u)
	attrs := make(map[string]interface{})

	if !u.Email.IsZero() {
		attrs["email"] = u.Email.String
	}

	if !u.Password.IsZero() {
		hash, err := CreateHash(u.Password.String)
		if err != nil {
			return &u, err
		}
		attrs["password_hash"] = hash
		u.PasswordHash = null.StringFrom(hash)
	}

	g := urs.db.Model(&d).Where("id = ?", d.ID).Update(attrs)
	if g.Error != nil {
		return &u, g.Error
	}

	updatedUser, err := urs.FindByID(ctx, d.ID)
	if err != nil {
		return &u, err
	}

	return updatedUser, nil
}

func (urs *userRepositorySqlite) DeleteUserByID(ctx router.ApiContext, id uuid.UUID) error {
	g := urs.db.Where("id = ?", id).Delete(&DbUser{})
	if g.Error != nil {
		return g.Error
	}
	return nil
}

func (urs *userRepositorySqlite) CheckPassword(login string, password string) (*User, error) {
	var d DbUser
	var u User

	g := urs.db.Where("login = ? AND password = ?", login, password).Take(&d)
	if g.Error != nil {
		return &u, g.Error
	}
	u = db2user(d)
	return &u, nil
}
