package config

import (
	"os"

	"github.com/joho/godotenv"

	use "github.com/Nemo08/NCTW/usecase"
)

type ConfigInterface interface {
	Get(param string) string
	IsSet(param string) bool
}

type appConfig struct {
	log use.LogInterface
}

func NewAppConfigLoader(l use.LogInterface) appConfig {
	err := godotenv.Load()
	if err != nil {
		l.LogError("Error loading .env file")

	}
	return appConfig{log: l}
}

func (ac appConfig) Get(param string) string {
	ac.log.LogMessage("Read ENV param ", param, " = ", os.Getenv(param))
	return os.Getenv(param)
}

func (ac appConfig) IsSet(param string) bool {
	_, set := os.LookupEnv(param)
	ac.log.LogMessage("Check ENV param ", param, " with result ", set)
	return set
}
