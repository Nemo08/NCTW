package main

import (
	"github.com/Pallinder/go-randomdata"
	"gopkg.in/guregu/null.v4"

	cfg "github.com/Nemo08/NCTW/infrastructure/config"
	db "github.com/Nemo08/NCTW/infrastructure/database"
	"github.com/Nemo08/NCTW/infrastructure/router"
	user "github.com/Nemo08/NCTW/services/user"
)

//FillDatbaseByUsers заполняет базу фейковыми данными
func FillDatbaseByUsers(uc *user.UserUsecaseStruct, c int) {
	for i := 0; i < c; i++ {
		prof := randomdata.GenerateProfile(1)
		newuser, _ := user.NewUser(
			null.NewString(prof.Name.First+prof.Name.Last+randomdata.Digits(3), true),
			null.NewString(prof.Login.Password, true),
			null.NewString(prof.Email+randomdata.Digits(3), true))
		uc.AddUser(router.ApiContext{}, newuser)
	}
}

func main() {
	//конфигуратор
	conf := cfg.NewAppConfigLoader()

	//база
	sqliterepo := db.NewSqliteRepository(conf)
	defer sqliterepo.Close()

	//создаем репозитории объектов
	userrepo := user.NewUserRepositorySqlite(sqliterepo.GetDB())
	//contrepo := repo.NewContactRepositorySqlite(logger, conf, sqliterepo.GetDB())

	//Автомиграция таблиц
	sqliterepo.Migrate(&user.DbUser{})

	//бизнес-логика
	ucase := user.NewUserUsecase(userrepo)
	//contcase := user.NewContactUsecase(logger, contrepo)
	FillDatbaseByUsers(ucase, 100)
}
