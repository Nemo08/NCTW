package database

import (
	"testing"

	"github.com/google/uuid"
	"gopkg.in/guregu/null.v4"

	cfg "github.com/Nemo08/NCTW/infrastructure/config"
	log "github.com/Nemo08/NCTW/infrastructure/logger"

	ent "github.com/Nemo08/NCTW/entity"
	repo "github.com/Nemo08/NCTW/infrastructure/repository"
	use "github.com/Nemo08/NCTW/usecase"
)

func TestNewSqliteRepository(t *testing.T) {
	logger := log.NewStdLogger()
	conf := cfg.NewCustomAppConfigLoader(logger)
	sqliterepo := NewSqliteRepository(logger, conf)
	defer sqliterepo.Close()

	sqliterepo.Migrate(&repo.DbUser{}, &repo.DbContact{}, &repo.DbBranch{})
	userrepo := repo.NewUserRepositorySqlite(logger, conf, sqliterepo.GetDB())
	ucase := use.NewUserUsecase(logger, userrepo)
	a := ent.User{
		ID:           uuid.New(),
		Login:        null.StringFrom("ЛОГин"),
		PasswordHash: null.StringFrom(""),
	}
	d, err := ucase.AddUser(a)
	u, err := ucase.Find("логин")
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
