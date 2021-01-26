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

type jsonUser struct {
	ID           uuid.UUID
	Login        string
	Password     string `json:",omitempty"`
	PasswordHash string `json:"-"`
	Email        string
}

//json2entity Json объект копируем в Entity
func json2entity(i jsonUser) ent.User {
	return ent.User{
		ID:           i.ID,
		Login:        i.Login,
		PasswordHash: i.PasswordHash,
		Email:        i.Email,
	}
}

//entity2json Entity объект копируем в Json
func entity2json(i ent.User) jsonUser {
	return jsonUser{
		ID:    i.ID,
		Login: i.Login,
		Email: i.Email,
	}
}

type userHTTPRouter struct {
	uc  use.UserUsecase
	log log.LogInterface
}

//NewUserHTTPRouter роутер пользователей
func NewUserHTTPRouter(l log.LogInterface, u use.UserUsecase, r *mux.Router) {
	l.LogMessage("Создаю роутер для user")

	us := userHTTPRouter{
		uc:  u,
		log: l,
	}

	subr := r.PathPrefix("/user").Subrouter()
	subr.HandleFunc("", us.Store).Methods("POST")
	//subr.HandleFunc("", us.GetAllUsers).Methods("GET")
	subr.HandleFunc("/{id:[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}}", us.GetUser).Methods("GET")
	subr.HandleFunc("/search/{query}", us.Find).Methods("GET")
	subr.HandleFunc("", us.Update).Methods("PUT")
	subr.HandleFunc("/{id:[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}}", us.Delete).Methods("DELETE")
}

func (ush *userHTTPRouter) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	ush.log.LogMessage("Http request to get all users")

	var u []*ent.User
	var jsusers []*jsonUser
	u, err := ush.uc.GetAllUsers()
	if err != nil {
		resp := Message("error:" + err.Error())
		Respond(w, resp, http.StatusInternalServerError)
		return
	}

	for _, d := range u {
		j := entity2json(*d)
		jsusers = append(jsusers, &j)
	}

	resp := Message("success")
	resp["data"] = jsusers
	Respond(w, resp, http.StatusOK)
}

func (ush *userHTTPRouter) GetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := uuid.Parse(params["id"])

	ush.log.LogMessage("Http request to get one user with id ", id)

	if err != nil {
		Respond(w, Message("Error while decoding request body:"+err.Error()), http.StatusBadRequest)
		return
	}

	var u *ent.User
	u, err = ush.uc.FindByID(id)
	if err != nil {
		resp := Message("error:" + err.Error())
		Respond(w, resp, http.StatusInternalServerError)
		return
	}
	resp := Message("success")
	j := entity2json(*u)
	resp["data"] = &j
	Respond(w, resp, http.StatusOK)
}

func (ush *userHTTPRouter) Find(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	q := params["query"]

	ush.log.LogMessage("Http request to find users with query ", q)

	if len(q) < 3 {
		resp := Message("Too short query string")
		Respond(w, resp, http.StatusBadRequest)
		return
	}

	var u []*ent.User
	u, err := ush.uc.Find(q)
	if err != nil {
		resp := Message("error:" + err.Error())
		Respond(w, resp, http.StatusInternalServerError)
		return
	}

	var jsusers []*jsonUser
	for _, d := range u {
		j := entity2json(*d)
		jsusers = append(jsusers, &j)
	}

	resp := Message("success")
	resp["data"] = jsusers
	Respond(w, resp, http.StatusOK)
}

func (ush *userHTTPRouter) Store(w http.ResponseWriter, r *http.Request) {
	ush.log.LogMessage("Http request to make one user")

	j := &jsonUser{}

	err := json.NewDecoder(r.Body).Decode(j)
	if err != nil {
		Respond(w, Message("Error while decoding request body: "+err.Error()), http.StatusBadRequest)
		return
	}

	u, err := ent.NewUser(j.Login, j.Password, j.Email)
	if err != nil {
		Respond(w, Message("Error while user store: "+err.Error()), http.StatusInternalServerError)
		return
	}

	u2, err := ush.uc.AddUser(u)
	if err != nil {
		Respond(w, Message("Error while user store: "+err.Error()), http.StatusInternalServerError)
		return
	}

	*j = entity2json(*u2)
	resp := Message("success")
	resp["data"] = j
	Respond(w, resp, http.StatusOK)
}

func (ush *userHTTPRouter) Update(w http.ResponseWriter, r *http.Request) {
	ush.log.LogMessage("Http request to update one user")

	j := &jsonUser{}
	err := json.NewDecoder(r.Body).Decode(j)
	if err != nil {
		Respond(w, Message("Error while decoding request body: "+err.Error()), http.StatusBadRequest)
		return
	}

	u, err := ush.uc.UpdateUser(json2entity(*j))
	if err != nil {
		Respond(w, Message("Error while user store: "+err.Error()), http.StatusInternalServerError)
		return
	}

	*j = entity2json(*u)
	resp := Message("success")
	resp["data"] = j
	Respond(w, resp, http.StatusOK)
}

func (ush *userHTTPRouter) Delete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := uuid.Parse(params["id"])

	ush.log.LogMessage("Http request to delete one user with id ", id)

	if err != nil {
		Respond(w, Message("Error while decoding request body:"+err.Error()), http.StatusBadRequest)
		return
	}

	err = ush.uc.DeleteUserByID(id)
	if err != nil {
		resp := Message("Error while delete user:" + err.Error())
		Respond(w, resp, http.StatusInternalServerError)
		return
	}

	Respond(w, Message("success"), http.StatusOK)
}
