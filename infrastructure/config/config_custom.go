package config

import (
	"os"
)

func NewCustomAppConfigLoader() appConfig {
	os.Setenv("DBTYPE", "sqlite3")
	os.Setenv("DBCONNECTIONSTRING", "file::memory:?cache=shared")
	os.Setenv("SERVEPORT", "8222")
	os.Setenv("SERVESTATIC", "true")
	os.Setenv("SERVESTATICPATH", "../../static")
	return appConfig{}
}
