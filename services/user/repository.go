package user

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"gopkg.in/guregu/null.v4"

	"github.com/Nemo08/NCTW/infrastructure/logger"
	repo "github.com/Nemo08/NCTW/infrastructure/repository"
	"github.com/Nemo08/NCTW/services/api"
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

type RepositorySqlite struct {
	db *gorm.DB
}

//NewRepositorySqlite создание объекта репозитория для User
func NewRepositorySqlite(db *gorm.DB) *RepositorySqlite {
	return &RepositorySqlite{
		db: db,
	}
}

func (urs *RepositorySqlite) Store(ctx api.Context, user User) (*User, error) {
	var d DbUser = user2db(user)
	d.ID = uuid.New()

	errSlice := urs.db.Create(&d).GetErrors()
	var estr string
	if len(errSlice) != 0 {

		for _, err := range errSlice {
			logger.Log.LogError("Error while user create", err)
			estr = estr + err.Error()
		}
		return &user, errors.New("Error while user create:" + estr)
	}
	u := db2user(d)
	return &u, nil
}

func (urs *RepositorySqlite) Get(ctx api.Context) ([]*User, int, error) {
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

func (urs *RepositorySqlite) FindByID(ctx api.Context, id uuid.UUID) (*User, error) {
	var d DbUser
	var u User

	g := urs.db.Where("id = ?", id).First(&d)
	if g.Error != nil {
		return &u, g.Error
	}

	u = db2user(d)
	return &u, nil
}

func (urs *RepositorySqlite) Find(ctx api.Context, q string) ([]*User, int, error) {
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

func (urs *RepositorySqlite) Update(ctx api.Context, u User) (*User, error) {
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

func (urs *RepositorySqlite) DeleteByID(ctx api.Context, id uuid.UUID) error {
	g := urs.db.Where("id = ?", id).Delete(&DbUser{})
	if g.Error != nil {
		return g.Error
	}
	return nil
}

func (urs *RepositorySqlite) CheckPassword(login string, password string) (*User, error) {
	var d DbUser
	var u User

	g := urs.db.Where("login = ? AND password = ?", login, password).Take(&d)
	if g.Error != nil {
		return &u, g.Error
	}
	u = db2user(d)
	return &u, nil
}
