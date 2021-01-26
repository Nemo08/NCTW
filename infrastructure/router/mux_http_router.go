package router

import (
	"github.com/gorilla/mux"

	log "github.com/Nemo08/NCTW/infrastructure/logger"
)

//MuxHTTPRouter структура базового роутера
type MuxHTTPRouter struct {
	router *mux.Router
}

//NewMuxHTTPRouter новый базовый роутер
func NewMuxHTTPRouter(l log.LogInterface) *MuxHTTPRouter {
	l.LogMessage("Создаю основной роутер")

	var us MuxHTTPRouter
	us.router = mux.NewRouter()
	us.router.StrictSlash(true)
	return &us
}

//GetRouter передает роутер
func (mr *MuxHTTPRouter) GetRouter() *mux.Router {
	return mr.router
}
