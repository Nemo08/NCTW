package user2

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"

	"nctw/infrastructure/core"
)

type Repository interface {
	Store(ctx core.ServiceContext, User User) error
	Get(ctx core.ServiceContext) error
	FindByID(ctx core.ServiceContext, id uuid.UUID) error
	Find(ctx core.ServiceContext, q string) error
	Update(ctx core.ServiceContext, User User) error
	DeleteByID(ctx core.ServiceContext, id uuid.UUID) error
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
func (urs *repositorySqlite) Store(ctx core.ServiceContext, user User) error {
	var d DbUser = user2db(user)
	d.ID = uuid.New()

	err := urs.db.
		Scopes(core.CtxLogger(ctx)).Debug().
		Create(&d).
		Error
	if err != nil {
		return err
	}
	ctx.ResponseData = db2user(d)
	return nil
}

//Get получает пользователя из базы
func (urs *repositorySqlite) Get(ctx core.ServiceContext) error {
	var users []*User
	var DbUsers []*DbUser
	var count int64
	ctx.Log.Info(ctx)

	urs.db.
		Scopes(core.CtxLogger(ctx)).Debug().
		Model(&DbUsers).
		Count(&count)
	g := urs.db.
		//Scopes(repo.Paginate(ctx), core.CtxLogger(ctx)).Debug().
		Find(&DbUsers)

	if g.Error != nil {
		return g.Error
	}

	for _, d := range DbUsers {
		e := db2user(*d)
		users = append(users, &e)
	}
	ctx.ResponseData = users
	return nil
}

//FindByID ищет пользователя в базе по ИД
func (urs *repositorySqlite) FindByID(ctx core.ServiceContext, id uuid.UUID) error {
	var d DbUser

	g := urs.db.
		Scopes(core.CtxLogger(ctx)).Debug().Where("id = ?", id).
		First(&d)
	if g.Error != nil {
		return g.Error
	}

	ctx.ResponseData = db2user(d)
	return nil
}

//Find ищет в базе емейл, логин и прочее, совпадающее с запросом
func (urs *repositorySqlite) Find(ctx core.ServiceContext, q string) error {
	var users []*User
	var DbUsers []*DbUser
	var count int64

	//считаем количество результатов в запросе
	urs.db.
		//Scopes(repo.Paginate(ctx), core.CtxLogger(ctx)).Debug().
		Where("(utflower(login) LIKE ?) OR (utflower(email) LIKE ?)",
			strings.ToLower("%"+q+"%"),
			strings.ToLower("%"+q+"%")).
		Find(&DbUsers).
		Count(&count)

	g := urs.db.
		//Scopes(repo.Paginate(ctx), core.CtxLogger(ctx)).
		Where("(utflower(login) LIKE ?) OR (utflower(email) LIKE ?)",
			strings.ToLower("%"+q+"%"),
			strings.ToLower("%"+q+"%")).
		Find(&DbUsers)
	if g.Error != nil {
		return g.Error
	}
	for _, d := range DbUsers {
		e := db2user(*d)
		users = append(users, &e)
	}
	ctx.ResponseData = users
	return nil
}

func (urs *repositorySqlite) Update(ctx core.ServiceContext, u User) error {
	d := user2db(u)

	if !u.Email.IsZero() {
		g := urs.db.
			Scopes(core.CtxLogger(ctx)).Debug().
			Model(&d).
			Where("id = ?", d.ID).
			Update("email", u.Email.String)
		if g.Error != nil {
			return g.Error
		}
	}

	if !u.PasswordHash.IsZero() {
		g := urs.db.
			Scopes(core.CtxLogger(ctx)).Debug().
			Model(&d).
			Where("id = ?", d.ID).
			Update("password_hash", u.PasswordHash.String)
		if g.Error != nil {
			return g.Error
		}
	}

	err := urs.FindByID(ctx, d.ID)
	if err != nil {
		return err
	}

	return nil
}

func (urs *repositorySqlite) DeleteByID(ctx core.ServiceContext, id uuid.UUID) error {
	g := urs.db.
		Scopes(core.CtxLogger(ctx)).Debug().
		Where("id = ?", id).
		Delete(&DbUser{})

	return g.Error
}
