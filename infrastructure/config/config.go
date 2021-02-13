package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/Nemo08/NCTW/infrastructure/logger"
)

type ConfigInterface interface {
	Get(param string) string
	IsSet(param string) bool
}

type appConfig struct {
}

func NewAppConfigLoader() appConfig {
	err := godotenv.Load()
	if err != nil {
		logger.Log.LogError("Ошибка загрузки .env файла")

	}
	return appConfig{}
}

func (ac appConfig) Get(param string) string {
	logger.Log.LogMessage("Читаю переменную окружения '", param, "' = ", os.Getenv(param))
	return os.Getenv(param)
}

func (ac appConfig) IsSet(param string) bool {
	_, set := os.LookupEnv(param)
	logger.Log.LogMessage("Проверяю переменную окружения '", param, "', она установлена")
	return set
}
