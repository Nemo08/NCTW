package user2

import (
	"nctw/infrastructure/config"
	"nctw/infrastructure/core"
	"nctw/infrastructure/database"
	"nctw/infrastructure/logger"
)

func UserService() core.Service {
	//логгер
	core.Log.SetLevel(logger.DebugLevel)
	core.Log.Info("Запуск сервиса user")

	//конфигуратор
	conf := config.NewAppConfigLoader(logger.Log)

	//база
	database := database.NewSqliteRepository(conf, logger.Log)
	defer database.Close()

	//создаем репозитории объектов
	userrepo := NewSqliteRepository(database.GetDB())

	usecase := NewUsecase(userrepo)

	s := core.NewService()
	s.NewCommandHandler("store",
		func(sc core.ServiceContext) error {
			err := usecase.Get(sc)
			return err
		},
		jsonUserInput{},
	)
	return s
}
