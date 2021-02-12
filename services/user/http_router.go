package user

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gopkg.in/guregu/null.v4"

	log "github.com/Nemo08/NCTW/infrastructure/logger"
)

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
	uc UserUsecase
}

//NewUserHTTPRouter роутер пользователей
func NewUserHTTPRouter(u UserUsecase, g *echo.Group) {
	log.LogMessage("Создаю роутер для user")

	us := userHTTPRouter{
		uc: u,
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
	log.LogMessage("Запрошены пользователи постранично")

	var u []*User
	var jsusers []*jsonUser

	limit, offset, err := ush.findLimits(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "error:"+err.Error())
	}

	u, count, err := ush.uc.GetUsers(limit, offset)
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
	log.LogMessage("Http request to get one user with id ", c.Param("id"))

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Error while decoding request body:"+err.Error())
	}

	var u *User
	u, err = ush.uc.FindByID(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "error:"+err.Error())
	}

	j := user2json(*u)
	return c.JSON(http.StatusOK, j)
}

func (ush *userHTTPRouter) Find(c echo.Context) (err error) {
	log.LogMessage("Http request to find users with query ", c.Param("query"))

	q := c.Param("query")
	if len(q) < 3 {
		return echo.NewHTTPError(http.StatusBadRequest, "Too short query string")
	}

	limit, offset, err := ush.findLimits(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "error:"+err.Error())
	}

	var u []*User
	u, count, err := ush.uc.Find(q, limit, offset)
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
	log.LogMessage("Http request to make one user")

	j := &jsonUser{}

	if err = c.Bind(j); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Error while decoding request body: "+err.Error())
	}

	u, err := NewUser(j.Login, j.Password, j.Email)
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
	log.LogMessage("Http request to update one user")

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
	log.LogMessage("Http request to delete one user with id ", c.Param("id"))

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

func (ush *userHTTPRouter) findLimits(c echo.Context) (l int, o int, e error) {
	if (c.QueryParam("limit") == "") || (c.QueryParam("offset") == "") {
		return 0, 0, errors.New("error: Not set limit or offset")
	}
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		return 0, 0, err
	}
	offset, err := strconv.Atoi(c.QueryParam("offset"))
	if err != nil {
		return 0, 0, err
	}
	return limit, offset, nil
}
