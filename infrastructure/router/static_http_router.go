package router

import (
	"github.com/labstack/echo/v4"

	"nctw/infrastructure/logger"
)

//NewStaticHTTPRouter роутер статики
func NewStaticHTTPRouter(e *echo.Echo) {
	logger.Log.Info("Создаю роутер статики")

	e.Static("/", "../../static")
}
