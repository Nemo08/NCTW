// main
package main

import (
	"fmt"
	"net/http"
	"os"

	cfg "github.com/Nemo08/NCTW/infrastructure/config"
	db "github.com/Nemo08/NCTW/infrastructure/database"
	log "github.com/Nemo08/NCTW/infrastructure/logger"
	repo "github.com/Nemo08/NCTW/infrastructure/repository"
	rout "github.com/Nemo08/NCTW/infrastructure/router"
	use "github.com/Nemo08/NCTW/usecase"
	"github.com/gorilla/mux"
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
	muxrouter := rout.NewMuxHTTPRouter(logger)
	apiV1Router := muxrouter.GetRouter().PathPrefix("/api/v1").Subrouter()
	rout.NewUserHTTPRouter(logger, ucase, apiV1Router)
	rout.NewContactHTTPRouter(logger, contcase, apiV1Router)
	rout.NewStaticHTTPRouter(logger, muxrouter.GetRouter())

	muxrouter.GetRouter().Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		t, err := route.GetPathTemplate()
		if err != nil {
			return err
		}
		m, _ := route.GetMethods()
		fmt.Println(t, m)
		return nil
	})

	if !conf.IsSet("SERVEPORT") {
		logger.LogError("Переменная окружения SERVEPORT для порта не установлена")
		os.Exit(1)
	}
	port := conf.Get("SERVEPORT")
	logger.LogMessage("Сервер запущен на порту " + port)
	err := http.ListenAndServe(":"+port, muxrouter.GetRouter())
	if err != nil {
		logger.LogError(err)
	}
}
