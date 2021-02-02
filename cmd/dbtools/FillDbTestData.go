package main

import (
	"github.com/Pallinder/go-randomdata"
	"gopkg.in/guregu/null.v4"

	ent "github.com/Nemo08/NCTW/entity"
	cfg "github.com/Nemo08/NCTW/infrastructure/config"
	db "github.com/Nemo08/NCTW/infrastructure/database"
	log "github.com/Nemo08/NCTW/infrastructure/logger"
	repo "github.com/Nemo08/NCTW/infrastructure/repository"
	use "github.com/Nemo08/NCTW/usecase"
)

//FillDatbaseByUsers заполняет базу фейковыми данными
func FillDatbaseByUsers(uc *use.UserUsecaseStruct, c int) {
	for i := 0; i < c; i++ {
		prof := randomdata.GenerateProfile(1)
		newuser, _ := ent.NewUser(
			null.NewString(prof.Name.First+prof.Name.Last+randomdata.Digits(3), true),
			null.NewString(prof.Login.Password, true),
			null.NewString(prof.Email+randomdata.Digits(3), true))
		uc.AddUser(newuser)
	}
}

func main() {
	//логгер
	logger := log.NewStdLogger()

	//конфигуратор
	conf := cfg.NewAppConfigLoader(logger)

	//база
	sqliterepo := db.NewSqliteRepository(logger, conf)
	defer sqliterepo.Close()

	//создаем репозитории объектов
	userrepo := repo.NewUserRepositorySqlite(logger, conf, sqliterepo.GetDB())
	//contrepo := repo.NewContactRepositorySqlite(logger, conf, sqliterepo.GetDB())

	//Автомиграция таблиц
	sqliterepo.Migrate(&repo.DbUser{}, &repo.DbContact{}, &repo.DbBranch{})

	//бизнес-логика
	ucase := use.NewUserUsecase(logger, userrepo)
	//contcase := use.NewContactUsecase(logger, contrepo)
	FillDatbaseByUsers(ucase, 100)
}
