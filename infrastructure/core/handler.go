package core

import "reflect"

type commandHandlerStruct struct {
	name, description string
	dataType          reflect.Type
	usecase           func(sc ServiceContext) error
}

func commandHandler(sc ServiceContext) error {
	return nil
}
