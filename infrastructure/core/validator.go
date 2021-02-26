package core

import (
	"database/sql/driver"
	"reflect"

	"github.com/go-playground/validator/v10"
	"gopkg.in/guregu/null.v4"
)

func dataValidate(sc ServiceContext, comm commandHandlerStruct) error {
	validate := validator.New()
	validate.RegisterCustomTypeFunc(validateValuer, null.String{}, null.Int{}, null.Bool{}, null.Time{}, null.Float{})
	return validate.Struct(sc.RequestData)
}

func validateValuer(field reflect.Value) interface{} {
	if valuer, ok := field.Interface().(driver.Valuer); ok {
		val, err := valuer.Value()
		if err == nil {
			return val
		}
	}
	return nil
}
