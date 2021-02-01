package router

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gopkg.in/guregu/null.v4"

	ent "github.com/Nemo08/NCTW/entity"
	log "github.com/Nemo08/NCTW/infrastructure/logger"
	use "github.com/Nemo08/NCTW/usecase"
)

type jsonContact struct {
	ID       uuid.UUID
	Position null.String //должность
}

//json2contact Json объект копируем в Entity
func json2contact(i jsonContact) ent.Contact {
	return ent.Contact{
		ID:       i.ID,
		Position: i.Position,
	}
}

//contact2json Entity объект копируем в Json
func contact2json(i ent.Contact) jsonContact {
	return jsonContact{
		ID:       i.ID,
		Position: i.Position,
	}
}

type contactHTTPRouter struct {
	uc  use.ContactUsecase
	log log.LogInterface
}

//NewContactHTTPRouter роутер для контактов
func NewContactHTTPRouter(l log.LogInterface, u use.ContactUsecase, g *echo.Group) {
	l.LogMessage("Создаю роутер для contact")

	var us contactHTTPRouter
	us.uc = u
	us.log = l

	subr := g.Group("/contact")
	subr.GET("", us.GetAllContacts)
	subr.POST("", us.Store)
	subr.GET("/:id", us.GetContact)
	subr.GET("/search/:query", us.Find)
	subr.PUT("", us.Update)
	subr.DELETE("/:id", us.Delete)
}

func (ush *contactHTTPRouter) GetAllContacts(c echo.Context) (err error) {
	ush.log.LogMessage("Http request to get all contacts")

	var u []*ent.Contact
	var jscontacts []*jsonContact

	u, err = ush.uc.GetAllContacts()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "error:"+err.Error())
	}

	for _, d := range u {
		j := contact2json(*d)
		jscontacts = append(jscontacts, &j)
	}

	return c.JSON(http.StatusOK, jscontacts)
}

func (ush *contactHTTPRouter) GetContact(c echo.Context) (err error) {
	ush.log.LogMessage("Http request to get one contact with id ", c.Param("id"))

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Error while decoding request body:"+err.Error())
	}

	var u *ent.Contact
	u, err = ush.uc.FindByID(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "error:"+err.Error())
	}

	j := contact2json(*u)
	return c.JSON(http.StatusOK, j)
}

func (ush *contactHTTPRouter) Find(c echo.Context) (err error) {
	ush.log.LogMessage("Http request to find contacts with query ", c.Param("query"))

	q := c.Param("query")
	if len(q) < 3 {
		return echo.NewHTTPError(http.StatusBadRequest, "Too short query string")
	}

	var u []*ent.Contact
	u, err = ush.uc.Find(q)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "error:"+err.Error())
	}

	var jscontacts []*jsonContact
	for _, d := range u {
		j := contact2json(*d)
		jscontacts = append(jscontacts, &j)
	}

	return c.JSON(http.StatusOK, jscontacts)
}

func (ush *contactHTTPRouter) Store(c echo.Context) (err error) {
	ush.log.LogMessage("Http request to make one contact")

	j := &jsonContact{}

	if err = c.Bind(j); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Error while decoding request body: "+err.Error())
	}

	u, err := ent.NewContact(j.Position)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error while contact store: "+err.Error())
	}

	u2, err := ush.uc.AddContact(u)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error while contact store: "+err.Error())
	}

	*j = contact2json(*u2)
	return c.JSON(http.StatusOK, j)
}

func (ush *contactHTTPRouter) Update(c echo.Context) (err error) {
	ush.log.LogMessage("Http request to update one contact")

	j := &jsonContact{}
	if err = c.Bind(j); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Error while decoding request body: "+err.Error())
	}

	u, err := ush.uc.UpdateContact(json2contact(*j))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error while contact store: "+err.Error())
	}

	*j = contact2json(*u)
	return c.JSON(http.StatusOK, j)
}

func (ush *contactHTTPRouter) Delete(c echo.Context) (err error) {
	ush.log.LogMessage("Http request to delete one contact with id ", c.Param("id"))

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Error while decoding request body:"+err.Error())
	}

	err = ush.uc.DeleteContactByID(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error while delete contact:"+err.Error())
	}

	return c.JSON(http.StatusOK, id)
}
