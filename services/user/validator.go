package user

import (
	"database/sql/driver"
	"reflect"

	"github.com/go-playground/validator/v10"
	"gopkg.in/guregu/null.v4"
)

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
