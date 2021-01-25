package router

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	ent "github.com/Nemo08/nctw/entity"
	use "github.com/Nemo08/nctw/usecase"
)

type contactHttpRouter struct {
	uc  use.ContactUsecase
	log use.LogInterface
}

func NewContactHttpRouter(l use.LogInterface, u use.ContactUsecase, r *mux.Router) {
	l.LogMessage("Set up contact router")

	var us contactHttpRouter
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

func (ush *contactHttpRouter) GetAllContacts(w http.ResponseWriter, r *http.Request) {
	ush.log.LogMessage("Http request to get all contacts")

	var contacts []*ent.Contact
	contacts, err := ush.uc.GetAllContacts()
	if err != nil {
		resp := Message(true, "error")
		Respond(w, resp)
		return
	}
	resp := Message(true, "success")
	resp["data"] = contacts
	Respond(w, resp)
}

func (ush *contactHttpRouter) GetContact(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := uuid.Parse(params["id"])

	ush.log.LogMessage("Http request to get one contact with id ", id)

	if err != nil {
		Respond(w, Message(false, "Error while decoding request body"))
		return
	}

	var Contact *ent.Contact
	Contact, err = ush.uc.FindById(id)
	if err != nil {
		resp := Message(true, "error")
		Respond(w, resp)
		return
	}
	resp := Message(true, "success")
	resp["data"] = &Contact
	Respond(w, resp)
}

func (ush *contactHttpRouter) Find(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	q := params["query"]

	ush.log.LogMessage("Http request to find contacts with query ", q)

	if len(q) < 3 {
		resp := Message(true, "Too short query string")
		Respond(w, resp)
		return
	}

	var contacts []*ent.Contact
	contacts, err := ush.uc.Find(q)
	if err != nil {
		resp := Message(true, "error")
		Respond(w, resp)
		return
	}
	resp := Message(true, "success")
	resp["data"] = contacts
	Respond(w, resp)
}

func (ush *contactHttpRouter) Store(w http.ResponseWriter, r *http.Request) {
	ush.log.LogMessage("Http request to make one contact")

	Contact := &ent.Contact{}
	err := json.NewDecoder(r.Body).Decode(Contact)
	if err != nil {
		Respond(w, Message(false, "Error while decoding request body: "+err.Error()))
		return
	}
	resp := Message(true, "success")
	resp["data"], err = ush.uc.AddContact(*Contact)
	if err != nil {
		Respond(w, Message(false, "Error while contact store: "+err.Error()))
		return
	}

	Respond(w, resp)
}

func (ush *contactHttpRouter) Update(w http.ResponseWriter, r *http.Request) {
	ush.log.LogMessage("Http request to update one contact")

	Contact := &ent.Contact{}
	err := json.NewDecoder(r.Body).Decode(Contact)
	if err != nil {
		Respond(w, Message(false, "Error while decoding request body: "+err.Error()))
		return
	}
	resp := Message(true, "success")
	resp["data"], err = ush.uc.UpdateContact(*Contact)
	if err != nil {
		Respond(w, Message(false, "Error while contact store: "+err.Error()))
		return
	}

	Respond(w, resp)
}

func (ush *contactHttpRouter) Delete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := uuid.Parse(params["id"])

	ush.log.LogMessage("Http request to delete one contact with id ", id)

	if err != nil {
		Respond(w, Message(false, "Error while decoding request body"))
		return
	}

	err = ush.uc.DeleteContactById(id)
	if err != nil {
		resp := Message(true, "Error while delete contact")
		Respond(w, resp)
		return
	}
	resp := Message(true, "success")
	Respond(w, resp)
}
