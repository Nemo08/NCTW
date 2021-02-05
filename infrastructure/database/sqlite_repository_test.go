package database

import (
	"testing"

	"github.com/google/uuid"
	"gopkg.in/guregu/null.v4"

	cfg "github.com/Nemo08/NCTW/infrastructure/config"

	ent "github.com/Nemo08/NCTW/entity"
	repo "github.com/Nemo08/NCTW/infrastructure/repository"
	use "github.com/Nemo08/NCTW/usecase"
)

func TestNewSqliteRepository(t *testing.T) {
	conf := cfg.NewCustomAppConfigLoader()
	sqliterepo := NewSqliteRepository(conf)
	defer sqliterepo.Close()

	sqliterepo.Migrate(&repo.DbUser{})
	userrepo := repo.NewUserRepositorySqlite(sqliterepo.GetDB())
	ucase := use.NewUserUsecase(userrepo)
	a := ent.User{
		ID:           uuid.New(),
		Login:        null.StringFrom("ЛОГин"),
		PasswordHash: null.StringFrom(""),
	}
	d, err := ucase.AddUser(a)
	u, _, err := ucase.Find("логин", 5, 0)
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
