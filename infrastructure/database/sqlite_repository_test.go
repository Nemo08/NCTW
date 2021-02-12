package database

import (
	"testing"

	"github.com/google/uuid"
	"gopkg.in/guregu/null.v4"

	cfg "github.com/Nemo08/NCTW/infrastructure/config"
	"github.com/Nemo08/NCTW/infrastructure/router"

	user "github.com/Nemo08/NCTW/services/user"
)

func TestNewSqliteRepository(t *testing.T) {
	conf := cfg.NewCustomAppConfigLoader()
	sqliterepo := NewSqliteRepository(conf)
	defer sqliterepo.Close()

	sqliterepo.Migrate(&user.DbUser{})
	userrepo := user.NewUserRepositorySqlite(sqliterepo.GetDB())
	ucase := user.NewUserUsecase(userrepo)
	a := user.User{
		ID:           uuid.New(),
		Login:        null.StringFrom("ЛОГин"),
		PasswordHash: null.StringFrom(""),
	}
	d, err := ucase.AddUser(router.ApiContext{}, a)
	if err != nil {
		t.Error("Ошибка в добавлении пользователя ", err.Error())
	}

	u, _, err := ucase.Find(router.ApiContext{}, "логин")
	if err != nil {
		t.Error("Ошибка в поиске пользователя ", err.Error())
	}
	if len(u) == 0 {
		t.Error("Ошибка в поиске пользователя - нет результатов")
		return
	}
	if u[0].ID != d.ID {
		t.Error("Ошибка в поиске пользователя - ID не совпали!", u[0])
	}
}
