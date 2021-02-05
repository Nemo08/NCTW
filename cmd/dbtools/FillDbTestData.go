package main

import (
	"github.com/Pallinder/go-randomdata"
	"gopkg.in/guregu/null.v4"

	ent "github.com/Nemo08/NCTW/entity"
	cfg "github.com/Nemo08/NCTW/infrastructure/config"
	db "github.com/Nemo08/NCTW/infrastructure/database"
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
	//конфигуратор
	conf := cfg.NewAppConfigLoader()

	//база
	sqliterepo := db.NewSqliteRepository(conf)
	defer sqliterepo.Close()

	//создаем репозитории объектов
	userrepo := repo.NewUserRepositorySqlite(sqliterepo.GetDB())
	//contrepo := repo.NewContactRepositorySqlite(logger, conf, sqliterepo.GetDB())

	//Автомиграция таблиц
	sqliterepo.Migrate(&repo.DbUser{})

	//бизнес-логика
	ucase := use.NewUserUsecase(userrepo)
	//contcase := use.NewContactUsecase(logger, contrepo)
	FillDatbaseByUsers(ucase, 100)
}
