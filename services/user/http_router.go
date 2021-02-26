package user

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
	"gopkg.in/guregu/null.v4"

	"github.com/Nemo08/NCTW/infrastructure/logger"
	"github.com/Nemo08/NCTW/services/api"
)

//jsonUserInput структура входящих данных
type jsonUserInput struct {
	ID       uuid.UUID   `json:"id" validate:"omitempty,uuid_rfc4122"`
	Login    null.String `json:"login" validate:"required,ascii"`
	Password null.String `json:",omitempty" validate:"required,ascii,lt=101,gt=4"`
	Email    null.String `json:"email" validate:"required,email"`
}

//jsonUserOutput структура исходящих данных
type jsonUserOutput struct {
	ID    uuid.UUID   `json:"id"`
	Login null.String `json:"login"`
	Email null.String `json:"email"`
}

type jsonUserUpdate struct {
	ID       string      `json:"id" validate:"required,uuid4"`
	Login    null.String `json:"login" validate:"omitempty,ascii"`
	Password null.String `json:",omitempty" validate:"omitempty,ascii,lt=101,gt=4"`
	Email    null.String `json:"email" validate:"omitempty,email"`
}

type userHTTPRouter struct {
	uc Usecase
}

//NewUserHTTPRouter роутер пользователей
func NewUserHTTPRouter(log logger.Logr, u Usecase, g *echo.Group) {
	log.Info("Создаю роутер для user")

	us := userHTTPRouter{
		uc: u,
	}

	subr := g.Group("/user")
	subr.POST("", us.Store)
	subr.GET("", us.GetUsers)
	subr.GET("/:id", us.GetUserByID)
	subr.GET("/search/:query", us.Find)
	subr.PUT("", us.Update)
	subr.DELETE("/:id", us.DeleteByID)
}

//GetUsers получаем пользователей
//Пагинация пробрасывается к базе в контексте
func (ush *userHTTPRouter) GetUsers(c echo.Context) (err error) {
	c.(api.Context).Log.Info("Запрошены пользователи постранично")

	//Нечего валидировать
	var u []*User
	var juo []*jsonUserOutput
	//Получаем пользователей
	u, err = ush.uc.Get(c.(api.Context))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "error:"+err.Error())
	}
	//Копируем данные в ответ с преобразованием
	if copier.Copy(&juo, &u) != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "error:"+err.Error())
	}
	//При отсутствии результатов отправляем "нет контента"
	if len(juo) == 0 {
		return c.NoContent(http.StatusOK)
	}
	//Передаем в ответ в заголовке количество возвращаемых пользователей
	c.Response().Header().Set("X-Total-Count", strconv.FormatInt(int64(len(u)), 10))
	return c.JSON(http.StatusOK, juo)
}

//GetUserByID получаем пользователя по ID
func (ush *userHTTPRouter) GetUserByID(c echo.Context) (err error) {
	c.(api.Context).Log.Info("Http request to get one user with id ", c.Param("id"))

	//Валидация ID
	err = IDValidate(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Validation error:"+err.Error())
	}
	//Делаем UUID из запроса
	id, _ := uuid.Parse(c.Param("id"))
	//Ищем
	user, err := ush.uc.FindByID(c.(api.Context), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "error:"+err.Error())
	}
	//При отсутствии результатов отправляем "нет контента"
	if user == nil {
		return c.NoContent(http.StatusOK)
	}
	//Копируем данные в ответ с преобразованием
	juo := jsonUserOutput{}
	if copier.Copy(&juo, &user) != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "error:"+err.Error())
	}
	//Передаем в ответ в заголовке количество возвращаемых пользователей
	c.Response().Header().Set("X-Total-Count", "1")
	return c.JSON(http.StatusOK, juo)
}

func (ush *userHTTPRouter) Find(c echo.Context) (err error) {
	c.(api.Context).Log.Info("Http request to find users with query ", c.Param("query"))

	//Достаем запрос
	q := c.Param("query")
	q, err = url.QueryUnescape(q)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Query unescape error")
	}
	//Валидация
	if len(q) < 3 {
		return echo.NewHTTPError(http.StatusBadRequest, "Too short query string")
	}
	//Ищем
	var users []*User
	users, err = ush.uc.Find(c.(api.Context), q)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "error:"+err.Error())
	}
	//При отсутствии результатов отправляем "нет контента"
	if len(users) == 0 {
		return c.NoContent(http.StatusOK)
	}
	//Копируем данные в ответ с преобразованием
	var juo []*jsonUserOutput
	if copier.Copy(&juo, &users) != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "error:"+err.Error())
	}
	//При отсутствии результатов отправляем "нет контента"
	if len(juo) == 0 {
		return c.NoContent(http.StatusOK)
	}
	//Пишем в ответ количество записей данных
	c.Response().Header().Set("X-Total-Count", strconv.FormatInt(int64(len(juo)), 10))
	return c.JSON(http.StatusOK, juo)
}

func (ush *userHTTPRouter) Store(c echo.Context) (err error) {
	c.(api.Context).Log.Info("Запрос на сохранение одного пользователя")
	//Достаем данные из запроса
	jui := &jsonUserInput{}
	if err = c.Bind(jui); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Error while decoding request body: "+err.Error())
	}
	//Валидация данных
	err = NewUserValidate(*jui)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Validation error:"+err.Error())
	}
	//Копируем данные из запроса с преобразованием
	user := User{}
	if copier.Copy(&user, &jui) != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "error:"+err.Error())
	}
	//Генерируем хэш пароля
	passwordHash, err := CreateHash(jui.Password.String)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error while user store: "+err.Error())
	}
	user.PasswordHash = null.StringFrom(passwordHash)
	//Пишем в базу, получаем с ИДом из базы
	user2, err := ush.uc.Store(c.(api.Context), user)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error while user store: "+err.Error())
	}
	//Копируем данные в ответ с преобразованием
	juo := jsonUserOutput{}
	if copier.Copy(&juo, &user2) != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error while user store:"+err.Error())
	}
	//Пишем в ответ количество записей данных
	c.Response().Header().Set("X-Total-Count", "1")
	return c.JSON(http.StatusOK, juo)
}

func (ush *userHTTPRouter) Update(c echo.Context) (err error) {
	c.(api.Context).Log.Info("Запрос на апдейт одного пользователя")
	//Достаем данные из запроса
	juu := jsonUserUpdate{}
	if err = c.Bind(&juu); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Error while decoding request body: "+err.Error())
	}
	//Валидация данных
	err = NewUserValidate(juu)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Validation error:"+err.Error())
	}
	//Копируем данные из запроса с преобразованием
	user := User{}
	if copier.Copy(&user, &juu) != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "error:"+err.Error())
	}
	//Если пароль не пустой - генерируем хэш пароля
	if !juu.Password.IsZero() {
		passwordHash, err := CreateHash(juu.Password.String)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Error while user store: "+err.Error())
		}
		user.PasswordHash = null.StringFrom(passwordHash)
	}
	//Сохраняем в базу и получаем обновленный
	user2, err := ush.uc.Update(c.(api.Context), user)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error while user store: "+err.Error())
	}
	//Копируем данные в ответ с преобразованием
	juo := jsonUserOutput{}
	if copier.Copy(&juo, &user2) != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "error:"+err.Error())
	}
	//Передаем в ответ в заголовке количество возвращаемых пользователей
	c.Response().Header().Set("X-Total-Count", "1")
	return c.JSON(http.StatusOK, juo)
}

//DeleteByID удаляет пользователя по ID
func (ush *userHTTPRouter) DeleteByID(c echo.Context) (err error) {
	c.(api.Context).Log.Info("Http request to delete one user with id ", c.Param("id"))
	//Валидация ID
	err = IDValidate(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "error:"+err.Error())
	}
	//Делаем UUID из запроса
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Error while decoding request body:"+err.Error())
	}
	//Удаляем по ID
	err = ush.uc.DeleteByID(c.(api.Context), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error while delete user:"+err.Error())
	}
	//Передаем в ответ в заголовке количество возвращаемых пользователей
	c.Response().Header().Set("X-Total-Count", "0")
	return c.JSON(http.StatusOK, id)
}
