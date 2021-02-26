package user

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"

	repo "github.com/Nemo08/NCTW/infrastructure/repository"
	"github.com/Nemo08/NCTW/services/api"
)

type Repository interface {
	Store(ctx api.Context, User User) (*User, error)
	Get(ctx api.Context) ([]*User, error)
	FindByID(ctx api.Context, id uuid.UUID) (*User, error)
	Find(ctx api.Context, q string) ([]*User, error)
	Update(ctx api.Context, User User) (*User, error)
	DeleteByID(ctx api.Context, id uuid.UUID) error
	CheckPassword(login string, password string) (*User, error)
}

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

type repositorySqlite struct {
	db *gorm.DB
}

//NewSqliteRepository создание объекта репозитория для User
func NewSqliteRepository(db *gorm.DB) *repositorySqlite {
	return &repositorySqlite{
		db: db,
	}
}

//Store сохраняет пользователя в базе
func (urs *repositorySqlite) Store(ctx api.Context, user User) (*User, error) {
	var d DbUser = user2db(user)
	d.ID = uuid.New()

	err := urs.db.
		Scopes(repo.CtxLogger(ctx)).Debug().
		Create(&d).
		Error
	if err != nil {
		return &user, err
	}
	u := db2user(d)
	return &u, nil
}

//Get получает пользователя из базы
func (urs *repositorySqlite) Get(ctx api.Context) ([]*User, error) {
	var users []*User
	var DbUsers []*DbUser
	var count int64
	ctx.Log.Info(ctx)

	urs.db.
		Scopes(repo.CtxLogger(ctx)).Debug().
		Model(&DbUsers).
		Count(&count)
	g := urs.db.
		Scopes(repo.Paginate(ctx), repo.CtxLogger(ctx)).Debug().
		Find(&DbUsers)

	if g.Error != nil {
		return users, g.Error
	}

	for _, d := range DbUsers {
		e := db2user(*d)
		users = append(users, &e)
	}

	return users, nil
}

//FindByID ищет пользователя в базе по ИД
func (urs *repositorySqlite) FindByID(ctx api.Context, id uuid.UUID) (*User, error) {
	var d DbUser
	var u User

	g := urs.db.
		Scopes(repo.CtxLogger(ctx)).Debug().Where("id = ?", id).
		First(&d)
	if g.Error != nil {
		return &u, g.Error
	}

	u = db2user(d)
	return &u, nil
}

//Find ищет в базе емейл, логин и прочее, совпадающее с запросом
func (urs *repositorySqlite) Find(ctx api.Context, q string) ([]*User, error) {
	var users []*User
	var DbUsers []*DbUser
	var count int64

	//считаем количество результатов в запросе
	urs.db.
		Scopes(repo.Paginate(ctx), repo.CtxLogger(ctx)).Debug().
		Where("(utflower(login) LIKE ?) OR (utflower(email) LIKE ?)",
			strings.ToLower("%"+q+"%"),
			strings.ToLower("%"+q+"%")).
		Find(&DbUsers).
		Count(&count)

	g := urs.db.
		Scopes(repo.Paginate(ctx), repo.CtxLogger(ctx)).
		Where("(utflower(login) LIKE ?) OR (utflower(email) LIKE ?)",
			strings.ToLower("%"+q+"%"),
			strings.ToLower("%"+q+"%")).
		Find(&DbUsers)
	if g.Error != nil {
		return users, g.Error
	}
	for _, d := range DbUsers {
		e := db2user(*d)
		users = append(users, &e)
	}
	return users, nil
}

func (urs *repositorySqlite) Update(ctx api.Context, u User) (*User, error) {
	d := user2db(u)

	if !u.Email.IsZero() {
		g := urs.db.
			Scopes(repo.CtxLogger(ctx)).Debug().
			Model(&d).
			Where("id = ?", d.ID).
			Update("email", u.Email.String)
		if g.Error != nil {
			return &u, g.Error
		}
	}

	if !u.PasswordHash.IsZero() {
		g := urs.db.
			Scopes(repo.CtxLogger(ctx)).Debug().
			Model(&d).
			Where("id = ?", d.ID).
			Update("password_hash", u.PasswordHash.String)
		if g.Error != nil {
			return &u, g.Error
		}
	}

	updatedUser, err := urs.FindByID(ctx, d.ID)
	if err != nil {
		return &u, err
	}

	return updatedUser, nil
}

func (urs *repositorySqlite) DeleteByID(ctx api.Context, id uuid.UUID) error {
	g := urs.db.
		Scopes(repo.CtxLogger(ctx)).Debug().
		Where("id = ?", id).
		Delete(&DbUser{})

	return g.Error
}

func (urs *repositorySqlite) CheckPassword(login string, password string) (*User, error) {
	var d DbUser
	var u User

	g := urs.db.Debug().
		Where("login = ? AND password = ?", login, password).
		Take(&d)
	if g.Error != nil {
		return &u, g.Error
	}
	u = db2user(d)
	return &u, nil
}
