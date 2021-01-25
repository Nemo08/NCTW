package router

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	ent "github.com/Nemo08/NCTW/entity"
	log "github.com/Nemo08/NCTW/infrastructure/logger"
	use "github.com/Nemo08/NCTW/usecase"
)

type userHttpRouter struct {
	uc  use.UserUsecase
	log log.LogInterface
}

func NewUserHttpRouter(l log.LogInterface, u use.UserUsecase, r *mux.Router) {
	l.LogMessage("Set up user router")

	var us userHttpRouter
	us.uc = u
	us.log = l

	subr := r.PathPrefix("/user").Subrouter()
	subr.HandleFunc("", us.GetAllUsers).Methods("GET")
	subr.HandleFunc("", us.Store).Methods("POST")
	subr.HandleFunc("/{id:[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}}", us.GetUser).Methods("GET")
	subr.HandleFunc("/search/{query}", us.Find).Methods("GET")
	subr.HandleFunc("", us.Update).Methods("PUT")
	subr.HandleFunc("/{id:[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}}", us.Delete).Methods("DELETE")
}

func (ush *userHttpRouter) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	ush.log.LogMessage("Http request to get all users")

	var users []*ent.User
	users, err := ush.uc.GetAllUsers()
	if err != nil {
		resp := Message(true, "error")
		Respond(w, resp)
		return
	}
	resp := Message(true, "success")
	resp["data"] = users
	Respond(w, resp)
}

func (ush *userHttpRouter) GetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := uuid.Parse(params["id"])

	ush.log.LogMessage("Http request to get one user with id ", id)

	if err != nil {
		Respond(w, Message(false, "Error while decoding request body"))
		return
	}

	var User *ent.User
	User, err = ush.uc.FindById(id)
	if err != nil {
		resp := Message(true, "error")
		Respond(w, resp)
		return
	}
	resp := Message(true, "success")
	resp["data"] = &User
	Respond(w, resp)
}

func (ush *userHttpRouter) Find(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	q := params["query"]

	ush.log.LogMessage("Http request to find users with query ", q)

	if len(q) < 3 {
		resp := Message(true, "Too short query string")
		Respond(w, resp)
		return
	}

	var users []*ent.User
	users, err := ush.uc.Find(q)
	if err != nil {
		resp := Message(true, "error")
		Respond(w, resp)
		return
	}
	resp := Message(true, "success")
	resp["data"] = users
	Respond(w, resp)
}

func (ush *userHttpRouter) Store(w http.ResponseWriter, r *http.Request) {
	ush.log.LogMessage("Http request to make one user")

	User := &ent.User{}
	err := json.NewDecoder(r.Body).Decode(User)
	if err != nil {
		Respond(w, Message(false, "Error while decoding request body: "+err.Error()))
		return
	}
	resp := Message(true, "success")
	resp["data"], err = ush.uc.AddUser(*User)
	if err != nil {
		Respond(w, Message(false, "Error while user store: "+err.Error()))
		return
	}

	Respond(w, resp)
}

func (ush *userHttpRouter) Update(w http.ResponseWriter, r *http.Request) {
	ush.log.LogMessage("Http request to update one user")

	User := &ent.User{}
	err := json.NewDecoder(r.Body).Decode(User)
	if err != nil {
		Respond(w, Message(false, "Error while decoding request body: "+err.Error()))
		return
	}
	resp := Message(true, "success")
	resp["data"], err = ush.uc.UpdateUser(*User)
	if err != nil {
		Respond(w, Message(false, "Error while user store: "+err.Error()))
		return
	}

	Respond(w, resp)
}

func (ush *userHttpRouter) Delete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := uuid.Parse(params["id"])

	ush.log.LogMessage("Http request to delete one user with id ", id)

	if err != nil {
		Respond(w, Message(false, "Error while decoding request body"))
		return
	}

	err = ush.uc.DeleteUserById(id)
	if err != nil {
		resp := Message(true, "Error while delete user")
		Respond(w, resp)
		return
	}
	resp := Message(true, "success")
	Respond(w, resp)
}
