// main
package main

import (
	"net/http"

	cfg "github.com/Nemo08/NCTW/infrastructure/config"
	db "github.com/Nemo08/NCTW/infrastructure/database"
	log "github.com/Nemo08/NCTW/infrastructure/logger"
	repo "github.com/Nemo08/NCTW/infrastructure/repository"
	rout "github.com/Nemo08/NCTW/infrastructure/router"
	use "github.com/Nemo08/NCTW/usecase"
)

func main() {
	//логгер
	logger := log.NewStdLogger()

	//конфигуратор
	conf := cfg.NewAppConfigLoader(logger)

	//база
	sqliterepo := db.NewSqliteRepository(logger, conf)
	urepo := repo.NewUserRepositorySqlite(logger, conf, sqliterepo.GetDB())
	contrepo := repo.NewContactRepositorySqlite(logger, conf, sqliterepo.GetDB())
	sqliterepo.Migrate(&repo.DbUser{}, &repo.DbContact{}, &repo.DbBranch{})
	defer sqliterepo.Close()

	//бизнес-логика
	ucase := use.NewUserUsecase(logger, urepo)
	contcase := use.NewContactUsecase(logger, contrepo)

	//роуты и сервер
	muxrouter := rout.NewMuxHttpRouter(logger)
	rout.NewUserHttpRouter(logger, ucase, muxrouter.GetRouter())
	rout.NewContactHttpRouter(logger, contcase, muxrouter.GetRouter())
	rout.NewStaticHttpRouter(logger, muxrouter.GetRouter())

	logger.LogMessage("Сервер запущен на порту 8222")
	err := http.ListenAndServe(":8222", muxrouter.GetRouter())
	if err != nil {
		logger.LogError(err)
	}
}
