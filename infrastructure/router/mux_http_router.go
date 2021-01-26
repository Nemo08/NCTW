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
	//us.router.Use(suffixMiddleware)
	/*
		m := http.NewServeMux()

		m.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
			fmt.Println(req.URL.Path)
			if req.URL.Path != "/" {
				req.URL.Path = strings.TrimSuffix(req.URL.Path, "/")
			}
			// do something with req
			us.router.ServeHTTP(w, req)
		})
	*/
	return &us
}

//GetRouter передает роутер
func (mr *MuxHTTPRouter) GetRouter() *mux.Router {
	return mr.router
}
