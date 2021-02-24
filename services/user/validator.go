package user

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gopkg.in/guregu/null.v4"
)

type NewValidUser struct {
	ID       uuid.UUID   `validate:"required,uuid4"`
	Login    null.String `validate:"required,alphanum"`
	Password null.String `validate:"required,alphaunicode"`
	Email    null.String `validate:"required,email"`
}

func IDValidate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		validate := validator.New()
		errs := validate.Var(c.Param("id"), "required,uuid4")

		if errs != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Error while decoding ID: "+errs.Error())
		}
		return next(c)
	}
}

func NewUserValidate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		j := &jsonUser{}

		if err := c.Bind(j); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Error while decoding request body: "+err.Error())
		}

		validate := validator.New()
		errs := validate.Struct(j)

		if errs != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Error while decoding ID: "+errs.Error())
		}
		return next(c)
	}
}
