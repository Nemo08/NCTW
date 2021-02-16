package user

import (
	"log"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gopkg.in/guregu/null.v4"

	"github.com/Nemo08/NCTW/infrastructure/logger"
	"github.com/Nemo08/NCTW/services/api"
)

type CustomContext struct {
	echo.Context
	log log.Logger
}
type jsonUser struct {
	ID           uuid.UUID   `json:"id"`
	Login        null.String `json:"login"`
	Password     null.String `json:",omitempty"`
	PasswordHash null.String `json:"-"`
	Email        null.String `json:"email"`
}

//json2user Json объект копируем в Entity
func json2user(i jsonUser) User {
	return User{
		ID:           i.ID,
		Login:        i.Login,
		PasswordHash: i.PasswordHash,
		Email:        i.Email,
		Password:     i.Password,
	}
}

//user2json Entity объект копируем в Json
func user2json(i User) jsonUser {
	return jsonUser{
		ID:    i.ID,
		Login: i.Login,
		Email: i.Email,
	}
}

type userHTTPRouter struct {
	uc  Usecase
	log logger.Logr
}

//NewUserHTTPRouter роутер пользователей
func NewUserHTTPRouter(log logger.Logr, u Usecase, g *echo.Group) {
	log.Info("Создаю роутер для user")

	us := userHTTPRouter{
		uc:  u,
		log: log,
	}

	subr := g.Group("/user")
	subr.POST("", us.Store)
	subr.GET("", us.GetUsers)
	subr.GET("/:id", us.GetUser)
	subr.GET("/search/:query", us.Find)
	subr.PUT("", us.Update)
	subr.DELETE("/:id", us.Delete)
}

func (ush *userHTTPRouter) GetUsers(c echo.Context) (err error) {
	ush.log.WithContext(c.(api.Context)).Info("Запрошены пользователи постранично")
	var u []*User
	var jsusers []*jsonUser

	u, count, err := ush.uc.Get(c.(api.Context))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "error:"+err.Error())
	}

	for _, d := range u {
		j := user2json(*d)
		jsusers = append(jsusers, &j)
	}

	c.Response().Header().Set("X-Total-Count", strconv.Itoa(count))
	return c.JSON(http.StatusOK, jsusers)
}

func (ush *userHTTPRouter) GetUser(c echo.Context) (err error) {
	ush.log.WithContext(c.(api.Context)).Info("Http request to get one user with id ", c.Param("id"))

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Error while decoding request body:"+err.Error())
	}

	var u *User
	u, err = ush.uc.FindByID(c.(api.Context), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "error:"+err.Error())
	}

	j := user2json(*u)
	return c.JSON(http.StatusOK, j)
}

func (ush *userHTTPRouter) Find(c echo.Context) (err error) {
	ush.log.WithContext(c.(api.Context)).Info("Http request to find users with query ", c.Param("query"))

	q := c.Param("query")
	if len(q) < 3 {
		return echo.NewHTTPError(http.StatusBadRequest, "Too short query string")
	}

	var u []*User
	u, count, err := ush.uc.Find(c.(api.Context), q)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "error:"+err.Error())
	}

	var jsusers []*jsonUser
	for _, d := range u {
		j := user2json(*d)
		jsusers = append(jsusers, &j)
	}

	c.Response().Header().Set("X-Total-Count", strconv.Itoa(count))
	if len(jsusers) == 0 {
		return c.NoContent(http.StatusOK)
	}
	return c.JSON(http.StatusOK, jsusers)
}

func (ush *userHTTPRouter) Store(c echo.Context) (err error) {
	ush.log.WithContext(c.(api.Context)).Info("Запрос на сохранение одного пользователя")

	j := &jsonUser{}

	if err = c.Bind(j); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Error while decoding request body: "+err.Error())
	}

	u, err := NewUser(j.Login, j.Password, j.Email)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error while user store: "+err.Error())
	}

	u2, err := ush.uc.Add(c.(api.Context), u)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error while user store: "+err.Error())
	}

	*j = user2json(*u2)
	return c.JSON(http.StatusOK, j)
}

func (ush *userHTTPRouter) Update(c echo.Context) (err error) {
	ush.log.WithContext(c.(api.Context)).Info("Запрос на апдейт одного пользователя")

	j := &jsonUser{}
	if err = c.Bind(j); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Error while decoding request body: "+err.Error())
	}

	u, err := ush.uc.Update(c.(api.Context), json2user(*j))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error while user store: "+err.Error())
	}

	*j = user2json(*u)
	return c.JSON(http.StatusOK, j)
}

func (ush *userHTTPRouter) Delete(c echo.Context) (err error) {
	ush.log.WithContext(c.(api.Context)).Info("Http request to delete one user with id ", c.Param("id"))

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Error while decoding request body:"+err.Error())
	}

	err = ush.uc.DeleteByID(c.(api.Context), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error while delete user:"+err.Error())
	}

	return c.JSON(http.StatusOK, id)
}
