package api

import (
	"github.com/labstack/echo/v4"
)

type Context struct {
	echo.Context
}

func CustomContext(h echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cc := Context{c}
		return h(cc)
	}
}
