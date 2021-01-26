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

type contactHTTPRouter struct {
	uc  use.ContactUsecase
	log log.LogInterface
}

//NewContactHTTPRouter роутер для контактов
func NewContactHTTPRouter(l log.LogInterface, u use.ContactUsecase, r *mux.Router) {
	l.LogMessage("Создаю роутер для contact")

	var us contactHTTPRouter
	us.uc = u
	us.log = l

	subr := r.PathPrefix("/contact").Subrouter()
	subr.HandleFunc("", us.GetAllContacts).Methods("GET")
	subr.HandleFunc("", us.Store).Methods("POST")
	subr.HandleFunc("/{id:[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}}", us.GetContact).Methods("GET")
	subr.HandleFunc("/search/{query}", us.Find).Methods("GET")
	subr.HandleFunc("", us.Update).Methods("PUT")
	subr.HandleFunc("/{id:[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}}", us.Delete).Methods("DELETE")
}

func (ush *contactHTTPRouter) GetAllContacts(w http.ResponseWriter, r *http.Request) {
	ush.log.LogMessage("Http request to get all contacts")

	var contacts []*ent.Contact
	contacts, err := ush.uc.GetAllContacts()
	if err != nil {
		resp := Message("error:"+err.Error())
		Respond(w, resp, http.StatusInternalServerError)
		return
	}
	resp := Message("success")
	resp["data"] = contacts
	Respond(w, resp, http.StatusOK)
}

func (ush *contactHTTPRouter) GetContact(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := uuid.Parse(params["id"])

	ush.log.LogMessage("Http request to get one contact with id ", id)

	if err != nil {
		Respond(w, Message("Error while decoding request body:"+err.Error()), http.StatusBadRequest)
		return
	}

	var Contact *ent.Contact
	Contact, err = ush.uc.FindByID(id)
	if err != nil {
		resp := Message("error:" + err.Error())
		Respond(w, resp, http.StatusInternalServerError)
		return
	}
	resp := Message("success")
	resp["data"] = &Contact
	Respond(w, resp, http.StatusOK)
}

func (ush *contactHTTPRouter) Find(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	q := params["query"]

	ush.log.LogMessage("Http request to find contacts with query", q)

	if len(q) < 3 {
		resp := Message("Too short query string")
		Respond(w, resp, http.StatusBadRequest)
		return
	}

	var contacts []*ent.Contact
	contacts, err := ush.uc.Find(q)
	if err != nil {
		resp := Message("error:" + err.Error())
		Respond(w, resp, http.StatusInternalServerError)
		return
	}
	resp := Message("success")
	resp["data"] = contacts
	Respond(w, resp, http.StatusOK)
}

func (ush *contactHTTPRouter) Store(w http.ResponseWriter, r *http.Request) {
	ush.log.LogMessage("Http request to make one contact")

	Contact := &ent.Contact{}
	err := json.NewDecoder(r.Body).Decode(Contact)
	if err != nil {
		Respond(w, Message("Error while decoding request body: "+err.Error()), http.StatusBadRequest)
		return
	}
	resp := Message("success")
	resp["data"], err = ush.uc.AddContact(*Contact)
	if err != nil {
		Respond(w, Message("Error while contact store: "+err.Error()), http.StatusInternalServerError)
		return
	}

	Respond(w, resp, http.StatusOK)
}

func (ush *contactHTTPRouter) Update(w http.ResponseWriter, r *http.Request) {
	ush.log.LogMessage("Http request to update one contact")

	Contact := &ent.Contact{}
	err := json.NewDecoder(r.Body).Decode(Contact)
	if err != nil {
		Respond(w, Message("Error while decoding request body: "+err.Error()), http.StatusBadRequest)
		return
	}
	resp := Message("success")
	resp["data"], err = ush.uc.UpdateContact(*Contact)
	if err != nil {
		Respond(w, Message("Error while contact store: "+err.Error()), http.StatusInternalServerError)
		return
	}

	Respond(w, resp, http.StatusOK)
}

func (ush *contactHTTPRouter) Delete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := uuid.Parse(params["id"])

	ush.log.LogMessage("Http request to delete one contact with id ", id)

	if err != nil {
		Respond(w, Message("Error while decoding request body:"+err.Error()), http.StatusBadRequest)
		return
	}

	err = ush.uc.DeleteContactByID(id)
	if err != nil {
		resp := Message("Error while delete contact:" + err.Error())
		Respond(w, resp, http.StatusNotFound)
		return
	}
	resp := Message("success")
	Respond(w, resp, http.StatusOK)
}
