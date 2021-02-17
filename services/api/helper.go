package api

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/Nemo08/NCTW/infrastructure/logger"
)

//Context структура для проброса "контекста" по цепочке запроса
type Context struct {
	echo.Context
	Log *logger.Logr
}

//CustomContext миддлварь для оборачивания контекста эхи в кастомный
func CustomContext(h echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		zl, _ := zap.NewProduction()
		defer zl.Sync()
		//l := zl.Sugar().With(zap.String("request_id", c.Response().Header().Get("X-Request-ID")))
		log := logger.Log
		cc := Context{
			c,
			log.WithField("request_id", c.Response().Header().Get("X-Request-ID")),
		}
		return h(cc)
	}
}
