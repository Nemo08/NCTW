package repository

import (
	"strconv"

	"gorm.io/gorm"

	"github.com/Nemo08/NCTW/services/api"
)

func Paginate(ctx api.Context) func(db *gorm.DB) *gorm.DB {
	limit, offset, error := findLimits(ctx)
	if error != nil {
		return func(db *gorm.DB) *gorm.DB {
			return db
		}
	}

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

//findLimits достает из контекста запроса лимиты, смещение
func findLimits(ctx api.Context) (l int, o int, e error) {
	limit := 0
	offset := 0
	/*
		if (c.QueryParam("limit") == "") || (c.QueryParam("offset") == "") {
			return 0, 0, "", errors.New("error: Not set limit or offset")
		}
	*/
	limit, err := strconv.Atoi(ctx.QueryParam("limit"))
	if err != nil {
		return limit, offset, err
	}
	offset, err = strconv.Atoi(ctx.QueryParam("offset"))
	if err != nil {
		return limit, offset, err
	}
	//order = ctx.QueryParam("offset")
	return limit, offset, nil
}
