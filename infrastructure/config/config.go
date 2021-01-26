package config

import (
	"os"

	log "github.com/Nemo08/NCTW/infrastructure/logger"
	"github.com/joho/godotenv"
)

type ConfigInterface interface {
	Get(param string) string
	IsSet(param string) bool
}

type appConfig struct {
	log log.LogInterface
}

func NewAppConfigLoader(l log.LogInterface) appConfig {
	err := godotenv.Load()
	if err != nil {
		l.LogError("Ошибка загрузки .env файла")

	}
	return appConfig{log: l}
}

func (ac appConfig) Get(param string) string {
	ac.log.LogMessage("Читаю переменную окружения ", param, " = ", os.Getenv(param))
	return os.Getenv(param)
}

func (ac appConfig) IsSet(param string) bool {
	_, set := os.LookupEnv(param)
	ac.log.LogMessage("Проверяю переменную окружения ", param, ", она равна установлена")
	return set
}
