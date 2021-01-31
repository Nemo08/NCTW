// main
package main

import (
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

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
	e := echo.New()
	e.HideBanner = true
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.RequestID())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	apiV1Router := e.Group("/api/v1")

	rout.NewUserHTTPRouter(logger, ucase, apiV1Router)
	rout.NewContactHTTPRouter(logger, contcase, apiV1Router)
	rout.NewStaticHTTPRouter(logger, e)

	//запуск сервера
	if !conf.IsSet("SERVEPORT") {
		logger.LogError("Переменная окружения SERVEPORT для порта не установлена")
		os.Exit(1)
	}
	port := conf.Get("SERVEPORT")
	logger.LogMessage("Сервер запущен на порту " + port)
	e.Logger.Fatal(e.Start(":" + port))
}
