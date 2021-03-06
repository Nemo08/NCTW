// main
package main

import (
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	cfg "github.com/Nemo08/NCTW/infrastructure/config"
	db "github.com/Nemo08/NCTW/infrastructure/database"
	logger "github.com/Nemo08/NCTW/infrastructure/logger"
	rout "github.com/Nemo08/NCTW/infrastructure/router"
	api "github.com/Nemo08/NCTW/services/api"
	user "github.com/Nemo08/NCTW/services/user"
)

var (
	gitTag, gitCommit, gitBranch string
)

func main() {
	//логгер
	logger.Log.SetLevel(logger.DebugLevel)
	logger.Log.Info("Запуск NCTW, git tag:'", gitTag, "', git commit:'", gitCommit, "', git branch:'", gitBranch, "'")

	//конфигуратор
	conf := cfg.NewAppConfigLoader(logger.Log)

	//база
	database := db.NewSqliteRepository(conf, logger.Log)
	defer database.Close()

	//создаем репозитории объектов
	userrepo := user.NewSqliteRepository(database.GetDB())

	//Автомиграция таблиц
	database.Migrate(&user.DbUser{})

	//бизнес-логика
	ucase := user.NewUsecase(userrepo)

	//роуты и сервер
	e := echo.New()
	e.HideBanner = true
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.RequestID())
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(logger.Log.EchoLogger())
	e.Use(api.CustomContext)
	apiV1Router := e.Group("/api/v1")

	user.NewUserHTTPRouter(logger.Log, ucase, apiV1Router)
	rout.NewStaticHTTPRouter(e)

	//запуск сервера
	if !conf.IsSet("SERVEPORT") {
		logger.Log.Error("Переменная окружения SERVEPORT для порта не установлена")
		os.Exit(1)
	}
	port := conf.Get("SERVEPORT")
	logger.Log.Info("Сервер запущен на порту " + port)
	e.Logger.Fatal(e.Start(":" + port))
}
