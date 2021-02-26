package user2

import (
	"database/sql/driver"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gopkg.in/guregu/null.v4"
)

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

func IDValidate(id string) error {
	validate := validator.New()
	return validate.Var(id, "required,uuid4")
}

func NewUserValidate(jui interface{}) error {
	validate := validator.New()
	validate.RegisterCustomTypeFunc(ValidateValuer, null.String{}, null.Int{}, null.Bool{}, null.Time{}, null.Float{})
	return validate.Struct(jui)
}

// ValidateValuer implements validator.CustomTypeFunc
func ValidateValuer(field reflect.Value) interface{} {
	if valuer, ok := field.Interface().(driver.Valuer); ok {
		val, err := valuer.Value()
		if err == nil {
			return val
		}
	}
	return nil
}
