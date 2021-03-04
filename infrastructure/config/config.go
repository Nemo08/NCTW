package config

import (
	"os"

	"github.com/joho/godotenv"

	"nctw/infrastructure/logger"
)

type ConfigInterface interface {
	Get(param string) string
	IsSet(param string) bool
}

type appConfig struct {
	l logger.Logr
}

func NewAppConfigLoader(logr logger.Logr) appConfig {
	err := godotenv.Load()
	if err != nil {
		logr.Error("Ошибка загрузки .env файла")

	}
	return appConfig{
		l: logr,
	}
}

func (ac appConfig) Get(param string) string {
	ac.l.Info("Читаю переменную окружения '", param, "' = ", os.Getenv(param))
	return os.Getenv(param)
}

func (ac appConfig) IsSet(param string) bool {
	_, set := os.LookupEnv(param)
	ac.l.Info("Проверяю переменную окружения '", param, "', она установлена")
	return set
}
