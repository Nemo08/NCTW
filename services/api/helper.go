package api

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

//Context структура для проброса "контекста" по цепочке запроса
type Context struct {
	echo.Context
	log *logrus.Entry
}

//CustomContext миддлварь для оборачивания контекста эхи в кастомный
func CustomContext(h echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cc := Context{
			c,
			logrus.WithFields(logrus.Fields{
				"id": c.Request().Header.Get("HeaderXRequestID"),
			}),
		}
		return h(cc)
	}
}
