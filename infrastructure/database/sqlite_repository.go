package database

import (
	"os"
	"strings"

	"gorm.io/gorm"

	cfg "github.com/Nemo08/NCTW/infrastructure/config"
	"github.com/Nemo08/NCTW/infrastructure/database/sq3driver"
	"github.com/Nemo08/NCTW/infrastructure/logger"
)

type sqliteRepository struct {
	db *gorm.DB
}

//utflower служебная функция для регистронезависимого поиска в sqlite
func utflower(s string) string {
	return strings.ToLower(s)
}

//NewSqliteRepository новый объект репозитория sqlite
func NewSqliteRepository(c cfg.ConfigInterface, log logger.Logr) *sqliteRepository {
	if !c.IsSet("DBTYPE") || !c.IsSet("DSN") {
		log.Error("Не установлены переменные окружения: DBTYPE или DSN")
		os.Exit(1)
	}

	dbtype := c.Get("DBTYPE")
	
	var db *gorm.DB
	var err error
	switch dbtype {
	case "sqlite3":
		db, err = gorm.Open(sq3driver.Open(c.Get("DSN")), &gorm.Config{Logger: logger.Log.GormLogger()})
	default:
		{
			log.Error("База ", c.Get("DBTYPE"), " не поддерживается")
			os.Exit(1)
		}
	}

	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	return &sqliteRepository{db: db}
}

func (sq *sqliteRepository) GetDB() *gorm.DB {
	return sq.db
}

func (sq *sqliteRepository) Close() {
	//sq.db.Close()
}

func (sq *sqliteRepository) Migrate(objs ...interface{}) {
	for _, obj := range objs {
		sq.db.AutoMigrate(obj)
	}
}
