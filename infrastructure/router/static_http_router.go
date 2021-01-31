package router

import (
	"github.com/labstack/echo/v4"

	log "github.com/Nemo08/NCTW/infrastructure/logger"
)

//NewStaticHTTPRouter роутер статики
func NewStaticHTTPRouter(l log.LogInterface, e *echo.Echo) {
	l.LogMessage("Создаю роутер статики")

	e.Static("/", "../../static")
}
