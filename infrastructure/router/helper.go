package router

import (
	"github.com/labstack/echo/v4"
)

type ApiContext struct {
	echo.Context
}

func CustomApiContext(h echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cc := ApiContext{c}
		return h(cc)
	}
}
