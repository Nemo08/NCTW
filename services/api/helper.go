package api

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"moul.io/zapgorm2"

	"github.com/Nemo08/NCTW/infrastructure/logger"
)

//Context структура для проброса "контекста" по цепочке запроса
type Context struct {
	echo.Context
	Log *logger.Logr
}

func (c *Context) GormLogger() zapgorm2.Logger {
	return zapgorm2.New(c.Log.Desugar())
}

//CustomContext миддлварь для оборачивания контекста эхи в кастомный
func CustomContext(h echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		zl, _ := zap.NewProduction()
		defer zl.Sync()

		log := logger.Log.WithField("request_id", c.Response().Header().Get("X-Request-ID"))
		cc := Context{
			c,
			log,
		}
		return h(cc)
	}
}
