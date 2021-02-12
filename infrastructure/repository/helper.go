package repository

import (
	"github.com/jinzhu/gorm"
)

func Paginate(limit, offset int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		switch {
		case limit > 100:
			limit = 100
		case limit <= 0:
			limit = 10
		}
		return db.Offset(offset).Limit(limit)
	}
}
