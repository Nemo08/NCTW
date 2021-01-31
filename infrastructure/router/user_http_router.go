package router

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	ent "github.com/Nemo08/NCTW/entity"
	log "github.com/Nemo08/NCTW/infrastructure/logger"
	use "github.com/Nemo08/NCTW/usecase"
)

type jsonUser struct {
	ID           uuid.UUID `json:"id"`
	Login        string    `json:"login"`
	Password     string    `json:",omitempty"`
	PasswordHash string    `json:"-"`
	Email        string    `json:"email"`
}

//json2user Json объект копируем в Entity
func json2user(i jsonUser) ent.User {
	return ent.User{
		ID:           i.ID,
		Login:        i.Login,
		PasswordHash: i.PasswordHash,
		Email:        i.Email,
	}
}

//user2json Entity объект копируем в Json
func user2json(i ent.User) jsonUser {
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
func NewUserHTTPRouter(l log.LogInterface, u use.UserUsecase, g *echo.Group) {
	l.LogMessage("Создаю роутер для user")

	us := userHTTPRouter{
		uc:  u,
		log: l,
	}

	subr := g.Group("/user")
	subr.POST("", us.Store)
	subr.GET("", us.GetAllUsers)
	subr.GET("/:id", us.GetUser)
	subr.GET("/search/:query", us.Find)
	subr.PUT("", us.Update)
	subr.DELETE("/:id", us.Delete)
}

func (ush *userHTTPRouter) GetAllUsers(c echo.Context) (err error) {
	ush.log.LogMessage("Запрошены все пользователи")

	var u []*ent.User
	var jsusers []*jsonUser

	u, err = ush.uc.GetAllUsers()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "error:"+err.Error())
	}

	for _, d := range u {
		j := user2json(*d)
		jsusers = append(jsusers, &j)
	}

	return c.JSON(http.StatusOK, jsusers)
}

func (ush *userHTTPRouter) GetUser(c echo.Context) (err error) {
	ush.log.LogMessage("Http request to get one user with id ", c.Param("id"))

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Error while decoding request body:"+err.Error())
	}

	var u *ent.User
	u, err = ush.uc.FindByID(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "error:"+err.Error())
	}

	j := user2json(*u)
	return c.JSON(http.StatusOK, j)
}

func (ush *userHTTPRouter) Find(c echo.Context) (err error) {
	ush.log.LogMessage("Http request to find users with query ", c.Param("query"))

	q := c.Param("query")
	if len(q) < 3 {
		return echo.NewHTTPError(http.StatusBadRequest, "Too short query string")
	}

	var u []*ent.User
	u, err = ush.uc.Find(q)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "error:"+err.Error())
	}

	var jsusers []*jsonUser
	for _, d := range u {
		j := user2json(*d)
		jsusers = append(jsusers, &j)
	}

	return c.JSON(http.StatusOK, jsusers)
}

func (ush *userHTTPRouter) Store(c echo.Context) (err error) {
	ush.log.LogMessage("Http request to make one user")

	j := &jsonUser{}

	if err = c.Bind(j); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Error while decoding request body: "+err.Error())
	}

	u, err := ent.NewUser(j.Login, j.Password, j.Email)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error while user store: "+err.Error())
	}

	u2, err := ush.uc.AddUser(u)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error while user store: "+err.Error())
	}

	*j = user2json(*u2)
	return c.JSON(http.StatusOK, j)
}

func (ush *userHTTPRouter) Update(c echo.Context) (err error) {
	ush.log.LogMessage("Http request to update one user")

	j := &jsonUser{}
	if err = c.Bind(j); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Error while decoding request body: "+err.Error())
	}

	u, err := ush.uc.UpdateUser(json2user(*j))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error while user store: "+err.Error())
	}

	*j = user2json(*u)
	return c.JSON(http.StatusOK, j)
}

func (ush *userHTTPRouter) Delete(c echo.Context) (err error) {
	ush.log.LogMessage("Http request to delete one user with id ", c.Param("id"))

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Error while decoding request body:"+err.Error())
	}

	err = ush.uc.DeleteUserByID(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error while delete user:"+err.Error())
	}

	return c.JSON(http.StatusOK, id)
}
