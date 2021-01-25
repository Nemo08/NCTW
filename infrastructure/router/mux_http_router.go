package router

import (
	"github.com/gorilla/mux"

	use "github.com/Nemo08/nctw/usecase"
)

type muxHttpRouter struct {
	router *mux.Router
}

func NewMuxHttpRouter(l use.LogInterface) *muxHttpRouter {
	l.LogMessage("Set up main router")

	var us muxHttpRouter
	us.router = mux.NewRouter()
	us.router.StrictSlash(true)
	return &us
}

func (mr *muxHttpRouter) GetRouter() *mux.Router {
	return mr.router
}
