package core

type commandHandlerStruct struct {
	name, description string
	dataType          interface{}
	usecase           func(sc ServiceContext) error
}

func commandHandler(sc ServiceContext) error {
	return nil
}
