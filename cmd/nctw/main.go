// main
package main

import (
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	cfg "github.com/Nemo08/NCTW/infrastructure/config"
	db "github.com/Nemo08/NCTW/infrastructure/database"
	"github.com/Nemo08/NCTW/infrastructure/logger"
	rout "github.com/Nemo08/NCTW/infrastructure/router"
	vld "github.com/Nemo08/NCTW/infrastructure/validator"
	api "github.com/Nemo08/NCTW/services/api"
	user "github.com/Nemo08/NCTW/services/user"
)

var (
	gitTag, gitCommit, gitBranch string
)

func main() {
	//логгер
	log := logger.NewLogger()
	log.Info("Запуск NCTW" + gitTag + gitCommit + gitBranch)

	//конфигуратор
	conf := cfg.NewAppConfigLoader()

	//валидатор
	vld.NewValidator()

	//база
	database := db.NewSqliteRepository(conf, log)
	defer database.Close()

	//создаем репозитории объектов
	userrepo := user.NewSqliteRepository(database.GetDB())

	//Автомиграция таблиц
	database.Migrate(&user.DbUser{})

	//бизнес-логика
	ucase := user.NewUsecase(log, userrepo)

	//роуты и сервер
	e := echo.New()
	e.HideBanner = true
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.RequestID())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(api.CustomContext)
	apiV1Router := e.Group("/api/v1")

	user.NewUserHTTPRouter(log, ucase, apiV1Router)
	rout.NewStaticHTTPRouter(e)

	//запуск сервера
	if !conf.IsSet("SERVEPORT") {
		log.Error("Переменная окружения SERVEPORT для порта не установлена")
		os.Exit(1)
	}
	port := conf.Get("SERVEPORT")
	log.Info("Сервер запущен на порту " + port)
	e.Logger.Fatal(e.Start(":" + port))
}
