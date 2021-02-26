package user2

import (
	"github.com/Nemo08/NCTW/infrastructure/config"
	"github.com/Nemo08/NCTW/infrastructure/core"
	"github.com/Nemo08/NCTW/infrastructure/database"
	"github.com/Nemo08/NCTW/infrastructure/logger"
	"github.com/Nemo08/NCTW/services/user2"
)

func UserService() {
	//логгер
	core.Log.SetLevel(logger.DebugLevel)
	core.Log.Info("Запуск сервиса user")

	//конфигуратор
	conf := config.NewAppConfigLoader(logger.Log)

	//база
	database := database.NewSqliteRepository(conf, logger.Log)
	defer database.Close()

	//создаем репозитории объектов
	userrepo := user2.NewSqliteRepository(database.GetDB())

	use := user2.NewUsecase(userrepo)

	s := core.NewService()
	s.NewCommandHandler("store",
		func(sc core.ServiceContext) error {
			_, err := use.Get(sc)
			return err
		},
		jsonUserInput{},
	)
}
