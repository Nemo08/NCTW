// main
package main

import (
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	cfg "github.com/Nemo08/NCTW/infrastructure/config"
	db "github.com/Nemo08/NCTW/infrastructure/database"
	"github.com/Nemo08/NCTW/infrastructure/logger"
	repo "github.com/Nemo08/NCTW/infrastructure/repository"
	rout "github.com/Nemo08/NCTW/infrastructure/router"
	use "github.com/Nemo08/NCTW/usecase"
)

func main() {
	//конфигуратор
	conf := cfg.NewAppConfigLoader()

	//база
	sqliterepo := db.NewSqliteRepository(conf)
	defer sqliterepo.Close()

	//создаем репозитории объектов
	userrepo := repo.NewUserRepositorySqlite(conf, sqliterepo.GetDB())

	//Автомиграция таблиц
	sqliterepo.Migrate(&repo.DbUser{})

	//бизнес-логика
	ucase := use.NewUserUsecase(userrepo)

	//роуты и сервер
	e := echo.New()
	e.HideBanner = true
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.RequestID())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	apiV1Router := e.Group("/api/v1")

	rout.NewUserHTTPRouter(ucase, apiV1Router)
	rout.NewStaticHTTPRouter(e)

	//запуск сервера
	if !conf.IsSet("SERVEPORT") {
		logger.LogError("Переменная окружения SERVEPORT для порта не установлена")
		os.Exit(1)
	}
	port := conf.Get("SERVEPORT")
	logger.LogMessage("Сервер запущен на порту " + port)
	e.Logger.Fatal(e.Start(":" + port))
}
