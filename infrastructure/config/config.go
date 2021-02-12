package config

import (
	"os"

	"github.com/joho/godotenv"

	log "github.com/Nemo08/NCTW/infrastructure/logger"
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
		log.LogError("Ошибка загрузки .env файла")

	}
	return appConfig{}
}

func (ac appConfig) Get(param string) string {
	log.LogMessage("Читаю переменную окружения '", param, "' = ", os.Getenv(param))
	return os.Getenv(param)
}

func (ac appConfig) IsSet(param string) bool {
	_, set := os.LookupEnv(param)
	log.LogMessage("Проверяю переменную окружения '", param, "', она установлена")
	return set
}
