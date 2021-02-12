package repository

import (
	"strconv"

	"github.com/jinzhu/gorm"

	"github.com/Nemo08/NCTW/infrastructure/router"
)

func Paginate(ctx router.ApiContext) func(db *gorm.DB) *gorm.DB {
	limit, offset, order, error := findLimits(ctx)
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
		return db.Offset(offset).Limit(limit).Order(order)
	}
}

//findLimits достает из контекста запроса лимиты, смещение
func findLimits(ctx router.ApiContext) (l int, o int, ord string, e error) {
	limit := 0
	offset := 0
	order := ""
	/*
		if (c.QueryParam("limit") == "") || (c.QueryParam("offset") == "") {
			return 0, 0, "", errors.New("error: Not set limit or offset")
		}
	*/
	limit, err := strconv.Atoi(ctx.QueryParam("limit"))
	if err != nil {
		return limit, offset, order, err
	}
	offset, err = strconv.Atoi(ctx.QueryParam("offset"))
	if err != nil {
		return limit, offset, order, err
	}
	order = ctx.QueryParam("offset")
	return limit, offset, order, nil
}
