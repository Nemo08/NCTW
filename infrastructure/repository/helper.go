package repository

import (
	"github.com/jinzhu/gorm"
)

func paginate(limit, offset int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if offset == 0 {
			offset = 1
		}

		switch {
		case limit > 100:
			limit = 100
		case limit <= 0:
			limit = 10
		}

		return db.Offset(offset).Limit(limit)
	}
}
