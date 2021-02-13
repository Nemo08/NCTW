package api

import (
	"github.com/labstack/echo/v4"
)

//Context структура для проброса "контекста" по цепочке запроса
type Context struct {
	echo.Context
}

//CustomContext миддлварь для оборачиванияконтекстаэхи в кастомный
func CustomContext(h echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cc := Context{c}
		return h(cc)
	}
}
