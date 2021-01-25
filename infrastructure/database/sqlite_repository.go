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

func utflower(s string) string {
	return strings.ToLower(s)
}

func NewSqliteRepository(l log.LogInterface, c cfg.ConfigInterface) *sqliteRepository {
	if !c.IsSet("DBTYPE") || !c.IsSet("DBCONNECTIONSTRING") {
		l.LogError("Unable to get config: DBTYPE or DBCONNECTIONSTRING")
		os.Exit(1)
	}

	dbtype := c.Get("DBTYPE")

	sql.Register("sqlite3_custom", &sqlite3.SQLiteDriver{
		ConnectHook: func(conn *sqlite3.SQLiteConn) error {
			if err := conn.RegisterFunc("utflower", utflower, true); err != nil {
				return err
			}

			return nil
		},
	})

	if c.Get("DBTYPE") == "sqlite3" {
		dbtype = "sqlite3_custom"
	}

	db, err := gorm.Open(dbtype, c.Get("DBCONNECTIONSTRING"))
	db.SetLogger(l)

	if err != nil {
		l.LogError(err)
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
