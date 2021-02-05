package database

import (
	"database/sql"
	"os"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/mattn/go-sqlite3"

	cfg "github.com/Nemo08/NCTW/infrastructure/config"
	log "github.com/Nemo08/NCTW/infrastructure/logger"
)

type sqliteRepository struct {
	db *gorm.DB
}

//utflower служебная функция для регистронезависимого поиска в sqlite
func utflower(s string) string {
	return strings.ToLower(s)
}

//NewSqliteRepository новый объект репозитория sqlite
func NewSqliteRepository(c cfg.ConfigInterface) *sqliteRepository {
	if !c.IsSet("DBTYPE") || !c.IsSet("DBCONNECTIONSTRING") {
		log.LogError("Не установлены переменные окружения: DBTYPE или DBCONNECTIONSTRING")
		os.Exit(1)
	}

	dbtype := c.Get("DBTYPE")
	if dbtype == "sqlite3" {
		dbtype = "sqlite3_custom"

		//регистрируем свой драйвер для добавления в sqlite функции utflower
		//для обеспечения регистронезависимого поиска
		sql.Register("sqlite3_custom", &sqlite3.SQLiteDriver{
			ConnectHook: func(conn *sqlite3.SQLiteConn) error {
				if err := conn.RegisterFunc("utflower", utflower, true); err != nil {
					return err
				}

				return nil
			},
		})
	}

	db, err := gorm.Open(dbtype, c.Get("DBCONNECTIONSTRING"))

	if err != nil {
		log.LogError(err)
		os.Exit(1)
	}

	//db.LogMode(true)
	return &sqliteRepository{db: db}
}

func (sq *sqliteRepository) GetDB() *gorm.DB {
	return sq.db
}

func (sq *sqliteRepository) Close() {
	sq.db.Close()
}

func (sq *sqliteRepository) Migrate(objs ...interface{}) {
	for _, obj := range objs {
		sq.db.AutoMigrate(obj)
	}
}
