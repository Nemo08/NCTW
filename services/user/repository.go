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
func NewSqliteRepository(db *gorm.DB) *RepositorySqlite {
	return &RepositorySqlite{
		db: db,
	}
}

func (urs *RepositorySqlite) Store(ctx api.Context, user User) (*User, error) {
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

func (urs *RepositorySqlite) Get(ctx api.Context) ([]*User, int64, error) {
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

	g := urs.db.
		Scopes(repo.CtxLogger(ctx)).Debug().Where("id = ?", id).
		First(&d)
	if g.Error != nil {
		return &u, g.Error
	}

	u = db2user(d)
	return &u, nil
}

func (urs *RepositorySqlite) Find(ctx api.Context, q string) ([]*User, int64, error) {
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

	if !u.Password.IsZero() {
		hash, err := CreateHash(u.Password.String)
		if err != nil {
			return &u, err
		}
		u.PasswordHash = null.StringFrom(hash)
		g := urs.db.
			Scopes(repo.CtxLogger(ctx)).Debug().
			Model(&d).
			Where("id = ?", d.ID).
			Update("password_hash", hash)
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

func (urs *RepositorySqlite) DeleteByID(ctx api.Context, id uuid.UUID) error {
	g := urs.db.
		Scopes(repo.CtxLogger(ctx)).Debug().
		Where("id = ?", id).
		Delete(&DbUser{})
	if g.Error != nil {
		return g.Error
	}
	return nil
}

func (urs *RepositorySqlite) CheckPassword(login string, password string) (*User, error) {
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
