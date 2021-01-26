// main
package main

import (
	"net/http"
	"os"

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
	defer sqliterepo.Close()

	//создаем репозитории объектов
	userrepo := repo.NewUserRepositorySqlite(logger, conf, sqliterepo.GetDB())
	contrepo := repo.NewContactRepositorySqlite(logger, conf, sqliterepo.GetDB())

	//Автомиграция таблиц
	sqliterepo.Migrate(&repo.DbUser{}, &repo.DbContact{}, &repo.DbBranch{})

	//бизнес-логика
	ucase := use.NewUserUsecase(logger, userrepo)
	contcase := use.NewContactUsecase(logger, contrepo)

	//роуты и сервер
	muxrouter := rout.NewMuxHttpRouter(logger)
	rout.NewUserHttpRouter(logger, ucase, muxrouter.GetRouter())
	rout.NewContactHttpRouter(logger, contcase, muxrouter.GetRouter())
	rout.NewStaticHttpRouter(logger, muxrouter.GetRouter())

	if !conf.IsSet("SERVEPORT") {
		logger.LogError("Переменная окружения SERVEPORT для порта не установлена")
		os.Exit(1)
	}
	err := http.ListenAndServe(":"+conf.Get("SERVEPORT"), muxrouter.GetRouter())
	logger.LogMessage("Сервер запущен на порту " + conf.Get("SERVEPORT"))
	if err != nil {
		logger.LogError(err)
	}
}
