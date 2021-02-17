package main

import (
	"github.com/Pallinder/go-randomdata"
	"gopkg.in/guregu/null.v4"

	cfg "github.com/Nemo08/NCTW/infrastructure/config"
	db "github.com/Nemo08/NCTW/infrastructure/database"
	"github.com/Nemo08/NCTW/infrastructure/logger"
	api "github.com/Nemo08/NCTW/services/api"
	user "github.com/Nemo08/NCTW/services/user"
)

//FillDatbaseByUsers заполняет базу фейковыми данными
func FillDatbaseByUsers(uc *user.UsecaseStruct, c int) {
	for i := 0; i < c; i++ {
		prof := randomdata.GenerateProfile(1)
		newuser, _ := user.NewUser(
			null.NewString(prof.Name.First+prof.Name.Last+randomdata.Digits(3), true),
			null.NewString(prof.Login.Password, true),
			null.NewString(prof.Email+randomdata.Digits(3), true))
		uc.Add(api.Context{}, newuser)
	}
}

func main() {
	//логгер
	log := logger.Log

	//конфигуратор
	conf := cfg.NewAppConfigLoader(log)

	//база
	sqliterepo := db.NewSqliteRepository(conf, log)
	defer sqliterepo.Close()

	//создаем репозитории объектов
	userrepo := user.NewSqliteRepository(sqliterepo.GetDB())
	//contrepo := repo.NewContactRepositorySqlite(logger, conf, sqliterepo.GetDB())

	//Автомиграция таблиц
	sqliterepo.Migrate(&user.DbUser{})

	//бизнес-логика
	ucase := user.NewUsecase(userrepo)
	//contcase := user.NewContactUsecase(logger, contrepo)
	FillDatbaseByUsers(ucase, 100)
	log.Info("Заливка базы тестовыми данными закончена")
}
