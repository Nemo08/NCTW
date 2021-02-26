package core

import (
	"gorm.io/gorm"
)

func CtxLogger(ctx ServiceContext) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Session(&gorm.Session{Logger: ctx.Log.GormLogger()})
	}
}
