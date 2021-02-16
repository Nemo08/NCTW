package api

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

//Context структура для проброса "контекста" по цепочке запроса
type Context struct {
	echo.Context
	Log *zap.SugaredLogger
}

//CustomContext миддлварь для оборачивания контекста эхи в кастомный
func CustomContext(h echo.HandlerFunc) echo.HandlerFunc {
	fmt.Println("THIS!!!!!!!!!!!!!!!!!!!!!!")
	return func(c echo.Context) error {
		zl, _ := zap.NewProduction()
		defer zl.Sync()
		cc := Context{
			c,
			zl.Sugar().With(zap.String("request_id", c.Response().Header().Get("X-Request-ID"))),
		}
		return h(cc)
	}
}
