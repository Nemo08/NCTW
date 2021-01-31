package router

import (
	"github.com/gorilla/mux"

	log "github.com/Nemo08/NCTW/infrastructure/logger"
)

type muxHttpRouter struct {
	router *mux.Router
}

func NewMuxHttpRouter(l log.LogInterface) *muxHttpRouter {
	l.LogMessage("Set up main router")

	var us muxHttpRouter
	us.router = mux.NewRouter()
	us.router.StrictSlash(true)
	return &us
}

func (mr *muxHttpRouter) GetRouter() *mux.Router {
	return mr.router
}
