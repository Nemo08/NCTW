package router

import (
	"net/http"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	use "github.com/Nemo08/nctw/usecase"
)

//NewStaticHttpRouter
func NewStaticHttpRouter(l use.LogInterface, r *mux.Router) {
	l.LogMessage("Set up static router")

	s := http.StripPrefix("/", noFoldersContent(http.FileServer(http.Dir("../../static"))))
	loggedRouter := handlers.LoggingHandler(l, s)

	r.PathPrefix("/").Handler(loggedRouter)
}

func noFoldersContent(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}
